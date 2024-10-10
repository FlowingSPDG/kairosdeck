package handlers

import (
	"context"

	"github.com/FlowingSPDG/kairosdeck/Source/backend/pi"
	"github.com/FlowingSPDG/streamdeck"
)

const (
	// uuids
	patchSceneUUID = "dev.flowingspdg.kairos.patchScene"
	patchMacroUUID = "dev.flowingspdg.kairos.patchMacro"
)

// Handlers
// TODO: Actionごとに細分化する
type Handlers struct {
	client   *streamdeck.Client
	settings *pi.PropertyInspectorStore
}

func SetupHandlers(client *streamdeck.Client) *Handlers {
	h := &Handlers{
		client:   client,
		settings: &pi.PropertyInspectorStore{},
	}

	actionPatchScene := client.Action(patchSceneUUID)
	actionPatchScene.RegisterHandler(streamdeck.WillAppear, h.PatchSceneWillAppear)
	actionPatchScene.RegisterHandler(streamdeck.WillDisappear, h.PatchSceneWillDisappear)
	actionPatchScene.RegisterHandler(streamdeck.KeyDown, h.PatchSceneKeyDown)
	actionPatchScene.RegisterHandler(streamdeck.SendToPlugin, h.PatchSceneSendToPlugin)

	actionPatchMacro := client.Action(patchMacroUUID)
	actionPatchMacro.RegisterHandler(streamdeck.WillAppear, h.PatchMacroWillAppear)
	actionPatchMacro.RegisterHandler(streamdeck.WillDisappear, h.PatchMacroWillDisappear)
	actionPatchMacro.RegisterHandler(streamdeck.KeyDown, h.PatchMacroKeyDown)
	actionPatchMacro.RegisterHandler(streamdeck.SendToPlugin, h.PatchMacroSendToPlugin)

	return h
}

func (h *Handlers) Run(ctx context.Context) error {
	return h.client.Run(ctx)
}
