package forms_exchange

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"Service_1Cv8/internal/constants"
)

var exchanger = Exchanger{}

type KeyContext string

type CustomModel interface {
	SetParameters([]interface{})
}

type CM CustomModel

type Exchanger struct {
	Models []CM
	Tea    *tea.Program
}

type Task struct {
	status      constants.Status
	title       string
	description string
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

func NewExchanger() Exchanger {

	fmt.Print("loading...")

	exchanger.Models = []CM{
		&FormBasic{},
		NewFormListToken(),
		NewFormKeyToken(),
		NewFormListQueues(),
		NewFormQueue(),
	}

	frm := exchanger.Models[constants.FormExchangeBasic].(*FormBasic)
	frm.SetParameters(nil)

	exchanger.Tea = tea.NewProgram(frm)
	return exchanger
}

func (c *Exchanger) Run() error {
	_, err := c.Tea.Run()
	return err
}
