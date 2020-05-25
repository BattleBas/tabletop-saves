package tts

// SaveFile describes a Tabletop Simulator save file
// More info: https://tabletopsimulator.gamepedia.com/Save_File_Format
type SaveFile struct {
	SaveName     string    `json:"SaveName"`
	TableURL     string    `json:"TableURL"`
	SkyURL       string    `json:"SkyURL"`
	FaceURL      string    `json:"FaceURL"`
	BackURL      string    `json:"BackURL"`
	ImageURL     string    `json:"ImageURL"`
	DiffuseURL   string    `json:"DiffuseURL"`
	NormalURL    string    `json:"NormalURL"`
	ObjectStates []Objects `json:"ObjectStates"`
}

// Objects describes the objects and their state used in the game
type Objects struct {
	Name string `json:"Names"`
}
