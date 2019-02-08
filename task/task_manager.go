package task

import "github.com/GaruGaru/duty/storage"

type Manager struct {
	Storage         storage.Storage
	RunningTasksMap map[string]ScheduledTask
	WorkPool        Pool
	Results         chan ScheduledTaskResult
}

func Initialize() {

}

func (m Manager) handleResults() error {
	for result := range m.Results {

		if err := m.Storage.Update(result.ScheduledTask, result.Status); err != nil {
			return err
		}

		if result.Status.Completed {
			delete(m.RunningTasksMap, result.ScheduledTask.ID)
		} else {
			m.RunningTasksMap[result.ScheduledTask.ID] = result.ScheduledTask
		}

	}

	return nil
}
