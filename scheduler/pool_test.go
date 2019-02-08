package scheduler

import (
	"fmt"
	"github.com/GaruGaru/duty/task"
	"testing"
)

type TestTask struct {
	Fail bool
}

func (t TestTask) Run() error {
	if t.Fail {
		return fmt.Errorf("test error")
	}
	return nil
}

func (TestTask) Type() string {
	return "test"
}

func TestPoolTaskSchedulingAndExecution(t *testing.T) {

	results := make(chan ScheduledTaskResult, 1)
	pool := NewWorkerPool(1, 1, results)

	defer pool.Stop()

	go pool.Start()

	scheduledTask := task.ScheduledTask{
		Task: TestTask{},
	}

	scheduled, err := pool.Schedule(scheduledTask)

	if !scheduled || err != nil {
		t.Fatalf("unable to schedule task")
	}

	result := <-results

	if result.ScheduledTask.ID != "" {
		t.Fatal("worker pool can't assign ID to task")
	}

	if result.Status.State != task.StatePending {
		t.Fatalf("expected task state %s but got %s", task.StatePending, result.Status.State)
	}

	result = <-results

	if result.Status.State != task.StateRunning {
		t.Fatalf("expected task state %s but got %s", task.StateRunning, result.Status.State)
	}

	result = <-results

	if result.Status.State != task.StateSuccess {
		t.Fatalf("expected task state %s but got %s", task.StateSuccess, result.Status.State)
	}

	if !result.Status.Completed {
		t.Fatalf("expected task state to be completed")
	}

	if !result.Status.Success {
		t.Fatalf("expected task state to be successful")
	}

}

func TestPoolFailingTaskSchedulingAndExecution(t *testing.T) {

	results := make(chan ScheduledTaskResult, 1)
	pool := NewWorkerPool(1, 1, results)

	defer pool.Stop()

	go pool.Start()

	scheduledTask := task.ScheduledTask{
		Task: TestTask{Fail: true},
	}

	scheduled, err := pool.Schedule(scheduledTask)

	if !scheduled || err != nil {
		t.Fatalf("unable to schedule task")
	}

	result := <-results

	if result.Status.State != task.StatePending {
		t.Fatalf("expected task state %s but got %s", task.StatePending, result.Status.State)
	}

	result = <-results

	if result.Status.State != task.StateRunning {
		t.Fatalf("expected task state %s but got %s", task.StateRunning, result.Status.State)
	}

	result = <-results

	if result.Status.State != task.StateError {
		t.Fatalf("expected task state %s but got %s", task.StateError, result.Status.State)
	}

	if !result.Status.Completed {
		t.Fatalf("expected task state to be completed")
	}

	if result.Status.Success {
		t.Fatalf("expected task state to NOT be successful")
	}
}
