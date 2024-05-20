package forms_exchange

import (
	"Service_1Cv8/internal/environment"
	"Service_1Cv8/internal/exchange"
	"Service_1Cv8/internal/files"
	"bytes"
	"context"
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/text/encoding/charmap"
	"io"
	"net/http"
	"strings"
	"time"

	"Service_1Cv8/internal/cli/charm"
	"Service_1Cv8/internal/cli/styles"
	"Service_1Cv8/internal/constants"
	"Service_1Cv8/internal/winsys"
)

const (
	itemsInputsFormQueue   = 1
	itemsAreasFormQueue    = 0
	itemsButtonFormQueue   = 4
	itemsCheckBoxFormQueue = 2
	itemsTableFormQueue    = 0
)

type buttonFormQueue struct {
	buttonCreate charm.BoolModel
	buttonMove   charm.BoolModel
	buttonClean  charm.BoolModel
	buttonDel    charm.BoolModel
	//buttonProlong charm.BoolModel
}

type inputsFormQueue struct {
	name       textinput.Model
	downloadAt textinput.Model
	uploadAt   textinput.Model
	counts     textinput.Model
	size       textinput.Model
}

type checkBoxFormQueue struct {
	soft charm.BoolModel
	hard charm.BoolModel
}

type areasFormQueue struct {
}

type FormQueue struct {
	focusIndex int

	inputs     inputsFormQueue
	areas      areasFormQueue
	buttons    buttonFormQueue
	checkBoxes checkBoxFormQueue

	message     string
	errorConnet bool

	exchangeQueueInfo exchange.ExchangeQueueInfo

	spinner    spinner.Model
	spinnering bool

	crawlElements charm.CrawlElements

	chanOut chan bool
}

func (m *FormQueue) refreshData() {
	m.inputs.name.SetValue(m.exchangeQueueInfo.Name)
	m.inputs.uploadAt.SetValue(m.exchangeQueueInfo.UploadAt.Format("2006-01-02T15:04:05"))
	m.inputs.downloadAt.SetValue(m.exchangeQueueInfo.DownloadAt.Format("2006-01-02T15:04:05"))
	m.inputs.counts.SetValue(m.exchangeQueueInfo.Messages)
	m.inputs.size.SetValue(m.exchangeQueueInfo.Size)

	if m.exchangeQueueInfo.TypeStorage == constants.Hard.String() {
		m.checkBoxes.soft.Value = false
		m.checkBoxes.hard.Value = true
	}
}

func (m *FormQueue) SetParameters(args []interface{}) {

	m.inputs.name.Placeholder = "Name"
	m.inputs.name.SetValue("")
	m.inputs.name.Focus()
	m.inputs.name.PromptStyle = styles.FocusedStyleFB
	m.inputs.name.TextStyle = styles.FocusedStyleFB

	m.inputs.uploadAt.Placeholder = "Upload At"
	m.inputs.uploadAt.SetValue("")

	m.inputs.downloadAt.Placeholder = "Download At"
	m.inputs.downloadAt.SetValue("")

	m.inputs.counts.Placeholder = "Counts"
	m.inputs.counts.SetValue("")

	m.inputs.size.Placeholder = "Size"
	m.inputs.size.SetValue("")

	m.checkBoxes.soft.Value = true
	m.checkBoxes.hard.Value = false

	m.focusIndex = 0

	for _, v := range args {
		switch v.(type) {
		case exchange.ExchangeQueueInfo:

			m.exchangeQueueInfo = v.(exchange.ExchangeQueueInfo)
			m.refreshData()

		}
	}

	m.FillFormQueueElements()

	m.chanOut = make(chan bool)
	ctx := context.Background()
	go m.wsQueueInfo(ctx)

	m.Init()
}

func (m *FormQueue) FillFormQueueElements() {
	ce := make(charm.CrawlElements)

	ce[0] = &m.inputs.name
	ce[1] = &m.checkBoxes.soft
	ce[2] = &m.checkBoxes.hard
	ce[3] = &m.buttons.buttonCreate
	ce[4] = &m.buttons.buttonClean
	ce[5] = &m.buttons.buttonMove
	ce[6] = &m.buttons.buttonDel

	m.crawlElements = ce
}

