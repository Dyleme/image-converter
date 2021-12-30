package model

// Information about image conversion.
type ConversionInfo struct {
	// Ration with which you will convert image.
	Ratio float32 `json:"ratio"`

	// Type to which you will convert image.
	Type string `json:"newType"`
}

// Information about image.
type ReuquestImageInfo struct {
	Type string
	URL  string
}

type ConvImageInfo struct {
	UserID  int
	OldImID int
	OldURL  string
	OldType string
	NewType string
	Ratio   float32
}

// ConverstionedImage is struct, which contains images and all needed
// information to convert images.
type ConverstionedImage struct {
	ReqID    int    `json:"reqID"`
	FileName string `json:"fileName"`
}
