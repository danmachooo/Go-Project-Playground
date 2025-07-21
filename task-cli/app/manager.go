package app

import (
	"task-cli/task"

)

func AddTask(task task.Task, tasks *[]task.Task) {
	*tasks = append(*tasks, task)
}

func MarkTask(taskNumber int, tasks []task.Task) (*task.Task, bool) {
	if taskNumber < 1 || taskNumber > len(tasks) {
		return nil, false
	}

	idx := taskNumber - 1
	tasks[idx].IsCompleted = !tasks[idx].IsCompleted
	return &tasks[idx], true
}

func RemoveTask(taskNumber int, tasks []task.Task) ([]task.Task, bool) {
	if taskNumber < 1 || taskNumber > len(tasks) {
		return nil, false
	}

	idx := taskNumber - 1
	newTasks := append(tasks[:idx], tasks[idx+1:]...)
	return newTasks, true
}