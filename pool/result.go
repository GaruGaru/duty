package pool

import "github.com/GaruGaru/duty/task"

type ScheduledTaskResult struct {
	ScheduledTask task.ScheduledTask
	Status        task.Status
}
