package scheduler

import (
	"fmt"
	"github.com/GaruGaru/duty/pool"
	"github.com/GaruGaru/duty/storage"
	"github.com/GaruGaru/duty/task"
	"github.com/satori/go.uuid"
)

type Manager struct {
	Storage         storage.Storage
	RunningTasksMap map[string]task.ScheduledTask
	WorkPool        pool.Pool
}

func NewTaskManager(storage storage.Storage) Manager {
	return Manager{
		Storage:         storage,
		RunningTasksMap: make(map[string]task.ScheduledTask, 1),
		WorkPool:        pool.New(pool.Options{}),
	}
}

func (m Manager) Init() {
	m.WorkPool.Init()
	m.WorkPool.ResultCallback = m.handleResults
}

func (m Manager) Close() {
	m.Storage.Close()
	m.WorkPool.Close()
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

func (m Manager) handleResults(result pool.ScheduledTaskResult) {

	if err := m.Storage.Update(result.ScheduledTask, result.Status); err != nil {
		panic(err)
	}

	if result.Status.Completed {
		delete(m.RunningTasksMap, result.ScheduledTask.ID)
	} else {
		m.RunningTasksMap[result.ScheduledTask.ID] = result.ScheduledTask
	}

}

func (m Manager) Execute(t task.Task) (pool.ScheduledTaskResult, error) {
	_, result, err := m.WorkPool.Execute(m.schedule(t))
	return result, err
}

func (m Manager) Enqueue(t task.Task) (task.ScheduledTask, error) {
	scheduledTask := m.schedule(t)

	scheduled := m.WorkPool.Enqueue(scheduledTask)

	if !scheduled {
		return task.ScheduledTask{}, fmt.Errorf("task rejected, pool is full")
	}

	return scheduledTask, nil
}

func (m Manager) schedule(t task.Task) task.ScheduledTask {
	return task.ScheduledTask{
		ID:   uuid.NewV4().String(),
		Type: t.Type(),
		Status: task.Status{
			State:     task.StateScheduled,
			Completed: false,
			Success:   false,
		},
		Task: t,
	}
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
