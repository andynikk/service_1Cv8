package cli

import (
	"Service_1Cv8/internal/constants"
)

type Task struct {
	status      constants.Status
	dbName      string
	uid         string
	title       string
	description string
}

type Checkbox struct {
	cursor   int
	choices  []string
	selected map[int]struct{}
}

func NewTask(status constants.Status, title, description string) Task {
	return Task{status: status, title: title, description: description}
}

func (t Task) UpdateTitle(val *string) {
	t.title = *val
}

func (t *Task) Next() {
	if t.status == constants.Services {
		t.status = constants.Blocking
	} else {
		t.status++
	}
}

// FilterValue the list.Item interface
func (t Task) FilterValue() string {
	return t.title
}

func (t Task) Title() string {
	return t.title
}

func (t Task) Description() string {
	return t.description
}
