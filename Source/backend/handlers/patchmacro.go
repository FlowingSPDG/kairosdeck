package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/FlowingSPDG/kairos-go"
	"github.com/FlowingSPDG/kairosdeck/Source/backend/pi"
	"github.com/FlowingSPDG/streamdeck"

	"golang.org/x/xerrors"
)

func (h *Handlers) PatchMacroWillAppear(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
	msg := fmt.Sprintf("Context %s received WillAppear with payload :%s", event.Context, event.Payload)
	client.LogMessage(ctx, msg)

	p := streamdeck.WillAppearPayload[pi.PatchMacroSetting]{}
	if err := json.Unmarshal(event.Payload, &p); err != nil {
		return xerrors.Errorf("Failed to Unmarshal JSON : %w", err)
	}

	if p.Settings.IsDefault() {
		p.Settings.Initialize()
	}

	if err := client.SetSettings(ctx, p.Settings); err != nil {
		return xerrors.Errorf("Failed to save setting : %w", err)
	}
	h.settings.PatchMacroSettings.Store(event.Context, &p.Settings)

	return nil
}

func (h *Handlers) PatchMacroWillDisappear(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
	msg := fmt.Sprintf("Deleting context %s", event.Context)
	client.LogMessage(ctx, msg)

	h.settings.PatchMacroSettings.Delete(event.Context)
	return client.SetSettings(ctx, map[string]any{})
}

func (h *Handlers) PatchMacroKeyDown(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
	msg := fmt.Sprintf("Context %s Keydown", event.Context)
	client.LogMessage(ctx, msg)

	s, ok := h.settings.PatchMacroSettings.Load(event.Context)
	if !ok {
		return fmt.Errorf("couldn't find settings for context %v", event.Context)
	}
	msg = fmt.Sprintf("Context %s Keydown with settings :%v", event.Context, s)
	client.LogMessage(ctx, msg)

	kr := kairos.NewKairosRestClient(s.Host, fmt.Sprint(s.Port), s.User, s.Password)
	if err := kr.PatchMacro(ctx, s.MacroUUID, s.State); err != nil {
		return xerrors.Errorf("Failed to patch scene : %w", err)
	}

	if err := client.ShowOk(ctx); err != nil {
		return xerrors.Errorf("Failed to execute ShowOk() : %w", err)
	}
	return nil
}

func (h *Handlers) PatchMacroSendToPlugin(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
	msg := fmt.Sprintf("SendToPlugin for Context %s with payload :%s", event.Context, event.Payload)
	client.LogMessage(ctx, msg)

	p := pi.ToPluginPayload{}
	if err := json.Unmarshal(event.Payload, &p); err != nil {
		return err
	}

	switch p.Action {
	case "refresh":
		return h.PatchMacroRefreshMacro(ctx, client, event)
	}

	return nil
}

func (h *Handlers) PatchMacroRefreshMacro(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
	s, ok := h.settings.PatchMacroSettings.Load(event.Context)
	if !ok {
		return fmt.Errorf("couldn't find settings for context %v", event.Context)
	}

	kr := kairos.NewKairosRestClient(s.Host, fmt.Sprint(s.Port), s.User, s.Password)
	macros, err := kr.GetMacros(ctx)
	if err != nil {
		return xerrors.Errorf("Failed to get Scenes() : %w", err)
	}
	s.Macros = macros

	if err := client.SetSettings(ctx, s); err != nil {
		return xerrors.Errorf("Failed to save setting : %w", err)
	}
	h.settings.PatchMacroSettings.Store(event.Context, s)

	return nil
}

func (h *Handlers) PatchMacroDidReceiveSettings(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
	msg := fmt.Sprintf("Received settings for Context %s with payload :%s", event.Context, event.Payload)
	client.LogMessage(ctx, msg)

	p := streamdeck.DidReceiveSettingsPayload[pi.PatchMacroSetting]{}
	if err := json.Unmarshal(event.Payload, &p); err != nil {
		return xerrors.Errorf("Failed to Unmarshal JSON : %w", err)
	}

	kr := kairos.NewKairosRestClient(p.Settings.Host, fmt.Sprint(p.Settings.Port), p.Settings.User, p.Settings.Password)
	macros, err := kr.GetMacros(ctx)
	if err != nil {
		return xerrors.Errorf("Failed to get Scenes() : %w", err)
	}
	p.Settings.Macros = macros

	if err := client.SetSettings(ctx, p.Settings); err != nil {
		return xerrors.Errorf("Failed to save setting : %w", err)
	}
	h.settings.PatchMacroSettings.Store(event.Context, &p.Settings)

	msg = fmt.Sprintf("Settings for Context %s overwritten. payload :%s", event.Context, event.Payload)
	client.LogMessage(ctx, msg)

	return nil
}
