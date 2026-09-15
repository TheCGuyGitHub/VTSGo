package vtsgo

import (
	"encoding/json"
	"fmt"
)

type BaseMessage struct {
	APIName     string          `json:"apiName"`    // "VTubeStudioPublicAPI"
	APIVersion  string          `json:"apiVersion"` // "1.0"
	RequestID   string          `json:"requestId"`  // Unique correlation ID
	MessageType string          `json:"messageType"`
	Data        json.RawMessage `json:"data,omitempty"`
}

type PluginInfo struct {
	PluginName      string `json:"pluginName"`
	PluginDeveloper string `json:"pluginDeveloper"`
	PluginIcon      string `json:"pluginIcon,omitempty"`
}

//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------

type TokenRequestPayload struct {
	PluginName      string `json:"pluginName"`
	PluginDeveloper string `json:"pluginDeveloper"`
	PluginIcon      string `json:"pluginIcon,omitempty"`
}

type TokenResponsePayload struct {
	AuthenticationToken string `json:"authenticationToken"`
	ErrorId             int    `json:"errorID"`
	Message             string `json:"message"`
}

//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------

type AuthPayload struct {
	PluginName      string `json:"pluginName"`
	PluginDeveloper string `json:"pluginDeveloper"`
	AuthToken       string `json:"authenticationToken"`
}

type AuthResponsePayload struct {
	Authenticated bool   `json:"authenticated"`
	Reason        string `json:"reason"`
}

//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------

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

//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------

type VTSError struct {
	ErrorID int    `json:"errorID"`
	Message string `json:"message"`
}

func (e *VTSError) Error() string {
	return fmt.Sprintf("vts api error [%d]: %s", e.ErrorID, e.Message)
}

//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------

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

//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------

type VTSFolderInfoResponse struct {
	Models      string `json:"models"`
	Backgrounds string `json:"backgrounds"`
	Items       string `json:"items"`
	Config      string `json:"config"`
	Logs        string `json:"logs"`
	Backup      string `json:"backup"`
}

//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------

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

//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------

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

//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------

type ModelLoadRequest struct {
	ModelID string `json:"modelID"`
}

type ModelLoadResponse struct {
	ModelID string `json:"modelID"`
}

//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------

type MoveModelRequest struct {
	TimeInSeconds            float64  `json:"timeInSeconds"`
	ValuesAreRelativeToModel bool     `json:"valuesAreRelativeToModel"`
	PositionX                *float64 `json:"positionX,omitempty"` // Range: -1000 to 1000
	PositionY                *float64 `json:"positionY,omitempty"` // Range: -1000 to 1000
	Rotation                 *float64 `json:"rotation,omitempty"`  // Range: -360 to 360
	Size                     *float64 `json:"size,omitempty"`      // Range: -100 to 100
}

func Float64Ptr(v float64) *float64 {
	return &v
}

//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------

type HotkeysInCurrentModelRequest struct {
	ModelID            string `json:"modelID,omitempty"`
	Live2DItemFileName string `json:"live2DItemFileName,omitempty"`
}

type HotkeyKeyCombination struct {
}

type Hotkey struct {
	Name             string                 `json:"name"`
	Type             string                 `json:"type"`
	Description      string                 `json:"description"`
	File             string                 `json:"file"`
	HotkeyID         string                 `json:"hotkeyID"`
	KeyCombination   []HotkeyKeyCombination `json:"keyCombination"`
	OnScreenButtonID int                    `json:"onScreenButtonID"`
}

type HotkeysInCurrentModelResponse struct {
	ModelLoaded      bool     `json:"modelLoaded"`
	ModelName        string   `json:"modelName"`
	ModelID          string   `json:"modelID"`
	AvailableHotkeys []Hotkey `json:"availableHotkeys"`
}

//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------

type HotkeyTriggerRequest struct {
	HotkeyID       string `json:"hotkeyID"`                 // ID or case-insensitive name
	ItemInstanceID string `json:"itemInstanceID,omitempty"` // Optional Live2D item target
}

