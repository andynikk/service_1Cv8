package cli

import (
	"Service_1Cv8/internal/winsys"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/text/encoding/charmap"
	"strings"

	"Service_1Cv8/internal/cli/styles"
	"Service_1Cv8/internal/constants"
	"Service_1Cv8/internal/repository"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type FormAddDB struct {
	focusIndex   int
	inputs       []textinput.Model
	areas        []textarea.Model
	buttonsFocus []bool

	cursorMode textinput.CursorMode

	spinner    spinner.Model
	spinnering bool

	uid string
	new bool

	message string

	db         repository.DataBases
	propertyDB repository.PropertyDB
}

func NewFormAddDB() *FormAddDB {

	m := FormAddDB{
		inputs:       make([]textinput.Model, 9),
		areas:        make([]textarea.Model, 1),
		buttonsFocus: make([]bool, 1),
	}

	var t textinput.Model
	var a textarea.Model

	for i := range m.inputs {
		t = textinput.New()
		m.inputs[i] = t
	}

	for i := range m.areas {
		a = textarea.New()
		a.CharLimit = 1000

		m.areas[i] = a
	}

	return &m

}

func (m *FormAddDB) SetParameters(args []interface{}) {

	m.db = repository.DataBases{}
	m.propertyDB = repository.PropertyDB{}
	m.uid = ""
	m.message = ""
	m.focusIndex = 0
	m.buttonsFocus[0] = false

	m.inputs[0].Placeholder = "Alias"
	m.inputs[0].SetValue("")
	m.inputs[0].Focus()
	m.inputs[0].PromptStyle = styles.FocusedStyleFB
	m.inputs[0].TextStyle = styles.FocusedStyleFB

	m.inputs[1].Placeholder = "Server"
	m.inputs[1].SetValue("")

	m.inputs[2].Placeholder = "Port"
	m.inputs[2].SetValue("")

	m.inputs[3].Placeholder = "Name on server"
	m.inputs[3].SetValue("")

	m.inputs[4].Placeholder = "Start blocking"
	m.inputs[4].SetValue("")
	if m.new {
		m.inputs[4].SetValue(client.Storage.Settings.StartBlock)
	}
	m.inputs[4].CharLimit = 19

	m.inputs[5].Placeholder = "Finish blocking"
	m.inputs[5].SetValue("")
	if m.new {
		m.inputs[5].SetValue(client.Storage.Settings.FinishBlock)
	}
	m.inputs[5].CharLimit = 19

	m.inputs[6].Placeholder = "Unlock key"
	m.inputs[6].SetValue("")
	if m.new {
		m.inputs[6].SetValue(client.Storage.Settings.KeyUnlock)
	}

	m.inputs[7].Placeholder = "Name user"
	m.inputs[7].SetValue("")
	if m.new {
		m.inputs[7].SetValue(client.Storage.Settings.NameUser)
	}

	m.inputs[8].Placeholder = "Password user"
	m.inputs[8].SetValue("")
	if m.new {
		m.inputs[8].SetValue(client.Storage.Settings.PasswordUser)
	}
	m.inputs[8].EchoMode = textinput.EchoPassword
	m.inputs[8].EchoCharacter = '•'

	m.areas[0].SetWidth(150)
	m.areas[0].Placeholder = "massage"

	for _, v := range args {
		switch v.(type) {
		case repository.DataBases:
			m.db = v.(repository.DataBases)

			if m.new {
				continue
			}

			m.uid = m.db.UID
			m.inputs[0].SetValue(m.db.Name)
			m.inputs[1].SetValue(m.db.Server)
			m.inputs[2].SetValue(m.db.Port)
			m.inputs[3].SetValue(m.db.NameOnServer)

		case repository.PropertyDB:
			m.propertyDB = v.(repository.PropertyDB)

			if m.new {
				continue
			}

			m.uid = m.db.UID
			m.inputs[4].SetValue(m.propertyDB.StartBlock)
			m.inputs[5].SetValue(m.propertyDB.FinishBlock)
			m.inputs[6].SetValue(m.propertyDB.KeyUnlock)
			m.inputs[7].SetValue(m.propertyDB.NameUser)
			m.inputs[8].SetValue(m.propertyDB.PasswordUser)

			m.areas[0].SetValue(m.propertyDB.Massage)
		}
	}

	if m.propertyDB.Massage == "" {
		m.areas[0].SetValue(constants.MASSAGE)
	}

	if m.uid == "" {
		m.uid = uuid.New().String()
	}
}

func (m *FormAddDB) Init() tea.Cmd {
	return nil
}

func (m *FormAddDB) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		strMsg := msg.String()
		switch strMsg {
		case "esc", "ctrl+q":
			m.uid = ""
			m.db = repository.DataBases{}
			m.propertyDB = repository.PropertyDB{}

			frm := client.Models[constants.FormMain].(*FormMain)
			frm.initLists(204, 31)
			frm.loaded = true

			return frm, nil
		case "ctrl+c", "f10":
			return m, tea.Quit
			// Change cursor mode
		case "ctrl+r":
			m.cursorMode++
			if m.cursorMode > textinput.CursorHide {
				m.cursorMode = textinput.CursorBlink
			}
			cmds := make([]tea.Cmd, len(m.inputs))
			for i := range m.inputs {
				cmds[i] = m.inputs[i].SetCursorMode(m.cursorMode)
			}
			return m, tea.Batch(cmds...)
		case "tab", "shift+tab", "up", "down":

			if m.focusIndex == len(m.inputs) &&
				(strMsg == "up" || strMsg == "down") {

				cmd := m.updateInputs(msg)
				return m, cmd
			}

			s := msg.String()

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > m.lenForm() {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = m.lenForm()
			}

			cmds := make([]tea.Cmd, 0)
			if m.focusIndex <= len(m.inputs)-1 {
				m.areas[len(m.areas)-1].Blur()
				m.buttonsFocus[0] = false

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
			} else if m.focusIndex > len(m.inputs)-1 && m.focusIndex <= len(m.inputs)+len(m.areas)-1 {
				cmds = make([]tea.Cmd, len(m.areas))
				for i := 0; i <= len(m.areas)-1; i++ {
					if i+len(m.inputs) == m.focusIndex {
						// Set focused state
						cmds[i] = m.areas[i].Focus()

						m.inputs[len(m.inputs)-1].Blur()
						m.inputs[len(m.inputs)-1].PromptStyle = styles.NoStyleFB
						m.inputs[len(m.inputs)-1].TextStyle = styles.NoStyleFB

						continue
					}
					// Remove focused state
					m.areas[i].Blur()
				}
			} else if m.focusIndex <= len(m.inputs)+len(m.areas) && m.focusIndex >= m.lenForm() {
				m.areas[len(m.areas)-1].Blur()
				cmds = make([]tea.Cmd, len(m.buttonsFocus))
				for i := 0; i <= len(m.buttonsFocus)-1; i++ {
					m.buttonsFocus[0] = true
				}
			}

			return m, tea.Batch(cmds...)
		case "enter":
			if m.buttonsFocus[0] {
				m.editDB()
			}
		}
	}

	cmd := m.updateInputs(msg)
	return m, cmd
	//return m, nil
}

