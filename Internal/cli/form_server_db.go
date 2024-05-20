package cli

import (
	"Service_1Cv8/internal/winsys"
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	"golang.org/x/text/encoding/charmap"
	"sort"
	"strings"

	OneCv8 "Service_1Cv8/internal/1cv8"
	"Service_1Cv8/internal/cli/styles"
	"Service_1Cv8/internal/constants"
	"Service_1Cv8/internal/repository"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
)

const (
	itemsInputsFormServerDB = 0
	itemsAreasFormServerDB  = 0
	itemsButtonFormServerDB = 0
	itemsTableFormServerDB  = 1
)

const (
	rowId int = iota
	rowFavourite
	rowControl
	rowPort
	rowDB
	rowDesc
	rowUser
	rowPassword
)

type inputsFormServerDB struct {
	nameUser     textinput.Model
	passwordUser textinput.Model
}

type FormServerDB struct {
	focusIndex       int
	focusIndexDialog int
	buttonsFocus     []bool

	inputs inputsFormServerDB

	rows    []table.Row
	itemsDB []OneCv8.ItemDB
	table   table.Model

	uid string
	new bool

	message string

	spinner    spinner.Model
	spinnering bool
	dialog     bool

	srv repository.Services
}

func (m *FormServerDB) SetParameters(args []interface{}) {

	for _, v := range args {
		switch v.(type) {
		case repository.Services:
			m.srv = v.(repository.Services)
		}
	}

	massageJSON := OneCv8.MassageJSON{
		NameServer:   m.srv.NameServer,
		NameUser:     m.srv.User,
		PasswordUser: m.srv.Password,
	}

	m.rows = []table.Row{}

	m.inputs.nameUser.Placeholder = "User"
	m.inputs.nameUser.SetValue(client.Storage.Settings.NameUser)

	m.inputs.passwordUser.Placeholder = "Password user"
	m.inputs.passwordUser.SetValue(client.Storage.Settings.PasswordUser)
	m.inputs.passwordUser.EchoMode = textinput.EchoPassword
	m.inputs.passwordUser.EchoCharacter = '•'

	pp := 0
	m.itemsDB, _ = OneCv8.ListDB(massageJSON)
	sort.Slice(m.itemsDB, func(i, j int) bool {
		return m.itemsDB[i].Name < m.itemsDB[j].Name
	})

	for k, v := range m.itemsDB {
		v.UID = m.presentInFavorites(m.srv.NameServer, v.MainPort, v.Name)
		f := "[ ]"
		if v.UID != "" {
			f = "[X]"
		}

		id := m.presentInControlDoubleCon(m.srv.NameServer, v.Name)
		c := "[ ]"
		u := ""
		p := ""
		if id != -1 {
			c = "[X]"
			u = client.Storage.BasesDoubleControl[id].User
			pwd := client.Storage.BasesDoubleControl[id].Password
			if pwd != "" {
				p = "•"
			}

			m.itemsDB[k].Password = pwd
			m.itemsDB[k].User = u
		}

		pp++

		rowT := table.Row{
			fmt.Sprintf("%d", pp),
			f,
			c,
			v.MainPort,
			v.Name,
			v.Descr,
			u,
			p,
		}
		m.rows = append(m.rows, rowT)
	}

}

func NewFormServerDB() *FormServerDB {

	newSpinner := spinner.New()
	newSpinner.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))
	newSpinner.Spinner = spinner.Pulse

	m := FormServerDB{
		buttonsFocus: make([]bool, 1),
		spinnering:   false,
		spinner:      newSpinner,
		rows:         []table.Row{},
		itemsDB:      []OneCv8.ItemDB{},
	}

	var t textinput.Model

	t = textinput.New()
	m.inputs.nameUser = t

	t = textinput.New()
	m.inputs.passwordUser = t

	m.createTable()
	return &m

}

func (m *FormServerDB) Init() tea.Cmd {
	return nil
}

