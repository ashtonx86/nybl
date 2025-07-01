package dependencies

import (
	"context"
	"fmt"
	"os"
	"strconv"

	gomail "gopkg.in/mail.v2"
)

const BasicTemplate = `
<html>
<header>
<style>
 %s
</style>
</header>
<body>
	%s
</body>
</html>
`

type MailSingleton struct {
	Dialer *gomail.Dialer
}

func NewMailSingleton() *MailSingleton {
	host := os.Getenv("EMAIL_SMTP_HOST")
	port := os.Getenv("EMAIL_SMTP_PORT")
	username := os.Getenv("EMAIL_SMTP_USERNAME")
	password := os.Getenv("EMAIL_SMTP_PASSWORD")

	portInt, err := strconv.Atoi(port)
	if err != nil {
		panic("EMAIL_SMTP_PORT must be an integer")
	}

	dialer := gomail.NewDialer(host, portInt, username, password)
	return &MailSingleton{
		Dialer: dialer,
	}
}
                                    
func (s *MailSingleton) Init(ctx context.Context) error {
	return nil 
}

func (s *MailSingleton) Stop(ctx context.Context) error {
	return nil
}
func (s *MailSingleton) SendMail(from string, to string, subject string, style string, body string) error {
	message := gomail.NewMessage()

	message.SetHeader("from", from)
	message.SetHeader("to", to)
	message.SetHeader("subject", subject)

	message.SetBody("text/html", s.MakeEmailHTML(style, body))

	return s.Dialer.DialAndSend(message)
}

func (s *MailSingleton) MakeEmailHTML(style string, body string) string {
	return fmt.Sprintf(BasicTemplate, style, body)
}