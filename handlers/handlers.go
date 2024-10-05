package handlers

import (
	"context"

	"github.com/FlowingSPDG/kairosdeck/pi"
	"github.com/FlowingSPDG/streamdeck"
)

const (
	// uuids
	patchSceneUUID = "dev.flowingspdg.kairos.patchScene"
)

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

	return h
}

func (h *Handlers) Run(ctx context.Context) error {
	return h.client.Run(ctx)
}
