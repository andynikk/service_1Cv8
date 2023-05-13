package cli

import (
	"Service_1Cv8/internal/cli/styles"
	"Service_1Cv8/internal/constants"
	"Service_1Cv8/internal/repository"
	"Service_1Cv8/internal/winsys"
	"errors"
	"fmt"
	"github.com/charmbracelet/bubbles/table"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	itemsInputsFormService   = 0
	itemsAreasFormService    = 0
	itemsButtonFormService   = 2
	itemsTableFormService    = 1
	itemsCheckBoxFormService = 4
)

type checkBoxFormService struct {
	stopService       bool
	stopServiceSelect bool

	delServerCache       bool
	delServerCacheSelect bool

	startService       bool
	startServiceSelect bool

	rebootServer       bool
	rebootServerSelect bool
}

type buttonsFormService struct {
	base1C  bool
	baseSQL bool
}

type FormServer struct {
	focusIndex int
	spinner    spinner.Model
	spinnering bool

	rows  []table.Row
	table table.Model

	checkBoxes checkBoxFormService
	buttons    buttonsFormService

	message string

	service repository.Services
}

func (m *FormServer) SetParameters(args []interface{}) {
	m.focusIndex = 0

	m.checkBoxes.stopService = false
	m.checkBoxes.stopServiceSelect = true

	m.checkBoxes.delServerCache = false
	m.checkBoxes.delServerCacheSelect = false

	m.checkBoxes.startService = false
	m.checkBoxes.startServiceSelect = false

	m.checkBoxes.rebootServer = false
	m.checkBoxes.rebootServerSelect = false

	m.buttons.baseSQL = false
	m.buttons.base1C = false

	for _, v := range args {
		switch v.(type) {
		case repository.Services:
			m.service = v.(repository.Services)
			m.fillRows()
		}
	}

	m.Init()
}

func NewFormServer() *FormServer {

	newSpinner := spinner.New()
	newSpinner.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))
	newSpinner.Spinner = spinner.Pulse

	m := &FormServer{
		focusIndex: 0,

		checkBoxes: checkBoxFormService{},
		buttons:    buttonsFormService{},

		spinnering: false,
		spinner:    newSpinner,
	}

	m.createTable()

	return m
}

func (m *FormServer) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m *FormServer) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "f10":
			return m, tea.Quit
		case "esc", "ctrl+q":
			return client.Models[constants.FormMain].(*FormMain), nil
		case "tab", "shift+tab", "up", "down":

			if m.focusIndex == 4 &&
				(msg.String() == "up" || msg.String() == "down") {

				table, cmd := m.table.Update(msg)

				m.table = table
				return m, cmd
			}

			s := msg.String()

			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > m.lenForm()-1 {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = m.lenForm() - 1
			}

			switch m.focusIndex {
			case 0: //stopService
				m.buttons.baseSQL = false

				//m.checkBoxes.stopService = true
				m.checkBoxes.stopServiceSelect = true

				//m.checkBoxes.delServerCache = false
				m.checkBoxes.delServerCacheSelect = false

			case 1: //delServerCache
				//m.checkBoxes.stopService = false
				m.checkBoxes.stopServiceSelect = false

				//m.checkBoxes.delServerCache = true
				m.checkBoxes.delServerCacheSelect = true

				//m.checkBoxes.startService = false
				m.checkBoxes.startServiceSelect = false
			case 2: //startService
				//m.checkBoxes.delServerCache = false
				m.checkBoxes.delServerCacheSelect = false

				//m.checkBoxes.startService = true
				m.checkBoxes.startServiceSelect = true

				//m.checkBoxes.rebootServer = false
				m.checkBoxes.rebootServerSelect = false
			case 3: //rebootServer
				//m.checkBoxes.startService = false
				m.checkBoxes.startServiceSelect = false

				//m.checkBoxes.rebootServer = true
				m.checkBoxes.rebootServerSelect = true

				m.table.Blur()
			case 4: //table
				//m.checkBoxes.rebootServer = false
				m.checkBoxes.rebootServerSelect = false

				m.table.Focus()

				m.buttons.base1C = false
			case 5: //base1C
				m.table.Blur()

				m.buttons.base1C = true

				m.buttons.baseSQL = false
			case 6: //baseSQL
				m.buttons.base1C = false

				m.buttons.baseSQL = true

				//m.checkBoxes.stopService = false
				m.checkBoxes.rebootServerSelect = false
			}
			return m, tea.Batch(cmds...)
		case "enter", " ":
			if m.checkBoxes.startServiceSelect ||
				m.checkBoxes.delServerCacheSelect ||
				m.checkBoxes.stopServiceSelect ||
				m.checkBoxes.rebootServerSelect {

				m.executeEnterFormServer()

				return m, tea.Batch(cmds...)
			}
			if m.buttons.base1C {
				model, cmd := m.executeFormServerDB()
				return model, cmd
			}
			if m.buttons.baseSQL {
				model, cmd := m.executeFormServerDBSQL()
				return model, cmd
			}
		}
	default:
		var cmd tea.Cmd

		if msg == nil {
			msg = spinner.TickMsg{
				ID:   0,
				Time: time.Now(),
			}
		}

		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, tea.Batch(cmds...)
}

