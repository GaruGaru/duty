package duty

import (
	"fmt"
	"github.com/GaruGaru/duty/pool"
	"github.com/GaruGaru/duty/storage"
	"github.com/GaruGaru/duty/task"
	"github.com/satori/go.uuid"
)

type Duty struct {
	Storage     storage.Storage
	WorkPool    pool.Pool
	StateKeeper StateKeeper
}

func NewTaskManager(storage storage.Storage) Duty {
	return Duty{
		Storage:     storage,
		StateKeeper: NewStateKeeper(),
		WorkPool:    pool.New(pool.Options{}),
	}
}

func (m Duty) Init() error {
	if err := m.reconcileStatus(); err != nil {
		return err
	}
	m.WorkPool.Init()
	m.WorkPool.ResultCallback = m.handleResults
	return nil
}

func (m Duty) Close() {
	m.Storage.Close()
	m.WorkPool.Close()
}

func (m Duty) reconcileStatus() error {
	storedTasks, err := m.Storage.ListAll()

	if err != nil {
		return err
	}

	for _, sTask := range storedTasks {
		if sTask.Status.State == task.StateRunning && !m.StateKeeper.IsRunning(sTask.ID) {
			sTask.Status = task.StatusError(fmt.Errorf("task terminated unexpectedly"))

			err := m.Storage.Store(sTask)

			if err != nil {
				return err
			}

			m.StateKeeper.RemoveRunningTask(sTask)
		}
	}

	return nil
}

func (m Duty) handleResults(result pool.ScheduledTaskResult) {

	if err := m.Storage.Update(result.ScheduledTask, result.Status); err != nil {
		panic(err)
	}

	if result.Status.Completed {
		m.StateKeeper.RemoveRunningTask(result.ScheduledTask)
	} else {
		m.StateKeeper.AddRunningTask(result.ScheduledTask)
	}

}

func (m Duty) Execute(t task.Task) (pool.ScheduledTaskResult, error) {
	_, result, err := m.WorkPool.Execute(m.schedule(t))
	m.handleResults(result)
	return result, err
}

func (m Duty) Enqueue(t task.Task) (task.ScheduledTask, error) {
	scheduledTask := m.schedule(t)

	scheduled := m.WorkPool.Enqueue(scheduledTask)

	if !scheduled {
		return task.ScheduledTask{}, fmt.Errorf("task rejected, pool is full")
	}

	return scheduledTask, nil
}

func (m Duty) schedule(t task.Task) task.ScheduledTask {
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
