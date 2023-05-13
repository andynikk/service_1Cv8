package cli

import (
	"Service_1Cv8/internal/winsys"
	"fmt"
	"golang.org/x/text/encoding/charmap"
	"strings"

	"Service_1Cv8/internal/cli/styles"
	"Service_1Cv8/internal/constants"
	"Service_1Cv8/internal/repository"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	itemsInputsFormAddSQLSrv = 4
	itemsAreasFormAddSQLSrv  = 0
	itemsButtonFormAddSQLSrv = 1
	itemsTableFormAddSQLSrv  = 0
)

type inputsFormAddSQLSrv struct {
	description textinput.Model
	name        textinput.Model
	user        textinput.Model
	password    textinput.Model
}

type buttonsFormAddSQLSrv struct {
	save bool
}

type FormAddSQLSrv struct {
	focusIndex int
	inputs     inputsFormAddSQLSrv
	buttons    buttonsFormAddSQLSrv

	uid string
	new bool

	message string

	srv repository.SQLServer
}

func (m *FormAddSQLSrv) SetParameters(args []interface{}) {

	m.uid = ""
	m.message = ""
	m.focusIndex = 0
	m.buttons.save = false

	m.inputs.description.Placeholder = "Description"
	m.inputs.description.SetValue("")
	m.inputs.description.Focus()
	m.inputs.description.PromptStyle = styles.FocusedStyleFB
	m.inputs.description.TextStyle = styles.FocusedStyleFB

	m.inputs.name.Placeholder = "Name server"
	m.inputs.name.SetValue("")

	m.inputs.user.Placeholder = "User"
	m.inputs.user.SetValue("")

	m.inputs.password.Placeholder = "Password"
	m.inputs.password.SetValue("")
	m.inputs.password.EchoMode = textinput.EchoPassword
	m.inputs.password.EchoCharacter = '•'

	for _, v := range args {
		switch v.(type) {
		case repository.SQLServer:
			m.srv = v.(repository.SQLServer)

			if !m.new {
				m.uid = m.srv.Name
			}

			m.inputs.description.SetValue(m.srv.Description)
			m.inputs.name.SetValue(m.srv.Name)
			m.inputs.user.SetValue(m.srv.User)
			m.inputs.password.SetValue(m.srv.Password)
		}
	}
}

func NewFormAddSQLSrv() *FormAddSQLSrv {

	m := FormAddSQLSrv{
		inputs:  inputsFormAddSQLSrv{},
		buttons: buttonsFormAddSQLSrv{},
	}

	var t textinput.Model

	t = textinput.New()
	m.inputs.description = t

	t = textinput.New()
	m.inputs.name = t

	t = textinput.New()
	m.inputs.user = t

	t = textinput.New()
	m.inputs.password = t

	return &m

}

func (m *FormAddSQLSrv) Init() tea.Cmd {
	return nil
}

func (m *FormAddSQLSrv) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	cmds := make([]tea.Cmd, 0)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		strMsg := msg.String()
		switch strMsg {
		case "esc", "ctrl+q":
			m.uid = ""
			m.srv = repository.SQLServer{}

			frm := client.Models[constants.FormAddService].(*FormAddService)

			frm.srv.SQLServers = append(frm.srv.SQLServers, m.srv)
			arg := []interface{}{frm.srv}
			frm.SetParameters(arg)

			return frm, nil
		case "ctrl+c", "f10":
			return m, tea.Quit
		case "tab", "shift+tab", "up", "down":

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
			case 0: //description
				m.buttons.save = false

				m.inputs.description.Focus()
				m.inputs.description.PromptStyle = styles.FocusedStyleFB
				m.inputs.description.TextStyle = styles.FocusedStyleFB

				m.inputs.name.Blur()
				m.inputs.name.PromptStyle = styles.NoStyleFB
				m.inputs.name.TextStyle = styles.NoStyleFB

			case 1: //name
				m.inputs.description.Blur()
				m.inputs.description.PromptStyle = styles.NoStyleFB
				m.inputs.description.TextStyle = styles.NoStyleFB

				m.inputs.name.Focus()
				m.inputs.name.PromptStyle = styles.FocusedStyleFB
				m.inputs.name.TextStyle = styles.FocusedStyleFB

				m.inputs.user.Blur()
				m.inputs.user.PromptStyle = styles.NoStyleFB
				m.inputs.user.TextStyle = styles.NoStyleFB
			case 2: //user
				m.inputs.name.Blur()
				m.inputs.name.PromptStyle = styles.NoStyleFB
				m.inputs.name.TextStyle = styles.NoStyleFB

				m.inputs.user.Focus()
				m.inputs.user.PromptStyle = styles.FocusedStyleFB
				m.inputs.user.TextStyle = styles.FocusedStyleFB

				m.inputs.password.Blur()
				m.inputs.password.PromptStyle = styles.NoStyleFB
				m.inputs.password.TextStyle = styles.NoStyleFB

			case 3: //password
				m.inputs.user.Blur()
				m.inputs.user.PromptStyle = styles.NoStyleFB
				m.inputs.user.TextStyle = styles.NoStyleFB

				m.inputs.password.Focus()
				m.inputs.password.PromptStyle = styles.FocusedStyleFB
				m.inputs.password.TextStyle = styles.FocusedStyleFB

				m.buttons.save = false

			case 4: //save
				m.inputs.password.Blur()
				m.inputs.password.PromptStyle = styles.NoStyleFB
				m.inputs.password.TextStyle = styles.NoStyleFB

				m.buttons.save = true

				m.inputs.description.Blur()
				m.inputs.description.PromptStyle = styles.NoStyleFB
				m.inputs.description.TextStyle = styles.NoStyleFB
			}
			return m, tea.Batch(cmds...)
		case "enter":
			model, cmd := m.clickButton()
			return model, cmd
		}
	}

	cmd = m.updateInputs(msg)
	return m, cmd
}

