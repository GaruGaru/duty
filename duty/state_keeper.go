package duty

import "github.com/GaruGaru/duty/task"

type StateKeeper struct {
	runningTasksMap map[string]task.ScheduledTask
}

func NewStateKeeper() StateKeeper {
	return StateKeeper{
		runningTasksMap: make(map[string]task.ScheduledTask, 0),
	}
}

func (sk *StateKeeper) AddRunningTask(task task.ScheduledTask) {
	sk.runningTasksMap[task.ID] = task
}

func (sk *StateKeeper) RemoveRunningTask(task task.ScheduledTask) {
	delete(sk.runningTasksMap, task.ID)
}

func (sk *StateKeeper) IsRunning(id string) bool {
	_, present := sk.runningTasksMap[id]
	return present
}

func (sk *StateKeeper) RunningTasks() []task.ScheduledTask {
	allTasks := make([]task.ScheduledTask, len(sk.runningTasksMap))
	for _, v := range sk.runningTasksMap {
		allTasks = append(allTasks, v)
	}
	return allTasks
}
