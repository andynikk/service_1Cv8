package cli

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"Service_1Cv8/internal/cli/charm"
	"Service_1Cv8/internal/cli/styles"
	"Service_1Cv8/internal/constants"
	"Service_1Cv8/internal/repository"
	"Service_1Cv8/internal/winsys"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	itemsInputsFormService   = 0
	itemsAreasFormService    = 0
	itemsButtonFormService   = 4
	itemsTableFormService    = 1
	itemsCheckBoxFormService = 4
)

type checkBoxFormService struct {
	stop   charm.BoolModel
	delete charm.BoolModel
	start  charm.BoolModel
	reboot charm.BoolModel

	//stopService       bool
	//stopServiceSelect bool
	//
	//delServerCache       bool
	//delServerCacheSelect bool
	//
	//startService       bool
	//startServiceSelect bool
	//
	//rebootServer       bool
	//rebootServerSelect bool
}

type buttonsFormService struct {
	refresh charm.BoolModel
	base1C  charm.BoolModel
	baseSQL charm.BoolModel
	token1C charm.BoolModel
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

	charm.CrawlElements
}

func (m *FormServer) SetParameters(args []interface{}) {
	m.focusIndex = 0

	m.buttons.baseSQL.Active = false
	m.buttons.baseSQL.Active = false
	m.buttons.token1C.Active = false

	for _, v := range args {
		switch v.(type) {
		case repository.Services:
			m.service = v.(repository.Services)
			m.fillRows()
		}
	}

	thisES := client.EventsService[m.service.NameServer]
	m.checkBoxes.stop.Active = false
	m.checkBoxes.stop.Value = thisES.Stop

	m.checkBoxes.delete.Active = false
	m.checkBoxes.delete.Value = thisES.Del

	m.checkBoxes.start.Active = false
	m.checkBoxes.start.Value = thisES.Start

	m.checkBoxes.reboot.Active = false
	m.checkBoxes.reboot.Value = thisES.Reboot

	m.buttons.refresh.Active = true

	m.FillFormServerElements()
	m.Init()
}

