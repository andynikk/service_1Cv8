package cli

import (
	"fmt"
	"golang.org/x/text/encoding/charmap"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
	//"github.com/jinzhu/copier"

	"Service_1Cv8/internal/cli/styles"
	"Service_1Cv8/internal/constants"
	"Service_1Cv8/internal/repository"
	"Service_1Cv8/internal/winsys"
)

const (
	itemsInputsFormAddService = 5
	itemsAreasFormAddService  = 0
	itemsButtonFormAddService = 1
	itemsTableFormAddService  = 1
)

type inputsFormAddService struct {
	serverName     textinput.Model
	serverIP       textinput.Model
	nameService    textinput.Model
	userServer     textinput.Model
	passwordServer textinput.Model
}

type buttonsFormAddService struct {
	save      bool
	addSQL    bool
	removeSQL bool
}

type FormAddService struct {
	focusIndex int
	inputs     inputsFormAddService
	buttons    buttonsFormAddService

	rows  []table.Row
	table table.Model

	uid string
	new bool

	message string

	spinner    spinner.Model
	spinnering bool

	srv repository.Services
}

func (m *FormAddService) SetParameters(args []interface{}) {
	m.uid = ""
	m.message = ""
	m.srv = repository.Services{}
	m.focusIndex = 0

	m.buttons.save = false
	m.buttons.addSQL = false
	m.buttons.removeSQL = false

	m.rows = []table.Row{}

	m.inputs.serverName.SetValue("")
	m.inputs.serverName.Focus()
	m.inputs.serverName.PromptStyle = styles.FocusedStyleFB

	m.inputs.serverIP.Placeholder = "IP server"
	m.inputs.serverIP.SetValue("")

	m.inputs.nameService.Placeholder = "Name service"
	m.inputs.nameService.SetValue("")

	m.inputs.userServer.Placeholder = "User"
	m.inputs.userServer.SetValue("")

	m.inputs.passwordServer.Placeholder = "Password"
	m.inputs.passwordServer.SetValue("")
	m.inputs.passwordServer.EchoMode = textinput.EchoPassword
	m.inputs.passwordServer.EchoCharacter = '•'

	for _, v := range args {
		switch v.(type) {
		case repository.Services:
			m.srv = v.(repository.Services)

			m.inputs.serverName.SetValue(m.srv.NameServer)
			m.inputs.serverIP.SetValue(m.srv.IP)
			m.inputs.nameService.SetValue(m.srv.NameService)
			m.inputs.userServer.SetValue(m.srv.User)
			m.inputs.passwordServer.SetValue(m.srv.Password)

			m.fillRows()

			if !m.new {
				m.uid = m.srv.UID
			}
		}
	}

	if m.uid == "" {
		m.uid = uuid.New().String()
	}
}

func NewFormAddService() *FormAddService {

	newSpinner := spinner.New()
	newSpinner.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))
	newSpinner.Spinner = spinner.Pulse

	m := FormAddService{
		inputs:     inputsFormAddService{},
		buttons:    buttonsFormAddService{},
		spinnering: false,
		spinner:    newSpinner,

		rows: []table.Row{},
	}

	var t textinput.Model

	t = textinput.New()
	m.inputs.serverName = t

	t = textinput.New()
	m.inputs.serverIP = t

	t = textinput.New()
	m.inputs.nameService = t

	t = textinput.New()
	m.inputs.userServer = t

	t = textinput.New()
	m.inputs.passwordServer = t

	m.createTable()
	return &m

}

func (m *FormAddService) Init() tea.Cmd {
	return nil
}

