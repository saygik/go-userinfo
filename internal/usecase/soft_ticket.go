package usecase

import (
	"fmt"
	"strings"
)

const softTicketIDMaxLen = 20

// normalizeSoftTicketID проверяет необязательный ID заявки GLPI.
// Пустое значение допустимо, непустое должно состоять только из цифр.
func normalizeSoftTicketID(value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", nil
	}
	if len(value) > softTicketIDMaxLen {
		return "", fmt.Errorf("ID заявки в GLPI не должен быть длиннее %d символов", softTicketIDMaxLen)
	}
	for _, r := range value {
		if r < '0' || r > '9' {
			return "", fmt.Errorf("ID заявки в GLPI должен содержать только цифры")
		}
	}
	return value, nil
}
