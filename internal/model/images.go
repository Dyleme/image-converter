package model

import "context"

type ConversionInfo struct {
	Ratio float32 `json:"ratio"`
	Type  string  `json:"newType"`
}

type Info struct {
	Width  int
	Height int
	Type   string
	URL    string
}

type ConversionData struct {
	Ctx       context.Context
	ImageInfo ConversionInfo `json:"imageInfo"`
	UserID    int            `json:"userID"`
	ReqID     int            `json:"reqID"`
	OldType   string         `json:"oldType"`
	Pic       []byte         `json:"pic"`
	FileName  string         `json:"fileName"`
}
