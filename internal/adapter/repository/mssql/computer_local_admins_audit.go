package mssql

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
)

func (r *Repository) ComputerLocalAdminsAudit(computer string, admins []string, isDomain bool) error {
	// 🔥 Заменили JSON на строку через запятую
	adminsJSON := strings.Join(admins, ",")

	lastRecord := &LastRecord{}
	err := r.db.SelectOne(lastRecord, `
        SELECT administrators, [date]
        FROM computerLocalAdminsAudit
        WHERE computer = ? AND domain = ?
        ORDER BY [date] DESC
        OFFSET 0 ROWS FETCH NEXT 1 ROWS ONLY
    `, computer, isDomain)

	if err == sql.ErrNoRows {
		_, err = r.db.Exec(`
            INSERT INTO computerLocalAdminsAudit (computer, administrators, domain, [date])
            VALUES (?, ?, ?, GETDATE())
        `, computer, adminsJSON, isDomain) // ← строка, не JSON
		if err != nil {
			return fmt.Errorf("insert error: %v", err)
		}
		fmt.Printf("✅ Новая запись: %s (domain=%t): %s\n", computer, isDomain, adminsJSON)
		return nil
	}

	if err != nil {
		return fmt.Errorf("select error: %v", err)
	}

	if lastRecord.Administrators != adminsJSON {
		_, err = r.db.Exec(`
            INSERT INTO computerLocalAdminsAudit (computer, administrators, domain, [date])
            VALUES (?, ?, ?, GETDATE())
        `, computer, adminsJSON, isDomain)
		if err != nil {
			return fmt.Errorf("insert changed error: %v", err)
		}
		fmt.Printf("✅ Обновлено: %s (domain=%t): %s → %s\n", computer, isDomain, lastRecord.Administrators, adminsJSON)
	}

	return nil
}

// Структура для gorp
type LastRecord struct {
	Administrators string    `db:"administrators"`
	Date           time.Time `db:"date"`
}