func NewFormQueue() *FormQueue {

	newSpinner := spinner.New()
	newSpinner.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))
	newSpinner.Spinner = spinner.Pulse

	m := FormQueue{
		buttons:    buttonFormQueue{},
		spinnering: false,
		spinner:    newSpinner,
	}

	var t textinput.Model

	t = textinput.New()
	m.inputs.name = t

	t = textinput.New()
	m.inputs.uploadAt = t

	t = textinput.New()
	m.inputs.downloadAt = t

	t = textinput.New()
	m.inputs.counts = t

	t = textinput.New()
	m.inputs.size = t

	return &m
}

func (m *FormQueue) Init() tea.Cmd {

	return m.spinner.Tick
}

func (m *FormQueue) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	cmds := make([]tea.Cmd, 0)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		strMsg := msg.String()
		switch strMsg {
		case "esc", "ctrl+q":

			if !m.errorConnet {
				m.chanOut <- false
			}

			model := exchanger.Models[constants.FormExchangeListQueues].(*FormListQueues)
			model.SetParameters(nil)

			return model.Update(nil)
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

			charm.TеraversingFormElements(m.crawlElements, m.focusIndex)

			return m, tea.Batch(cmds...)
		case "enter", " ":
			switch m.focusIndex {
			case 1:
				m.checkBoxes.soft.Value = true
				m.checkBoxes.hard.Value = false

				m.exchangeQueueInfo.TypeStorage = constants.Soft.String()
			case 2:
				m.checkBoxes.soft.Value = false
				m.checkBoxes.hard.Value = true

				m.exchangeQueueInfo.TypeStorage = constants.Hard.String()
			case 3:
				m.createQueue()
			case 4:
				m.clearMessages()
			case 5:
				m.moveMessages()
			case 6:
				m.deleteQueue()
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

		if !m.errorConnet {
			m.chanOut <- true
		}

		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	cmd = m.updateInputs(msg)
	return m, cmd
}

func (m *FormQueue) View() string {
	var b strings.Builder

	b.WriteRune('\n')
	b.WriteString(fmt.Sprintf("%s Queue %s\n", styles.ShortLine, styles.Line))
	b.WriteRune('\n')

	b.WriteString(fmt.Sprintf("%s%s\n", styles.СaptionStyleFB.Render("Name:"), m.inputs.name.View()))

	b.WriteString(fmt.Sprintf("%s %s %s %s\n",
		styles.СaptionStyleFB.Render("Upload At:"), m.inputs.uploadAt.View(),
		styles.СaptionStyleFB.Render("Download At:"), m.inputs.downloadAt.View()))

	b.WriteString(fmt.Sprintf("%s %s %s %s\n",
		styles.СaptionStyleFB.Render("Messages:"), m.inputs.counts.View(),
		styles.СaptionStyleFB.Render("Size (Kb):"), m.inputs.size.View()))

	b.WriteRune('\n')
	b.WriteString("Type")
	b.WriteRune('\n')

	cursor := " "
	checked := " "
	if m.checkBoxes.soft.Active {
		cursor = styles.FocusedStyleFB.Render(">")
	}
	if m.checkBoxes.soft.Value {
		checked = "x"
	}
	b.WriteString(fmt.Sprintf("%s [%s] %s\n", cursor, checked, "Soft"))

	cursor = " "
	checked = " "
	if m.checkBoxes.hard.Active {
		cursor = styles.FocusedStyleFB.Render(">")
	}
	if m.checkBoxes.hard.Value {
		checked = "x"
	}
	b.WriteString(fmt.Sprintf("%s [%s] %s\n", cursor, checked, "Hard"))

	b.WriteRune('\n')
	if m.buttons.buttonCreate.Active {
		b.WriteString(fmt.Sprintf("[ %s ]", styles.FocusedStyleFB.Render("edit/created")))
	} else {
		b.WriteString(fmt.Sprintf("[ %s ]", styles.BlurredStyleFB.Render("edit/created")))
	}
	b.WriteString(" ")

	if m.buttons.buttonClean.Active {
		b.WriteString(fmt.Sprintf("[ %s ]", styles.FocusedStyleFB.Render("cleaning")))
	} else {
		b.WriteString(fmt.Sprintf("[ %s ]", styles.BlurredStyleFB.Render("cleaning")))
	}
	b.WriteString(" ")

	if m.buttons.buttonMove.Active {
		b.WriteString(fmt.Sprintf("[ %s ]", styles.FocusedStyleFB.Render("move")))
	} else {
		b.WriteString(fmt.Sprintf("[ %s ]", styles.BlurredStyleFB.Render("move")))
	}
	b.WriteString(" ")

	if m.buttons.buttonDel.Active {
		b.WriteString(fmt.Sprintf("[ %s ]", styles.FocusedStyleFB.Render("delete")))
	} else {
		b.WriteString(fmt.Sprintf("[ %s ]", styles.BlurredStyleFB.Render("delete")))
	}
	b.WriteString("\n")

	s := ""
	if m.spinnering {
		s = m.spinner.View()
	}
	b.WriteString(fmt.Sprintf("\n\n %s ", s))

	statusVal := styles.StatusText.Copy().Width(styles.Width).Render(m.message)
	b.WriteString(statusVal)

	basement := fmt.Sprintf("\n\n Press %s too exit to the previous page | Press %s to quit\n",
		styles.GreenFg("[ ESC ]"), styles.GreenFg("[ F10 ]"))
	b.WriteString(basement)

	return b.String()
}

func (m *FormQueue) lenForm() int {

	return itemsInputsFormQueue + itemsAreasFormQueue + itemsButtonFormQueue +
		itemsCheckBoxFormQueue + itemsTableFormQueue
}

func (m *FormQueue) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, itemsInputsFormQueue+itemsAreasFormQueue)

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

	m.inputs.name, _ = m.inputs.name.Update(msg)
	//m.inputs.timeLive, _ = m.inputs.timeLive.Update(msg)
	//m.areas.key, _ = m.areas.key.Update(msg)
	//m.areas.public, _ = m.areas.public.Update(msg)
	//m.areas.token, _ = m.areas.token.Update(msg)

	return tea.Batch(cmds...)
}

