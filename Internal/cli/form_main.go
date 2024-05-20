package cli

import (
	"Service_1Cv8/internal/cli/charm"
	"Service_1Cv8/internal/repository"
	"fmt"
	"log"
	"sort"
	"strings"

	"Service_1Cv8/internal/cli/styles"
	"Service_1Cv8/internal/constants"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type FormMain struct {
	loaded   bool
	focused  constants.Status
	lists    []list.Model
	err      error
	quitting bool
}

type listKeyMap struct {
	toggleHelpMenu key.Binding
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

func NewFormMain() *FormMain {
	return &FormMain{}
}

func (m *FormMain) SetParameters(args []interface{}) {
	//m.initLists(0, 0)
}

func (m *FormMain) MoveToNext() tea.Msg {
	selectedItem := m.lists[m.focused].SelectedItem()
	selectedTask := selectedItem.(charm.Task)
	m.lists[selectedTask.Status].RemoveItem(m.lists[m.focused].Index())
	selectedTask.Next()
	m.lists[selectedTask.Status].InsertItem(len(m.lists[selectedTask.Status].Items())-1, list.Item(selectedTask))
	return nil
}

func (m *FormMain) DeleteCurrent() tea.Msg {
	if len(m.lists[m.focused].VisibleItems()) > 0 {
		selectedTask := m.lists[m.focused].SelectedItem().(charm.Task)
		m.lists[selectedTask.Status].RemoveItem(m.lists[m.focused].Index())
	}
	return nil
}

func (m *FormMain) Next() {
	if m.focused == constants.Services {
		m.focused = constants.Blocking
	} else {
		m.focused++
	}
}

func (m *FormMain) Prev() {
	if m.focused == constants.Blocking {
		m.focused = constants.Services
	} else {
		m.focused--
	}
}

func (m *FormMain) initLists(width, height int) {

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
	m.lists = []list.Model{defaultList, defaultList}
	m.lists[constants.Blocking].Title = "Data bases"

	listItem := []list.Item{}
	arrDB := client.Storage.DB
	sort.Slice(arrDB, func(i, j int) bool {
		return arrDB[i].Name < arrDB[j].Name
	})

	for _, db := range arrDB {
		propertyDB, _ := repository.GetPropertiesDB(client.Storage.PropertyDB, db.UID)

		title := fmt.Sprintf("• %s", db.Name)
		if !propertyDB.Block {
			title = styles.GreenFg(title)
		} else {
			title = styles.RedFg(title)
		}

		description := fmt.Sprintf("Srvr=\"%s\";Ref=\"%s\";", db.Server, db.NameOnServer)
		listItem = append(listItem, charm.Task{Status: constants.Blocking, DBName: db.Name,
			TaskTitle: fmt.Sprintf("%s", title), TaskDescription: description, Uid: db.UID})
	}
	m.lists[constants.Blocking].SetItems(listItem)

	// Init services
	m.lists[constants.Services].Title = "Servers"

	listServersItem := []list.Item{}
	storageService := client.Storage.Services
	sort.Slice(storageService, func(i, j int) bool {
		return storageService[i].NameServer < storageService[j].NameServer
	})

	for _, service := range storageService {

		listServersItem = append(listServersItem, charm.Task{Status: constants.Services, DBName: "",
			TaskTitle: service.NameServer, TaskDescription: service.NameService, Uid: service.UID})
	}
	m.lists[constants.Services].SetItems(listServersItem)
}

func (m *FormMain) Init() tea.Cmd {
	return nil
}

func (m *FormMain) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !m.loaded {
			styles.ColumnStyle.Width(msg.Width / constants.Divisor)
			styles.FocusedStyle.Width(msg.Width / constants.Divisor)
			styles.ColumnStyle.Height(msg.Height - constants.Divisor)
			styles.FocusedStyle.Height(msg.Height - constants.Divisor)
			//m.initLists(msg.Width, msg.Height)
			m.initLists(210, 35)
			m.loaded = true
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "f10":

			m.quitting = true
			return m, tea.Quit

		case "left", "h":

			m.Prev()

		case "right", "l":

			m.Next()

		case "enter":

			model, cmd := m.executeEnterFormMain()
			return model, cmd

		case "ctrl+n", "insert":

			if m.focused == 0 {
				frm := client.Models[constants.FormAddDB].(*FormAddDB)
				frm.new = true
				frm.SetParameters(nil)
				return frm.Update(nil)
			}
			if m.focused == 1 {
				frm := client.Models[constants.FormAddService].(*FormAddService)
				frm.new = true
				frm.SetParameters(nil)
				return frm.Update(nil)
			}

		case "ctrl+e", "ctrl+o":

			model, cmd := m.executeEditDB(false)
			return model, cmd

		case "ctrl+c":

			model, cmd := m.executeEditDB(true)
			return model, cmd

		case "ctrl+d", "delete":

			model, _ := m.executeDelDB()
			return model, m.DeleteCurrent

		case "ctrl+s", "ctrl+w":

			frm := client.Models[constants.FormSetting].(*FormSetting)
			frm.SetParameters(nil)

			return frm.Update(nil)

		case "f5":
			frm := client.Models[constants.FormCopyFile].(*FormCopyFile)
			frm.SetParameters(nil)

			return frm.Update(nil)
		case "esc":

			m.quitting = false
			return m, nil

		case "ctrl+t", "ctrl+g":

			frm := client.Models[constants.FormTgMsg].(*FormTgMsg)
			frm.SetParameters(nil)

			return frm.Update(nil)
		}
	case charm.Task:
		task := msg
		return m, m.lists[task.Status].InsertItem(len(m.lists[task.Status].Items()), task)
	}
	var cmd tea.Cmd
	m.lists[m.focused], cmd = m.lists[m.focused].Update(msg)
	return m, cmd
}

