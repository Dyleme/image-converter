package model

import (
	"time"
)

// Sruct to put it in requests database.
type Request struct {
	ID             int       `json:"id"`
	OpStatus       string    `json:"status"`
	RequestTime    time.Time `json:"requestTime"`
	CompletionTime time.Time `json:"completionTime,omitempty"`
	OriginalID     int       `json:"originalID"`
	ProcessedID    int       `json:"processedID"`
	Ratio          float32   `json:"ratio"`
	OriginalType   string    `json:"originalType"`
	ProcessedType  string    `json:"processedType"`
}
