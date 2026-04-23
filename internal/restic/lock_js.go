//go:build js && wasm

package restic

import "os/user"

// UidGidInt always returns 0 on js/wasm because the browser runtime does not
// expose numeric OS user identifiers.
func UidGidInt(_ *user.User) (uid, gid uint32, err error) {
	return 0, 0, nil
}

// checkProcess intentionally avoids PID-based stale-lock detection in the
// browser runtime where there is no meaningful host process model.
func (l *Lock) processExists() bool {
	return true
}
