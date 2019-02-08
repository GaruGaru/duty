package task

import (
	"fmt"
	"github.com/GaruGaru/duty/storage"
)

type Manager struct {
	Storage         storage.Storage
	RunningTasksMap map[string]ScheduledTask
	WorkPool        Pool
	Results         chan ScheduledTaskResult
}

func NewTaskManager(storage storage.Storage) Manager {
	results := make(chan ScheduledTaskResult)
	return Manager{
		Storage:         storage,
		RunningTasksMap: make(map[string]ScheduledTask, 0),
		WorkPool:        NewWorkerPool(10, 10, results),
		Results:         results,
	}
}

func (m Manager) Initialize() error {

	tasks, err := m.Storage.ListAll()

	if err != nil {
		return err
	}

	for _, task := range tasks {
		if !task.Status.Completed {
			task.Status = StatusError(fmt.Errorf("task terminated unexpectedly"))
			if err := m.Storage.Store(task); err != nil {
				return err
			}
		}
	}

	go m.WorkPool.Start()

	return nil
}

func (m Manager) handleResults() error {
	for result := range m.Results {

		if err := m.Storage.Update(result.ScheduledTask, result.Status); err != nil {
			return err
		}

		if result.Status.Completed {
			delete(m.RunningTasksMap, result.ScheduledTask.ID)
		} else {
			m.RunningTasksMap[result.ScheduledTask.ID] = result.ScheduledTask
		}

	}

	return nil
}

func (m Manager) Schedule(task ScheduledTask) (bool, error) {
	return m.WorkPool.Schedule(task)
}
