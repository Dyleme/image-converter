package storage

import (
	"strings"

	"github.com/google/uuid"
)

// generateName is function which can be used to
// generate uniqe names for files, fith the same extension.
func generateName(filename string) string {
	name := uuid.NewString()

	dotPos := strings.LastIndex(filename, ".")
	if dotPos != -1 {
		name += filename[dotPos:]
	}

	return name
}
