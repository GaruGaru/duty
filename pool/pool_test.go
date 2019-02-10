package pool

import (
	"fmt"
	"github.com/GaruGaru/duty/scheduler"
	"github.com/GaruGaru/duty/task"
	"github.com/satori/go.uuid"
	"testing"
)

type TestTask struct {
	Fn func() error
}

func (t TestTask) Run() error {
	return t.Fn()
}

func (TestTask) Type() string {
	return "test"
}

func schedule(t task.Task) task.ScheduledTask {
	return task.ScheduledTask{
		ID:     uuid.NewV4().String(),
		Type:   t.Type(),
		Status: task.StatusPending,
		Task:   t,
	}
}

func TestPoolSyncTaskExecution(t *testing.T) {

	pool := New(Options{})

	notified, err := pool.Execute(schedule(TestTask{
		Fn: func() error {
			return nil
		},
	}))

	if err != nil {
		t.Fatalf("unexpected error executing task: %s", err)
	}

	if notified {
		t.Fatalf("worker pool should not notify task result since result channel has not been defined")
	}

}

func TestPoolSyncTaskExecutionWithWithError(t *testing.T) {

	pool := New(Options{})

	notified, err := pool.Execute(schedule(TestTask{
		Fn: func() error {
			return fmt.Errorf("test error")
		},
	}))

	if err == nil {
		t.Fatalf("expected error executing task")
	}

	if notified {
		t.Fatalf("worker pool should not notify task result since result channel has not been defined")
	}

}

func TestPoolSyncTaskExecutionResultChannel(t *testing.T) {

	results := make(chan scheduler.ScheduledTaskResult, 5)

	pool := New(Options{
		Results: results,
	})

	notified, err := pool.Execute(schedule(TestTask{
		Fn: func() error {
			return nil
		},
	}))

	if err != nil {
		t.Fatalf("unexpected error executing task")
	}

	if !notified {
		t.Fatalf("worker pool should notify task result since result channel has been defined")
	}

	result := <-results

	if result.Status.State != task.StateRunning {
		t.Fatalf("expected task state to be %s but got %s", task.StateRunning, result.Status.State)
	}

	result = <-results

	if result.Status.State != task.StateSuccess {
		t.Fatalf("expected task state to be %s but got %s", task.StateSuccess, result.Status.State)
	}

	if !result.Status.Completed {
		t.Fatal("expected task to be completed")
	}

	if !result.Status.Success {
		t.Fatal("expected task to be successful")
	}

}

func TestPoolSyncTaskExecutionWithErrorResultChannel(t *testing.T) {

	results := make(chan scheduler.ScheduledTaskResult, 5)

	pool := New(Options{
		Results: results,
	})

	notified, err := pool.Execute(schedule(TestTask{
		Fn: func() error {
			return fmt.Errorf("test error")
		},
	}))

	if err == nil {
		t.Fatalf("expected error executing task")
	}

	if !notified {
		t.Fatalf("worker pool should notify task result since result channel has been defined")
	}

	result := <-results

	if result.Status.State != task.StateRunning {
		t.Fatalf("expected task state to be %s but got %s", task.StateRunning, result.Status.State)
	}

	result = <-results

	if result.Status.State != task.StateError {
		t.Fatalf("expected task state to be %s but got %s", task.StateSuccess, result.Status.State)
	}

	if !result.Status.Completed {
		t.Fatal("expected task to be completed")
	}

	if result.Status.Success {
		t.Fatal("expected task to be failed")
	}

	if result.Status.Message == "" {
		t.Fatal("expected error message if the task is failed")
	}

}
