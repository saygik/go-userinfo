package usecase

import "errors"

func (u *UseCase) Error(err string) error {
	return errors.New(err)
}
