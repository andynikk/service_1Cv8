package forms_exchange

import (
	"Service_1Cv8/internal/cli/charm"
	"Service_1Cv8/internal/exchange"
	"context"
	"fmt"
	"golang.org/x/text/encoding/charmap"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"Service_1Cv8/internal/cli/styles"
	"Service_1Cv8/internal/constants"
	"Service_1Cv8/internal/winsys"
)

const (
	itemsInputsFormListQueues = 0
	itemsAreasFormListQueues  = 0
	itemsButtonFormListQueues = 0
	itemsTableFormListQueues  = 1
)

type buttonFormListQueues struct {
	buttonTransferToDisk charm.BoolModel
}

type FormListQueues struct {
	focusIndex int

	rows    []table.Row
	itemsDB []exchange.ExchangeQueueInfo
	table   table.Model
	buttons buttonFormListQueues

	message     string
	errorConnet bool

	spinner    spinner.Model
	spinnering bool

	chanOut chan bool

	charm.CrawlElements
}

func (m *FormListQueues) SetParameters(args []interface{}) {

	m.rows = []table.Row{}

	m.chanOut = make(chan bool)
	ctx := context.Background()

	go m.wsQueuesInfo(ctx)

	m.FillFormTgMsgElements()

}

func (m *FormListQueues) FillFormTgMsgElements() {
	ce := make(charm.CrawlElements)

	i := -1
	increaseI := func() int { i++; return i }

	//ce[increaseI()] = &m.buttons.buttonTransferToDisk
	ce[increaseI()] = &m.table

	m.CrawlElements = ce

	m.Init()
}

func NewFormListQueues() *FormListQueues {

	newSpinner := spinner.New()
	newSpinner.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))
	newSpinner.Spinner = spinner.Pulse

	m := FormListQueues{
		spinnering: false,
		spinner:    newSpinner,
		rows:       []table.Row{},
	}

	m.createTable()
	return &m

}

func (m *FormListQueues) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m *FormListQueues) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		strMsg := msg.String()
		switch strMsg {
		case "esc", "ctrl+q":

			if !m.errorConnet {
				m.chanOut <- false
			}

			frm := exchanger.Models[constants.FormExchangeBasic].(*FormBasic)
			frm.SetParameters(nil)

			return frm, nil

		case "ctrl+", "f10":

			return m, tea.Quit

		case "tab", "shift+tab", "up", "down":
			s := msg.String()

			if m.table.Focused() &&
				(s == "up" || s == "down") {

				m.table, cmd = m.table.Update(msg)
				return m, cmd
			}

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

			cmds := make([]tea.Cmd, 0)
			return m, tea.Batch(cmds...)
		case "ctrl+n", "insert":
			frm := exchanger.Models[constants.FormExchangeQueues].(*FormQueue)
			frm.SetParameters(nil)

			return frm.Update(nil)

		case "enter", " ":
			if m.table.Focused() {
				selectedRow := m.table.Cursor()
				eq := exchange.ExchangeQueueInfo{
					Name: m.rows[selectedRow][1],
				}

				frm := exchanger.Models[constants.FormExchangeQueues].(*FormQueue)

				arg := []interface{}{eq}
				frm.SetParameters(arg)

				return frm.Update(nil)
			}
			if m.buttons.buttonTransferToDisk.Active {

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

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m *FormListQueues) View() string {
	var b strings.Builder

	b.WriteRune('\n')
	title := fmt.Sprintf("%s - Tokens - %s", styles.Line, styles.Line)

	b.WriteString(fmt.Sprintf(" %s %s\n", title, styles.ShortLine))
	b.WriteRune('\n')

	m.table.SetRows(m.rows)
	m.table.SetHeight(styles.Min(20, len(m.rows)))
	b.WriteString(styles.BaseStyle.Render(m.table.View()))
	m.table.Focus()

	b.WriteString("\n")

	s := ""
	if m.spinnering {
		s = m.spinner.View()
	}
	b.WriteString(fmt.Sprintf("\n\n %s ", s))

	statusVal := styles.StatusText.Copy().Width(styles.Width).Render(m.message)
	b.WriteString(statusVal)

	podval := fmt.Sprintf("\n\n Press %s to exit main menu | Press %s to quit\n",
		styles.GreenFg("[ ESC ]"), styles.GreenFg("[ F10 ]"))

	b.WriteString(podval)

	return b.String()

}

func (m *FormListQueues) lenForm() int {
	return itemsInputsFormListQueues + itemsAreasFormListQueues + itemsButtonFormListQueues + itemsTableFormListQueues
}

func (m *FormListQueues) createTable() {
	columns := []table.Column{
		{Title: "n/o", Width: 4},
		{Title: "Queues", Width: 30},
		{Title: "Type storage", Width: 12},
		{Title: "DownloadAt", Width: 20},
		{Title: "UploadAt", Width: 20},
		{Title: "Count", Width: 12},
		{Title: "Size", Width: 20},
	}

	rows := []table.Row{}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(15),
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

func (m *FormListQueues) createdToken(server, port, db string) string {
	//for _, v := range exchange.Storage.DatabaseTokens {
	//	if v.Server == server && v.NameOnServer == db && v.Port == port {
	//		return v.UID
	//	}
	//}

	return ""
}

func (m *FormListQueues) addToken(selectedRow int) error {

	return nil
}

func (m *FormListQueues) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, 2)

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

	return tea.Batch(cmds...)
}
