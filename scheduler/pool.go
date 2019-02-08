package scheduler

import (
	"github.com/GaruGaru/duty/task"
)

type ScheduledTaskResult struct {
	ScheduledTask task.ScheduledTask
	Status        task.Status
}

type Pool struct {
	BaseSize       int
	Tasks          chan task.ScheduledTask
	SignalsChannel chan interface{}
	WorkersCount   int32
	ResultsChan    chan ScheduledTaskResult
}

func NewWorkerPool(size int, backPressure int, resChan chan ScheduledTaskResult) Pool {
	return Pool{
		BaseSize:       size,
		Tasks:          make(chan task.ScheduledTask, backPressure),
		SignalsChannel: make(chan interface{}),
		WorkersCount:   0,
		ResultsChan:    resChan,
	}
}

func (wp Pool) Start() {
	for i := 0; i < wp.BaseSize; i++ {
		go wp.startWorker()
	}

	<-wp.SignalsChannel
}

func (wp Pool) Stop() {
	wp.SignalsChannel <- true
}

func (wp Pool) Schedule(scheduledTask task.ScheduledTask) (bool, error) {
	select {
	case wp.Tasks <- scheduledTask:
		wp.ResultsChan <- ScheduledTaskResult{
			ScheduledTask: scheduledTask,
			Status:        task.StatusCreated,
		}
		return true, nil
	default:
		return false, nil
	}
}

func (wp *Pool) startWorker() {
	for {
		select {
		case <-wp.SignalsChannel:
			close(wp.Tasks)
			close(wp.SignalsChannel)
			return
		case ctask := <-wp.Tasks:

			wp.ResultsChan <- ScheduledTaskResult{
				ScheduledTask: ctask,
				Status:        task.StatusRunning,
			}

			err := ctask.Task.Run()

			if err != nil {
				wp.ResultsChan <- ScheduledTaskResult{
					ScheduledTask: ctask,
					Status:        task.StatusError(err),
				}
			} else {
				wp.ResultsChan <- ScheduledTaskResult{
					ScheduledTask: ctask,
					Status:        task.StatusSuccess,
				}
			}

		}
	}
}
