package enum

type Status int8

const (
	NotStarted       Status = 0
	ComputesCreating Status = 1
	ComputesCreated  Status = 2
	Launched         Status = 3
	Completed        Status = 4
	Error            Status = 5
)
