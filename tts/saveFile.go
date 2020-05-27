package tts

// SaveFile describes a Tabletop Simulator save file
// More info: https://tabletopsimulator.gamepedia.com/Save_File_Format
type SaveFile struct {
	SaveName     string    `json:"SaveName"`
	TableURL     string    `json:"TableURL"`
	SkyURL       string    `json:"SkyURL"`
	ObjectStates []Objects `json:"ObjectStates"`
}

// Objects describes the objects and their state used in the game
type Objects struct {
	Name        string          `json:"Names"`
	CustomMesh  Meshes          `json:"CustomMesh"`
	CustomImage Image           `json:"CustomImage"`
	CustomDeck  map[string]Deck `json:"CustomDeck"`
}

// Meshes describes a custom mesh
type Meshes struct {
	MeshURL     string `json:"MeshURL"`
	DiffuseURL  string `json:"DiffuseURL"`
	NormalURL   string `json:"NormalURL"`
	ColliderURL string `json:"ColliderURL"`
}

// Image describes a custom image
type Image struct {
	ImageURL string `json:"ImageURL"`
}

// Deck describes a custom deck
type Deck struct {
	FaceURL string `json:"FaceURL"`
	BackURL string `json:"BackURL"`
}
