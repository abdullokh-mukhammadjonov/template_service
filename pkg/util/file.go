package util

import (
	"os"
	"path/filepath"
)

func CreateDolders(p string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(p), 0770); err != nil {
		return nil, err
	}
	return os.Create(p)
}
