package usecase

import (
	"fmt"
	"mime"
	"net/smtp"
)

// smtpServer data to smtp server
type smtpServer struct {
	host string
	port string
}

// Address URI to smtp server
func (s *smtpServer) Address() string {
	return s.host + ":" + s.port
}
func (u *UseCase) SendMail(toMail string, body string) error {
	// Sender data.
	from := "UserInfo@brnv.rw"
	// Receiver email address.
	to := []string{toMail}
	// smtp server configuration.
	smtpServer := smtpServer{host: "mail.brnv.rw", port: "25"}
	// Message.
	subject := "Срок действия вашей учетной записи истекает"
	encodedSubject := encodeSubject(subject)
	//mime := "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"

	header := make(map[string]string)
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/plain; charset=\"utf-8\""
	header["From"] = from
	header["Subject"] = encodedSubject

	// Формирование полного сообщения
	msg := ""
	for k, v := range header {
		msg += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	msg += "\r\n" + body

	// msg := []byte(encodedSubject + mime + "\n" + body)
	// Authentication.
	//	auth := smtp.PlainAuth("", from, password, smtpServer.host)

	// Sending email.
	err := smtp.SendMail(smtpServer.Address(), nil, from, to, []byte(msg))
	if err != nil {
		return err
	}

	return nil
}

func encodeSubject(subject string) string {
	// Кодируем тему письма в соответствии с RFC 2047
	return mime.QEncoding.Encode("utf-8", subject)
}