func (m *FormMain) View() string {
	var b strings.Builder

	//bb := styles.CursorModeHelpStyle.Render(" - ") +
	//	styles.CursorModeHelpStyleWhite.Render("[del/ctrl+d delete]") +
	//	styles.CursorModeHelpStyle.Render(" - ") +
	//	styles.CursorModeHelpStyleWhite.Render("[ctrl+e  edit current]") +
	//	styles.CursorModeHelpStyle.Render(" - ") +
	//	styles.CursorModeHelpStyleWhite.Render("[ctrl+c add a new one based current]") +
	//	styles.CursorModeHelpStyle.Render(" - ") +
	//	styles.CursorModeHelpStyleWhite.Render("[insert/ctrl+n add a new]") +
	//	styles.CursorModeHelpStyle.Render(" - ") +
	//	styles.CursorModeHelpStyleWhite.Render("[ctrl+s settings]")

	if m.quitting {
		return ""
	}
	if m.loaded {
		todoView := m.lists[constants.Blocking].View()
		inProgView := m.lists[constants.Services].View()
		switch m.focused {
		case constants.Services:

			return lipgloss.JoinHorizontal(
				lipgloss.Left,
				styles.ColumnStyle.Render(todoView),
				styles.FocusedStyle.Render(inProgView))

		default:

			items := m.lists[constants.Blocking].Items()
			storageDB := client.Storage.DB
			for kk, v := range items {
				for _, val := range storageDB {
					if v.(charm.Task).DBName != val.Name {
						continue
					}

					propertyDB, _ := repository.GetPropertiesDB(client.Storage.PropertyDB, val.UID)

					title := fmt.Sprintf("• %s", val.Name)
					if !propertyDB.Block {
						title = styles.GreenFg(title)
					} else {
						title = styles.RedFg(title)
					}
					items[kk] = charm.Task{Status: constants.Blocking, DBName: val.Name, TaskTitle: fmt.Sprintf("%s", title),
						TaskDescription: fmt.Sprintf("Srvr=\"%s\";Ref=\"%s\";", val.Server, val.NameOnServer),
						Uid:             val.UID}
				}
			}

			m.lists[constants.Blocking].SetItems(items)

			panels := lipgloss.JoinHorizontal(
				lipgloss.Left,
				styles.FocusedStyle.Render(todoView),
				styles.ColumnStyle.Render(inProgView),
			)
			b.WriteString(panels)
		}
	} else {
		return "loading..."
	}

	//b.WriteString("\n" + bb)
	return b.String()
}

