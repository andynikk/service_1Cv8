package cli

import (
	"Service_1Cv8/internal/cli/charm"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"Service_1Cv8/internal/cli/styles"
	"Service_1Cv8/internal/constants"
	"Service_1Cv8/internal/mssql"
	"Service_1Cv8/internal/repository"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	itemsInputsFormSQLDB = 0
	itemsAreasFormSQLDB  = 0
	itemsButtonFormSQLDB = 4
	itemsTableFormSQLDB  = 1
)

type buttonsFormSQLDB struct {
	selectAll   charm.BoolModel
	unSelectAll charm.BoolModel
	shrink      charm.BoolModel
	restore     charm.BoolModel
	config      charm.BoolModel
}

type FormSQLDB struct {
	focusIndex int

	buttons buttonsFormSQLDB

	rows    []table.Row
	itemsDB []mssql.DB
	table   table.Model

	uid string
	new bool

	message string

	spinner    spinner.Model
	spinnering bool

	srv repository.SQLServer

	charm.CrawlElements
}

func (m *FormSQLDB) SetParameters(args []interface{}) {
	m.focusIndex = 0

	m.buttons.selectAll.Active = true
	m.buttons.unSelectAll.Active = false
	m.buttons.shrink.Active = false
	m.buttons.config.Active = false

	for _, v := range args {
		switch v.(type) {
		case repository.SQLServer:
			m.srv = v.(repository.SQLServer)
		}
	}

	c := mssql.ConnectSQLSetting{
		Server:   m.srv.Name,
		User:     m.srv.User,
		Password: m.srv.Password,
		Database: "master",
	}

	dbs, err := c.GetDatabasesOnServer()
	if err != nil {
		return
	}

	m.rows = []table.Row{}

	m.itemsDB = dbs
	sort.Slice(m.itemsDB, func(i, j int) bool {
		return m.itemsDB[i].Name < m.itemsDB[j].Name
	})

	for _, v := range m.itemsDB {
		f := "[ ]"
		if v.Mark != "" {
			f = "[X]"
		}

		rowT := table.Row{
			fmt.Sprintf("%d", v.ID),
			f,
			v.Name,
			v.State,
			v.RecoveryModel,
		}
		m.rows = append(m.rows, rowT)
	}

	m.Init()
	m.FillFormServerElements()
}

func NewFormSQLDB() *FormSQLDB {

	newSpinner := spinner.New()
	newSpinner.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))
	newSpinner.Spinner = spinner.Pulse

	m := FormSQLDB{
		spinnering: false,
		spinner:    newSpinner,
		rows:       []table.Row{},
		itemsDB:    []mssql.DB{},
	}

	m.createTable()
	return &m

}

func (m *FormSQLDB) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m *FormSQLDB) FillFormServerElements() {
	ce := make(charm.CrawlElements)

	ce[0] = &m.buttons.selectAll
	ce[1] = &m.buttons.unSelectAll
	ce[2] = &m.table
	ce[3] = &m.buttons.shrink
	ce[4] = &m.buttons.config

	m.CrawlElements = ce
}

func (m *FormSQLDB) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			if m.focusIndex == 2 &&
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
			model, cmd := m.executeFormSQLDB()
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

