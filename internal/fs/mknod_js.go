//go:build js && wasm

package fs

import "github.com/restic/restic/internal/errors"

func mknod(path string, mode uint32, dev uint64) error {
	_ = path
	_ = mode
	_ = dev
	return errors.New("device nodes are not supported on js/wasm")
}
