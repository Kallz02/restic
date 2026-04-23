//go:build !js

package restic

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/restic/restic/internal/debug"
)

// listen for incoming SIGHUP and ignore
var ignoreSIGHUP sync.Once

func init() {
	ignoreSIGHUP.Do(func() {
		go func() {
			c := make(chan os.Signal, 1)
			signal.Notify(c, syscall.SIGHUP)
			for s := range c {
				debug.Log("Signal received: %v\n", s)
			}
		}()
	})
}
