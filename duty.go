package duty

import (
	"github.com/GaruGaru/duty/scheduler"
	"github.com/GaruGaru/duty/storage"
	"github.com/GaruGaru/duty/task"
	"github.com/satori/go.uuid"
)

type Duty struct {
	Manager scheduler.Manager
}

func New(store storage.Storage) Duty {
	return Duty{
		Manager: scheduler.NewTaskManager(store),
	}
}

func (d Duty) Schedule(t task.Task) (bool, error) {
	scheduledTask := task.ScheduledTask{
		ID:   uuid.NewV4().String(),
		Type: t.Type(),
		Status: task.Status{
			State:     task.StateScheduled,
			Completed: false,
			Success:   false,
		},
		Task: t,
	}

	scheduled, err := d.Manager.Schedule(scheduledTask)

	if !scheduled {
		return false, err
	}

	return true, err
}
