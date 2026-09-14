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

// RawEnvelope allows extracting metadata (like requestId and messageType)
// without breaking when VTS sends non-standard or error payloads.
type RawEnvelope struct {
	APIName     string          `json:"apiName"`
	APIVersion  string          `json:"apiVersion"`
	Timestamp   int64           `json:"timestamp"`
	MessageType string          `json:"messageType"`
	RequestID   string          `json:"requestID"`
	Data        json.RawMessage `json:"data"`
}

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

		// Dispatch raw message byte slice based on requestID
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
