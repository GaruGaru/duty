package storage

import (
	"github.com/GaruGaru/duty/task"
	"testing"
)

func TestInMemoryStoreScheduledTask(t *testing.T) {

	db := NewMemoryStorage()

	scheduledTask := task.ScheduledTask{
		ID:   "0000-0000-0000-0000",
		Type: "test",
		Status: task.Status{
			State:     "SCHEDULED",
			Completed: false,
			Success:   false,
		},
	}

	err := db.Store(scheduledTask)

	if err != nil {
		t.Fatal(err)
		t.FailNow()
	}

}

func TestInMemoryStoreGetScheduledTask(t *testing.T) {

	db := NewMemoryStorage()

	taskID := "0000-0000-0000-0000"

	scheduledTask := task.ScheduledTask{
		ID:   taskID,
		Type: "test",
		Status: task.Status{
			State:     "SCHEDULED",
			Completed: false,
			Success:   false,
		},
	}

	err := db.Store(scheduledTask)

	if err != nil {
		t.Fatal(err)
		t.FailNow()
	}

	retrivedTask, err := db.Status(taskID)

	if err != nil {
		t.Fatal(err)
		t.FailNow()
	}

	if retrivedTask.ID != taskID {
		t.Fatalf("expecting id %s but got %s", retrivedTask.ID, taskID)
		t.FailNow()
	}

	if retrivedTask.Type != "test" {
		t.Fatalf("expecting task arg: %s but got %s", "test", retrivedTask.Type)
		t.FailNow()
	}

}
