//go:build js && wasm

package fs

import (
	"os"
)

func extendedStat(fi os.FileInfo) *ExtendedFileInfo {
	modTime := fi.ModTime()
	return &ExtendedFileInfo{
		Name:       fi.Name(),
		Mode:       fi.Mode(),
		Size:       fi.Size(),
		ModTime:    modTime,
		AccessTime: modTime,
		ChangeTime: modTime,
	}
}

func (*ExtendedFileInfo) RecallOnDataAccess() (bool, error) {
	return false, nil
}
