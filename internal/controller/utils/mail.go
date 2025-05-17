package utils

import (
	"fmt"
	"os"

	gomail "gopkg.in/mail.v2"
)

type emailSender struct {
	email   string
	message *gomail.Message
}

var EmailSender = new(emailSender)

func (data *emailSender) New(email string) error {
	data.email = email
	if data.message == nil {
		data.message = gomail.NewMessage()
	}
	data.message.SetHeader("From", email)
	return nil
}

func (data *emailSender) Body() error {
	data.message.SetHeader("Subject", "Testing Class email via Gmail SMTP")
	data.message.SetBody("text/html", `
		<html>
			<body>
				<h1>Test Email</h1>
				<p><b>Hello!</b> TEST MAIL.</p>
				<p>Thanks,<br>-----------------</p>
			</body>
		</html>
	`)
	return nil
}

func (data *emailSender) Send(email string) error {
	data.message.SetHeader("To", email)

	dialer := gomail.NewDialer("smtp.gmail.com", 587, data.email, os.Getenv("PASS")) // secure pass later
	if err := dialer.DialAndSend(data.message); err != nil {
		return err
	}

	fmt.Println("Email sent successfully!")
	return nil
}