func (m *FormSQLDB) View() string {
	var b strings.Builder

	b.WriteRune('\n')

	seporator := styles.CursorModeHelpStyle.Render(" - ")
	title := fmt.Sprintf("%s Server: %s %s %s %s ", styles.ShortLine, m.srv.Description, seporator, m.srv.Name,
		seporator)

	b.WriteString(fmt.Sprintf(" %s %s\n", title, styles.ShortLine))
	b.WriteRune('\n')

	if m.buttons.selectAll.Active {
		b.WriteString(fmt.Sprintf("[ %s ]", styles.FocusedStyleFB.Render("Select all")))
	} else {
		b.WriteString(fmt.Sprintf("[ %s ]", styles.BlurredStyleFB.Render("Select all")))
	}
	b.WriteString(" ")
	if m.buttons.unSelectAll.Active {
		b.WriteString(fmt.Sprintf("[ %s ]", styles.FocusedStyleFB.Render("Unselect all")))
	} else {
		b.WriteString(fmt.Sprintf("[ %s ]", styles.BlurredStyleFB.Render("Unselect all")))
	}
	b.WriteString("\n")

	m.table.SetRows(m.rows)
	m.table.SetHeight(styles.Min(15, len(m.rows)))
	b.WriteString(styles.BaseStyle.Render(m.table.View()))
	m.table.Focus()

	b.WriteString("\n")

	if m.buttons.shrink.Active {
		b.WriteString(fmt.Sprintf("[ %s ]", styles.FocusedStyleFB.Render("Shrink DB")))
	} else {
		b.WriteString(fmt.Sprintf("[ %s ]", styles.BlurredStyleFB.Render("Shrink DB")))
	}

	if m.buttons.config.Active {
		b.WriteString(fmt.Sprintf("[ %s ]", styles.FocusedStyleFB.Render("Save config DB")))
	} else {
		b.WriteString(fmt.Sprintf("[ %s ]", styles.BlurredStyleFB.Render("Save config DB")))
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

func (m *FormSQLDB) lenForm() int {
	return itemsInputsFormSQLDB + itemsAreasFormSQLDB + itemsButtonFormSQLDB + itemsTableFormSQLDB
}

func (m *FormSQLDB) createTable() {
	columns := []table.Column{
		{Title: "id", Width: 4},
		{Title: "Mark", Width: 4},
		{Title: "Name", Width: 25},
		{Title: "Recovery model", Width: 15},
		{Title: "State", Width: 10},
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

func (m *FormSQLDB) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, itemsInputsFormServerDB)

	return tea.Batch(cmds...)
}

func (m *FormSQLDB) executeFormSQLDB() (tea.Model, tea.Cmd) {

	if m.focusIndex == 0 {
		for k, _ := range m.itemsDB {
			m.itemsDB[k].Mark = "[X]"
			m.rows[k][1] = "[X]"
		}

		cmds := make([]tea.Cmd, 0)
		return m, tea.Batch(cmds...)
	}

	if m.focusIndex == 1 {
		for k, _ := range m.itemsDB {
			m.itemsDB[k].Mark = ""
			m.rows[k][1] = "[ ]"
		}

		cmds := make([]tea.Cmd, 0)
		return m, tea.Batch(cmds...)
	}

	if m.focusIndex == 2 {
		selectedRow := m.table.Cursor()
		if m.rows[selectedRow][1] == "[X]" {
			m.rows[selectedRow][1] = "[ ]"
			m.itemsDB[selectedRow].Mark = ""

		} else {
			m.rows[selectedRow][1] = "[X]"
			m.itemsDB[selectedRow].Mark = "[X]"
		}

		cmds := make([]tea.Cmd, 0)
		return m, tea.Batch(cmds...)
	}

	if m.focusIndex == 3 {

		c := mssql.ConnectSQLSetting{
			Server:   m.srv.Name,
			User:     m.srv.User,
			Password: m.srv.Password,
			Database: "",
		}

		chanRes := make(chan string)
		go c.ShrinkDatabases(m.itemsDB, chanRes)
		go m.reviewFormSQLDB(chanRes)

		cmds := make([]tea.Cmd, 0)
		return m, tea.Batch(cmds...)
	}

	if m.focusIndex == 4 {
		cmds := make([]tea.Cmd, 0)

		c := mssql.ConnectSQLSetting{
			Server:   m.srv.Name,
			User:     m.srv.User,
			Password: m.srv.Password,
			Database: "uh_3_2",
		}

		arrConfig, err := c.GetConfig1C()
		if err != nil {
			m.message = err.Error()
			return m, tea.Batch(cmds...)
		}

		for k, v := range arrConfig {
			_ = k
			_ = v

			err = os.WriteFile(fmt.Sprintf("./xml/%s.xml", v.FileName), v.BinaryData, os.FileMode(0644))
			if err != nil {
				m.message = err.Error()
				return m, tea.Batch(cmds...)
			}
		}

		return m, tea.Batch(cmds...)
	}

	return m, nil
}

func (m *FormSQLDB) reviewFormSQLDB(chanRes chan string) {
	m.spinnering = true
	for {
		res, ok := <-chanRes
		if !ok {

			break
		}

		m.message = res
	}

	m.message = "OK shrink"
	m.spinnering = false
}
