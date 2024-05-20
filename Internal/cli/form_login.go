package cli

import (
	"fmt"
	"strings"
	"time"

	"Service_1Cv8/internal/cli/styles"
	"Service_1Cv8/internal/constants"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	itemsFormLogin = 4
)

type FormLogin struct {
	focusIndex   int
	inputs       []textinput.Model
	buttonsFocus []bool
	message      string

	spinner    spinner.Model
	spinnering bool
}

func (m *FormLogin) SetParameters(args []interface{}) {
	m.inputs[0].Placeholder = "Login"
	m.inputs[0].SetValue("")
	m.inputs[0].Focus()
	m.inputs[0].PromptStyle = styles.FocusedStyleFB
	m.inputs[0].TextStyle = styles.FocusedStyleFB
	m.inputs[0].CharLimit = 19

	m.inputs[1].Placeholder = "Password"
	m.inputs[1].SetValue("")
	m.inputs[1].EchoMode = textinput.EchoPassword
	m.inputs[1].EchoCharacter = '•'
}

func NewFormLogin() *FormLogin {

	newSpinner := spinner.New()
	newSpinner.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))
	newSpinner.Spinner = spinner.Pulse

	m := FormLogin{
		inputs:       make([]textinput.Model, 2),
		buttonsFocus: make([]bool, 2),
		spinnering:   false,
		spinner:      newSpinner,
	}

	var t textinput.Model

	for i := range m.inputs {
		t = textinput.New()
		m.inputs[i] = t
	}

	return &m

}

func (m *FormLogin) Init() tea.Cmd {
	return nil
}

func (m *FormLogin) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

		case "tab", "shift+tab": //, "up", "down":

			s := msg.String()

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > itemsFormLogin-1 {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = itemsFormLogin - 1
			}

			cmds := make([]tea.Cmd, 0)
			if m.focusIndex <= len(m.inputs)-1 {
				cmds = make([]tea.Cmd, len(m.inputs))
				for i := 0; i <= len(m.inputs)-1; i++ {
					if i == m.focusIndex {
						// Set focused state
						cmds[i] = m.inputs[i].Focus()
						m.inputs[i].PromptStyle = styles.FocusedStyleFB
						m.inputs[i].TextStyle = styles.FocusedStyleFB
						continue
					}
					// Remove focused state
					m.inputs[i].Blur()
					m.inputs[i].PromptStyle = styles.NoStyleFB
					m.inputs[i].TextStyle = styles.NoStyleFB
				}
				m.buttonsFocus[0] = false
				m.buttonsFocus[1] = false
			} else if m.focusIndex == 2 {
				cmds = make([]tea.Cmd, 0)
				//for i := 0; i <= len(m.buttonsFocus)-1; i++ {

				m.inputs[len(m.inputs)-1].Blur()
				m.inputs[len(m.inputs)-1].PromptStyle = styles.NoStyleFB
				m.inputs[len(m.inputs)-1].TextStyle = styles.NoStyleFB

				m.buttonsFocus[0] = true
				m.buttonsFocus[1] = false
			} else if m.focusIndex == 3 {
				cmds = make([]tea.Cmd, 0)
				//for i := 0; i <= len(m.buttonsFocus)-1; i++ {

				m.buttonsFocus[0] = false
				m.buttonsFocus[1] = true
				//}
			}

			return m, tea.Batch(cmds...)
		case " ", "enter":

			frm, err := m.executeFormServer()
			if err == nil {
				model := frm.(*FormMain)
				model.SetParameters(nil)
				return model.Update(tea.WindowSizeMsg{130, 30})
			}
			m.message = err.Error()
		}
	}

	cmd := m.updateInputs(msg)
	return m, cmd
}

func (m *FormLogin) View() string {
	var b strings.Builder

	b.WriteString("Login\n\n")

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	b.WriteString("\n\n\n")

	if m.buttonsFocus[0] {
		b.WriteString(fmt.Sprintf("[ %s ]", styles.FocusedStyleFB.Render("sign in")))
	} else {
		b.WriteString(fmt.Sprintf("[ %s ]", styles.BlurredStyleFB.Render("sign in")))
	}
	b.WriteString(" ")
	if m.buttonsFocus[1] {
		b.WriteString(fmt.Sprintf("[ %s ]", styles.FocusedStyleFB.Render("sign up")))
	} else {
		b.WriteString(fmt.Sprintf("[ %s ]", styles.BlurredStyleFB.Render("sign up")))
	}

	s := ""
	if m.spinnering {
		s = m.spinner.View()
	}
	b.WriteString(fmt.Sprintf("\n\n %s ", s))

	statusVal := styles.StatusText.Copy().Width(styles.Width).Render(m.message)
	b.WriteString(statusVal)

	b.WriteString(styles.CursorModeHelpStyleFB.Render("\n\n\n" + " F10 exit program"))

	return b.String()

}

func (m *FormLogin) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m *FormLogin) executeFormServer() (CM, error) {

	for k, v := range m.buttonsFocus {
		switch k {
		case 0:
			if v {
				//go m.goSignIn()
				return client.Models[constants.FormMain], nil
			}
		case 1:
			if v {

			}
		}
	}
	return nil, nil
}

func (m *FormLogin) goSignIn() {

	m.spinnering = true
	now15 := time.Now().Add(15 * time.Second)
	for time.Now().Before(now15) {

	}
	m.spinnering = false

	pref := "sign in"
	m.message = fmt.Sprintf("%s/%s %s", m.inputs[0].View(), m.inputs[1].View(), pref)
	//if err != nil {
	//	m.error = err.Error()
	//}
}
