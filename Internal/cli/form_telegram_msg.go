package cli

import (
	"Service_1Cv8/internal/cli/charm"
	"Service_1Cv8/internal/cli/styles"
	"Service_1Cv8/internal/compression"
	"Service_1Cv8/internal/constants"
	"Service_1Cv8/internal/environment"
	"Service_1Cv8/internal/iron"
	"Service_1Cv8/internal/telegram"
	"Service_1Cv8/internal/winsys"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"golang.org/x/text/encoding/charmap"
	"net/http"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	itemsInputsFormTgMsg   = 2
	itemsAreasFormTgMsg    = 1
	itemsButtonFormTgMsg   = 3
	itemsCheckBoxFormTgMsg = 0
	itemsTableFormTgMsg    = 0
)

type buttonFormTgMsg struct {
	buttonDefault charm.BoolModel
	buttonSend    charm.BoolModel
	buttonNullify charm.BoolModel
}

type inputsFormTgMsg struct {
	id  textinput.Model
	api textinput.Model
}

type areasFormTgMsg struct {
	massage textarea.Model
}

type FormTgMsg struct {
	focusIndex int
	inputs     inputsFormTgMsg
	areas      areasFormTgMsg
	buttons    buttonFormTgMsg
	message    string

	spinner    spinner.Model
	spinnering bool

	charm.CrawlElements
}

func (m *FormTgMsg) SetParameters(args []interface{}) {
	m.buttons.buttonDefault.Active = true

	m.inputs.id.Placeholder = "Telegram ID"
	m.inputs.id.SetValue(client.Storage.Settings.TgID)

	m.inputs.api.Placeholder = "Telegram API"
	m.inputs.api.SetValue(client.Storage.Settings.TgAPI)

	m.areas.massage.Placeholder = "Telegram message"
	m.areas.massage.SetValue(client.Storage.Settings.PathExe1C)
	m.areas.massage.Focus()

	m.FillFormTgMsgElements()
}

func (m *FormTgMsg) FillFormTgMsgElements() {
	ce := make(charm.CrawlElements)

	i := -1
	increaseI := func() int { i++; return i }

	ce[increaseI()] = &m.buttons.buttonDefault
	ce[increaseI()] = &m.inputs.id
	ce[increaseI()] = &m.inputs.api
	ce[increaseI()] = &m.areas.massage
	ce[increaseI()] = &m.buttons.buttonNullify
	ce[increaseI()] = &m.buttons.buttonSend

	m.CrawlElements = ce

	m.Init()
}

func NewFormTgMsg() *FormTgMsg {

	newSpinner := spinner.New()
	newSpinner.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))
	newSpinner.Spinner = spinner.Pulse

	m := FormTgMsg{
		buttons:    buttonFormTgMsg{buttonDefault: charm.BoolModel{}, buttonNullify: charm.BoolModel{}, buttonSend: charm.BoolModel{}},
		spinnering: false,
		spinner:    newSpinner,
	}

	var a textarea.Model
	var t textinput.Model

	t = textinput.New()
	m.inputs.id = t

	t = textinput.New()
	m.inputs.api = t

	a = textarea.New()
	a.CharLimit = 2000
	m.areas.massage = a

	return &m
}

func (m *FormTgMsg) Init() tea.Cmd {
	m.defaultMsg()
	return nil
}

func (m *FormTgMsg) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

			if m.areas.massage.Focused() &&
				(strMsg == "up" || strMsg == "down") {
				cmd = m.updateInputs(msg)
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
		case "enter":
			if m.buttons.buttonDefault.Active {
				m.defaultMsg()
			}
			if m.buttons.buttonSend.Active {
				m.sendMsg()
			}
			if m.buttons.buttonNullify.Active {
				m.nullifyDataMsg()
			}
		}
	}

	cmd = m.updateInputs(msg)
	return m, cmd
}