func (m *FormServerDB) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		strMsg := msg.String()
		switch strMsg {
		case "esc", "ctrl+q":
			if m.dialog {
				m.dialog = false

				return m, nil
			}
			m.uid = ""

			frm := client.Models[constants.FormMain].(*FormMain)
			frm.initLists(204, 31)

			return frm, nil
		case "ctrl+c", "f10":

			return m, tea.Quit

		case "tab", "shift+tab", "up", "down":
			s := msg.String()

			if m.dialog {
				if s == "up" || s == "shift+tab" {
					m.focusIndexDialog--
				} else {
					m.focusIndexDialog++
				}

				if m.focusIndexDialog > 3 {
					m.focusIndexDialog = 0
				} else if m.focusIndexDialog < 0 {
					m.focusIndexDialog = 3
				}

				switch m.focusIndexDialog {
				case 0: //user
					m.inputs.nameUser.Focus()
					m.inputs.nameUser.PromptStyle = styles.FocusedStyleFB
					m.inputs.nameUser.TextStyle = styles.FocusedStyleFB

					m.inputs.passwordUser.Blur()
					m.inputs.passwordUser.PromptStyle = styles.NoStyleFB
					m.inputs.passwordUser.TextStyle = styles.NoStyleFB
				case 1: //passwordUser
					m.inputs.nameUser.Blur()
					m.inputs.nameUser.PromptStyle = styles.NoStyleFB
					m.inputs.nameUser.TextStyle = styles.NoStyleFB

					m.inputs.passwordUser.Focus()
					m.inputs.passwordUser.PromptStyle = styles.FocusedStyleFB
					m.inputs.passwordUser.TextStyle = styles.FocusedStyleFB
				default:
					m.inputs.nameUser.Blur()
					m.inputs.nameUser.PromptStyle = styles.NoStyleFB
					m.inputs.nameUser.TextStyle = styles.NoStyleFB

					m.inputs.passwordUser.Blur()
					m.inputs.passwordUser.PromptStyle = styles.NoStyleFB
					m.inputs.passwordUser.TextStyle = styles.NoStyleFB
				}
			} else {

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
			}

			cmds := make([]tea.Cmd, 0)
			return m, tea.Batch(cmds...)

		case " ":
			selectedRow := m.table.Cursor()
			if m.dialog {
				if m.focusIndexDialog == 2 {

					m.rows[selectedRow][rowUser] = m.inputs.nameUser.Value()
					m.rows[selectedRow][rowPassword] = "•"
					m.rows[selectedRow][rowControl] = "[X]"

					m.itemsDB[selectedRow].User = m.inputs.nameUser.Value()
					m.itemsDB[selectedRow].Password = m.inputs.passwordUser.Value()

					_ = m.addInControlDoubleCon(selectedRow)
				}

				m.dialog = false
				cmds := make([]tea.Cmd, 0)
				return m, tea.Batch(cmds...)
			}

			if m.rows[selectedRow][rowControl] == "[X]" {
				m.rows[selectedRow][rowControl] = "[ ]"

				if err := m.delFromControlDoubleCon(selectedRow); err != nil {
					m.message = err.Error()
				}

			} else {
				m.rows[selectedRow][rowControl] = "[X]"
				m.message = "OK add control double!"

				if err := m.addInControlDoubleCon(selectedRow); err != nil {
					m.message = err.Error()
				}
			}

			cmds := make([]tea.Cmd, 0)
			return m, tea.Batch(cmds...)
		case "enter":

			selectedRow := m.table.Cursor()
			if m.rows[selectedRow][rowFavourite] == "[X]" {
				m.rows[selectedRow][rowFavourite] = "[ ]"

				if err := m.delFromFavorites(selectedRow); err != nil {
					m.message = err.Error()
				}

			} else {
				m.rows[selectedRow][rowFavourite] = "[X]"
				m.message = "OK add favorites!"

				if err := m.addInFavorites(selectedRow); err != nil {
					m.message = err.Error()
				}
			}

			cmds := make([]tea.Cmd, 0)
			return m, tea.Batch(cmds...)
		case "f2":
			selectedRow := m.table.Cursor()

			m.inputs.nameUser.Focus()
			m.inputs.nameUser.PromptStyle = styles.FocusedStyleFB
			m.inputs.nameUser.TextStyle = styles.FocusedStyleFB
			m.focusIndexDialog = 0

			m.dialog = !m.dialog

			if m.dialog {
				m.inputs.nameUser.SetValue(m.itemsDB[selectedRow].User)
				if m.inputs.nameUser.Value() == "" {
					m.inputs.nameUser.SetValue(client.Storage.NameUser)
				}
				m.inputs.passwordUser.SetValue(m.itemsDB[selectedRow].Password)
				if m.inputs.passwordUser.Value() == "" {
					m.inputs.passwordUser.SetValue(client.Storage.PasswordUser)
				}
			}

			m.table, cmd = m.table.Update(msg)
			return m, cmd
		}
	}

	if m.dialog {
		cmd = m.updateInputs(msg)
	}
	if !m.dialog {
		m.table, cmd = m.table.Update(msg)
	}
	return m, cmd
}

