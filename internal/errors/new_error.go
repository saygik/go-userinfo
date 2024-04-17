package errors

import "errors"

func Error(err string) error {
	return errors.New(err)
}
