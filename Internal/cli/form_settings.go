package cli

import (
	"Service_1Cv8/internal/cli/styles"
	"Service_1Cv8/internal/constants"
	"Service_1Cv8/internal/repository"
	"Service_1Cv8/internal/winsys"
	"fmt"
	"github.com/charmbracelet/bubbles/textarea"
	"golang.org/x/text/encoding/charmap"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	itemsInputsFormSetting   = 12
	itemsAreasFormSetting    = 1
	itemsButtonFormSetting   = 1
	itemsCheckBoxFormSetting = 0
	itemsTableFormSetting    = 0
)

type buttonFormSetting struct {
	buttonSave bool
}

type inputsFormSetting struct {
	path1C          textinput.Model
	nameUser        textinput.Model
	passwordUser    textinput.Model
	startBlock      textinput.Model
	finishBlock     textinput.Model
	keyUnlock       textinput.Model
	pathStorage     textinput.Model
	pathCopy        textinput.Model
	httpPort        textinput.Model
	intervalService textinput.Model
	killProcessKb   textinput.Model
	controlServer   textinput.Model
}

type areasFormSetting struct {
	massage textarea.Model
}

type FormSetting struct {
	focusIndex int
	inputs     inputsFormSetting
	areas      areasFormSetting
	buttons    buttonFormSetting
	message    string

	spinner    spinner.Model
	spinnering bool
}

func (m *FormSetting) SetParameters(args []interface{}) {
	m.inputs.path1C.Placeholder = "Path 1C exe"
	m.inputs.path1C.SetValue(client.Storage.Settings.PathExe1C)
	m.inputs.path1C.Focus()
	m.inputs.path1C.PromptStyle = styles.FocusedStyleFB
	m.inputs.path1C.TextStyle = styles.FocusedStyleFB

	m.inputs.startBlock.Placeholder = "Start block"
	m.inputs.startBlock.SetValue(client.Storage.Settings.StartBlock)

	m.inputs.finishBlock.Placeholder = "Finish block"
	m.inputs.finishBlock.SetValue(client.Storage.Settings.FinishBlock)

	m.areas.massage.Placeholder = "Massage block"
	m.areas.massage.SetValue(client.Storage.Settings.Massage)
	if m.areas.massage.Value() == "" {
		m.areas.massage.SetValue(constants.MASSAGE)
	}
	m.areas.massage.SetWidth(250)
	m.areas.massage.SetHeight(5)

	m.inputs.keyUnlock.Placeholder = "Key unlock"
	m.inputs.keyUnlock.SetValue(client.Storage.Settings.KeyUnlock)

	m.inputs.nameUser.Placeholder = "Name user"
	m.inputs.nameUser.SetValue(client.Storage.Settings.NameUser)

	m.inputs.passwordUser.Placeholder = "Password user"
	m.inputs.passwordUser.SetValue(client.Storage.Settings.PasswordUser)
	m.inputs.passwordUser.EchoMode = textinput.EchoPassword
	m.inputs.passwordUser.EchoCharacter = '•'

	m.inputs.pathStorage.Placeholder = "Storage path"
	m.inputs.pathStorage.SetValue(client.Storage.Settings.PathStorage)

	m.inputs.pathCopy.Placeholder = "Copy path"
	m.inputs.pathCopy.SetValue(client.Storage.Settings.PathCopy)

	m.inputs.httpPort.Placeholder = "HTTP port"
	m.inputs.httpPort.SetValue(client.Storage.Settings.HTTPPort)

	m.inputs.intervalService.Placeholder = "interval service"
	m.inputs.intervalService.SetValue(client.Storage.Settings.IntervalService)

	m.inputs.killProcessKb.Placeholder = "Kill process (Kb)"
	m.inputs.killProcessKb.SetValue(client.Storage.Settings.KillProcessKb)

	m.inputs.controlServer.Placeholder = "Control server"
	m.inputs.controlServer.SetValue(client.Storage.Settings.ControlServer)
}

func NewFormSetting() *FormSetting {

	newSpinner := spinner.New()
	newSpinner.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))
	newSpinner.Spinner = spinner.Pulse

	m := FormSetting{
		buttons:    buttonFormSetting{buttonSave: false},
		spinnering: false,
		spinner:    newSpinner,
	}

	var t textinput.Model
	var a textarea.Model

	t = textinput.New()
	m.inputs.path1C = t

	t = textinput.New()
	m.inputs.keyUnlock = t

	t = textinput.New()
	m.inputs.startBlock = t
	m.inputs.startBlock.CharLimit = 19

	t = textinput.New()
	m.inputs.finishBlock = t
	m.inputs.finishBlock.CharLimit = 19

	t = textinput.New()
	m.inputs.nameUser = t

	t = textinput.New()
	m.inputs.passwordUser = t

	t = textinput.New()
	m.inputs.pathStorage = t

	t = textinput.New()
	m.inputs.pathCopy = t

	t = textinput.New()
	m.inputs.httpPort = t

	t = textinput.New()
	m.inputs.intervalService = t

	t = textinput.New()
	m.inputs.killProcessKb = t

	t = textinput.New()
	m.inputs.controlServer = t

	a = textarea.New()
	a.CharLimit = 2000
	m.areas.massage = a

	return &m
}

