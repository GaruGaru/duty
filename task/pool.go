package task

import (
	"github.com/sirupsen/logrus"
)

type ScheduledTaskResult struct {
	ScheduledTask ScheduledTask
	Status        Status
}

type Pool struct {
	BaseSize       int
	Tasks          chan ScheduledTask
	SignalsChannel chan interface{}
	WorkersCount   int32
	ResultsChan    chan ScheduledTaskResult
}

func NewWorkerPool(size int, backPressure int, resChan chan ScheduledTaskResult) Pool {
	return Pool{
		BaseSize:       size,
		Tasks:          make(chan ScheduledTask, backPressure),
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

	logrus.Warn("Received stop signal")
}

func (wp Pool) Stop() {
	wp.SignalsChannel <- true
}

func (wp Pool) Schedule(scheduledTask ScheduledTask) (bool, error) {
	select {
	case wp.Tasks <- scheduledTask:
		wp.ResultsChan <- ScheduledTaskResult{
			ScheduledTask: scheduledTask,
			Status:        StatusCreated,
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
		case task := <-wp.Tasks:

			wp.ResultsChan <- ScheduledTaskResult{
				ScheduledTask: task,
				Status:        StatusRunning,
			}

			err := task.Task.Run()

			if err != nil {
				wp.ResultsChan <- ScheduledTaskResult{
					ScheduledTask: task,
					Status:        StatusError(err),
				}
			} else {
				wp.ResultsChan <- ScheduledTaskResult{
					ScheduledTask: task,
					Status:        StatusSuccess,
				}
			}

		}
	}
}
