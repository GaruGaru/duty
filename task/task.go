package task

type Task interface {
	Run() error
	Type() string
}

type ScheduledTask struct {
	ID     string
	Type   string
	Status Status
	Task   Task `json:"-"`
}
