package enum

type EmailFor int8

const (
	Both EmailFor = 0
	Job  EmailFor = 1
	Step EmailFor = 2
)

type EmailType int8

const (
	All     EmailType = 0
	SUCCESS EmailType = 1
	ERROR   EmailType = 2
)