type HotkeyTriggerResponse struct {
	HotkeyID string `json:"hotkeyID"`
}

//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------

type ExpressionStateRequest struct {
	Details        bool   `json:"details"`                  // Returns hotkeys & params if true
	ExpressionFile string `json:"expressionFile,omitempty"` // Leave empty for all expressions
}

type ExpressionHotkey struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type ExpressionParameter struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

type Expression struct {
	Name                       string                `json:"name"`
	File                       string                `json:"file"`
	Active                     bool                  `json:"active"`
	DeactivateWhenKeyIsLetGo   bool                  `json:"deactivateWhenKeyIsLetGo"`
	AutoDeactivateAfterSeconds bool                  `json:"autoDeactivateAfterSeconds"`
	SecondsRemaining           float64               `json:"secondsRemaining"`
	UsedInHotkeys              []ExpressionHotkey    `json:"usedInHotkeys"`
	Parameters                 []ExpressionParameter `json:"parameters"`
}

type ExpressionStateResponse struct {
	ModelLoaded bool         `json:"modelLoaded"`
	ModelName   string       `json:"modelName"`
	ModelID     string       `json:"modelID"`
	Expressions []Expression `json:"expressions"`
}

//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------

type ExpressionActivationRequest struct {
	ExpressionFile string  `json:"expressionFile"` // Must end in .exp3.json
	FadeTime       float64 `json:"fadeTime"`       // Clamped between 0.0 and 2.0 (default 0.25)
	Active         bool    `json:"active"`
}

//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------

type ArtMeshGroup struct {
	GroupID                  string   `json:"groupID"`
	GroupName                string   `json:"groupName"`
	NumberOfArtMeshesInGroup int      `json:"numberOfArtMeshesInGroup"`
	ArtMeshNames             []string `json:"artMeshNames"`
}

type ArtMeshListResponse struct {
	ModelLoaded           bool           `json:"modelLoaded"`
	NumberOfArtMeshNames  int            `json:"numberOfArtMeshNames"`
	NumberOfArtMeshTags   int            `json:"numberOfArtMeshTags"`
	ArtMeshNames          []string       `json:"artMeshNames"`
	ArtMeshTags           []string       `json:"artMeshTags"`
	NumberOfArtMeshGroups int            `json:"numberOfArtMeshGroups"`
	ArtMeshGroups         []ArtMeshGroup `json:"artMeshGroups"`
}

//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------

type ArtMeshAtPositionRequest struct {
	X         float64 `json:"x"`                   // Range: -1000 to 1000
	Y         float64 `json:"y"`                   // Range: -1000 to 1000
	Visualize float64 `json:"visualize,omitempty"` // Seconds indicator dot stays visible (0.0 to 5.0)
}

type Vector2 struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type ArtMeshHitInfo struct {
	ModelID       string  `json:"modelID"`
	ArtMeshID     string  `json:"artMeshID"`
	Angle         float64 `json:"angle"`
	Size          float64 `json:"size"`
	VertexID1     int     `json:"vertexID1"`
	VertexID2     int     `json:"vertexID2"`
	VertexID3     int     `json:"vertexID3"`
	VertexWeight1 float64 `json:"vertexWeight1"`
	VertexWeight2 float64 `json:"vertexWeight2"`
	VertexWeight3 float64 `json:"vertexWeight3"`
}

type ArtMeshHit struct {
	ArtMeshOrder int            `json:"artMeshOrder"`
	IsMasked     bool           `json:"isMasked"`
	HitInfo      ArtMeshHitInfo `json:"hitInfo"`
}