func (m *FormAddService) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	cmds := make([]tea.Cmd, 0)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		strMsg := msg.String()
		switch strMsg {
		case "esc", "ctrl+q":
			m.uid = ""

			frm := client.Models[constants.FormMain].(*FormMain)
			frm.initLists(204, 31)

			return frm, nil
		case "ctrl+c", "f10":

			return m, tea.Quit

		case "tab", "shift+tab", "up", "down":
			if m.focusIndex == 5 &&
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
			case 0: //serverName
				m.buttons.save = false

				m.inputs.serverName.Focus()
				m.inputs.serverName.PromptStyle = styles.FocusedStyleFB
				m.inputs.serverName.TextStyle = styles.FocusedStyleFB

				m.inputs.serverIP.Blur()
				m.inputs.serverIP.PromptStyle = styles.NoStyleFB
				m.inputs.serverIP.TextStyle = styles.NoStyleFB

			case 1: //serverIP
				m.inputs.serverName.Blur()
				m.inputs.serverName.PromptStyle = styles.NoStyleFB
				m.inputs.serverName.TextStyle = styles.NoStyleFB

				m.inputs.serverIP.Focus()
				m.inputs.serverIP.PromptStyle = styles.FocusedStyleFB
				m.inputs.serverIP.TextStyle = styles.FocusedStyleFB

				m.inputs.nameService.Blur()
				m.inputs.nameService.PromptStyle = styles.NoStyleFB
				m.inputs.nameService.TextStyle = styles.NoStyleFB
			case 2: //nameService
				m.inputs.serverIP.Blur()
				m.inputs.serverIP.PromptStyle = styles.NoStyleFB
				m.inputs.serverIP.TextStyle = styles.NoStyleFB

				m.inputs.nameService.Focus()
				m.inputs.nameService.PromptStyle = styles.FocusedStyleFB
				m.inputs.nameService.TextStyle = styles.FocusedStyleFB

				m.inputs.userServer.Blur()
				m.inputs.userServer.PromptStyle = styles.NoStyleFB
				m.inputs.userServer.TextStyle = styles.NoStyleFB

			case 3: //userServer
				m.inputs.nameService.Blur()
				m.inputs.nameService.PromptStyle = styles.NoStyleFB
				m.inputs.nameService.TextStyle = styles.NoStyleFB

				m.inputs.userServer.Focus()
				m.inputs.userServer.PromptStyle = styles.FocusedStyleFB
				m.inputs.userServer.TextStyle = styles.FocusedStyleFB

				m.inputs.passwordServer.Blur()
				m.inputs.passwordServer.PromptStyle = styles.NoStyleFB
				m.inputs.passwordServer.TextStyle = styles.NoStyleFB

			case 4: //passwordServer
				m.inputs.userServer.Blur()
				m.inputs.userServer.PromptStyle = styles.NoStyleFB
				m.inputs.userServer.TextStyle = styles.NoStyleFB

				m.inputs.passwordServer.Focus()
				m.inputs.passwordServer.PromptStyle = styles.FocusedStyleFB
				m.inputs.passwordServer.TextStyle = styles.FocusedStyleFB

				m.table.Blur()
				m.table.SetHeight(1)
			case 5: //table
				m.inputs.passwordServer.Blur()
				m.inputs.passwordServer.PromptStyle = styles.NoStyleFB
				m.inputs.passwordServer.TextStyle = styles.NoStyleFB

				m.table.Focus()
				if len(m.rows) > 0 {
					m.table.SetHeight(styles.Min(20, len(m.rows)))
				} else {
					m.table.SetHeight(1)
				}

				m.buttons.save = false

			case 6: //save

				m.table.Blur()
				m.table.SetHeight(1)

				m.buttons.save = true

				m.inputs.serverName.Blur()
				m.inputs.serverName.PromptStyle = styles.NoStyleFB
				m.inputs.serverName.TextStyle = styles.NoStyleFB
			}
			return m, tea.Batch(cmds...)
		case "enter":

			m.clickButton()

		case "insert", "ctrl+n":

			model, cmd := m.addSQLServer()
			return model, cmd

		case "delete", "ctrl+d":

			if m.focusIndex == 5 {
				m.delSQLServer()
				return m, nil
			}
		}
	}

	cmd = m.updateInputs(msg)
	return m, cmd
}

func (m *FormAddService) View() string {
	var b strings.Builder

	b.WriteRune('\n')

	seporator := styles.CursorModeHelpStyle.Render(" - ")
	title := fmt.Sprintf("Service: %s %s %s %s %s ", m.srv.NameServer, seporator, m.srv.IP,
		seporator, strings.ToUpper(m.srv.NameService))

	if m.new {
		title = "Service: NEW "
	}

	b.WriteString(fmt.Sprintf(" %s %s\n", title, styles.ShortLine))
	b.WriteRune('\n')

	b.WriteString(fmt.Sprintf("%s %s %s %s\n",
		styles.СaptionStyleFB.Render("Name server:"), m.inputs.serverName.View(),
		styles.СaptionStyleFB.Render("IP server:"), m.inputs.serverIP.View()))

	b.WriteString(fmt.Sprintf("%s %s\n", styles.СaptionStyleFB.Render("Name service:"),
		m.inputs.nameService.View()))

	b.WriteString(fmt.Sprintf("%s %s %s %s\n",
		styles.СaptionStyleFB.Render("User:"), m.inputs.userServer.View(),
		styles.СaptionStyleFB.Render("Password:"), m.inputs.passwordServer.View()))

	b.WriteString(fmt.Sprintf("\n%s SQL servrs: %s", styles.ShortLine, styles.Line))

	b.WriteRune('\n')
	m.table.SetRows(m.rows)
	m.table.SetHeight(styles.Min(20, len(m.rows)))
	b.WriteString(styles.BaseStyle.Render(m.table.View()))

	b.WriteString("\n\n\n")

	if m.buttons.save {
		b.WriteString(fmt.Sprintf("[ %s ]", styles.FocusedStyleFB.Render("save")))
	} else {
		b.WriteString(fmt.Sprintf("[ %s ]", styles.BlurredStyleFB.Render("save")))
	}
	b.WriteString(" ")

	s := ""
	if m.spinnering {
		s = m.spinner.View()
	}
	b.WriteString(fmt.Sprintf("\n\n %s ", s))

	statusVal := styles.StatusText.Copy().Width(styles.Width).Render(m.message)
	b.WriteString(statusVal)

	podval := fmt.Sprintf("\n\n\n %s %s | %s %s\n %s %s | %s %s",
		styles.GreenFg("[ INS, CTRL+N ]"), "add new SQL Server",
		styles.GreenFg("[ DEL, CTRL+D ]"), "delete current SQL Server",
		styles.GreenFg("[ ESC ]"), "exit main menu",
		styles.GreenFg("[ F10 ]"), "exit program")

	b.WriteString(styles.CursorModeHelpStyleFB.Render(podval))

	return b.String()

}

