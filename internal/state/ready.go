package state

import (
	"sync/atomic"
	"time"
)

var initialized atomic.Bool
var fillingRedis atomic.Bool
var startedAt = time.Now()
var lastADFill atomic.Int64

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

// StartedAt возвращает время запуска процесса приложения.
func StartedAt() time.Time {
	return startedAt
}

// SetLastADFill запоминает время последнего успешного заполнения кеша AD.
func SetLastADFill(t time.Time) {
	lastADFill.Store(t.Unix())
}

// LastADFill возвращает время последнего успешного заполнения кеша AD.
// Нулевое значение означает, что за время работы процесса заполнение ещё не завершалось.
func LastADFill() time.Time {
	unix := lastADFill.Load()
	if unix == 0 {
		return time.Time{}
	}
	return time.Unix(unix, 0)
}
