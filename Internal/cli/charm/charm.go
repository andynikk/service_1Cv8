package charm

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"

	"Service_1Cv8/internal/cli/styles"
	"Service_1Cv8/internal/constants"
)

type Task struct {
	Status          constants.Status
	OrdinalNumber   int
	DBName          string
	Uid             string
	TaskTitle       string
	TaskDescription string
}

type Checkbox struct {
	Cursor   int
	Choices  []string
	Selected map[int]struct{}
}

type CrawlElements map[int]interface{}

type BoolModel struct {
	Value  bool
	Active bool
}

func NewTask(status constants.Status, title, description string) Task {
	return Task{Status: status, TaskTitle: title, TaskDescription: description}
}

func (t Task) UpdateTitle(val *string) {
	t.TaskTitle = *val
}

func (t *Task) Next() {
	if t.Status == constants.Services {
		t.Status = constants.Blocking
	} else {
		t.Status++
	}
}

// FilterValue the list.Item interface
func (t Task) FilterValue() string {
	return t.TaskTitle
}

func (t Task) Title() string {
	return t.TaskTitle
}

func (t Task) Description() string {
	return t.TaskDescription
}

func TеraversingFormElements(ce CrawlElements, focusIndex int) {
	for kElement, vElement := range ce {
		switch vElement.(type) {
		case *textinput.Model:
			model := vElement.(*textinput.Model)
			model.Blur()
			model.PromptStyle = styles.NoStyleFB
			model.TextStyle = styles.NoStyleFB

			if focusIndex == kElement {
				model.Focus()
				model.PromptStyle = styles.FocusedStyleFB
				model.TextStyle = styles.FocusedStyleFB
			}
		case *textarea.Model:
			model := vElement.(*textarea.Model)
			model.Blur()
			if focusIndex == kElement {
				model.Focus()
			}
		case *table.Model:
			model := vElement.(*table.Model)
			model.Blur()
			if focusIndex == kElement {
				model.Focus()
			}
		case *BoolModel:
			model := vElement.(*BoolModel)
			model.Active = false
			if focusIndex == kElement {
				model.Active = true
			}
		}
	}
}