func (m *FormServerDB) View() string {
	var b strings.Builder

	b.WriteRune('\n')

	seporator := styles.CursorModeHelpStyle.Render(" - ")
	title := fmt.Sprintf("%s Server: %s %s %s %s ", styles.ShortLine, m.srv.NameServer, seporator, m.srv.IP,
		seporator)

	if m.new {
		title = "Service: NEW "
	}

	b.WriteString(fmt.Sprintf(" %s %s\n", title, styles.ShortLine))
	b.WriteRune('\n')

	// Dialog.
	subtle := lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	dialogBoxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#874BFD")).
		Padding(1, 0).
		BorderTop(true).
		BorderLeft(true).
		BorderRight(true).
		BorderBottom(true)

	buttonStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFF7DB")).
		Background(lipgloss.Color("#888B7E")).
		Padding(0, 3).
		MarginTop(1)

	activeButtonStyle := buttonStyle.Copy().
		Foreground(lipgloss.Color("#FFF7DB")).
		Background(lipgloss.Color("#F25D94")).
		MarginRight(2).
		Underline(true)

	width := 96

	//{
	okButton := buttonStyle.Render("Yes")
	cancelButton := buttonStyle.Render("Cancel")
	if m.focusIndexDialog == 2 {
		okButton = activeButtonStyle.Render("Yes")
	}
	if m.focusIndexDialog == 3 {
		cancelButton = activeButtonStyle.Render("Cancel")
	}

	inputs := fmt.Sprintf("\n%s %s\n%s %s\n",
		styles.СaptionStyleFB.Render("User:"), m.inputs.nameUser.View(),
		styles.СaptionStyleFB.Render("Password:"), m.inputs.passwordUser.View())

	selectedRow := m.table.Cursor()
	if len(m.rows) != 0 {
		question := lipgloss.NewStyle().Width(50).Align(lipgloss.Center).
			Render(fmt.Sprintf("Enter the database (%s) name and password", m.rows[selectedRow][rowDB]))
		buttons := lipgloss.JoinHorizontal(lipgloss.Top, okButton, " ", cancelButton)
		ui := lipgloss.JoinVertical(lipgloss.Center, question, inputs, buttons)
		dialog := lipgloss.Place(width, 9,
			lipgloss.Center, lipgloss.Center,
			dialogBoxStyle.Render(ui),
			lipgloss.WithWhitespaceChars("猫咪"),
			lipgloss.WithWhitespaceForeground(subtle),
		)
		m.table.SetRows(m.rows)
		m.table.SetHeight(styles.Min(20, len(m.rows)))
		if m.dialog {
			b.WriteString(dialog + "\n\n")
		} else {
			b.WriteString(styles.BaseStyle.Render(m.table.View()))
			m.table.Focus()
		}
	}

	//}

	podval := fmt.Sprintf("\n\n Press %s add/remove favourites | Press %s add/remove control double connections\n"+
		" Press %s to exit main menu | Press %s to quit\n",
		styles.GreenFg("[ ENTER ]"), styles.GreenFg("[ SPACE ]"), styles.GreenFg("[ ESC ]"), styles.GreenFg("[ F10 ]"))
	b.WriteString(podval)

	return b.String()

}

func (m *FormServerDB) lenForm() int {
	return itemsInputsFormServerDB + itemsAreasFormServerDB + itemsButtonFormServerDB + itemsTableFormServerDB
}

