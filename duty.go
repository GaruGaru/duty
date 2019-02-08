package duty

import (
	"github.com/GaruGaru/duty/task"
	"github.com/satori/go.uuid"
)

type Duty struct {
	WorkPool task.Pool
}

func (d Duty) Schedule(t task.Task) (bool, error) {
	scheduledTask := task.ScheduledTask{
		ID:   uuid.NewV4().String(),
		Type: t.Type(),
		Status: task.Status{
			State:     "SCHEDULED",
			Completed: false,
			Success:   false,
		},
		Task: t,
	}

	scheduled, err := d.WorkPool.Schedule(scheduledTask)

	if !scheduled {
		return false, err
	}

	return true, err
}
