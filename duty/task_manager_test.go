package duty

import (
	"github.com/GaruGaru/duty/storage"
	"github.com/GaruGaru/duty/task"
	"github.com/satori/go.uuid"
	"testing"
)

func TestTaskManagerInitWithNoTasks(t *testing.T) {

	store := storage.NewMemoryStorage()

	manager := NewTaskManager(store)

	err := manager.Init()

	if err != nil {
		t.Fatal(err)
	}

}

func TestTaskManagerStatusReconcilePreviousRunningTaskAfterInit(t *testing.T) {

	taskID := uuid.NewV4().String()

	store := storage.NewMemoryStorage()

	err := store.Store(task.ScheduledTask{
		ID:     taskID,
		Type:   "test",
		Status: task.StatusRunning,
	})

	if err != nil {
		t.Fatal(err)
	}

	manager := NewTaskManager(store)

	err = manager.Init()

	if err != nil {
		t.Fatal(err)
	}

	status, err := store.Status(taskID)

	if err != nil {
		t.Fatal(err)
	}

	if status.Status.State != task.StateError {
		t.Fatalf("expected status to be %s but was %s", task.StateError, status.Status.State)
	}

}
