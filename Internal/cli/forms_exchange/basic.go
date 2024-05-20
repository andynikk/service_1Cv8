package forms_exchange

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"Service_1Cv8/internal/cli/charm"
	"Service_1Cv8/internal/cli/styles"
	"Service_1Cv8/internal/constants"
)

type FormBasic struct {
	loaded   bool
	focused  constants.Status
	lists    []list.Model
	err      error
	quitting bool
}

type listKeyMap struct {
	toggleHelpMenu key.Binding
}

func (f *FormBasic) SetParameters(args []interface{}) {

}

func (f *FormBasic) Init() tea.Cmd {
	return nil
}

func (f *FormBasic) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !f.loaded {
			styles.ColumnStyle.Width(msg.Width / 1)
			styles.FocusedStyle.Width(msg.Width / 1)
			styles.ColumnStyle.Height(msg.Height - 1)
			styles.FocusedStyle.Height(msg.Height - 1)
			//m.initLists(msg.Width, msg.Height)
			f.initLists(210, 35)
			f.loaded = true
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "f10", "ctrl+q":

			f.quitting = true
			return f, tea.Quit

		case "enter":

			model, cmd := f.executeEnterFormMain()
			if model != nil {
				return model, cmd
			}

		case "esc":

			f.quitting = false
			return f, nil

		case "ctrl+n", "insert":

			//if f.focused == 0 {
			//	frm := client.Models[constants.FormAddDB].(*FormAddDB)
			//	frm.new = true
			//	frm.SetParameters(nil)
			//	return frm.Update(nil)
			//}
			//if m.focused == 1 {
			//	frm := client.Models[constants.FormAddService].(*FormAddService)
			//	frm.new = true
			//	frm.SetParameters(nil)
			//	return frm.Update(nil)
			//}

		}
	case charm.Task:
		task := msg
		return f, f.lists[task.Status].InsertItem(len(f.lists[task.Status].Items()), task)
	}
	var cmd tea.Cmd
	f.lists[f.focused], cmd = f.lists[f.focused].Update(msg)
	return f, cmd
}

func (f *FormBasic) View() string {
	var b strings.Builder

	if f.quitting {
		return ""
	}
	if f.loaded {
		todoView := f.lists[0].View()

		items := f.lists[0].Items()

		title := styles.GreenFg(fmt.Sprintf("1. %s", "Keys, tokens"))
		items[0] = charm.Task{
			Status: 0, OrdinalNumber: 1, TaskTitle: title,
			TaskDescription: fmt.Sprintf("Access. Public, private keys. Tokins"),
		}

		title = styles.GreenFg(fmt.Sprintf("2. %s", "Queues"))
		items[1] = charm.Task{
			Status: 0, OrdinalNumber: 2, TaskTitle: title,
			TaskDescription: fmt.Sprintf("Queues. Public, private keys. Tokins"),
		}

		panels := lipgloss.JoinHorizontal(
			lipgloss.Left,
			styles.FocusedStyle.Render(todoView),
		)
		b.WriteString(panels)

	} else {
		return "loading..."
	}

	return b.String()
}

func (f *FormBasic) initLists(width, height int) {

	defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), width/constants.Divisor, height/2-10)
	defaultList.SetHeight(25)
	defaultList.SetShowHelp(true)
	defaultList.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			newKeyMapSetting().toggleHelpMenu,
			newKeyMapAddNew().toggleHelpMenu,
			newKeyMapAddCopying().toggleHelpMenu,
			newKeyMapEdit().toggleHelpMenu,
			newKeyMapDel().toggleHelpMenu,
			newKeyCopyFile().toggleHelpMenu,
		}
	}
	f.lists = []list.Model{defaultList}
	f.lists[0].Title = "Data bases"

	listItem := []list.Item{}

	title := styles.GreenFg(fmt.Sprintf("1. %s", "Keys, tokens"))
	description := fmt.Sprintf("Access. Public, private keys. Tokins")
	listItem = append(listItem, charm.Task{
		Status:          0,
		TaskTitle:       fmt.Sprintf("%s", title),
		TaskDescription: description,
	})

	title = styles.GreenFg(fmt.Sprintf("2. %s", "Queues"))
	description = fmt.Sprintf("Queues. Public, private keys. Tokins")
	listItem = append(listItem, charm.Task{
		Status:          0,
		TaskTitle:       fmt.Sprintf("%s", title),
		TaskDescription: description,
	})

	f.lists[0].SetItems(listItem)
	f.lists[0].Title = "Main selection menu"

}

func (f *FormBasic) executeEnterFormMain() (tea.Model, tea.Cmd) {

	i := f.lists[0].SelectedItem().(charm.Task)
	if i.OrdinalNumber == 1 {
		args := []interface{}{}
		model := exchanger.Models[constants.FormExchangeListToken].(*FormListToken)
		model.SetParameters(args)

		return model.Update(nil)
	}

	if i.OrdinalNumber == 2 {
		args := []interface{}{}
		model := exchanger.Models[constants.FormExchangeListQueues].(*FormListQueues)
		model.SetParameters(args)

		return model.Update(nil)
	}

	return nil, nil
}

func newKeyMapAddCopying() *listKeyMap {
	return &listKeyMap{
		toggleHelpMenu: key.NewBinding(
			key.WithKeys(styles.GreenFg("ctrl+c")),
			key.WithHelp(styles.GreenFg("ctrl+c"), styles.GreenFg("Add by copying")),
		),
	}
}

func newKeyMapEdit() *listKeyMap {
	return &listKeyMap{
		toggleHelpMenu: key.NewBinding(
			key.WithKeys(styles.GreenFg("ctrl+o")),
			key.WithHelp(styles.GreenFg("ctrl+o"), styles.GreenFg("Edit current")),
		),
	}
}

func newKeyMapAddNew() *listKeyMap {
	return &listKeyMap{
		toggleHelpMenu: key.NewBinding(
			key.WithKeys(styles.GreenFg("ins, ctrl+n")),
			key.WithHelp(styles.GreenFg("ins, ctrl+n"), styles.GreenFg("Add new")),
		),
	}
}

func newKeyMapDel() *listKeyMap {
	return &listKeyMap{
		toggleHelpMenu: key.NewBinding(
			key.WithKeys(styles.GreenFg("del, ctrl+d")),
			key.WithHelp(styles.GreenFg("del, ctrl+d"), styles.GreenFg("Del current")),
		),
	}
}

func newKeyCopyFile() *listKeyMap {
	return &listKeyMap{
		toggleHelpMenu: key.NewBinding(
			key.WithKeys(styles.GreenFg("F5")),
			key.WithHelp(styles.GreenFg("F5"), styles.GreenFg("Cope backup file")),
		),
	}
}

func newKeyMapSetting() *listKeyMap {
	return &listKeyMap{
		toggleHelpMenu: key.NewBinding(
			key.WithKeys(styles.GreenFg("ctrl+s")),
			key.WithHelp(styles.GreenFg("ctrl+s"), styles.GreenFg("Default settings")),
		),
	}
}