func (m *FormServerDB) createTable() {
	columns := []table.Column{
		{Title: "n/o", Width: 4},
		{Title: "Favourite", Width: 9},
		{Title: "Control double", Width: 9},
		{Title: "Port", Width: 4},
		{Title: "Name DB", Width: 20},
		{Title: "Description DB", Width: 25},
		{Title: "User", Width: 15},
		{Title: "Password", Width: 15},
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

func (m *FormServerDB) presentInFavorites(server, port, db string) string {
	for _, v := range client.Storage.DB {
		if v.Server == server && v.NameOnServer == db && v.Port == port {
			return v.UID
		}
	}
	return ""
}

func (m *FormServerDB) presentInControlDoubleCon(server, db string) int {
	for k, v := range client.Storage.BasesDoubleControl {
		if v.Server == server && v.Name == db {
			return k
		}
	}
	return -1
}

func (m *FormServerDB) addInFavorites(selectedRow int) error {
	uid := m.presentInFavorites(m.srv.NameServer, m.rows[selectedRow][rowPort], m.rows[selectedRow][rowDB])
	if uid == "" {
		uid = uuid.New().String()
	}

	setting := client.Storage.Settings
	db := repository.DataBases{
		Name:         m.rows[selectedRow][rowDB],
		Server:       m.srv.NameServer,
		NameOnServer: m.rows[selectedRow][rowDB],
		Port:         m.rows[selectedRow][rowPort],
		UID:          uid,
	}

	pDB, k := repository.GetPropertiesDB(client.Storage.PropertyDB, uid)
	if k == -1 {
		pDB = repository.PropertyDB{
			NameUser:     setting.NameUser,
			PasswordUser: setting.PasswordUser,
			StartBlock:   setting.StartBlock,
			FinishBlock:  setting.FinishBlock,
			KeyUnlock:    setting.KeyUnlock,
			Massage:      setting.Massage,
			UID:          uid,
		}
	}
	err := client.EditDB(uid, db, pDB)
	if err != nil {
		return err
	}

	for _, v := range m.itemsDB {
		v.UID = uid
	}
	return nil
}

func (m *FormServerDB) delFromFavorites(selectedRow int) error {

	uid := m.presentInFavorites(m.srv.NameServer, m.rows[selectedRow][rowPort], m.rows[selectedRow][rowDB])
	if uid == "" {
		return nil
	}

	_, keyDB := repository.GetDB(client.Storage.DB, uid)
	_, keyPDB := repository.GetPropertiesDB(client.Storage.PropertyDB, uid)
	dataDBJSON := &client.Storage.DataDBJSON

	if keyDB != -1 {
		dataDBJSON.DB[keyDB] = dataDBJSON.DB[len(dataDBJSON.DB)-1]
		dataDBJSON.DB[len(dataDBJSON.DB)-1] = repository.DataBases{}
		dataDBJSON.DB = dataDBJSON.DB[:len(dataDBJSON.DB)-1]
	}

	if keyPDB != -1 {
		dataDBJSON.PropertyDB[keyPDB] = dataDBJSON.PropertyDB[len(dataDBJSON.PropertyDB)-1]
		dataDBJSON.PropertyDB[len(dataDBJSON.PropertyDB)-1] = repository.PropertyDB{}
		dataDBJSON.PropertyDB = dataDBJSON.PropertyDB[:len(dataDBJSON.PropertyDB)-1]
	}

	err := client.Storage.SetPudgelData()
	if err != nil {
		return err
	}

	return nil
}

func (m *FormServerDB) addInControlDoubleCon(selectedRow int) error {
	var uid string

	m.itemsDB[selectedRow].User = m.inputs.nameUser.Value()
	m.itemsDB[selectedRow].Password = m.inputs.passwordUser.Value()

	id := m.presentInControlDoubleCon(m.srv.NameServer, m.rows[selectedRow][rowDB])
	if id == -1 {
		uid = uuid.New().String()
	} else {
		uid = client.Storage.BasesDoubleControl[id].UID
	}

	db := repository.BasesDoubleControl{
		Name:     m.itemsDB[selectedRow].Name,
		Server:   m.srv.NameServer,
		User:     m.itemsDB[selectedRow].User,
		Password: m.itemsDB[selectedRow].Password,
		UID:      uid,
	}

	err := client.EditControlDoubleConDB(uid, db)
	if err != nil {
		return err
	}

	return nil
}

func (m *FormServerDB) delFromControlDoubleCon(selectedRow int) error {

	id := m.presentInControlDoubleCon(m.srv.NameServer, m.rows[selectedRow][rowDB])
	if id == -1 {
		return nil
	}

	dataDBJSON := client.Storage

	dataDBJSON.BasesDoubleControl[id] = dataDBJSON.BasesDoubleControl[len(dataDBJSON.BasesDoubleControl)-1]
	dataDBJSON.BasesDoubleControl[len(dataDBJSON.BasesDoubleControl)-1] = repository.BasesDoubleControl{}
	dataDBJSON.BasesDoubleControl = dataDBJSON.BasesDoubleControl[:len(dataDBJSON.BasesDoubleControl)-1]

	err := client.Storage.SetPudgelData()
	if err != nil {
		return err
	}

	return nil
}

func (m *FormServerDB) updateInputs(msg tea.Msg) tea.Cmd {
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

	m.inputs.nameUser, _ = m.inputs.nameUser.Update(msg)
	m.inputs.passwordUser, _ = m.inputs.passwordUser.Update(msg)

	return tea.Batch(cmds...)
}
