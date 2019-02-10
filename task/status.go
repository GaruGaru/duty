package task

var (
	StateScheduled = "SCHEDULED"
	StatePending   = "PENDING"
	StateRunning   = "RUNNING"
	StateSuccess   = "SUCCESS"
	StateError     = "ERROR"
)

type Status struct {
	State     string
	Completed bool
	Success   bool
	Message   string
}

var StatusPending = Status{
	State:     StatePending,
	Completed: false,
	Success:   false,
}

var StatusRunning = Status{
	State:     StateRunning,
	Completed: false,
	Success:   false,
}

var StatusSuccess = Status{
	State:     StateSuccess,
	Completed: true,
	Success:   true,
}

func StatusError(err error) Status {
	return Status{
		State:     StateError,
		Completed: true,
		Success:   false,
		Message:   err.Error(),
	}
}
