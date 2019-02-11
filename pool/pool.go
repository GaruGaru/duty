package pool

import (
	"fmt"
	"github.com/GaruGaru/duty/task"
	"sync"
)

var (
	DefaultWorkers   = 10
	DefaultQueueSize = 10
)

type Pool struct {
	TaskQueue      chan task.ScheduledTask
	ResultsCh      chan ScheduledTaskResult
	WorkersCount   int
	ResultCallback func(ScheduledTaskResult)
	wg             *sync.WaitGroup
}

type Options struct {
	Workers        int
	Results        chan ScheduledTaskResult
	QueueSize      int
	ResultCallback func(ScheduledTaskResult)
}

func New(opt Options) Pool {

	if opt.QueueSize <= 0 {
		opt.QueueSize = DefaultQueueSize
	}

	if opt.Workers <= 0 {
		opt.Workers = DefaultWorkers
	}

	if opt.ResultCallback == nil {
		opt.ResultCallback = func(result ScheduledTaskResult) {}
	}

	return Pool{
		TaskQueue:      make(chan task.ScheduledTask, opt.QueueSize*opt.Workers),
		ResultsCh:      opt.Results,
		WorkersCount:   opt.Workers,
		ResultCallback: opt.ResultCallback,
		wg:             &sync.WaitGroup{},
	}

}

func (p Pool) Init() {
	p.wg.Add(1)
	for i := 0; i < p.WorkersCount; i++ {
		go func() {
			for t := range p.TaskQueue {
				notified, _, _ := p.Execute(t)
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

func (p Pool) Execute(t task.ScheduledTask) (bool, ScheduledTaskResult, error) {
	p.notifyTaskResult(t, task.StatusRunning)

	err := t.Task.Run()

	status := task.StatusSuccess

	if err != nil {
		status = task.StatusError(err)
	}

	notified := p.notifyTaskResult(t, status)

	return notified, ScheduledTaskResult{
		Status:        status,
		ScheduledTask: t,
	}, err
}

func (p Pool) Close() {
	close(p.ResultsCh)
	p.wg.Done()
}

func (p Pool) Wait() {
	p.wg.Wait()
}

func (p Pool) notifyTaskResult(task task.ScheduledTask, status task.Status) bool {

	result := ScheduledTaskResult{
		Status:        status,
		ScheduledTask: task,
	}

	p.ResultCallback(result)

	select {
	case p.ResultsCh <- result:
		return true
	default:
		return false
	}

}
