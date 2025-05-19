package mail

import (
	"fmt"
	"os"

	// "golang.org/x/text/message"
	gomail "gopkg.in/mail.v2"
)

type emailSender struct {
	email string
	// pass  string
}

var EmailSender = new(emailSender)

func (data *emailSender) NewEmail(email string) {
	data.email = email
}

func (data *emailSender) Send(message *gomail.Message) error {
	dialer := gomail.NewDialer("smtp.gmail.com", 587, data.email, os.Getenv("PASS")) // secure pass later
	if err := dialer.DialAndSend(message); err != nil {
		return err
	}

	fmt.Println("Email sent successfully!")
	return nil
}
