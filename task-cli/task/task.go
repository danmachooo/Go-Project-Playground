package task

import (
	"github.com/google/uuid"
)

type Task struct {
	Id          string
	Name        string
	Description string
	IsCompleted bool
}

type Data struct {
	Tasks []Task `json:"tasks"`
}

func NewTask(name, description string) Task {
	return Task{
		Id:          uuid.New().String(),
		Name:        name,
		Description: description,
		IsCompleted: false,
	}
}