func (m *FormTgMsg) View() string {
	var b strings.Builder

	b.WriteRune('\n')
	b.WriteString(fmt.Sprintf("%s Telegram message %s\n", styles.ShortLine, styles.Line))
	b.WriteRune('\n')

	if m.buttons.buttonDefault.Active {
		b.WriteString(fmt.Sprintf("[ %s ]", styles.FocusedStyleFB.Render("default msg")))
	} else {
		b.WriteString(fmt.Sprintf("[ %s ]", styles.BlurredStyleFB.Render("default msg")))
	}
	b.WriteRune('\n')

	b.WriteString(fmt.Sprintf("%s %s\n%s %s\n",
		styles.СaptionStyleFB.Render("Telegram ID:"), m.inputs.id.View(),
		styles.СaptionStyleFB.Render("Telegram API:"), m.inputs.api.View()))

	b.WriteString(fmt.Sprintf("%s\n%s\n", styles.СaptionStyleFB.Render("Massage:"), m.areas.massage.View()))

	if m.buttons.buttonNullify.Active {
		b.WriteString(fmt.Sprintf("[ %s ]", styles.FocusedStyleFB.Render("nullify")))
	} else {
		b.WriteString(fmt.Sprintf("[ %s ]", styles.BlurredStyleFB.Render("nullify")))
	}
	b.WriteString(" ")

	if m.buttons.buttonSend.Active {
		b.WriteString(fmt.Sprintf("[ %s ]", styles.FocusedStyleFB.Render("send")))
	} else {
		b.WriteString(fmt.Sprintf("[ %s ]", styles.BlurredStyleFB.Render("send")))
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

func (m *FormTgMsg) lenForm() int {
	return itemsInputsFormTgMsg + itemsAreasFormTgMsg + itemsButtonFormTgMsg +
		itemsTableFormTgMsg + itemsCheckBoxFormTgMsg
}

func (m *FormTgMsg) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, itemsInputsFormTgMsg+itemsAreasFormTgMsg)

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

	m.inputs.id, _ = m.inputs.id.Update(msg)
	m.inputs.api, _ = m.inputs.api.Update(msg)
	m.areas.massage, _ = m.areas.massage.Update(msg)

	return tea.Batch(cmds...)
}

func (m *FormTgMsg) defaultMsg() {
	txtUpdateDB := fmt.Sprintf("**1. Обновленные базы:**  *_%s_*",
		strings.Join(client.PerformedActions.UpdateDB, ", "))
	txtRestartService := fmt.Sprintf("**2. Серверный кэш почищен:**  *_%s_*",
		strings.Join(client.PerformedActions.RestartService, ", "))
	txtClearCash := fmt.Sprintf("**3. Службы перезапущены:** *_%s_*",
		strings.Join(client.PerformedActions.ClearCash, ", "))
	txtRebutServer := fmt.Sprintf("**4. Сервера перегружены:** *_%s_*",
		strings.Join(client.PerformedActions.RebutServer, ", "))

	m.areas.massage.SetValue(fmt.Sprintf("%s\n%s\n%s\n%s",
		txtUpdateDB, txtRestartService, txtClearCash, txtRebutServer))
}

func (m *FormTgMsg) nullifyDataMsg() {
	client.PerformedActions = PerformedActions{}
}

func (m *FormTgMsg) sendMsg() {

	macAdress, err := iron.GetMacAddr()
	if err != nil {
		m.message = fmt.Sprintf("%s: %s", "Get MAC adress", err.Error())
		return
	}
	snHDD, err := iron.GetDiskDrivers()
	if err != nil {
		m.message = fmt.Sprintf("%s: %s", "Get HDD s/n", err.Error())
		return
	}

	tgEmoji := telegram.TgEmoji{"\U00002611", " Готово!"}
	i, err := strconv.ParseInt(m.inputs.id.Value(), 10, 64)
	if err != nil {
		m.message = err.Error()
		return
	}

	tgMsg := telegram.TgMsg{i, m.inputs.api.Value(), m.areas.massage.Value(), tgEmoji}
	jsonMsg, err := json.Marshal(tgMsg)
	if err != nil {
		m.message = err.Error()
		return
	}

	gzip, err := compression.Compress(jsonMsg)
	if err != nil {
		m.message = err.Error()
		return
	}

	url := environment.ServerUrlApi() + "/tg/send"
	//url := "tk-test-app.telematika.local:7171/tg/send"
	req, err := http.NewRequest("POST", url, bytes.NewReader(gzip))
	if err != nil {
		m.message = fmt.Sprintf("%d. %s", 1, err.Error())
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Authentication-Key", fmt.Sprintf("%s%s", macAdress, snHDD))
	defer req.Body.Close()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		m.message = fmt.Sprintf("%d. %s", 2, err.Error())
		return
	}
	defer resp.Body.Close()

	m.message = resp.Status
}