func (m *FormAddSQLSrv) View() string {
	var b strings.Builder

	nameSrv := m.srv.Name
	if m.new {
		nameSrv = "NEW"
	}

	b.WriteString(fmt.Sprintf("%s Data base %s %s\n", styles.ShortLine, nameSrv, styles.Line))
	b.WriteString(fmt.Sprintf("%s %s %s %s\n",
		styles.СaptionStyleFB.Render("Description:"), m.inputs.description.View(),
		styles.СaptionStyleFB.Render("Name:"), m.inputs.name.View()))

	b.WriteString(fmt.Sprintf("%s %s %s %s\n",
		styles.СaptionStyleFB.Render("User:"), m.inputs.user.View(),
		styles.СaptionStyleFB.Render("Password:"), m.inputs.password.View()))

	b.WriteString("\n\n\n")

	if m.buttons.save {
		b.WriteString(fmt.Sprintf("[ %s ]", styles.FocusedStyleFB.Render("save")))
	} else {
		b.WriteString(fmt.Sprintf("[ %s ]", styles.BlurredStyleFB.Render("save")))
	}
	b.WriteString(" ")

	s := ""

	b.WriteString(s)

	b.WriteString(fmt.Sprintf("\n %s ", s))

	statusVal := styles.StatusText.Copy().Width(styles.Width).Render(m.message)
	b.WriteString(statusVal)

	b.WriteString(styles.CursorModeHelpStyleFB.Render("\n\n" + " ESC exit main menu | F10 exit program"))

	return b.String()

}

func (m *FormAddSQLSrv) lenForm() int {
	return itemsInputsFormAddSQLSrv + itemsAreasFormAddSQLSrv + itemsButtonFormAddSQLSrv + itemsTableFormAddSQLSrv
}

func (m *FormAddSQLSrv) updateInputs(msg tea.Msg) tea.Cmd {
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

	m.inputs.description, _ = m.inputs.description.Update(msg)
	m.inputs.name, _ = m.inputs.name.Update(msg)
	m.inputs.user, _ = m.inputs.user.Update(msg)
	m.inputs.password, _ = m.inputs.password.Update(msg)

	return tea.Batch(cmds...)
}

func (m *FormAddSQLSrv) clickButton() (tea.Model, tea.Cmd) {

	if m.buttons.save {
		model, _ := m.saveSQLServer()
		return model.Update(nil)
	}

	m.message = "edit OK"

	return nil, nil
}

func (m *FormAddSQLSrv) saveSQLServer() (tea.Model, tea.Cmd) {

	m.srv.Description = m.inputs.description.Value()
	m.srv.Name = m.inputs.name.Value()
	m.srv.User = m.inputs.user.Value()
	m.srv.Password = m.inputs.password.Value()

	model := client.Models[constants.FormAddService].(*FormAddService)
	model.srv.SQLServers = append(model.srv.SQLServers, m.srv)
	model.fillRows()

	m.message = "edit OK"

	return model.Update(nil)

}