type ArtMeshAtPositionResponse struct {
	ModelLoaded     bool         `json:"modelLoaded"`
	LoadedModelID   string       `json:"loadedModelID"`
	LoadedModelName string       `json:"loadedModelName"`
	ModelWasHit     bool         `json:"modelWasHit"`
	CheckedPosition Vector2      `json:"checkedPosition"`
	WindowSize      Vector2      `json:"windowSize"`
	ArtMeshHitCount int          `json:"artMeshHitCount"`
	ArtMeshHits     []ArtMeshHit `json:"artMeshHits"`
}

//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------

type ColorTint struct {
	ColorR                    uint8    `json:"colorR"`                              // Range: 0 to 255
	ColorG                    uint8    `json:"colorG"`                              // Range: 0 to 255
	ColorB                    uint8    `json:"colorB"`                              // Range: 0 to 255
	ColorA                    uint8    `json:"colorA"`                              // Range: 0 to 255
	MixWithSceneLightingColor *float64 `json:"mixWithSceneLightingColor,omitempty"` // Range: 0.0 to 1.0 (default: 1.0)
}

type ArtMeshMatcher struct {
	TintAll             bool     `json:"tintAll"`
	ArtMeshNumber       []int    `json:"artMeshNumber,omitempty"`
	NameExact           []string `json:"nameExact,omitempty"`
	NameContains        []string `json:"nameContains,omitempty"`
	TagExact            []string `json:"tagExact,omitempty"`
	TagContains         []string `json:"tagContains,omitempty"`
	ArtMeshGroupIDExact []string `json:"artMeshGroupIDExact,omitempty"`
}

type ColorTintRequest struct {
	ColorTint      ColorTint      `json:"colorTint"`
	ArtMeshMatcher ArtMeshMatcher `json:"artMeshMatcher"`
}

type ColorTintResponse struct {
	MatchedArtMeshes int `json:"matchedArtMeshes"`
}

//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------

type CapturePartColor struct {
	Active bool  `json:"active"`
	ColorR uint8 `json:"colorR"` // Range: 0 to 255
	ColorG uint8 `json:"colorG"` // Range: 0 to 255
	ColorB uint8 `json:"colorB"` // Range: 0 to 255
}

type SceneColorOverlayInfoResponse struct {
	Active            bool             `json:"active"`
	ItemsIncluded     bool             `json:"itemsIncluded"`
	IsWindowCapture   bool             `json:"isWindowCapture"`
	BaseBrightness    int              `json:"baseBrightness"` // Range: 0 to 100
	ColorBoost        int              `json:"colorBoost"`     // Range: 0 to 100
	Smoothing         int              `json:"smoothing"`      // Range: 0 to 60
	ColorOverlayR     int              `json:"colorOverlayR"`  // Range: 0 to 459
	ColorOverlayG     int              `json:"colorOverlayG"`  // Range: 0 to 459
	ColorOverlayB     int              `json:"colorOverlayB"`  // Range: 0 to 459
	ColorAvgR         uint8            `json:"colorAvgR"`      // Range: 0 to 255
	ColorAvgG         uint8            `json:"colorAvgG"`      // Range: 0 to 255
	ColorAvgB         uint8            `json:"colorAvgB"`      // Range: 0 to 255
	LeftCapturePart   CapturePartColor `json:"leftCapturePart"`
	MiddleCapturePart CapturePartColor `json:"middleCapturePart"`
	RightCapturePart  CapturePartColor `json:"rightCapturePart"`
}

//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------

type FaceFoundResponse struct {
	Found bool `json:"found"`
}

//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------

type TrackingParameter struct {
	Name         string  `json:"name"`
	AddedBy      string  `json:"addedBy"`
	Value        float64 `json:"value"`
	Min          float64 `json:"min"`
	Max          float64 `json:"max"`
	DefaultValue float64 `json:"defaultValue"`
}

type InputParameterListResponse struct {
	ModelLoaded       bool                `json:"modelLoaded"`
	ModelName         string              `json:"modelName"`
	ModelID           string              `json:"modelID"`
	CustomParameters  []TrackingParameter `json:"customParameters"`
	DefaultParameters []TrackingParameter `json:"defaultParameters"`
}

