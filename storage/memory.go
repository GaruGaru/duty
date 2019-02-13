package storage

import (
	"fmt"
	"github.com/GaruGaru/duty/task"
)

type Memory struct {
	Tasks map[string]task.ScheduledTask
}

func NewMemoryStorage() Memory {
	return Memory{
		Tasks: make(map[string]task.ScheduledTask, 0),
	}
}

func (m Memory) Store(task task.ScheduledTask) error {
	m.Tasks[task.ID] = task
	return nil
}

func (m Memory) Update(task task.ScheduledTask, status task.Status) error {
	oldTask := m.Tasks[task.ID]
	oldTask.Status = status
	m.Tasks[task.ID] = oldTask
	return nil
}

func (m Memory) Status(id string) (task.ScheduledTask, error) {
	ctask, found := m.Tasks[id]
	if !found {
		return task.ScheduledTask{}, fmt.Errorf("task not found with id %s", id)
	}
	return ctask, nil
}

func (m Memory) ListByType(types string) ([]task.ScheduledTask, error) {
	filteredTasks := make([]task.ScheduledTask, 0)
	for _, v := range m.Tasks {
		if v.Type == types {
			filteredTasks = append(filteredTasks, v)
		}
	}
	return filteredTasks, nil
}

func (m Memory) ListAll() ([]task.ScheduledTask, error) {
	allTasks := make([]task.ScheduledTask, len(m.Tasks))
	for _, v := range m.Tasks {
		allTasks = append(allTasks, v)
	}
	return allTasks, nil
}

func (m Memory) Delete(id string) (bool, error) {
	_, found := m.Tasks[id]
	delete(m.Tasks, id)
	return found, nil
}

func (m Memory) Exists(id string) (bool, error) {
	_, p := m.Tasks[id]
	return p, nil
}

func (m Memory) Close() {
	m.Tasks = nil
}
