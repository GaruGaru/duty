package scheduler

import (
	"github.com/GaruGaru/duty/storage"
	"testing"
)

func TestTaskManager(t *testing.T) {

	store := storage.NewMemoryStorage()

	taskManager := NewTaskManager(store)



}
