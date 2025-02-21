package usecase

import "net/smtp"

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
	msg := []byte("From:" + from + "\r\n" +
		"Subject: Срок действия вашей учетной записи истекает\r\n\r\n" +
		body + "\r\n")
	// Authentication.
	//	auth := smtp.PlainAuth("", from, password, smtpServer.host)

	// Sending email.
	err := smtp.SendMail(smtpServer.Address(), nil, from, to, msg)
	if err != nil {
		return err
	}

	return nil
}
