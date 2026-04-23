//go:build js && wasm

package fs

import "github.com/restic/restic/internal/data"

func lchown(_ string, _ *data.Node, _ bool) error {
	return nil
}

func nodeRestoreGenericAttributes(node *data.Node, _ string, warn func(msg string)) error {
	return data.HandleAllUnknownGenericAttributesFound(node.GenericAttributes, warn)
}

func nodeFillGenericAttributes(_ *data.Node, _ string, _ *ExtendedFileInfo) error {
	return nil
}

func nodeRestoreExtendedAttributes(_ *data.Node, _ string, _ func(xattrName string) bool) error {
	return nil
}

func nodeFillExtendedAttributes(_ *data.Node, _ string, _ bool, _ func(format string, args ...any)) error {
	return nil
}

func utimesNano(path string, atime, mtime int64, typ data.NodeType) error {
	_ = path
	_ = atime
	_ = mtime
	_ = typ
	return nil
}