func (m *FormServer) FillFormServerElements() {
	ce := make(charm.CrawlElements)

	i := -1
	increaseI := func() int { i++; return i }

	ce[increaseI()] = &m.buttons.refresh
	ce[increaseI()] = &m.checkBoxes.stop
	ce[increaseI()] = &m.checkBoxes.delete
	ce[increaseI()] = &m.checkBoxes.start
	ce[increaseI()] = &m.checkBoxes.reboot
	ce[increaseI()] = &m.table
	ce[increaseI()] = &m.buttons.base1C
	ce[increaseI()] = &m.buttons.baseSQL
	ce[increaseI()] = &m.buttons.token1C

	m.CrawlElements = ce
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

			charm.TеraversingFormElements(m.CrawlElements, m.focusIndex)

			return m, tea.Batch(cmds...)
		case "enter", " ":
			if m.checkBoxes.start.Active {
				m.executeEnterFormServer(constants.EvStart)
				return m, tea.Batch(cmds...)
			}
			if m.checkBoxes.delete.Active {
				m.executeEnterFormServer(constants.EvDelete)
				return m, tea.Batch(cmds...)
			}
			if m.checkBoxes.stop.Active {
				m.executeEnterFormServer(constants.EvStop)
				return m, tea.Batch(cmds...)
			}
			if m.checkBoxes.reboot.Active {
				m.executeEnterFormServer(constants.EvReboot)
				return m, tea.Batch(cmds...)
			}
			if m.buttons.base1C.Active {
				model, cmd := m.executeFormServerDB()
				return model, cmd
			}
			if m.buttons.baseSQL.Active {
				model, cmd := m.executeFormServerDBSQL()
				return model, cmd
			}
			if m.buttons.token1C.Active {
				model, cmd := m.executeFormListToken()
				return model, cmd
			}
			if m.buttons.refresh.Active {
				m.executeRefresh()
				return m, tea.Batch(cmds...)
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

	if m.buttons.refresh.Active {
		b.WriteString(fmt.Sprintf("[ %s ]", styles.FocusedStyleFB.Render("Refresh")))
	} else {
		b.WriteString(fmt.Sprintf("[ %s ]", styles.BlurredStyleFB.Render("Refresh")))
	}
	b.WriteString("\n\n")

	cursor := " "
	checked := " "
	if m.checkBoxes.stop.Active {
		cursor = styles.FocusedStyleFB.Render(">")
	}
	if m.checkBoxes.stop.Value {
		checked = "x"
	}
	b.WriteString(fmt.Sprintf("%s [%s] %s\n", cursor, checked, "Stop service"))

	cursor = " "
	checked = " "
	if m.checkBoxes.delete.Active {
		cursor = styles.FocusedStyleFB.Render(">")
	}
	if m.checkBoxes.delete.Value {
		checked = "x"
	}
	b.WriteString(fmt.Sprintf("%s [%s] %s\n", cursor, checked, "Delete server cache"))

	cursor = " "
	checked = " "
	if m.checkBoxes.start.Active {
		cursor = styles.FocusedStyleFB.Render(">")
	}
	if m.checkBoxes.start.Value {
		checked = "x"
	}
	b.WriteString(fmt.Sprintf("%s [%s] %s\n", cursor, checked, "Start service"))

	cursor = " "
	checked = " "
	if m.checkBoxes.reboot.Active {
		cursor = styles.FocusedStyleFB.Render(">")
	}
	if m.checkBoxes.reboot.Value {
		checked = "x"
	}
	b.WriteString(fmt.Sprintf("%s [%s] %s\n", cursor, checked, "Reboot server"))

	b.WriteRune('\n')
	m.table.SetRows(m.rows)
	m.table.SetHeight(styles.Min(20, len(m.rows)))
	b.WriteString(styles.BaseStyle.Render(m.table.View()))

	b.WriteString("\n\n\n")

	b.WriteString("\n\n")

	if m.buttons.base1C.Active {
		b.WriteString(fmt.Sprintf("[ %s ]", styles.FocusedStyleFB.Render("1C server databases")))
	} else {
		b.WriteString(fmt.Sprintf("[ %s ]", styles.BlurredStyleFB.Render("1C server databases")))
	}
	b.WriteString(" ")

	if m.buttons.baseSQL.Active {
		b.WriteString(fmt.Sprintf("[ %s ]", styles.FocusedStyleFB.Render("SQL server databases")))
	} else {
		b.WriteString(fmt.Sprintf("[ %s ]", styles.BlurredStyleFB.Render("SQL server databases")))
	}
	b.WriteString(" ")

	if m.buttons.token1C.Active {
		b.WriteString(fmt.Sprintf("[ %s ]", styles.FocusedStyleFB.Render("Token's")))
	} else {
		b.WriteString(fmt.Sprintf("[ %s ]", styles.BlurredStyleFB.Render("Token's")))
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

	err := winsys.StopService(m.service.NameServer, m.service.NameService)

	m.message = "OK stop service"
	if err != nil {
		m.message = err.Error()
	}

	m.spinnering = false
	m.checkBoxes.stop.Value = err == nil
	if m.checkBoxes.stop.Value &&
		client.PerformedActions.RestartService.Find(m.service.NameServer) == -1 {

		client.PerformedActions.RestartService = append(client.PerformedActions.RestartService, m.service.NameServer)
	}

	m.setEventValue(constants.EvStop)
}

func (m *FormServer) setEventValue(event constants.EventService) {
	es := client.EventsService
	if es == nil {
		es = make(map[string]EventsService)
	}

	thisES := es[m.service.NameServer]

	switch event {
	case constants.EvStop:
		thisES.Stop = m.checkBoxes.stop.Value
	case constants.EvDelete:
		thisES.Del = m.checkBoxes.delete.Value
	case constants.EvStart:
		thisES.Start = m.checkBoxes.start.Value
	case constants.EvReboot:
		thisES.Reboot = m.checkBoxes.reboot.Value
	}

	es[m.service.NameServer] = thisES
	client.EventsService = es
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
	m.checkBoxes.delete.Value = err == nil

	if m.checkBoxes.delete.Value &&
		client.PerformedActions.ClearCash.Find(m.service.NameServer) == -1 {

		client.PerformedActions.ClearCash = append(client.PerformedActions.ClearCash, m.service.NameServer)
	}

	m.setEventValue(constants.EvDelete)
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
	m.checkBoxes.start.Value = err == nil

	m.setEventValue(constants.EvStart)
}

func (m *FormServer) goRebootWindows() {
	m.message = ""
	m.spinnering = true

	err := winsys.RebootRemoteWindows(m.service.NameServer)
	pref := "reboot"

	m.message = fmt.Sprintf("OK %s", pref)
	if err != nil {
		m.message = err.Error()
	}

	m.spinnering = false
	m.checkBoxes.reboot.Value = err == nil

	if m.checkBoxes.reboot.Value &&
		client.PerformedActions.RebutServer.Find(m.service.NameServer) == -1 {

		client.PerformedActions.RebutServer = append(client.PerformedActions.RebutServer, m.service.NameServer)
	}

	m.setEventValue(constants.EvReboot)
}

func (m *FormServer) executeEnterFormServer(event constants.EventService) {

	switch event {
	case constants.EvStop:
		go m.goStopService()
	case constants.EvDelete:
		go m.goClearServerCache()
	case constants.EvStart:
		go m.goStartService()
	case constants.EvReboot:
		go m.goRebootWindows()
	}

	if m.buttons.base1C.Active {
		args := []interface{}{nil}
		model := client.Models[constants.FormServerDB].(*FormServerDB)
		model.SetParameters(args)

		model.Update(nil)
	}

	if m.buttons.baseSQL.Active {

	}
}

func (m *FormServer) executeFormServerDB() (tea.Model, tea.Cmd) {

	args := []interface{}{m.service}
	model := client.Models[constants.FormServerDB].(*FormServerDB)
	model.SetParameters(args)

	return model.Update(nil)
}

func (m *FormServer) executeFormListToken() (tea.Model, tea.Cmd) {

	//args := []interface{}{m.service}
	//model := client.Models[constants.FormListToken].(*FormListToken)
	//model.SetParameters(args)
	//
	//return model.Update(nil)

	return nil, nil
}

func (m *FormServer) executeRefresh() {
	es := client.EventsService
	if es == nil {
		es = make(map[string]EventsService)
	}

	thisES := es[m.service.NameServer]

	thisES.Stop = false
	thisES.Del = false
	thisES.Start = false
	thisES.Reboot = false

	es[m.service.NameServer] = thisES
	client.EventsService = es

	m.checkBoxes.stop.Value = false
	m.checkBoxes.delete.Value = false
	m.checkBoxes.start.Value = false
	m.checkBoxes.reboot.Value = false

	m.buttons.refresh.Active = false

	m.checkBoxes.stop.Active = true
	m.focusIndex = m.focusIndex + 1
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