func (m *FormAddService) lenForm() int {
	return itemsInputsFormAddService + itemsAreasFormAddService + itemsButtonFormAddService + itemsTableFormAddService
}

func (m *FormAddService) createTable() {
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

func (m *FormAddService) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, itemsInputsFormAddService+itemsAreasFormAddService)

	switch msg.(type) {
	case tea.KeyMsg:
		_, ok := winsys.SubstitutionRune(msg.(tea.KeyMsg).Runes)
		if !ok {
			for k, v := range msg.(tea.KeyMsg).Runes {

				dec := charmap.Windows1251.DecodeByte(byte(v))
				msg.(tea.KeyMsg).Runes[k] = dec
			}
		}
	}

	m.inputs.serverName, _ = m.inputs.serverName.Update(msg)
	m.inputs.serverIP, _ = m.inputs.serverIP.Update(msg)
	m.inputs.nameService, _ = m.inputs.nameService.Update(msg)
	m.inputs.userServer, _ = m.inputs.userServer.Update(msg)
	m.inputs.passwordServer, _ = m.inputs.passwordServer.Update(msg)

	return tea.Batch(cmds...)
}

func (m *FormAddService) clickButton() {

	if m.buttons.save {
		m.saveService()
	}

}

func (m *FormAddService) saveService() {
	srv, keySrv := repository.GetService(client.Storage.Services, m.uid)
	if keySrv == -1 {
		srv = repository.Services{}
	}
	srv.NameServer = m.inputs.serverName.Value()
	srv.IP = m.inputs.serverIP.Value()
	srv.NameService = m.inputs.nameService.Value()
	srv.User = m.inputs.userServer.Value()
	srv.Password = m.inputs.passwordServer.Value()
	srv.UID = m.uid
	srv.SQLServers = m.srv.SQLServers

	//storage := &client.Storage
	if keySrv != -1 {
		client.Storage.Services[keySrv] = client.Storage.Services[len(client.Storage.Services)-1]
		client.Storage.Services[len(client.Storage.Services)-1] = repository.Services{}
		client.Storage.Services = client.Storage.Services[:len(client.Storage.Services)-1]
	}
	client.Storage.Services = append(client.Storage.Services, srv)

	err := client.Storage.SetPudgelData()
	if err != nil {
		m.message = err.Error()
	}
	m.message = "save service OK"
}

func (m *FormAddService) addSQLServer() (tea.Model, tea.Cmd) {
	model := client.Models[constants.FormAddSQLSrv].(*FormAddSQLSrv)
	model.SetParameters([]interface{}{})

	return model.Update(nil)
}

func (m *FormAddService) delSQLServer() {

	selectedRow := m.table.Cursor()
	name := m.rows[selectedRow][1]

	keySrv := 0
	for k, v := range m.srv.SQLServers {
		if v.Name == name {
			keySrv = k
			break
		}
	}

	m.srv.SQLServers[keySrv] = m.srv.SQLServers[len(m.srv.SQLServers)-1]
	m.srv.SQLServers[len(m.srv.SQLServers)-1] = repository.SQLServer{}
	m.srv.SQLServers = m.srv.SQLServers[:len(m.srv.SQLServers)-1]

	m.fillRows()
}

func (m *FormAddService) fillRows() {
	m.rows = []table.Row{}

	for _, sqlsrv := range m.srv.SQLServers {
		rowT := table.Row{
			sqlsrv.Description,
			sqlsrv.Name,
			sqlsrv.User,
			"•",
		}
		m.rows = append(m.rows, rowT)
	}
}
