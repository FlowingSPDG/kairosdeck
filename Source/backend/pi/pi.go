package pi

import (
	"reflect"

	"github.com/FlowingSPDG/kairos-go/internals/objects"
)

type BaseSetting struct {
	// TODO: Add this to "Account" section on StreamDeck app.
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type PatchSceneSetting struct {
	BaseSetting
	SceneUUID string  `json:"scene_uuid"`
	LayerUUID string  `json:"layer_uuid"`
	A         *string `json:"a"`
	B         *string `json:"b"`

	// External sources
	Scenes []*objects.SceneR `json:"scenes"`
}

func (p PatchSceneSetting) IsDefault() bool {
	return reflect.DeepEqual(p, PatchSceneSetting{})
}

type PatchMacroSetting struct {
	BaseSetting
	MacroUUID string `json:"macro_uuid"`
	State     string `json:"state"`

	// External sources
	// TODO: Support multi-host
	Macros []*objects.MacroR `json:"macros"`
}

func (p PatchMacroSetting) IsDefault() bool {
	return reflect.DeepEqual(p, PatchMacroSetting{})
}
func (p *PatchMacroSetting) Initialize() {
	p.BaseSetting = BaseSetting{
		Host: "localhost",
		Port: 1234,
		User: "Kairos",
		Password: "Kairos",
	}
	p.MacroUUID = ""
	p.State = "play"
	p.Macros = make([]*objects.MacroR, 0)
}

type PropertyInspectorStore struct {
	PatchSceneSettings SettingStore[PatchSceneSetting]
	PatchMacroSettings SettingStore[PatchMacroSetting]
}
