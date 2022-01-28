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

// ConvImageInfo is a struct with request info to get from ConvRepo.
type ConvImageInfo struct {
	UserID  int
	OldImID int
	OldURL  string
	OldType string
	NewType string
	Ratio   float32
}

// RequestToProcess is struct, which contains request id and name of converted image.
type RequestToProcess struct {
	ReqID    int    `json:"reqID"`
	FileName string `json:"fileName"`
}
