package mssql

import (
	"database/sql"
	"errors"
)

func sqlRowAffectedErrorWrapper(res sql.Result, err error) error {
	if err != nil {
		return err
	}
	if res != nil {
		ra, err1 := res.RowsAffected()
		if err1 != nil {
			return err1
		}
		if ra < 1 {
			return errors.New("add row to sql not affected")
		}
		return nil
	}
	return nil
}
