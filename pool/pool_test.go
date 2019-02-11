package pool

import (
	"fmt"
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

	notified, status, err := pool.Execute(schedule(TestTask{
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

	if !status.Status.Success {
		t.Fatal("expected state to be successful")
	}

	if !status.Status.Completed {
		t.Fatal("expected state to be completed")
	}

}

func TestPoolSyncTaskExecutionWithWithError(t *testing.T) {

	pool := New(Options{})

	notified, status, err := pool.Execute(schedule(TestTask{
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

	if status.Status.Success {
		t.Fatal("expected state to be not successful")
	}

	if !status.Status.Completed {
		t.Fatal("expected state to be completed")
	}

}

func TestPoolSyncTaskExecutionResultChannel(t *testing.T) {

	results := make(chan ScheduledTaskResult, 5)

	pool := New(Options{
		Results: results,
	})

	notified, status, err := pool.Execute(schedule(TestTask{
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

	if !status.Status.Success {
		t.Fatal("expected state to be successful")
	}

	if !status.Status.Completed {
		t.Fatal("expected state to be completed")
	}

	testChannelResults(t, results, PoolTestOutcome{
		Steps:      []string{task.StateRunning, task.StateSuccess},
		Successful: true,
	})
}

func TestPoolSyncTaskExecutionWithErrorResultChannel(t *testing.T) {

	results := make(chan ScheduledTaskResult, 5)

	pool := New(Options{
		Results: results,
	})

	notified, status, err := pool.Execute(schedule(TestTask{
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

	if status.Status.Success {
		t.Fatal("expected state to be not successful")
	}

	if !status.Status.Completed {
		t.Fatal("expected state to be completed")
	}

	testChannelResults(t, results, PoolTestOutcome{
		Steps:      []string{task.StateRunning, task.StateError},
		Successful: false,
	})
}

func TestPoolAsyncTaskExecution(t *testing.T) {

	results := make(chan ScheduledTaskResult, 5)

	pool := New(Options{
		Results: results,
	})

	pool.Init()

	scheduled := pool.Enqueue(schedule(TestTask{
		Fn: func() error {
			return nil
		},
	}))

	if !scheduled {
		t.Fatal("unable to schedule task")
	}

	testChannelResults(t, results, PoolTestOutcome{
		Steps:      []string{task.StatePending, task.StateRunning, task.StateSuccess},
		Successful: true,
	})

}

func TestPoolAsyncTaskExecutionWithError(t *testing.T) {

	results := make(chan ScheduledTaskResult, 5)

	pool := New(Options{
		Results: results,
	})

	pool.Init()

	scheduled := pool.Enqueue(schedule(TestTask{
		Fn: func() error {
			return fmt.Errorf("test error")
		},
	}))

	if !scheduled {
		t.Fatal("unable to schedule task")
	}

	testChannelResults(t, results, PoolTestOutcome{
		Steps:      []string{task.StatePending, task.StateRunning, task.StateError},
		Successful: false,
	})

}

type PoolTestOutcome struct {
	Steps      []string
	Successful bool
}

func testChannelResults(t *testing.T, results chan ScheduledTaskResult, testConfig PoolTestOutcome) {

	var outcome task.Status
	for _, expectedState := range testConfig.Steps {
		result := <-results

		if result.Status.State != expectedState {
			t.Fatalf("expected task state to be %s but got %s", expectedState, result.Status.State)
		}

		outcome = result.Status

	}

	if !outcome.Completed {
		t.Fatal("expected task to be completed after execution")
	}

	if outcome.Success != testConfig.Successful {
		t.Fatalf("expected task success to be %t but was %t", testConfig.Successful, outcome.Success)
	}

	if !outcome.Success && outcome.Message == "" {
		t.Fatal("expecting task result to have a message in case of fal")
	}

}
