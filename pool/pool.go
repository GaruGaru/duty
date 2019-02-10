package pool

import (
	"fmt"
	"github.com/GaruGaru/duty/scheduler"
	"github.com/GaruGaru/duty/task"
)

var (
	DefaultWorkers   = 10
	DefaultQueueSize = 10
)

type Pool struct {
	TaskQueue    chan task.ScheduledTask
	ResultsCh    chan scheduler.ScheduledTaskResult
	WorkersCount int
}

type Options struct {
	Workers   int
	Results   chan scheduler.ScheduledTaskResult
	QueueSize int
}

func New(opt Options) Pool {

	if opt.QueueSize <= 0 {
		opt.QueueSize = DefaultQueueSize
	}

	if opt.Workers <= 0 {
		opt.Workers = DefaultWorkers
	}

	return Pool{
		TaskQueue:    make(chan task.ScheduledTask, opt.QueueSize*opt.Workers),
		ResultsCh:    opt.Results,
		WorkersCount: opt.Workers,
	}

}

func (p Pool) Init() {
	for i := 0; i < p.WorkersCount; i++ {
		go func() {
			for t := range p.TaskQueue {
				notified, _ := p.Execute(t)
				if !notified {
					fmt.Println("unable to notify async task result")
				}
			}
		}()
	}
}

func (p Pool) Enqueue(t task.ScheduledTask) bool {
	select {
	case p.TaskQueue <- t:
		p.notifyTaskResult(t, task.StatusPending)
		return true
	default:
		return false
	}
}

func (p Pool) Execute(t task.ScheduledTask) (bool, error) {
	p.notifyTaskResult(t, task.StatusRunning)

	err := t.Task.Run()

	status := task.StatusSuccess

	if err != nil {
		status = task.StatusError(err)
	}

	notified := p.notifyTaskResult(t, status)

	return notified, err
}

func (p Pool) Close() {
	close(p.ResultsCh)
}

func (p Pool) notifyTaskResult(task task.ScheduledTask, status task.Status) bool {
	select {
	case p.ResultsCh <- scheduler.ScheduledTaskResult{
		Status:        status,
		ScheduledTask: task,
	}:
		return true
	default:
		return false
	}
}
