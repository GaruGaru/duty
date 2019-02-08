package storage

import (
	"encoding/json"
	"github.com/GaruGaru/duty/task"
)

func Serialize(task task.ScheduledTask) ([]byte, error) {
	result, err := json.Marshal(task)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func Deserialize(payload []byte) (task.ScheduledTask, error) {
	var result task.ScheduledTask
	err := json.Unmarshal(payload, &result)
	return result, err
}