//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------

type ParameterValueRequest struct {
	Name string `json:"name"`
}

type ParameterValueResponse struct {
	Name         string  `json:"name"`
	AddedBy      string  `json:"addedBy"`
	Value        float64 `json:"value"`
	Min          float64 `json:"min"`
	Max          float64 `json:"max"`
	DefaultValue float64 `json:"defaultValue"`
}

//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------

type Live2DParameter struct {
	Name         string  `json:"name"`
	Value        float64 `json:"value"`
	Min          float64 `json:"min"`
	Max          float64 `json:"max"`
	DefaultValue float64 `json:"defaultValue"`
}

type Live2DParameterListResponse struct {
	ModelLoaded bool              `json:"modelLoaded"`
	ModelName   string            `json:"modelName"`
	ModelID     string            `json:"modelID"`
	Parameters  []Live2DParameter `json:"parameters"`
}

//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------

type ParameterCreationRequest struct {
	ParameterName string  `json:"parameterName"`         // Alphanumeric, 4-32 chars, unique
	Explanation   string  `json:"explanation,omitempty"` // Short description (<256 chars)
	Min           float64 `json:"min"`                   // Range: -1000000 to 1000000
	Max           float64 `json:"max"`                   // Range: -1000000 to 1000000
	DefaultValue  float64 `json:"defaultValue"`          // Range: -1000000 to 1000000
}

type ParameterCreationResponse struct {
	ParameterName string `json:"parameterName"`
}

//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------

type ParameterDeletionRequest struct {
	ParameterName string `json:"parameterName"`
}

type ParameterDeletionResponse struct {
	ParameterName string `json:"parameterName"`
}

//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------

type ParameterValue struct {
	ID     string   `json:"id"`
	Value  float64  `json:"value"`            // Range: -1000000 to 1000000
	Weight *float64 `json:"weight,omitempty"` // Range: 0.0 to 1.0 (defaults to 1.0)
}

// InjectParameterDataRequest payload for feeding values into default or custom parameters
type InjectParameterDataRequest struct {
	FaceFound       *bool            `json:"faceFound,omitempty"` // Overrides tracking lost state if provided
	Mode            string           `json:"mode,omitempty"`      // "set" (default) or "add"
	ParameterValues []ParameterValue `json:"parameterValues"`
}

//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------

type PhysicsGroup struct {
	GroupID            string  `json:"groupID"`
	GroupName          string  `json:"groupName"`
	StrengthMultiplier float64 `json:"strengthMultiplier"` // Range: 0.0 to 2.0
	WindMultiplier     float64 `json:"windMultiplier"`     // Range: 0.0 to 2.0
}

type GetCurrentModelPhysicsResponse struct {
	ModelLoaded                  bool           `json:"modelLoaded"`
	ModelName                    string         `json:"modelName"`
	ModelID                      string         `json:"modelID"`
	ModelHasPhysics              bool           `json:"modelHasPhysics"`
	PhysicsSwitchedOn            bool           `json:"physicsSwitchedOn"`
	UsingLegacyPhysics           bool           `json:"usingLegacyPhysics"`
	PhysicsFPSSetting            int            `json:"physicsFPSSetting"` // 30, 60, 120, or -1 (app default)
	BaseStrength                 int            `json:"baseStrength"`      // Range: 0 to 100
	BaseWind                     int            `json:"baseWind"`          // Range: 0 to 100
	APIPhysicsOverrideActive     bool           `json:"apiPhysicsOverrideActive"`
	APIPhysicsOverridePluginName string         `json:"apiPhysicsOverridePluginName"`
	PhysicsGroups                []PhysicsGroup `json:"physicsGroups"`
}

//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------

