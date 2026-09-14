package vtsgo

import (
	"encoding/json"
	"fmt"
)

// BaseMessage matches the standard envelope required by VTS JSON-RPC.
type BaseMessage struct {
	APIName     string          `json:"apiName"`    // "VTubeStudioPublicAPI"
	APIVersion  string          `json:"apiVersion"` // "1.0"
	RequestID   string          `json:"requestId"`  // Unique correlation ID
	MessageType string          `json:"messageType"`
	Data        json.RawMessage `json:"data,omitempty"`
}

// PluginInfo metadata sent during auth requests
type PluginInfo struct {
	PluginName      string `json:"pluginName"`
	PluginDeveloper string `json:"pluginDeveloper"`
	PluginIcon      string `json:"pluginIcon,omitempty"`
}

// TokenRequestPayload asks VTS to prompt the user for a new token
type TokenRequestPayload struct {
	PluginName      string `json:"pluginName"`
	PluginDeveloper string `json:"pluginDeveloper"`
	PluginIcon      string `json:"pluginIcon,omitempty"`
}

// TokenResponsePayload contains the generated token from VTS
type TokenResponsePayload struct {
	AuthenticationToken string `json:"authenticationToken"`
	ErrorId             int    `json:"errorID"`
	Message             string `json:"message"`
}

// AuthPayload presents the token to VTS to authorize the session
type AuthPayload struct {
	PluginName      string `json:"pluginName"`
	PluginDeveloper string `json:"pluginDeveloper"`
	AuthToken       string `json:"authenticationToken"`
}

// AuthResponsePayload contains the result of the auth attempt
type AuthResponsePayload struct {
	Authenticated bool   `json:"authenticated"`
	Reason        string `json:"reason"`
}

type RequestSessionStatus struct {
	ApiName     string `json:"apiName"`
	ApiVersion  string `json:"apiVersion"`
	RequestID   string `json:"requestID"`
	MessageType string `json:"messageType"`
}

type SessionStatusResponse struct {
	ApiName     string `json:"apiName"`
	ApiVersion  string `json:"apiVersion"`
	Timestamp   int64  `json:"timestamp"`
	MessageType string `json:"messageType"`
	RequestID   string `json:"requestID"`
	Data        struct {
		Active                      bool   `json:"active"`
		VTubeStudioVersion          string `json:"vTubeStudioVersion"`
		CurrentSessionAuthenticated bool   `json:"currentSessionAuthenticated"`
	} `json:"data"`
}

type VTSError struct {
	ErrorID int    `json:"errorID"`
	Message string `json:"message"`
}

func (e *VTSError) Error() string {
	return fmt.Sprintf("vts api error [%d]: %s", e.ErrorID, e.Message)
}

type StatisticsResponse struct {
	Uptime              int    `json:"uptime"`
	Framerate           int    `json:"framerate"`
	VTubeStudioVersion  string `json:"vTubeStudioVersion"`
	AllowedPlugins      int    `json:"allowedPlugins"`
	ConnectedPlugins    int    `json:"connectedPlugins"`
	StartedWithSteam    bool   `json:"startedWithSteam"`
	WindowWidth         int    `json:"windowWidth"`
	WindowHeight        int    `json:"windowHeight"`
	WindowsIsFullscreen bool   `json:"windowIsFullscreen"`
}

type VTSFolderInfoResponse struct {
	Models      string `json:"models"`
	Backgrounds string `json:"backgrounds"`
	Items       string `json:"items"`
	Config      string `json:"config"`
	Logs        string `json:"logs"`
	Backup      string `json:"backup"`
}

type ModelPosition struct {
	PositionX float64 `json:"positionX"`
	PositionY float64 `json:"positionY"`
	Rotation  float64 `json:"rotation"`
	Size      float64 `json:"size"`
}

type CurrentModelResponse struct {
	ModelLoaded              bool          `json:"modelLoaded"`
	ModelName                string        `json:"modelName"`
	ModelID                  string        `json:"modelID"`
	VTSModelName             string        `json:"vtsModelName"`
	VTSModelIconName         string        `json:"vtsModelIconName"`
	Live2DModelName          string        `json:"live2DModelName"`
	ModelLoadTime            int64         `json:"modelLoadTime"`
	TimeSinceModelLoaded     int64         `json:"timeSinceModelLoaded"`
	NumberOfLive2DParameters int           `json:"numberOfLive2DParameters"`
	NumberOfLive2DArtmeshes  int           `json:"numberOfLive2DArtmeshes"`
	HasPhysicsFile           bool          `json:"hasPhysicsFile"`
	NumberOfTextures         int           `json:"numberOfTextures"`
	TextureResolution        int           `json:"textureResolution"`
	ModelPosition            ModelPosition `json:"modelPosition"`
}

type AvailableModel struct {
	ModelLoaded      bool   `json:"modelLoaded"`
	ModelName        string `json:"modelName"`
	ModelID          string `json:"modelID"`
	VTSModelName     string `json:"vtsModelName"`
	VTSModelIconName string `json:"vtsModelIconName"`
}

type AvailableModelsResponse struct {
	NumberOfModels  int              `json:"numberOfModels"`
	AvailableModels []AvailableModel `json:"availableModels"`
}
