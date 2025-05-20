package mail

import (
	"fmt"
	"os"
	"sync"

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

var mu sync.Mutex

func (data *emailSender) Send(message *gomail.Message) {
	mu.Lock()
	dialer := gomail.NewDialer("smtp.gmail.com", 587, data.email, os.Getenv("PASS")) // secure pass later
	if err := dialer.DialAndSend(message); err != nil {
		fmt.Println(fmt.Errorf("%w", err))
	} else {
		fmt.Println("Email sent successfully!")
	}
	mu.Unlock()
}
