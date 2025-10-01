package state

import "sync/atomic"

var initialized atomic.Bool

func IsInitialized() bool {
	return initialized.Load()
}

func SetInitialized() {
	initialized.Store(true)
}
