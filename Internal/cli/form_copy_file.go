package cli

import (
	"Service_1Cv8/internal/files"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"Service_1Cv8/internal/cli/styles"
	"Service_1Cv8/internal/constants"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	itemsInputsFormCopyFile = 0
	itemsAreasFormCopyFile  = 0
	itemsButtonFormCopyFile = 1
	itemsTableFormCopyFile  = 1
)

type buttonsFormCopyFile struct {
	copy bool
}

type FormCopyFile struct {
	focusIndex int

	buttons buttonsFormCopyFile

	rows       []table.Row
	itemsTable []files.Files
	table      table.Model

	uid string
	new bool

	message string

	spinner    spinner.Model
	spinnering bool
}

func (m *FormCopyFile) SetParameters(args []interface{}) {
	m.focusIndex = 0
	m.buttons.copy = false
	m.table.Focus()

	itemsTable, err := files.ListFile(client.Storage.PathStorage)
	if err != nil {
		return
	}

	ch := 0
	m.itemsTable = itemsTable
	sort.Slice(m.itemsTable, func(i, j int) bool {
		return m.itemsTable[i].Name < m.itemsTable[j].Name
	})

	for _, v := range m.itemsTable {

		ch++
		rowT := table.Row{
			fmt.Sprintf("%d", ch),
			v.Name,
			fmt.Sprintf("%s", v.SizeStrings),
		}
		m.rows = append(m.rows, rowT)
	}

	m.Init()
}

func NewFormCopyFile() *FormCopyFile {

	newSpinner := spinner.New()
	newSpinner.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))
	newSpinner.Spinner = spinner.Pulse

	m := FormCopyFile{
		spinnering: false,
		spinner:    newSpinner,
		rows:       []table.Row{},
		itemsTable: []files.Files{},
	}

	m.createTable()
	return &m

}

func (m *FormCopyFile) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m *FormCopyFile) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		strMsg := msg.String()
		switch strMsg {
		case "esc", "ctrl+q":
			m.uid = ""

			frm := client.Models[constants.FormServer].(*FormServer)
			return frm, nil

		case "ctrl+c", "f10":

			return m, tea.Quit

		case "tab", "shift+tab", "up", "down":
			if m.focusIndex == 0 &&
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
			case 0:
				m.table.Focus()
				m.buttons.copy = false
			case 1:
				m.table.Blur()
				m.buttons.copy = true
			}
			return m, tea.Batch(cmds...)
		case "enter", " ":
			model, cmd := m.executeForm()
			return model, cmd
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

	return m, nil
}

func (m *FormCopyFile) View() string {
	var b strings.Builder

	b.WriteRune('\n')

	seporator := styles.CursorModeHelpStyle.Render(" - ")
	title := fmt.Sprintf("%s Server: %s %s %s %s ", styles.ShortLine, client.Storage.PathStorage, seporator,
		client.Storage.PathCopy, seporator)

	b.WriteString(fmt.Sprintf(" %s %s\n", title, styles.ShortLine))
	b.WriteRune('\n')

	m.table.SetRows(m.rows)
	m.table.SetHeight(styles.Min(15, len(m.rows)))
	b.WriteString(styles.BaseStyle.Render(m.table.View()))
	m.table.Focus()

	b.WriteString("\n")

	if m.buttons.copy {
		b.WriteString(fmt.Sprintf("[ %s ]", styles.FocusedStyleFB.Render("Copy file")))
	} else {
		b.WriteString(fmt.Sprintf("[ %s ]", styles.BlurredStyleFB.Render("Copy file")))
	}
	b.WriteString("\n")

	s := " "
	if m.spinnering {
		s = m.spinner.View()
	}
	b.WriteString(fmt.Sprintf("\n %s ", s))

	statusVal := styles.StatusText.Copy().Width(styles.Width).Render(m.message)
	b.WriteString(statusVal + "\n")

	//b.WriteString(styles.CursorModeHelpStyleFB.Render("\n\n\n" + " ESC exit form service | F10 exit program"))
	podval := fmt.Sprintf("\n\n Press %s to exit main menu | Press %s to quit\n",
		styles.GreenFg("[ ESC ]"), styles.GreenFg("[ F10 ]"))
	b.WriteString(podval)

	return b.String()

}

func (m *FormCopyFile) lenForm() int {
	return itemsInputsFormCopyFile + itemsAreasFormCopyFile + itemsButtonFormCopyFile + itemsTableFormCopyFile
}

func (m *FormCopyFile) createTable() {
	columns := []table.Column{
		{Title: "id", Width: 4},
		{Title: "Name", Width: 60},
		{Title: "Size", Width: 100},
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

func (m *FormCopyFile) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, itemsInputsFormServerDB)

	return tea.Batch(cmds...)
}

func (m *FormCopyFile) executeForm() (tea.Model, tea.Cmd) {

	if m.focusIndex == 1 {

		selectedRow := m.table.Cursor()
		pathSource := fmt.Sprintf("%s\\%s", client.Storage.PathStorage, m.rows[selectedRow][1])
		pathReceiver := fmt.Sprintf("%s\\%s", client.Storage.PathCopy, m.rows[selectedRow][1])

		newFile, err := os.Create(pathReceiver)
		if err != nil {
			m.message = err.Error()
			return m, nil
		}

		chanOut := make(chan files.VC)
		chanInfo := make(chan int64)

		file, _ := os.Open(pathSource)
		fi, _ := file.Stat()
		totalBytes := fi.Size()

		go files.ReadFile(pathSource, chanOut, chanInfo)
		go m.reviewForm(chanInfo, totalBytes)

		go files.GoWriteFile(newFile, chanOut)
		go files.GoWriteFile(newFile, chanOut)
		go files.GoWriteFile(newFile, chanOut)
		go files.GoWriteFile(newFile, chanOut)
		go files.GoWriteFile(newFile, chanOut) //4min - 16 777 344; 4 min - 8 388 67
		//go files.GoWriteFile(newFile, chanOut)
		//go files.GoWriteFile(newFile, chanOut) //4 min - 8 388 672; 5min - 16 777 344

		cmds := make([]tea.Cmd, 0)
		return m, tea.Batch(cmds...)
	}

	return m, nil
}

func (m *FormCopyFile) reviewForm(chanInfo chan int64, totalBytes int64) {

	timeStart := time.Now()
	m.spinnering = true

	totalBytesStr := files.GroupSeparator(fmt.Sprintf("%d", totalBytes))
	for {
		res, ok := <-chanInfo
		if !ok {

			break
		}

		proc := int(res * 100 / totalBytes)

		m.message = fmt.Sprintf("%s from %s (%d%%)", files.GroupSeparator(fmt.Sprintf("%d", res)),
			totalBytesStr, proc)
	}
	timeFinish := time.Now()

	m.message = fmt.Sprintf("OK copy (%d min)", int(timeFinish.Sub(timeStart)/time.Minute))

	m.spinnering = false
}