type PhysicsOverride struct {
	ID              string  `json:"id"`              // Physics group ID, or empty string if setBaseValue is true
	Value           float64 `json:"value"`           // Multiplier (0.0 to 2.0) OR Base Value (0 to 100)
	SetBaseValue    bool    `json:"setBaseValue"`    // True to set base value, false to set group multiplier
	OverrideSeconds float64 `json:"overrideSeconds"` // Duration in seconds (0.5 to 5.0)
}

type SetCurrentModelPhysicsRequest struct {
	StrengthOverrides []PhysicsOverride `json:"strengthOverrides,omitempty"`
	WindOverrides     []PhysicsOverride `json:"windOverrides,omitempty"`
}

//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------

type ItemListRequest struct {
	IncludeAvailableSpots       bool   `json:"includeAvailableSpots"`
	IncludeItemInstancesInScene bool   `json:"includeItemInstancesInScene"`
	IncludeAvailableItemFiles   bool   `json:"includeAvailableItemFiles"` // Warning: Reads disk, may lag VTS
	OnlyItemsWithFileName       string `json:"onlyItemsWithFileName,omitempty"`
	OnlyItemsWithInstanceID     string `json:"onlyItemsWithInstanceID,omitempty"`
}

type ItemInstance struct {
	ItemFileName               string  `json:"itemFileName"`
	ItemInstanceID             string  `json:"itemInstanceID"`
	Order                      int     `json:"order"`
	CustomOrder                int     `json:"customOrder"`
	Flipped                    bool    `json:"flipped"`
	Locked                     bool    `json:"locked"`
	Smoothing                  float64 `json:"smoothing"`
	Framerate                  float64 `json:"framerate"`
	PinnedToModel              bool    `json:"pinnedToModel"`
	PinnedToModelID            string  `json:"pinnedToModelID"`
	PinnedToArtMeshID          string  `json:"pinnedToArtMeshID"`
	Group                      string  `json:"group"`
	SceneLightingMultiplyColor string  `json:"sceneLightingMultiplyColor"`
	PositionX                  float64 `json:"positionX"`
	PositionY                  float64 `json:"positionY"`
	Rotation                   float64 `json:"rotation"`
	Size                       float64 `json:"size"`
}

type ItemFile struct {
	FileName string `json:"fileName"`
	Type     string `json:"type"`
	Size     int64  `json:"size"`
}

type ItemListResponse struct {
	ItemsInSceneCount       int            `json:"itemsInSceneCount"`
	TotalItemsAllowed       int            `json:"totalItemsAllowed"`
	AvailableSpots          []int          `json:"availableSpots"`
	ItemInstancesInScene    []ItemInstance `json:"itemInstancesInScene"`
	AvailableItemFilesCount int            `json:"availableItemFilesCount"`
	AvailableItemFiles      []ItemFile     `json:"availableItemFiles"`
}

//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------

type ItemLoadRequest struct {
	FileName                              string  `json:"fileName"`
	PositionX                             float64 `json:"positionX,omitempty"`
	PositionY                             float64 `json:"positionY,omitempty"`
	Size                                  float64 `json:"size,omitempty"`
	Rotation                              float64 `json:"rotation,omitempty"`
	FadeTime                              float64 `json:"fadeTime,omitempty"`
	Order                                 int     `json:"order,omitempty"`
	FailIfOrderTaken                      bool    `json:"failIfOrderTaken"`
	Smoothing                             float64 `json:"smoothing,omitempty"`
	Censored                              bool    `json:"censored"`
	Flipped                               bool    `json:"flipped"`
	Locked                                bool    `json:"locked"`
	UnloadWhenPluginDisconnects           bool    `json:"unloadWhenPluginDisconnects"`
	CustomDataBase64                      string  `json:"customDataBase64,omitempty"`
	CustomDataAskUserFirst                bool    `json:"customDataAskUserFirst"`
	CustomDataSkipAskingUserIfWhitelisted bool    `json:"customDataSkipAskingUserIfWhitelisted"`
	CustomDataAskTimer                    int     `json:"customDataAskTimer,omitempty"` // Seconds before timing out (-1 for default)
}

