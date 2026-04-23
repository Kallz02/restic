//go:build js && wasm

package fs

import "os"

// The browser runtime does not expose Unix open flags. Keep the portable
// values and treat the unsupported flags as no-ops for compile-time support.
const (
	O_RDONLY    int = os.O_RDONLY
	O_WRONLY    int = os.O_WRONLY
	O_RDWR      int = os.O_RDWR
	O_APPEND    int = os.O_APPEND
	O_CREATE    int = os.O_CREATE
	O_EXCL      int = os.O_EXCL
	O_SYNC      int = os.O_SYNC
	O_TRUNC     int = os.O_TRUNC
	O_NONBLOCK  int = 0
	O_NOFOLLOW  int = 0
	O_DIRECTORY int = 0
)
