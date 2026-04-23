//go:build js && wasm

package fs

import (
	"errors"
	"os"
)

func fixpath(name string) string {
	return name
}

func TempFile(dir, prefix string) (*os.File, error) {
	return nil, errors.New("temporary files are not supported on js/wasm")
}

func isNotSupported(err error) bool {
	return false
}

func chmod(_ string, _ os.FileMode) error {
	return nil
}