type ItemLoadResponse struct {
	ItemInstanceID string `json:"itemInstanceID"`
	FileName       string `json:"fileName"`
}

//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------

type ItemUnloadRequest struct {
	UnloadAllInScene                              bool     `json:"unloadAllInScene"`
	UnloadAllLoadedByThisPlugin                   bool     `json:"unloadAllLoadedByThisPlugin"`
	AllowUnloadingItemsLoadedByUserOrOtherPlugins bool     `json:"allowUnloadingItemsLoadedByUserOrOtherPlugins"`
	InstanceIDs                                   []string `json:"instanceIDs,omitempty"`
	FileNames                                     []string `json:"fileNames,omitempty"`
}

type UnloadedItem struct {
	InstanceID string `json:"instanceID"`
	FileName   string `json:"fileName"`
}

type ItemUnloadResponse struct {
	UnloadedItems []UnloadedItem `json:"unloadedItems"`
}

//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------

type ItemAnimationControlRequest struct {
	ItemInstanceID        string  `json:"itemInstanceID"`
	Framerate             float64 `json:"framerate,omitempty"`  // FPS (0.1 to 120), -1 to leave unchanged
	Frame                 int     `json:"frame,omitempty"`      // Target frame index, -1 to leave unchanged
	Brightness            float64 `json:"brightness,omitempty"` // Range: 0.0 to 1.0, -1 to leave unchanged
	Opacity               float64 `json:"opacity,omitempty"`    // Range: 0.0 to 1.0, -1 to leave unchanged
	SetAutoStopFrames     bool    `json:"setAutoStopFrames"`
	AutoStopFrames        []int   `json:"autoStopFrames,omitempty"` // Max 1024 frame indices
	SetAnimationPlayState bool    `json:"setAnimationPlayState"`
	AnimationPlayState    bool    `json:"animationPlayState"`
}

type ItemAnimationControlResponse struct {
	Frame            int  `json:"frame"`            // Current frame index
	AnimationPlaying bool `json:"animationPlaying"` // True if animation is actively playing
}

//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------

type ItemMoveTarget struct {
	ItemInstanceID string  `json:"itemInstanceID"`
	TimeInSeconds  float64 `json:"timeInSeconds"`       // Range: 0 to 30
	FadeMode       string  `json:"fadeMode,omitempty"`  // "linear", "easeIn", "easeOut", "easeBoth", "overshoot", "zip"
	PositionX      float64 `json:"positionX,omitempty"` // Set <= -1000 to ignore
	PositionY      float64 `json:"positionY,omitempty"` // Set <= -1000 to ignore
	Size           float64 `json:"size,omitempty"`      // Set <= -1000 to ignore
	Rotation       float64 `json:"rotation,omitempty"`  // Set <= -1000 to ignore
	Order          int     `json:"order,omitempty"`     // Set <= -1000 to ignore
	SetFlip        bool    `json:"setFlip"`
	Flip           bool    `json:"flip"`
	UserCanStop    bool    `json:"userCanStop"`
}

type ItemMoveRequest struct {
	ItemsToMove []ItemMoveTarget `json:"itemsToMove"`
}

type MovedItemStatus struct {
	ItemInstanceID string `json:"itemInstanceID"`
	Success        bool   `json:"success"`
	ErrorID        int    `json:"errorID"` // -1 on success
}

type ItemMoveResponse struct {
	MovedItems []MovedItemStatus `json:"movedItems"`
}

//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------

