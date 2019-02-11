package duty

import (
	"github.com/GaruGaru/duty/scheduler"
	"github.com/GaruGaru/duty/storage"
)

type Duty struct {
	Manager scheduler.Manager
}

func New(store storage.Storage) Duty {
	return Duty{
		Manager: scheduler.NewTaskManager(store),
	}
}
