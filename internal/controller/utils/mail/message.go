package mail

import (
	"fmt"
	"strings"

	gomail "gopkg.in/mail.v2"
	// "sigs.k8s.io/controller-runtime/pkg/reconcile"

	kaasv1 "github.com/faithByte/kaas/api/v1"
	"github.com/faithByte/kaas/internal/controller/pods"
	"github.com/faithByte/kaas/internal/controller/utils"
	enum "github.com/faithByte/kaas/internal/controller/utils/enums"
	"github.com/faithByte/kaas/internal/controller/utils/interfaces"
)

type Message struct {
	message  *gomail.Message
	jobBody  string
	stepBody string
}

var Messages = make(map[string]*Message)

func NewMessage(uid, email string) {
	Messages[uid] = new(Message)
	Messages[uid].NewMessage(email)
}

func (data *Message) NewMessage(email string) error {
	if data.message == nil {
		data.message = gomail.NewMessage()
	}
	data.message.SetHeader("From", EmailSender.email)
	data.message.SetHeader("To", email)
	data.jobBody += "<html>\n<body>\n"
	return nil
}

func (data *Message) JobEndedMessage(job *kaasv1.JobSteps) error {
	if (job.Spec.Email.For != enum.Both) && (job.Spec.Email.For != enum.Job) {
		return nil
	}
	data.message.SetHeader("Subject", fmt.Sprintf("Job <%s> ended", job.Name))
	data.jobBody += "\n</body>\n</html>"
	data.message.SetBody("text/html", data.jobBody)
	return EmailSender.Send(data.message)
}

func (data *Message) StepMessage(job *kaasv1.JobSteps, stepType interfaces.Type, r utils.ReconcilerData) error {
	step := stepType.GetStepData()
	body := ""

	if (job.Spec.Email.For == enum.Both) || (job.Spec.Email.For == enum.Job) {
		body = pods.GetLogs(stepType.GetPodName(), job.Namespace, &r)
		body = strings.ReplaceAll(body, "\n", "<br />")
		data.jobBody += fmt.Sprintf("<h1>Step %s: </h1>\n<p>%s</p><hr>", step.Name, body)
	}

	if (job.Spec.Email.For != enum.Both) && (job.Spec.Email.For != enum.Step) {
		return nil
	}

	if body == "" {
		body = pods.GetLogs(stepType.GetPodName(), job.Namespace, &r)
		body = strings.ReplaceAll(body, "\n", "<br />")
	}

	data.message.SetHeader("Subject", fmt.Sprintf("Step <%s> ended successffuly", step.Name))
	data.message.SetBody("text/html", fmt.Sprintf("<html>\n<body>\n<h1>%s logs:</h1>\n<p>%s</p>\n</body>\n</html>", step.Name, body))
	return EmailSender.Send(data.message)
}