type ItemSortRequest struct {
	ItemInstanceID        string `json:"itemInstanceID"`
	FrontOn               bool   `json:"frontOn"`
	BackOn                bool   `json:"backOn"`
	SetSplitPoint         string `json:"setSplitPoint"`                   // "Unchanged", "UseArtMeshID"
	SetFrontOrder         string `json:"setFrontOrder"`                   // "Unchanged", "UseArtMeshID", "UseSpecialID"
	SetBackOrder          string `json:"setBackOrder"`                    // "Unchanged", "UseArtMeshID", "UseSpecialID"
	SplitAt               string `json:"splitAt,omitempty"`               // ArtMesh ID in item model
	WithinModelOrderFront string `json:"withinModelOrderFront,omitempty"` // ArtMesh ID in main model OR "FullyInFront" / "FullyInBack"
	WithinModelOrderBack  string `json:"withinModelOrderBack,omitempty"`  // ArtMesh ID in main model OR "FullyInBack"
}

type ItemSortResponse struct {
	ItemInstanceID                    string `json:"itemInstanceID"`
	ModelLoaded                       bool   `json:"modelLoaded"`
	ModelID                           string `json:"modelID"`
	ModelName                         string `json:"modelName"`
	LoadedModelHadRequestedFrontLayer bool   `json:"loadedModelHadRequestedFrontLayer"`
	LoadedModelHadRequestedBackLayer  bool   `json:"loadedModelHadRequestedBackLayer"`
}

//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------

type ArtMeshSelectionRequest struct {
	TextOverride          string   `json:"textOverride,omitempty"`    // Top text prompt (4-1024 chars)
	HelpOverride          string   `json:"helpOverride,omitempty"`    // Help button (?) popup text (4-1024 chars)
	RequestedArtMeshCount int      `json:"requestedArtMeshCount"`     // Required count, or <= 0 for arbitrary (>= 1)
	ActiveArtMeshes       []string `json:"activeArtMeshes,omitempty"` // Pre-selected ArtMesh IDs
}

type ArtMeshSelectionResponse struct {
	Success           bool     `json:"success"`           // True if user clicked "OK", false if canceled
	ActiveArtMeshes   []string `json:"activeArtMeshes"`   // Selected ArtMesh IDs
	InactiveArtMeshes []string `json:"inactiveArtMeshes"` // Unselected ArtMesh IDs in the current model
}

//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------

type ItemPinInfo struct {
	ModelID       string  `json:"modelID,omitempty"`   // Empty string defaults to currently loaded model
	ArtMeshID     string  `json:"artMeshID,omitempty"` // Empty string picks a random ArtMesh
	Angle         float64 `json:"angle"`
	Size          float64 `json:"size"`
	VertexID1     int     `json:"vertexID1,omitempty"`     // Required if VertexPinType is "Provided"
	VertexID2     int     `json:"vertexID2,omitempty"`     // Required if VertexPinType is "Provided"
	VertexID3     int     `json:"vertexID3,omitempty"`     // Required if VertexPinType is "Provided"
	VertexWeight1 float64 `json:"vertexWeight1,omitempty"` // Weights must sum to 1.0 (Barycentric coordinates)
	VertexWeight2 float64 `json:"vertexWeight2,omitempty"`
	VertexWeight3 float64 `json:"vertexWeight3,omitempty"`
}

type ItemPinRequest struct {
	Pin             bool         `json:"pin"`
	ItemInstanceID  string       `json:"itemInstanceID"`
	AngleRelativeTo string       `json:"angleRelativeTo,omitempty"` // "RelativeToWorld", "RelativeToCurrentItemRotation", "RelativeToModel", "RelativeToPinPosition"
	SizeRelativeTo  string       `json:"sizeRelativeTo,omitempty"`  // "RelativeToWorld", "RelativeToCurrentItemSize"
	VertexPinType   string       `json:"vertexPinType,omitempty"`   // "Provided", "Center", "Random"
	PinInfo         *ItemPinInfo `json:"pinInfo,omitempty"`         // Omit or set to nil when unpinning (pin: false)
}

type ItemPinResponse struct {
	IsPinned       bool   `json:"isPinned"`
	ItemInstanceID string `json:"itemInstanceID"`
	ItemFileName   string `json:"itemFileName"`
}
