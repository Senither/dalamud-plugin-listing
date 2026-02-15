package state

import (
	"os"
	"path/filepath"
	"strings"
)

const cacheDirEnv = "APP_CACHE_DIR"

func cachePath(filename string) string {
	dir := strings.TrimSpace(os.Getenv(cacheDirEnv))
	if dir == "" {
		return filename
	}

	return filepath.Join(dir, filename)
}