func (m *FormAddDB) View() string {
	var b strings.Builder

	nameDB := m.db.Name
	if m.new {
		nameDB = "new"
	}

	b.WriteString(fmt.Sprintf("%s Data base %s %s\n", styles.ShortLine, nameDB, styles.Line))
	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
		if i == 3 {
			b.WriteString(fmt.Sprintf("%s Property DB %s %s\n", styles.ShortLine, nameDB, styles.Line))
		}
	}

	b.WriteRune('\n')

	for i := range m.areas {
		b.WriteString(m.areas[i].View())
		if i < len(m.areas)-1 {
			b.WriteRune('\n')
		}
	}

	b.WriteRune('\n')

	s := ""

	b.WriteString(s)

	cursor := ""
	checked := "[x]"
	choice := m.db.Name
	if m.db.Name == "" || m.new {
		choice = "new"
	}

	b.WriteString(fmt.Sprintf("%s %s %s\n", cursor, checked, choice))

	b.WriteString("\n\n\n")

	if m.buttonsFocus[0] {
		b.WriteString(fmt.Sprintf("[ %s ]", styles.FocusedStyleFB.Render("add/edit")))
	} else {
		b.WriteString(fmt.Sprintf("[ %s ]", styles.BlurredStyleFB.Render("add/edit")))
	}

	if m.spinnering {
		s = m.spinner.View()
	}
	b.WriteString(fmt.Sprintf("\n %s ", s))

	statusVal := styles.StatusText.Copy().Width(styles.Width).Render(m.message)
	b.WriteString(statusVal)

	b.WriteString(styles.CursorModeHelpStyleFB.Render("\n\n" + " ESC exit main menu | F10 exit program"))

	return b.String()

}

