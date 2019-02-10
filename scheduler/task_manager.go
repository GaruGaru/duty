package scheduler

import (
	"fmt"
	"github.com/GaruGaru/duty/storage"
	"github.com/GaruGaru/duty/task"
	"github.com/satori/go.uuid"
)

type ResultCallback func(ScheduledTaskResult)

type Manager struct {
	Storage         storage.Storage
	RunningTasksMap map[string]task.ScheduledTask
	WorkPool        Pool
	Results         chan ScheduledTaskResult
	ResultCallback  ResultCallback
}

func NewTaskManager(storage storage.Storage) Manager {
	results := make(chan ScheduledTaskResult)
	return Manager{
		Storage:         storage,
		RunningTasksMap: make(map[string]task.ScheduledTask, 1),
		WorkPool:        NewWorkerPool(10, 10, results),
		Results:         results,
		ResultCallback:  func(result ScheduledTaskResult) {},
	}
}

func (m Manager) Init() error {

	go m.WorkPool.Start()

	if err := m.handleResults(); err != nil {
		return err
	}

	return nil
}

func (m Manager) OnTaskResult(callback ResultCallback) {
	m.ResultCallback = callback
}

func (m Manager) Close()  {
	m.Storage.Close()
	m.WorkPool.Stop()
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

		m.ResultCallback(result)

	}

	return nil
}

func (m Manager) Schedule(t task.Task) (task.ScheduledTask, error) {
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

	scheduled, err := m.WorkPool.Schedule(scheduledTask)

	if !scheduled {
		return task.ScheduledTask{}, fmt.Errorf("task rejected, pool is full")
	}

	return scheduledTask, err
}

func (m Manager) AllTasks() ([]task.ScheduledTask, error) {
	return m.Storage.ListAll()
}

func (m Manager) RunningTasks() []task.ScheduledTask {
	allTasks := make([]task.ScheduledTask, len(m.RunningTasksMap))
	for _, v := range m.RunningTasksMap {
		allTasks = append(allTasks, v)
	}
	return allTasks
}
