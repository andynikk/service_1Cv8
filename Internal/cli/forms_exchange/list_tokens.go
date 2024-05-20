package forms_exchange

import (
	"Service_1Cv8/internal/iron"
	"errors"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/text/encoding/charmap"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"Service_1Cv8/internal/cli/styles"
	"Service_1Cv8/internal/constants"
	"Service_1Cv8/internal/repository"
	"Service_1Cv8/internal/token"
	"Service_1Cv8/internal/winsys"
)

const (
	itemsInputsFormListToken = 0
	itemsAreasFormListToken  = 0
	itemsButtonFormListToken = 0
	itemsTableFormListToken  = 1
)

const (
	rowLTID int = iota
	rowLTKey
	rowLTValid
)

type FormListToken struct {
	focusIndex int

	rows    []table.Row
	itemsDB []token.ClaimStore
	table   table.Model

	message string

	spinner    spinner.Model
	spinnering bool
}

func (m *FormListToken) SetParameters(args []interface{}) {

	m.rows = []table.Row{}
	m.itemsDB, _ = repository.GetTokens()

	pp := 0
	for _, v := range m.itemsDB {

		//tkn, _ := jwt.Parse(string(v.Value), func(token *jwt.Token) (interface{}, error) {
		//	return v.Secret, nil
		//})

		tkn, err := jwt.Parse(string(v.Value), func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("there was an error")
			}
			return v.Secret, nil
		})

		//tkn, err := jwt.Parse(v, func(token *jwt.Token) (interface{}, error) {
		//	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		//		return nil, errors.New("there was an error")
		//	}
		//	return nil, errors.New("token.Method not type jwt.SigningMethodHMAC")
		//})
		if err != nil && err.Error() != "Token is expired" {
			continue
		}
		valid := "[ ]"
		if tkn.Valid {
			valid = "[X]"
		}

		pp++
		rowT := table.Row{
			fmt.Sprintf("%d", pp),
			tkn.Claims.(jwt.MapClaims)["key"].(string),
			valid,
		}
		m.rows = append(m.rows, rowT)
	}

	m.Init()

}

func NewFormListToken() *FormListToken {

	newSpinner := spinner.New()
	newSpinner.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))
	newSpinner.Spinner = spinner.Pulse

	m := FormListToken{
		spinnering: false,
		spinner:    newSpinner,
		rows:       []table.Row{},
	}

	m.createTable()
	return &m

}

func (m *FormListToken) Init() tea.Cmd {
	return nil
}

func (m *FormListToken) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		strMsg := msg.String()
		switch strMsg {
		case "esc", "ctrl+q":

			frm := exchanger.Models[constants.FormExchangeBasic].(*FormBasic)
			frm.initLists(204, 31)

			return frm, nil
		case "ctrl+c", "f10":

			return m, tea.Quit

		case "tab", "shift+tab", "up", "down":
			s := msg.String()

			if s == "up" || s == "down" {
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

			cmds := make([]tea.Cmd, 0)
			return m, tea.Batch(cmds...)
		case "ctrl+g", "insert":
			frm := exchanger.Models[constants.FormExchangeKeyToken].(*FormKeyToken)
			frm.SetParameters(nil)

			return frm.Update(nil)
		case "ctrl+u":
			m.getKeyLocalMachine()

			return m, tea.Batch(make([]tea.Cmd, 0)...)
		case "enter", " ":

			selectedRow := m.table.Cursor()

			cs, _ := repository.GetToken(m.rows[selectedRow][rowLTKey])

			arg := []interface{}{cs}
			frm := exchanger.Models[constants.FormExchangeKeyToken].(*FormKeyToken)

			frm.SetParameters(arg)

			return frm.Update(nil)

			//selectedRow := m.table.Cursor()
			//if m.rows[selectedRow][rowLTToken] == "[X]" {
			//	m.rows[selectedRow][rowLTToken] = "[ ]"
			//
			//	//if err := m.delToken(selectedRow); err != nil {
			//	//	m.message = err.Error()
			//	//}
			//
			//} else {
			//	m.rows[selectedRow][rowLTToken] = "[X]"
			//	m.message = "OK add token!"
			//
			//	if err := m.addToken(selectedRow); err != nil {
			//		m.message = err.Error()
			//	}
			//}

			//cmds := make([]tea.Cmd, 0)
			//return m, tea.Batch(cmds...)
		}
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m *FormListToken) View() string {
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

func (m *FormListToken) lenForm() int {
	return itemsInputsFormListToken + itemsAreasFormListToken + itemsButtonFormListToken + itemsTableFormListToken
}

func (m *FormListToken) createTable() {
	columns := []table.Column{
		{Title: "n/o", Width: 4},
		{Title: "Key", Width: 70},
		{Title: "Valid", Width: 8},
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

func (m *FormListToken) createdToken(server, port, db string) string {
	//for _, v := range exchange.Storage.DatabaseTokens {
	//	if v.Server == server && v.NameOnServer == db && v.Port == port {
	//		return v.UID
	//	}
	//}

	return ""
}

func (m *FormListToken) addToken(selectedRow int) error {
	//uid := m.createdToken(m.srv.NameServer, "", "")
	//if uid == "" {
	//	uid = uuid.New().String()
	//}

	//db := repository.DatabaseTokens{
	//	Name:         m.rows[selectedRow][rowLTID],
	//	Server:       m.srv.NameServer,
	//	NameOnServer: m.rows[selectedRow][rowLTID],
	//	Port:         m.rows[selectedRow][rowLTID],
	//	UID:          uid,
	//}
	//
	//_ = db

	return nil
}

func (m *FormListToken) updateInputs(msg tea.Msg) tea.Cmd {
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

func (m *FormListToken) getKeyLocalMachine() {
	macAdress, err := iron.GetMacAddr()
	if err != nil {
		m.message = fmt.Sprintf("%s: %s", "Get MAC adress", err.Error())
		return
	}
	snHDD, err := iron.GetDiskDrivers()
	if err != nil {
		m.message = fmt.Sprintf("%s: %s", "Get MAC adress", err.Error())
		return
	}

	appKey := fmt.Sprintf("%s%s", macAdress, snHDD)
	err = clipboard.WriteAll(string(appKey))
	if err != nil {
		m.message = err.Error()
	}
	m.message = "OK! Write to clipboard"
}
