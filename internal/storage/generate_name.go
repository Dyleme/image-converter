package storage

import (
	"strings"

	"github.com/google/uuid"
)

func generateName(filename string) string {
	name := uuid.NewString()

	dotPos := strings.LastIndex(filename, ".")
	if dotPos != -1 {
		name += filename[dotPos:]
	}

	return name
}
