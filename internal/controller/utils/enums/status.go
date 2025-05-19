package enum

type Phase int8

const (
	NotStarted       Phase = 0
	ComputesCreating Phase = 1
	ComputesCreated  Phase = 2
	Launched         Phase = 3
	Completed        Phase = 4
	Error            Phase = 5
)

type Status int8

const (
	All     Status = 0
	Success Status = 4 // == phase completed
	Failure Status = 5 // == phase Error
)
