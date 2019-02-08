package scheduler

import (
	"fmt"
	"github.com/GaruGaru/duty/storage"
	"github.com/GaruGaru/duty/task"
)

type Manager struct {
	Storage         storage.Storage
	RunningTasksMap map[string]task.ScheduledTask
	WorkPool        Pool
	Results         chan ScheduledTaskResult
}

func NewTaskManager(storage storage.Storage) Manager {
	results := make(chan ScheduledTaskResult)
	return Manager{
		Storage:         storage,
		RunningTasksMap: make(map[string]task.ScheduledTask, 1),
		WorkPool:        NewWorkerPool(10, 10, results),
		Results:         results,
	}
}

func (m Manager) Initialize() error {

	tasks, err := m.Storage.ListAll()

	if err != nil {
		return err
	}

	for _, ctask := range tasks {
		if !ctask.Status.Completed {
			ctask.Status = task.StatusError(fmt.Errorf("task terminated unexpectedly"))
			if err := m.Storage.Store(ctask); err != nil {
				return err
			}
		}
	}

	go m.WorkPool.Start()

	return nil
}

func (m Manager) Cleanup() error {
	tasks, err := m.Storage.ListAll()

	if err != nil {
		return err
	}

	for _, ctask := range tasks {
		if ctask.Status.Completed {
			ctask.Status = task.StatusError(fmt.Errorf("task terminated unexpectedly"))
			if err := m.Storage.Store(ctask); err != nil {
				return err
			}
		}
	}

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

func (m Manager) Schedule(task task.ScheduledTask) (bool, error) {
	return m.WorkPool.Schedule(task)
}