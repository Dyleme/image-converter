package model

// Information about image conversion.
type ConversionInfo struct {
	// Ration with which you will convert image.
	Ratio float32 `json:"ratio"`

	// Type to which you will convert image.
	Type string `json:"newType"`
}

// Information about image.
type Info struct {
	Width  int
	Height int
	Type   string
	URL    string
}
