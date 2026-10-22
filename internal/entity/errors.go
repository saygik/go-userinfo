package entity

import "errors"

// ErrRefreshTokenRejected - провайдер авторизации отверг refresh-токен:
// он истёк, отозван или уже был использован. Сессию нужно завершить.
var ErrRefreshTokenRejected = errors.New("refresh token rejected")
