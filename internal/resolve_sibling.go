package internal

import (
	"path/filepath"
	"strings"
)

func ResolveSibling(filename, extension string) string {
	return strings.TrimSuffix(filename, filepath.Ext(filename)) + extension
}
