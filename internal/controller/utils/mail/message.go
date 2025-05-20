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
	jobBody          string
	runningProcesses utils.AtomicCounter
}

var Messages = make(map[string]*Message)

func NewMessage(uid string, jobSpec kaasv1.JobStepsSpec) {
	if jobSpec.Email.For == enum.None {
		return
	}
	Messages[uid] = new(Message)
	Messages[uid].jobBody = "<html>\n<body>\n"
}

func (data *Message) NewMessage(email string) *gomail.Message {
	message := gomail.NewMessage()
	message.SetHeader("From", EmailSender.email)
	message.SetHeader("To", email)
	return message
}

func JobEmail(job kaasv1.JobSteps) {
	if job.Spec.Email.For == enum.None {
		return
	}

	uid := string(job.GetUID())

	for Messages[uid].runningProcesses.Value() != 0 {
	}

	if (job.Spec.Email.For != enum.Both) && (job.Spec.Email.For != enum.Job) {
		delete(Messages, uid)
		return
	}
	Messages[uid].jobEndedMessage(job)
	delete(Messages, uid)
}

func (data *Message) jobEndedMessage(job kaasv1.JobSteps) {
	// protected on JobEmail
	message := data.NewMessage(job.Spec.Email.Email)
	message.SetHeader("Subject", fmt.Sprintf("Job <%s> ended", job.Name))
	data.jobBody += "\n</body>\n</html>"
	message.SetBody("text/html", data.jobBody)
	EmailSender.Send(message)
}

func (data *Message) StepMessage(job kaasv1.JobSteps, stepType interfaces.Type, r utils.ReconcilerData) {
	if job.Spec.Email.For == enum.None {
		return
	}
	data.runningProcesses.Increment()
	step := stepType.GetStepData()
	body := ""

	if (job.Spec.Email.For == enum.Both) || (job.Spec.Email.For == enum.Job) {
		body = pods.GetLogs(stepType.GetPodName(), job.Namespace, &r)
		body = strings.ReplaceAll(body, "\n", "<br />")
		data.jobBody += fmt.Sprintf("<h1>Step %s: </h1>\n<p>%s</p><hr>", step.Name, body)
	}

	if body == "" {
		body = pods.GetLogs(stepType.GetPodName(), job.Namespace, &r)
		body = strings.ReplaceAll(body, "\n", "<br />")
	}

	message := data.NewMessage(job.Spec.Email.Email)
	if stepType.GetPhase() == enum.Completed {
		message.SetHeader("Subject", fmt.Sprintf("Step <%s> ended successffuly", step.Name))
	} else {
		message.SetHeader("Subject", fmt.Sprintf("Step <%s> failed", step.Name))
	}
	message.SetBody("text/html", fmt.Sprintf("<html>\n<body>\n<h1>%s logs:</h1>\n<p>%s</p>\n</body>\n</html>", step.Name, body))
	EmailSender.Send(message)
	data.runningProcesses.Decrement()
}