func (m *FormServer) View() string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("Server %s. Srvice %s\n\n", m.service.NameServer, m.service.NameService))

	cursor := " "
	checked := " "
	if m.checkBoxes.stopServiceSelect {
		cursor = styles.FocusedStyleFB.Render(">")
	}
	if m.checkBoxes.stopService {
		checked = "x"
	}
	b.WriteString(fmt.Sprintf("%s [%s] %s\n", cursor, checked, "Stop service"))

	cursor = " "
	checked = " "
	if m.checkBoxes.delServerCacheSelect {
		cursor = styles.FocusedStyleFB.Render(">")
	}
	if m.checkBoxes.delServerCache {
		checked = "x"
	}
	b.WriteString(fmt.Sprintf("%s [%s] %s\n", cursor, checked, "Delete server cache"))

	cursor = " "
	checked = " "
	if m.checkBoxes.startServiceSelect {
		cursor = styles.FocusedStyleFB.Render(">")
	}
	if m.checkBoxes.startService {
		checked = "x"
	}
	b.WriteString(fmt.Sprintf("%s [%s] %s\n", cursor, checked, "Start service"))

	cursor = " "
	checked = " "
	if m.checkBoxes.rebootServerSelect {
		cursor = styles.FocusedStyleFB.Render(">")
	}
	if m.checkBoxes.rebootServer {
		checked = "x"
	}
	b.WriteString(fmt.Sprintf("%s [%s] %s\n", cursor, checked, "Reboot server"))

	b.WriteRune('\n')
	m.table.SetRows(m.rows)
	m.table.SetHeight(styles.Min(20, len(m.rows)))
	b.WriteString(styles.BaseStyle.Render(m.table.View()))

	b.WriteString("\n\n\n")

	b.WriteString("\n\n")

	if m.buttons.base1C {
		b.WriteString(fmt.Sprintf("[ %s ]", styles.FocusedStyleFB.Render("1C server databases")))
	} else {
		b.WriteString(fmt.Sprintf("[ %s ]", styles.BlurredStyleFB.Render("1C server databases")))
	}
	b.WriteString(" ")

	if m.buttons.baseSQL {
		b.WriteString(fmt.Sprintf("[ %s ]", styles.FocusedStyleFB.Render("SQL server databases")))
	} else {
		b.WriteString(fmt.Sprintf("[ %s ]", styles.BlurredStyleFB.Render("SQL server databases")))
	}
	b.WriteString(" ")

	s := " "
	if m.spinnering {
		s = m.spinner.View()
	}
	b.WriteString(fmt.Sprintf("\n\n %s ", s))

	statusVal := styles.StatusText.Copy().Width(styles.Width).Render(m.message)
	b.WriteString(statusVal + "\n")

	podval := fmt.Sprintf("\n\n Press %s to exit main menu | Press %s to quit\n",
		styles.GreenFg("[ ESC ]"), styles.GreenFg("[ F10 ]"))

	b.WriteString(podval)

	return b.String()
}

