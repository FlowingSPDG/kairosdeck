package pi

type BaseSetting struct {
	IP       string `json:"ip"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type PatchSceneSetting struct {
	BaseSetting
	SceneUUID string
	LayerUUID string
	A         *string
	B         *string
}

type PropertyInspectorStore struct {
	PatchSceneSettings SettingStore[PatchSceneSetting]
}
