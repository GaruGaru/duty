package storage

import "github.com/GaruGaru/duty/task"

type Storage interface {
	Store(task task.ScheduledTask) error

	Update(task.ScheduledTask, task.Status) error

	Status(id string) (task.ScheduledTask, error)

	ListByType(types string) ([]task.ScheduledTask, error)

	ListAll() ([]task.ScheduledTask, error)

	Delete(id string) (bool, error)

	Close()
}
