package vtsgo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type VTSGo struct {
	websocketConn *websocket.Conn
	tokenPath     string
	pluginInfo    *PluginInfo

	mu      sync.RWMutex
	pending map[string]chan []byte // maps requestID to raw JSON response bytes
}

func New(tokenPath string, pluginInfo PluginInfo) *VTSGo {
	return &VTSGo{
		pending:    make(map[string]chan []byte),
		tokenPath:  tokenPath,
		pluginInfo: &pluginInfo,
	}
}

func (c *VTSGo) Connect(ctx context.Context, addr string) error {
	if addr == "" {
		addr = "ws://localhost:8001"
	}

	conn, _, err := websocket.DefaultDialer.DialContext(ctx, addr, nil)
	if err != nil {
		return fmt.Errorf("vts connect error: %w", err)
	}
	c.websocketConn = conn

	go c.readLoop()
	return nil
}

func (c *VTSGo) readLoop() {
	defer c.websocketConn.Close()

	for {
		_, rawMsg, err := c.websocketConn.ReadMessage()
		if err != nil {
			return
		}

		//log.Printf("DEBUG: %v", string(rawMsg))

		var env RawEnvelope
		if err := json.Unmarshal(rawMsg, &env); err != nil {
			continue
		}

		c.mu.Lock()
		ch, exists := c.pending[env.RequestID]
		if exists {
			delete(c.pending, env.RequestID)
		}
		c.mu.Unlock()

		if exists && ch != nil {
			ch <- rawMsg
		}
	}
}

func (c *VTSGo) Request(ctx context.Context, msgType string, payload any, responseData any) error {
	reqID := uuid.NewString()

	dataBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	msg := map[string]any{
		"apiName":     "VTubeStudioPublicAPI",
		"apiVersion":  "1.0",
		"requestID":   reqID,
		"messageType": msgType,
		"data":        json.RawMessage(dataBytes),
	}

	// 1. Fetch the full raw JSON response via RequestRaw
	var rawResp json.RawMessage
	if err := c.RequestRaw(ctx, reqID, msg, &rawResp); err != nil {
		return err
	}

	// 2. Unmarshal the raw envelope to inspect messageType and extract data
	var env RawEnvelope
	if err := json.Unmarshal(rawResp, &env); err != nil {
		return fmt.Errorf("failed to parse response envelope: %w", err)
	}

	// 3. Handle VTS error responses centrally
	if env.MessageType == "APIError" {
		var vtsErr VTSError
		if err := json.Unmarshal(env.Data, &vtsErr); err == nil && vtsErr.ErrorID != 0 {
			return &vtsErr
		}
		return fmt.Errorf("vts returned APIError: %s", string(env.Data))
	}

	// 4. Unmarshal ONLY the inner "data" payload into responseData
	if responseData != nil && len(env.Data) > 0 {
		if err := json.Unmarshal(env.Data, responseData); err != nil {
			return fmt.Errorf("failed to unmarshal data payload: %w", err)
		}
	}

	return nil
}

func (c *VTSGo) RequestRaw(ctx context.Context, reqID string, payload any, responseData any) error {
	if c.pending == nil {
		return errors.New("client pending map is uninitialized")
	}

	resChan := make(chan []byte, 1)

	c.mu.Lock()
	c.pending[reqID] = resChan
	c.mu.Unlock()

	if err := c.websocketConn.WriteJSON(payload); err != nil {
		c.mu.Lock()
		delete(c.pending, reqID)
		c.mu.Unlock()
		return err
	}

	select {
	case <-ctx.Done():
		c.mu.Lock()
		delete(c.pending, reqID)
		c.mu.Unlock()
		return ctx.Err()
	case rawResp := <-resChan:
		if responseData != nil {
			// Unmarshal raw WS response directly into target struct
			return json.Unmarshal(rawResp, responseData)
		}
		return nil
	}
}