func (m *FormAddDB) lenForm() int {
	return (len(m.inputs) + len(m.areas) + len(m.buttonsFocus)) - 1
}

func (m *FormAddDB) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs)+len(m.areas))

	switch msg.(type) {
	case tea.KeyMsg:

		_, ok := winsys.SubstitutionRune(msg.(tea.KeyMsg).Runes)
		//ok := false
		//for k, v := range msg.(tea.KeyMsg).Runes {
		//	b, ok := charmap.CodePage866.EncodeRune(v)
		//	if ok {
		//		msg.(tea.KeyMsg).Runes[k] = winsys.Convert_CP866_To_unicode(b)
		//	}
		//	m.message = string(b) + " " + string(msg.(tea.KeyMsg).Runes[k])
		//}
		if !ok {
			for k, v := range msg.(tea.KeyMsg).Runes {

				dec := charmap.Windows1251.DecodeByte(byte(v))
				//m.message = fmt.Sprintf("%d - %d", v, dec)
				msg.(tea.KeyMsg).Runes[k] = dec
			}
		}
	}

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	for i := range m.areas {
		m.areas[i], cmds[i] = m.areas[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m *FormAddDB) editDB() {

	db, keyDB := repository.GetDB(client.Storage.DB, m.uid)
	if keyDB == -1 {
		db = repository.DataBases{}
	}
	db.Name = m.inputs[0].Value()
	db.Server = m.inputs[1].Value()
	db.Port = m.inputs[2].Value()
	db.NameOnServer = m.inputs[3].Value()
	db.UID = m.uid

	propertyDB, keyPDB := repository.GetPropertiesDB(client.Storage.PropertyDB, m.uid)
	if keyPDB == -1 {
		propertyDB = repository.PropertyDB{}
	}
	propertyDB.StartBlock = m.inputs[4].Value()
	propertyDB.FinishBlock = m.inputs[5].Value()
	propertyDB.KeyUnlock = m.inputs[6].Value()
	propertyDB.NameUser = m.inputs[7].Value()
	propertyDB.PasswordUser = m.inputs[8].Value()
	propertyDB.Massage = m.areas[0].Value()
	propertyDB.UID = m.uid

	dataDBJSON := &client.Storage.DataDBJSON
	if keyDB != -1 {
		dataDBJSON.DB[keyDB] = dataDBJSON.DB[len(dataDBJSON.DB)-1]
		dataDBJSON.DB[len(dataDBJSON.DB)-1] = repository.DataBases{}
		dataDBJSON.DB = dataDBJSON.DB[:len(dataDBJSON.DB)-1]
	}
	dataDBJSON.DB = append(dataDBJSON.DB, db)

	if keyPDB != -1 {
		dataDBJSON.PropertyDB[keyPDB] = dataDBJSON.PropertyDB[len(dataDBJSON.PropertyDB)-1]
		dataDBJSON.PropertyDB[len(dataDBJSON.PropertyDB)-1] = repository.PropertyDB{}
		dataDBJSON.PropertyDB = dataDBJSON.PropertyDB[:len(dataDBJSON.PropertyDB)-1]
	}
	dataDBJSON.PropertyDB = append(dataDBJSON.PropertyDB, propertyDB)

	err := client.Storage.SetPudgelData()
	if err != nil {
		m.message = err.Error()
		return
	}

	m.message = "edit OK"
}
