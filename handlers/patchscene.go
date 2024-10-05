package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/FlowingSPDG/kairos-go"
	"github.com/FlowingSPDG/kairosdeck/pi"
	"github.com/FlowingSPDG/streamdeck"

	"golang.org/x/xerrors"
)

func (h *Handlers) PatchSceneWillAppear(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
	p := streamdeck.WillAppearPayload[pi.PatchSceneSetting]{}
	if err := json.Unmarshal(event.Payload, &p); err != nil {
		return err
	}

	s, _ := h.settings.PatchSceneSettings.LoadOrStore(event.Context, &pi.PatchSceneSetting{})
	_ = s
	return nil
}

func (h *Handlers) PatchSceneWillDisappear(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
	h.settings.PatchSceneSettings.Delete(event.Context)
	return client.SetSettings(ctx, map[string]any{})
}

func (h *Handlers) PatchSceneKeyDown(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
	s, ok := h.settings.PatchSceneSettings.Load(event.Context)
	if !ok {
		return fmt.Errorf("couldn't find settings for context %v", event.Context)
	}

	kr := kairos.NewKairosRestClient(s.IP, fmt.Sprint(s.Port), s.User, s.Password)
	if err := kr.PatchScene(ctx, s.SceneUUID, s.LayerUUID, s.A, s.B, nil); err != nil {
		return xerrors.Errorf("Failed to patch scene : %w", err)
	}

	if err := client.SetSettings(ctx, s); err != nil {
		return xerrors.Errorf("Failed to save setting : %w", err)
	}

	if err := client.ShowOk(ctx); err != nil {
		return xerrors.Errorf("Failed to execute ShowOk() : %w", err)
	}
	return nil
}