func (m *FormSetting) Init() tea.Cmd {
	return nil
}

func (m *FormSetting) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	cmds := make([]tea.Cmd, 0)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		strMsg := msg.String()
		switch strMsg {
		case "esc", "ctrl+q":

			frm := client.Models[constants.FormMain].(*FormMain)
			frm.SetParameters(nil)

			return frm, nil
		case "ctrl+c", "f10":

			return m, tea.Quit

		case "tab", "shift+tab", "up", "down":

			if m.focusIndex == 3 &&
				(strMsg == "up" || strMsg == "down") {

				if m.focusIndex == 3 {
					cmd = m.updateInputs(msg)
					return m, cmd
				}
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
			case 0: //path1C
				m.buttons.buttonSave = false

				m.inputs.path1C.Focus()
				m.inputs.path1C.PromptStyle = styles.FocusedStyleFB
				m.inputs.path1C.TextStyle = styles.FocusedStyleFB

				m.inputs.startBlock.Blur()
				m.inputs.startBlock.PromptStyle = styles.NoStyleFB
				m.inputs.startBlock.TextStyle = styles.NoStyleFB

			case 1: //startBlock
				m.inputs.path1C.Blur()
				m.inputs.path1C.PromptStyle = styles.NoStyleFB
				m.inputs.path1C.TextStyle = styles.NoStyleFB

				m.inputs.startBlock.Focus()
				m.inputs.startBlock.PromptStyle = styles.FocusedStyleFB
				m.inputs.startBlock.TextStyle = styles.FocusedStyleFB

				m.inputs.finishBlock.Blur()
				m.inputs.finishBlock.PromptStyle = styles.NoStyleFB
				m.inputs.finishBlock.TextStyle = styles.NoStyleFB
			case 2: //finishBlock
				m.inputs.startBlock.Blur()
				m.inputs.startBlock.PromptStyle = styles.NoStyleFB
				m.inputs.startBlock.TextStyle = styles.NoStyleFB

				m.inputs.finishBlock.Focus()
				m.inputs.finishBlock.PromptStyle = styles.FocusedStyleFB
				m.inputs.finishBlock.TextStyle = styles.FocusedStyleFB

				m.areas.massage.Blur()

			case 3: //massage
				m.inputs.finishBlock.Blur()
				m.inputs.finishBlock.PromptStyle = styles.NoStyleFB
				m.inputs.finishBlock.TextStyle = styles.NoStyleFB

				m.areas.massage.Focus()

				m.inputs.keyUnlock.Blur()
				m.inputs.keyUnlock.PromptStyle = styles.NoStyleFB
				m.inputs.keyUnlock.TextStyle = styles.NoStyleFB

			case 4: //keyUnlock
				m.areas.massage.Blur()

				m.inputs.keyUnlock.Focus()
				m.inputs.keyUnlock.PromptStyle = styles.FocusedStyleFB
				m.inputs.keyUnlock.TextStyle = styles.FocusedStyleFB

				m.buttons.buttonSave = false

			case 5: //user
				m.inputs.keyUnlock.Blur()
				m.inputs.keyUnlock.PromptStyle = styles.NoStyleFB
				m.inputs.keyUnlock.TextStyle = styles.NoStyleFB

				m.inputs.nameUser.Focus()
				m.inputs.nameUser.PromptStyle = styles.FocusedStyleFB
				m.inputs.nameUser.TextStyle = styles.FocusedStyleFB

				m.inputs.passwordUser.Blur()
				m.inputs.passwordUser.PromptStyle = styles.NoStyleFB
				m.inputs.passwordUser.TextStyle = styles.NoStyleFB

			case 6: // password
				m.inputs.nameUser.Blur()
				m.inputs.nameUser.PromptStyle = styles.NoStyleFB
				m.inputs.nameUser.TextStyle = styles.NoStyleFB

				m.inputs.passwordUser.Focus()
				m.inputs.passwordUser.PromptStyle = styles.FocusedStyleFB
				m.inputs.passwordUser.TextStyle = styles.FocusedStyleFB

				m.buttons.buttonSave = false

			case 7: //pathStorage
				m.inputs.passwordUser.Blur()
				m.inputs.passwordUser.PromptStyle = styles.NoStyleFB
				m.inputs.passwordUser.TextStyle = styles.NoStyleFB

				m.inputs.pathStorage.Focus()
				m.inputs.pathStorage.PromptStyle = styles.FocusedStyleFB
				m.inputs.pathStorage.TextStyle = styles.FocusedStyleFB

				m.inputs.pathCopy.Blur()
				m.inputs.pathCopy.PromptStyle = styles.NoStyleFB
				m.inputs.pathCopy.TextStyle = styles.NoStyleFB
			case 8: //table
				m.inputs.pathStorage.Blur()
				m.inputs.pathStorage.PromptStyle = styles.NoStyleFB
				m.inputs.pathStorage.TextStyle = styles.NoStyleFB

				m.inputs.pathCopy.Focus()
				m.inputs.pathCopy.PromptStyle = styles.FocusedStyleFB
				m.inputs.pathCopy.TextStyle = styles.FocusedStyleFB

				m.inputs.httpPort.Blur()
				m.inputs.httpPort.PromptStyle = styles.NoStyleFB
				m.inputs.httpPort.TextStyle = styles.NoStyleFB
			case 9: //table
				m.inputs.pathCopy.Blur()
				m.inputs.pathCopy.PromptStyle = styles.NoStyleFB
				m.inputs.pathCopy.TextStyle = styles.NoStyleFB

				m.inputs.httpPort.Focus()
				m.inputs.httpPort.PromptStyle = styles.FocusedStyleFB
				m.inputs.httpPort.TextStyle = styles.FocusedStyleFB

				m.inputs.intervalService.Blur()
				m.inputs.intervalService.PromptStyle = styles.NoStyleFB
				m.inputs.intervalService.TextStyle = styles.NoStyleFB
			case 10: //table
				m.inputs.httpPort.Blur()
				m.inputs.httpPort.PromptStyle = styles.NoStyleFB
				m.inputs.httpPort.TextStyle = styles.NoStyleFB

				m.inputs.intervalService.Focus()
				m.inputs.intervalService.PromptStyle = styles.FocusedStyleFB
				m.inputs.intervalService.TextStyle = styles.FocusedStyleFB

				m.inputs.killProcessKb.Blur()
				m.inputs.killProcessKb.PromptStyle = styles.NoStyleFB
				m.inputs.killProcessKb.TextStyle = styles.NoStyleFB
			case 11: //table
				m.inputs.intervalService.Blur()
				m.inputs.intervalService.PromptStyle = styles.NoStyleFB
				m.inputs.intervalService.TextStyle = styles.NoStyleFB

				m.inputs.killProcessKb.Focus()
				m.inputs.killProcessKb.PromptStyle = styles.FocusedStyleFB
				m.inputs.killProcessKb.TextStyle = styles.FocusedStyleFB

				m.inputs.controlServer.Blur()
				m.inputs.controlServer.PromptStyle = styles.NoStyleFB
				m.inputs.controlServer.TextStyle = styles.NoStyleFB
			case 12: //table
				m.inputs.killProcessKb.Blur()
				m.inputs.killProcessKb.PromptStyle = styles.NoStyleFB
				m.inputs.killProcessKb.TextStyle = styles.NoStyleFB

				m.inputs.controlServer.Focus()
				m.inputs.controlServer.PromptStyle = styles.FocusedStyleFB
				m.inputs.controlServer.TextStyle = styles.FocusedStyleFB

				m.buttons.buttonSave = false
			case 13: //table
				m.inputs.controlServer.Blur()
				m.inputs.controlServer.PromptStyle = styles.NoStyleFB
				m.inputs.controlServer.TextStyle = styles.NoStyleFB

				m.buttons.buttonSave = true

				m.inputs.path1C.Blur()
				m.inputs.path1C.PromptStyle = styles.NoStyleFB
				m.inputs.path1C.TextStyle = styles.NoStyleFB
			}
			return m, tea.Batch(cmds...)
		case "enter":

			m.saveSettings()

		}
	}

	cmd = m.updateInputs(msg)
	return m, cmd
}