func (c *VTSGo) GetSessionStatus(ctx context.Context) (*SessionStatusResponse, error) {
	reqID := uuid.NewString()

	req := RequestSessionStatus{
		ApiName:     "VTubeStudioPublicAPI",
		ApiVersion:  "1.0",
		RequestID:   reqID,
		MessageType: "APIStateRequest",
	}

	var resp SessionStatusResponse
	if err := c.RequestRaw(ctx, reqID, req, &resp); err != nil {
		return nil, fmt.Errorf("failed APIStateRequest: %w", err)
	}

	return &resp, nil
}

func (c *VTSGo) loadToken() (string, error) {
	data, err := os.ReadFile(c.tokenPath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (c *VTSGo) saveToken(token string) error {
	dir := filepath.Dir(c.tokenPath)
	if dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	return os.WriteFile(c.tokenPath, []byte(token), 0600)
}

func (c *VTSGo) RequestAuthToken(ctx context.Context) (string, error) {
	req := TokenRequestPayload{
		PluginName:      c.pluginInfo.PluginName,
		PluginDeveloper: c.pluginInfo.PluginDeveloper,
		PluginIcon:      c.pluginInfo.PluginIcon,
	}

	var resp TokenResponsePayload
	if err := c.Request(ctx, "AuthenticationTokenRequest", req, &resp); err != nil {
		return "", err
	}

	// Check if VTS returned an error code
	if resp.ErrorId != 0 {
		return "", fmt.Errorf("vts error [%d]: %s", resp.ErrorId, resp.Message)
	}

	if resp.AuthenticationToken == "" {
		return "", errors.New("vts returned an empty authentication token")
	}

	return resp.AuthenticationToken, nil
}

func (c *VTSGo) Authenticate(ctx context.Context) error {
	if c.tokenPath == "" {
		c.tokenPath = "vts_token.txt"
	}

	token, err := c.loadToken()

	if err != nil || token == "" {
		// Token request can take longer because the user has to click "Allow" in VTS.
		// Use a 30-second deadline if none is present on ctx.
		reqCtx := ctx
		if _, hasDeadline := ctx.Deadline(); !hasDeadline {
			var cancel context.CancelFunc
			reqCtx, cancel = context.WithTimeout(ctx, 30*time.Second)
			defer cancel()
		}

		token, err = c.RequestAuthToken(reqCtx)
		if err != nil {
			return fmt.Errorf("failed to obtain auth token: %w", err)
		}

		if err := c.saveToken(token); err != nil {
			fmt.Printf("warning: failed to save auth token to disk: %v\n", err)
		}
	}

	// 3. Authenticate the current WebSocket session using the token
	authReq := AuthPayload{
		PluginName:      c.pluginInfo.PluginName,
		PluginDeveloper: c.pluginInfo.PluginDeveloper,
		AuthToken:       token,
	}

	var authResp AuthResponsePayload
	if err := c.Request(ctx, "AuthenticationRequest", authReq, &authResp); err != nil {
		return fmt.Errorf("authentication request failed: %w", err)
	}

	if !authResp.Authenticated {
		_ = os.Remove(c.tokenPath)
		return fmt.Errorf("vts rejected authentication: %s", authResp.Reason)
	}

	return nil
}

//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------

func (c *VTSGo) GetCurrentModel(ctx context.Context) (*CurrentModelResponse, error) {
	var resp CurrentModelResponse
	err := c.Request(ctx, "CurrentModelRequest", nil, &resp)
	return &resp, err
}

func (c *VTSGo) GetAvailableModels(ctx context.Context) (*AvailableModelsResponse, error) {
	var resp AvailableModelsResponse
	err := c.Request(ctx, "AvailableModelsRequest", nil, &resp)
	return &resp, err
}

func (c *VTSGo) LoadModel(ctx context.Context, req ModelLoadRequest) (*ModelLoadResponse, error) {
	var resp ModelLoadResponse
	err := c.Request(ctx, "ModelLoadRequest", req, &resp)
	return &resp, err
}

func (c *VTSGo) MoveModel(ctx context.Context, req MoveModelRequest) error {
	return c.Request(ctx, "MoveModelRequest", req, nil)
}

//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------

func (c *VTSGo) GetHotkeys(ctx context.Context, req HotkeysInCurrentModelRequest) (*HotkeysInCurrentModelResponse, error) {
	var resp HotkeysInCurrentModelResponse
	err := c.Request(ctx, "HotkeysInCurrentModelRequest", req, &resp)
	return &resp, err
}

func (c *VTSGo) TriggerHotkey(ctx context.Context, hotkeyID string) (*HotkeyTriggerResponse, error) {
	var resp HotkeyTriggerResponse
	req := HotkeyTriggerRequest{HotkeyID: hotkeyID}
	err := c.Request(ctx, "HotkeyTriggerRequest", req, &resp)
	return &resp, err
}

func (c *VTSGo) GetExpressionState(ctx context.Context, req ExpressionStateRequest) (*ExpressionStateResponse, error) {
	var resp ExpressionStateResponse
	err := c.Request(ctx, "ExpressionStateRequest", req, &resp)
	return &resp, err
}

func (c *VTSGo) SetExpressionState(ctx context.Context, req ExpressionActivationRequest) error {
	return c.Request(ctx, "ExpressionActivationRequest", req, nil)
}

//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------

func (c *VTSGo) GetArtMeshList(ctx context.Context) (*ArtMeshListResponse, error) {
	var resp ArtMeshListResponse
	err := c.Request(ctx, "ArtMeshListRequest", nil, &resp)
	return &resp, err
}

func (c *VTSGo) GetArtMeshesAtPosition(ctx context.Context, req ArtMeshAtPositionRequest) (*ArtMeshAtPositionResponse, error) {
	var resp ArtMeshAtPositionResponse
	err := c.Request(ctx, "ArtMeshAtPositionRequest", req, &resp)
	return &resp, err
}

func (c *VTSGo) SetColorTint(ctx context.Context, req ColorTintRequest) (*ColorTintResponse, error) {
	var resp ColorTintResponse
	err := c.Request(ctx, "ColorTintRequest", req, &resp)
	return &resp, err
}

func (c *VTSGo) RequestArtMeshSelection(ctx context.Context, req ArtMeshSelectionRequest) (*ArtMeshSelectionResponse, error) {
	var resp ArtMeshSelectionResponse
	err := c.Request(ctx, "ArtMeshSelectionRequest", req, &resp)
	return &resp, err
}

//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------

func (c *VTSGo) GetInputParameterList(ctx context.Context) (*InputParameterListResponse, error) {
	var resp InputParameterListResponse
	err := c.Request(ctx, "InputParameterListRequest", nil, &resp)
	return &resp, err
}

func (c *VTSGo) GetLive2DParameterList(ctx context.Context) (*Live2DParameterListResponse, error) {
	var resp Live2DParameterListResponse
	err := c.Request(ctx, "Live2DParameterListRequest", nil, &resp)
	return &resp, err
}

func (c *VTSGo) GetParameterValue(ctx context.Context, name string) (*ParameterValueResponse, error) {
	var resp ParameterValueResponse
	req := ParameterValueRequest{Name: name}
	err := c.Request(ctx, "ParameterValueRequest", req, &resp)
	return &resp, err
}

func (c *VTSGo) CreateCustomParameter(ctx context.Context, req ParameterCreationRequest) (*ParameterCreationResponse, error) {
	var resp ParameterCreationResponse
	err := c.Request(ctx, "ParameterCreationRequest", req, &resp)
	return &resp, err
}

func (c *VTSGo) DeleteCustomParameter(ctx context.Context, name string) (*ParameterDeletionResponse, error) {
	var resp ParameterDeletionResponse
	req := ParameterDeletionRequest{ParameterName: name}
	err := c.Request(ctx, "ParameterDeletionRequest", req, &resp)
	return &resp, err
}

func (c *VTSGo) InjectParameterData(ctx context.Context, req InjectParameterDataRequest) error {
	return c.Request(ctx, "InjectParameterDataRequest", req, nil)
}

//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------

func (c *VTSGo) GetCurrentModelPhysics(ctx context.Context) (*GetCurrentModelPhysicsResponse, error) {
	var resp GetCurrentModelPhysicsResponse
	err := c.Request(ctx, "GetCurrentModelPhysicsRequest", nil, &resp)
	return &resp, err
}

func (c *VTSGo) SetCurrentModelPhysics(ctx context.Context, req SetCurrentModelPhysicsRequest) error {
	return c.Request(ctx, "SetCurrentModelPhysicsRequest", req, nil)
}

//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------

func (c *VTSGo) GetItemList(ctx context.Context, req ItemListRequest) (*ItemListResponse, error) {
	var resp ItemListResponse
	err := c.Request(ctx, "ItemListRequest", req, &resp)
	return &resp, err
}

func (c *VTSGo) LoadItem(ctx context.Context, req ItemLoadRequest) (*ItemLoadResponse, error) {
	var resp ItemLoadResponse
	err := c.Request(ctx, "ItemLoadRequest", req, &resp)
	return &resp, err
}

func (c *VTSGo) UnloadItem(ctx context.Context, req ItemUnloadRequest) (*ItemUnloadResponse, error) {
	var resp ItemUnloadResponse
	err := c.Request(ctx, "ItemUnloadRequest", req, &resp)
	return &resp, err
}

func (c *VTSGo) MoveItem(ctx context.Context, req ItemMoveRequest) (*ItemMoveResponse, error) {
	var resp ItemMoveResponse
	err := c.Request(ctx, "ItemMoveRequest", req, &resp)
	return &resp, err
}

func (c *VTSGo) ControlItemAnimation(ctx context.Context, req ItemAnimationControlRequest) (*ItemAnimationControlResponse, error) {
	var resp ItemAnimationControlResponse
	err := c.Request(ctx, "ItemAnimationControlRequest", req, &resp)
	return &resp, err
}

func (c *VTSGo) SortItem(ctx context.Context, req ItemSortRequest) (*ItemSortResponse, error) {
	var resp ItemSortResponse
	err := c.Request(ctx, "ItemSortRequest", req, &resp)
	return &resp, err
}

func (c *VTSGo) PinItem(ctx context.Context, req ItemPinRequest) (*ItemPinResponse, error) {
	var resp ItemPinResponse
	err := c.Request(ctx, "ItemPinRequest", req, &resp)
	return &resp, err
}

//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------

func (c *VTSGo) GetStatistics(ctx context.Context) (*StatisticsResponse, error) {
	var resp StatisticsResponse
	err := c.Request(ctx, "StatisticsRequest", nil, &resp)
	return &resp, err
}

func (c *VTSGo) GetVTSFolderInfo(ctx context.Context) (*VTSFolderInfoResponse, error) {
	var resp VTSFolderInfoResponse
	err := c.Request(ctx, "VTSFolderInfoRequest", nil, &resp)
	return &resp, err
}

func (c *VTSGo) GetSceneColorOverlayInfo(ctx context.Context) (*SceneColorOverlayInfoResponse, error) {
	var resp SceneColorOverlayInfoResponse
	err := c.Request(ctx, "SceneColorOverlayInfoRequest", nil, &resp)
	return &resp, err
}

func (c *VTSGo) CheckFaceFound(ctx context.Context) (*FaceFoundResponse, error) {
	var resp FaceFoundResponse
	err := c.Request(ctx, "FaceFoundRequest", nil, &resp)
	return &resp, err
}
