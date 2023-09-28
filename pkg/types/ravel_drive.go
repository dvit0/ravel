package types

type RavelDriveSpec struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
}

type RavelDrive struct {
	Id       string `json:"id"`
	Internal bool   `json:"internal"`
	*RavelDriveSpec
}
