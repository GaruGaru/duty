package scheduler

import (
	"github.com/GaruGaru/duty/storage"
	"testing"
)

func TestTaskManager(t *testing.T) {

	store := storage.NewMemoryStorage()

	taskManager := NewTaskManager(store)

	_, err := taskManager.Schedule(TestTask{})

	taskManager.Close()

	if err != nil{
		t.Fatal(err)
	}

	err = taskManager.Init()

	if err != nil{
		t.Fatal(err)
	}


}
