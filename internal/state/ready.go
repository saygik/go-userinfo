package state

import "sync/atomic"

var initialized atomic.Bool
var fillingRedis atomic.Bool

func IsInitialized() bool {
	return initialized.Load()
}

func SetInitialized() {
	initialized.Store(true)
}

func IsFillingRedis() bool {
	return fillingRedis.Load()
}

func SetFillingRedis(filling bool) {
	fillingRedis.Store(filling)
}
