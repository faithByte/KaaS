package enum

type EmailFor string

const (
	None EmailFor = ""
	Job  EmailFor = "job"
	Step EmailFor = "step"
	Both EmailFor = "all"
)