func (m *FormServer) goStopService() {

	m.message = ""
	m.spinnering = true

	//now := time.Now()
	//now15 := now.Add(time.Second * 15)
	//
	//m.spinner.Spinner = spinner.Pulse
	//m.spinner.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	//
	//err := errors.New("---")
	//for time.Now().Before(now15) {
	//	//m.spinner.Update(smsg)
	//}

	err := winsys.StopService(m.service.NameServer, m.service.NameService)

	m.message = "OK stop service"
	if err != nil {
		m.message = err.Error()
	}

	m.spinnering = false
	m.checkBoxes.stopService = err == nil
}

func (m *FormServer) lenForm() int {
	return itemsInputsFormService + itemsAreasFormService + itemsButtonFormService + itemsTableFormService +
		itemsCheckBoxFormService
}

func (m *FormServer) createTable() {
	columns := []table.Column{
		{Title: "Desc", Width: 15},
		{Title: "Name", Width: 45},
		{Title: "User", Width: 10},
		{Title: "Password", Width: 8},
	}

	rows := []table.Row{}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)

	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)

	t.SetStyles(s)

	m.table = t
}

func (m *FormServer) goClearServerCache() {

	m.message = ""
	m.spinnering = true

	err := errors.New("")
	now := time.Now()
	now10 := now.Add(time.Second * 10)
	for time.Now().Before(now10) {
		err = nil
		err = winsys.ClearServerCache(m.service.NameServer)

		if err == nil {
			break
		}
	}

	m.message = "OK clear server cache"
	if err != nil {
		m.message = err.Error()
	}

	m.spinnering = false
	m.checkBoxes.delServerCache = err == nil

}

func (m *FormServer) goStartService() {

	m.message = ""
	m.spinnering = true

	err := winsys.StartService(m.service.IP, m.service.NameServer, m.service.NameService)

	m.message = "OK start service"
	if err != nil {
		m.message = err.Error()
	}

	m.spinnering = false
	m.checkBoxes.startService = err == nil
}

func (m *FormServer) goRebootWindows() {
	m.message = ""
	m.spinnering = true

	err := winsys.RebootWindows()
	pref := "reboot"

	m.message = fmt.Sprintf("OK %s", pref)
	if err != nil {
		m.message = err.Error()
	}

	m.spinnering = false
	m.checkBoxes.rebootServer = err == nil
}

func (m *FormServer) executeEnterFormServer() {

	switch m.focusIndex {
	case 0:
		go m.goStopService()
	case 1:
		go m.goClearServerCache()
	case 2:
		go m.goStartService()
	case 3:
		go m.goRebootWindows()
	}

	if m.buttons.base1C {
		args := []interface{}{nil}
		model := client.Models[constants.FormServerDB].(*FormServerDB)
		model.SetParameters(args)

		model.Update(nil)
	}

	if m.buttons.baseSQL {

	}
}

func (m *FormServer) executeFormServerDB() (tea.Model, tea.Cmd) {

	args := []interface{}{m.service}
	model := client.Models[constants.FormServerDB].(*FormServerDB)
	model.SetParameters(args)

	return model.Update(nil)
}

func (m *FormServer) executeFormServerDBSQL() (tea.Model, tea.Cmd) {

	args := []interface{}{}
	model := client.Models[constants.FormSQLDB].(*FormSQLDB)

	selectedRow := m.table.Cursor()
	name := m.rows[selectedRow][1]
	for _, v := range m.service.SQLServers {
		if v.Name == name {
			args = []interface{}{v}
		}
	}

	model.SetParameters(args)
	return model.Update(nil)
}

func (m *FormServer) fillRows() {
	m.rows = []table.Row{}

	for _, sqlsrv := range m.service.SQLServers {
		rowT := table.Row{
			sqlsrv.Description,
			sqlsrv.Name,
			sqlsrv.User,
			"*",
		}
		//m.message = sqlsrv.Password
		m.rows = append(m.rows, rowT)
	}
}
