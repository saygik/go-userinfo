package usecase

import "github.com/saygik/go-userinfo/internal/state"

func (u *UseCase) IsAppInitialized() bool {
	return state.IsInitialized()
}
