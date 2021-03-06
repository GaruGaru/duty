package duty

import (
	"fmt"
	"github.com/GaruGaru/duty/pool"
	"github.com/GaruGaru/duty/storage"
	"github.com/GaruGaru/duty/task"
	"github.com/satori/go.uuid"
)

type Duty struct {
	Storage        storage.Storage
	WorkPool       pool.Pool
	StateKeeper    StateKeeper
	ResultCallback func(pool.ScheduledTaskResult)
}

type Options struct {
	ResultCallback func(pool.ScheduledTaskResult)
	Workers        int
	QueueSize      int
}

var Default = Options{}

var AssignID = ""

func New(storage storage.Storage, opt Options) Duty {
	return Duty{
		Storage:     storage,
		StateKeeper: NewStateKeeper(),
		WorkPool: pool.New(pool.Options{
			Workers:   opt.Workers,
			QueueSize: opt.QueueSize,
		}),
	}
}

func (m Duty) Init() error {
	if err := m.reconcileStatus(); err != nil {
		return err
	}
	m.WorkPool.ResultCallback = m.handleResults
	m.WorkPool.Init()
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

	if m.ResultCallback != nil {
		m.ResultCallback(result)
	}

}

func (m Duty) Execute(t task.Task) (pool.ScheduledTaskResult, error) {
	_, result, err := m.WorkPool.Execute(m.schedule(AssignID, t))
	m.handleResults(result)
	return result, err
}

func (m Duty) Enqueue(t task.Task) (task.ScheduledTask, error) {
	return m.EnqueueWithID(AssignID, t)
}

func (m Duty) EnqueueWithID(id string, t task.Task) (task.ScheduledTask, error) {
	scheduledTask := m.schedule(id, t)

	scheduled := m.WorkPool.Enqueue(scheduledTask)

	if !scheduled {
		return task.ScheduledTask{}, fmt.Errorf("task rejected, pool is full")
	}

	return scheduledTask, nil
}

func (m Duty) Get(id string) (task.ScheduledTask, error) {
	return m.Storage.Status(id)
}

func (m Duty) Tasks(id string) ([]task.ScheduledTask, error) {
	return m.Storage.ListAll()
}

func (m Duty) schedule(id string, t task.Task) task.ScheduledTask {

	if id == AssignID {
		id = uuid.NewV4().String()
	}

	return task.ScheduledTask{
		ID:   id,
		Type: t.Type(),
		Status: task.Status{
			State:     task.StateScheduled,
			Completed: false,
			Success:   false,
		},
		Task: t,
	}
}
