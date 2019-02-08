package storage

import (
	"github.com/GaruGaru/duty/task"
	"io/ioutil"
	"os"
	"testing"
)

func testStorage() (*BoltDB, string) {

	file, err := ioutil.TempFile("", "bolt-")

	if err != nil {
		panic(err)
	}

	defer os.Remove(file.Name())

	db, err := NewBoltDB(file.Name())

	if err != nil {
		panic(err)
	}

	return db, file.Name()
}

func TestStoreScheduledTask(t *testing.T) {

	db, file := testStorage()

	defer os.Remove(file)

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

func TestStoreGetScheduledTask(t *testing.T) {

	db, file := testStorage()

	defer os.Remove(file)

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

	retrivedTask, err := db.Status( taskID)

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