func (m *FormMain) executeEnterFormMain() (tea.Model, tea.Cmd) {
	if m.focused == 0 {
		uid := m.lists[m.focused].SelectedItem().(charm.Task).Uid
		db, _ := repository.GetDB(client.Storage.DB, uid)
		propertyDB, _ := repository.GetPropertiesDB(client.Storage.PropertyDB, uid)

		args := []interface{}{db, propertyDB}
		model := client.Models[constants.FormBlock].(*FormBlock)
		model.SetParameters(args)

		return model.Update(nil)
	}

	if m.focused == 1 {
		uid := m.lists[m.focused].SelectedItem().(charm.Task).Uid

		srv, _ := repository.GetService(client.Storage.Services, uid)

		args := []interface{}{srv}
		model := client.Models[constants.FormServer].(*FormServer)
		model.SetParameters(args)

		return model.Update(nil)
	}

	return client.Models[constants.FormMain].(*FormMain), nil
}

func (m *FormMain) executeEditDB(new bool) (tea.Model, tea.Cmd) {

	if m.focused == 0 {

		uid := m.lists[m.focused].SelectedItem().(charm.Task).Uid
		db, _ := repository.GetDB(client.Storage.DB, uid)
		propertyDB, _ := repository.GetPropertiesDB(client.Storage.PropertyDB, uid)

		args := []interface{}{db, propertyDB}
		model := client.Models[constants.FormAddDB].(*FormAddDB)
		model.new = new
		model.SetParameters(args)

		return model.Update(nil)
	}

	if m.focused == 1 {

		uid := m.lists[m.focused].SelectedItem().(charm.Task).Uid

		srv, _ := repository.GetService(client.Storage.Services, uid)

		args := []interface{}{srv}
		model := client.Models[constants.FormAddService].(*FormAddService)
		model.new = new
		model.SetParameters(args)

		return model.Update(nil)

	}

	return client.Models[constants.FormMain].(*FormMain), nil
}

func (m *FormMain) executeDelDB() (tea.Model, tea.Cmd) {

	if m.focused == 0 {

		uid := m.lists[m.focused].SelectedItem().(charm.Task).Uid
		_, keyDB := repository.GetDB(client.Storage.DB, uid)
		_, keyPDB := repository.GetPropertiesDB(client.Storage.PropertyDB, uid)

		dataDBJSON := &client.Storage.DataDBJSON

		if keyDB != -1 {
			dataDBJSON.DB[keyDB] = dataDBJSON.DB[len(dataDBJSON.DB)-1]
			dataDBJSON.DB[len(dataDBJSON.DB)-1] = repository.DataBases{}
			dataDBJSON.DB = dataDBJSON.DB[:len(dataDBJSON.DB)-1]
		}

		if keyPDB != -1 {
			dataDBJSON.PropertyDB[keyPDB] = dataDBJSON.PropertyDB[len(dataDBJSON.PropertyDB)-1]
			dataDBJSON.PropertyDB[len(dataDBJSON.PropertyDB)-1] = repository.PropertyDB{}
			dataDBJSON.PropertyDB = dataDBJSON.PropertyDB[:len(dataDBJSON.PropertyDB)-1]
		}

		model := client.Models[constants.FormMain].(*FormMain)
		model.SetParameters(nil)

		err := client.Storage.SetPudgelData()
		if err != nil {
			log.Println(err.Error())
		}

		return model.Update(nil)
	}

	if m.focused == 1 {

		uid := m.lists[m.focused].SelectedItem().(charm.Task).Uid
		_, keyS := repository.GetService(client.Storage.Services, uid)
		storage := client.Storage

		if keyS != -1 {
			storage.Services[keyS] = storage.Services[len(storage.Services)-1]
			storage.Services[len(storage.Services)-1] = repository.Services{}
			storage.Services = storage.Services[:len(storage.Services)-1]
		}

		model := client.Models[constants.FormMain].(*FormMain)
		model.SetParameters(nil)

		err := client.Storage.SetPudgelData()
		if err != nil {
			log.Println(err.Error())
		}

		return model.Update(nil)
	}

	return client.Models[constants.FormMain].(*FormMain), nil
}