func (m *FormSetting) View() string {
	var b strings.Builder

	b.WriteRune('\n')
	b.WriteString(fmt.Sprintf("%s Setting %s\n", styles.ShortLine, styles.Line))
	b.WriteRune('\n')

	b.WriteString(fmt.Sprintf("%s: %s\n", styles.СaptionStyleFB.Render("Patch 1C:"), m.inputs.path1C.View()))

	b.WriteString(fmt.Sprintf("%s %s %s %s\n",
		styles.СaptionStyleFB.Render("Time start block:"), m.inputs.startBlock.View(),
		styles.СaptionStyleFB.Render("Time finish block:"), m.inputs.finishBlock.View()))

	b.WriteString(fmt.Sprintf("%s\n%s\n", styles.СaptionStyleFB.Render("Massage:"), m.areas.massage.View()))

	b.WriteString(fmt.Sprintf("%s %s\n", styles.СaptionStyleFB.Render("Key unlock:"), m.inputs.keyUnlock.View()))

	b.WriteString(fmt.Sprintf("%s %s %s %s\n",
		styles.СaptionStyleFB.Render("User:"), m.inputs.nameUser.View(),
		styles.СaptionStyleFB.Render("Password:"), m.inputs.passwordUser.View()))

	b.WriteString(fmt.Sprintf("%s %s\n%s %s\n",
		styles.СaptionStyleFB.Render("Storage path:"), m.inputs.pathStorage.View(),
		styles.СaptionStyleFB.Render("Copy path:"), m.inputs.pathCopy.View()))

	b.WriteRune('\n')
	b.WriteString(fmt.Sprintf("%s Server setting %s\n", styles.ShortLine, styles.Line))
	b.WriteRune('\n')

	b.WriteString(fmt.Sprintf("%s %s %s %s %s %s %s %s\n",
		styles.СaptionStyleFB.Render("HTTP port:"), m.inputs.httpPort.View(),
		styles.СaptionStyleFB.Render("Interval service:"), m.inputs.intervalService.View(),
		styles.СaptionStyleFB.Render("Kill process (Kb):"), m.inputs.killProcessKb.View(),
		styles.СaptionStyleFB.Render("Control server:"), m.inputs.controlServer.View()))

	b.WriteRune('\n')

	if m.buttons.buttonSave {
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

	b.WriteString(styles.CursorModeHelpStyleFB.Render("\n\n\n" + " ESC exit main menu | F10 exit program"))

	return b.String()

}

func (m *FormSetting) lenForm() int {
	return itemsInputsFormSetting + itemsAreasFormSetting + itemsButtonFormSetting +
		itemsTableFormSetting + itemsCheckBoxFormSetting
}

func (m *FormSetting) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, itemsInputsFormSetting+itemsAreasFormSetting)

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

	m.inputs.keyUnlock, _ = m.inputs.keyUnlock.Update(msg)
	m.inputs.nameUser, _ = m.inputs.nameUser.Update(msg)
	m.inputs.passwordUser, _ = m.inputs.passwordUser.Update(msg)
	m.inputs.startBlock, _ = m.inputs.startBlock.Update(msg)
	m.inputs.finishBlock, _ = m.inputs.finishBlock.Update(msg)
	m.inputs.path1C, _ = m.inputs.path1C.Update(msg)
	m.inputs.pathStorage, _ = m.inputs.pathStorage.Update(msg)
	m.inputs.pathCopy, _ = m.inputs.pathCopy.Update(msg)
	m.inputs.httpPort, _ = m.inputs.httpPort.Update(msg)
	m.inputs.intervalService, _ = m.inputs.intervalService.Update(msg)
	m.inputs.killProcessKb, _ = m.inputs.killProcessKb.Update(msg)
	m.inputs.controlServer, _ = m.inputs.controlServer.Update(msg)

	m.areas.massage, _ = m.areas.massage.Update(msg)

	return tea.Batch(cmds...)
}

func (m *FormSetting) saveSettings() {

	client.Storage.Settings = repository.Settings{
		PathExe1C:       m.inputs.path1C.Value(),
		PathStorage:     m.inputs.pathStorage.Value(),
		PathCopy:        m.inputs.pathCopy.Value(),
		NameUser:        m.inputs.nameUser.Value(),
		PasswordUser:    m.inputs.passwordUser.Value(),
		StartBlock:      m.inputs.startBlock.Value(),
		FinishBlock:     m.inputs.finishBlock.Value(),
		KeyUnlock:       m.inputs.keyUnlock.Value(),
		HTTPPort:        m.inputs.httpPort.Value(),
		IntervalService: m.inputs.intervalService.Value(),
		KillProcessKb:   m.inputs.killProcessKb.Value(),
		ControlServer:   m.inputs.controlServer.Value(),
		Massage:         m.areas.massage.Value(),
	}

	err := client.Storage.WriteYamlData()
	if err != nil {
		m.message = err.Error()
	}
	m.message = "OK save settings"

}