func (m *FormQueue) createQueue() {

	url := environment.ServerUrlApi()

	addressPost := fmt.Sprintf("%s/addqueue/%s", url, m.inputs.name.Value())
	req, err := http.NewRequest("GET", addressPost, bytes.NewReader(nil))
	if err != nil {
		m.message = err.Error()
		return
	}

	req.Header.Set("Authentication-Key", "NativeСlient")
	defer req.Body.Close()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		m.message = err.Error()
		return
	}
	defer resp.Body.Close()

	bytBody, _ := io.ReadAll(resp.Body)
	m.message = fmt.Sprintf("%s; %s", resp.Status, string(bytBody))
}

func (m *FormQueue) deleteQueue() {

	url := environment.ServerUrlApi()

	addressPost := fmt.Sprintf("%s/delqueue/%s", url, m.inputs.name.Value())
	req, err := http.NewRequest("GET", addressPost, bytes.NewReader(nil))
	if err != nil {
		m.message = err.Error()
		return
	}

	req.Header.Set("Authentication-Key", "NativeСlient")
	defer req.Body.Close()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		m.message = err.Error()
		return
	}
	defer resp.Body.Close()

	_ = files.DelFolder(m.inputs.name.Value())
	m.message = resp.Status
}

func (m *FormQueue) moveMessages() {

	m.message = "Move OK!"
}

func (m *FormQueue) clearMessages() {

	url := environment.ServerUrlApi()

	addressPost := fmt.Sprintf("%s/clearqueue/%s", url, m.inputs.name.Value())
	req, err := http.NewRequest("GET", addressPost, bytes.NewReader(nil))
	if err != nil {
		m.message = err.Error()
		return
	}

	req.Header.Set("Authentication-Key", "NativeСlient")
	defer req.Body.Close()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		m.message = err.Error()
		return
	}
	defer resp.Body.Close()

	m.message = resp.Status
}
