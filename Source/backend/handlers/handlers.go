package handlers

import (
	"context"

	"github.com/FlowingSPDG/kairosdeck/Source/backend/pi"
	"github.com/FlowingSPDG/streamdeck"
)

const (
	// uuids
	patchSceneUUID = "dev.flowingspdg.kairos.patch-scene"
	patchMacroUUID = "dev.flowingspdg.kairos.patch-macro"
)

// Handlers
// TODO: Actionごとに細分化する
type Handlers struct {
	client   *streamdeck.Client
	settings *pi.PropertyInspectorStore
}

func SetupHandlers(ctx context.Context, params streamdeck.RegistrationParams) *Handlers {
	h := &Handlers{
		client: streamdeck.NewClient(ctx, params),
		settings: &pi.PropertyInspectorStore{
			PatchSceneSettings: pi.NewSettingStore[pi.PatchSceneSetting](),
			PatchMacroSettings: pi.NewSettingStore[pi.PatchMacroSetting](),
		},
	}

	actionPatchScene := h.client.Action(patchSceneUUID)
	actionPatchScene.RegisterHandler(streamdeck.WillAppear, h.PatchSceneWillAppear)
	actionPatchScene.RegisterHandler(streamdeck.WillDisappear, h.PatchSceneWillDisappear)
	actionPatchScene.RegisterHandler(streamdeck.KeyDown, h.PatchSceneKeyDown)
	actionPatchScene.RegisterHandler(streamdeck.SendToPlugin, h.PatchSceneSendToPlugin)

	actionPatchMacro := h.client.Action(patchMacroUUID)
	actionPatchMacro.RegisterHandler(streamdeck.WillAppear, h.PatchMacroWillAppear)
	actionPatchMacro.RegisterHandler(streamdeck.WillDisappear, h.PatchMacroWillDisappear)
	actionPatchMacro.RegisterHandler(streamdeck.KeyDown, h.PatchMacroKeyDown)
	actionPatchMacro.RegisterHandler(streamdeck.SendToPlugin, h.PatchMacroSendToPlugin)
	actionPatchMacro.RegisterHandler(streamdeck.DidReceiveSettings, h.PatchMacroDidReceiveSettings)

	return h
}

func (h *Handlers) Run(ctx context.Context) error {
	return h.client.Run(ctx)
}
