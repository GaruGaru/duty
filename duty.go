package duty

import (
	"github.com/GaruGaru/duty/scheduler"
	"github.com/GaruGaru/duty/storage"
	"github.com/GaruGaru/duty/task"
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
	return d.Manager.Schedule(t)
}

func (d Duty) RunningTasks()  {
	d.Manager.RunningTasks()
}
