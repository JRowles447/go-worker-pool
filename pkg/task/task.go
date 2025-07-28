package task

import (
	"fmt"
	"math/rand/v2"
)

type Task struct {
	Name       string
	RetryCount int
	MaxRetries int
}

type SystemError struct {
	Message   string
	Task      *Task
	Retriable bool
}

func (err *SystemError) Error() string {
	return fmt.Sprintf("error %s", err.Message)
}

// simulate task failure
func SimulateError(task *Task) error {
	res := rand.Float64()
	if res < .5 {
		return nil
	} else if res < .9 {
		return &SystemError{Message: "Task failed with retriable error", Task: task, Retriable: true}
	} else {
		return &SystemError{Message: "Encountered non-retriable error", Task: task, Retriable: false}
	}
}
