package cli

import (
	"Service_1Cv8/internal/cli/charm"
	"Service_1Cv8/internal/winsys"
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
	"github.com/go-ole/go-ole"
	"golang.org/x/text/encoding/charmap"
	"strings"
	"time"

	OneCv8 "Service_1Cv8/internal/1cv8"
	"Service_1Cv8/internal/cli/styles"
	"Service_1Cv8/internal/constants"
	"Service_1Cv8/internal/repository"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	itemsInputsFormBlock   = 5
	itemsAreasFormBlock    = 1
	itemsButtonFormBlock   = 0
	itemsCheckBoxFormBlock = 1
	itemsTableFormBlock    = 0
)

type inputsFormBlock struct {
	nameUser     textinput.Model
	passwordUser textinput.Model
	startBlock   textinput.Model
	finishBlock  textinput.Model
	keyUnlock    textinput.Model
}

type areasFormBlock struct {
	massage textarea.Model
}

type checkboxesFormBlock struct {
	blocked         charm.Checkbox
	blockedSelected bool
}

type FormBlock struct {
	focusIndex int
	inputs     inputsFormBlock
	areas      areasFormBlock
	checkboxes checkboxesFormBlock

	cursorMode textinput.CursorMode
	error      string

	uid string

	spinner    spinner.Model
	spinnering bool

	db         repository.DataBases
	propertyDB repository.PropertyDB
}

func (m *FormBlock) SetParameters(args []interface{}) {
	for _, v := range args {
		switch v.(type) {
		case repository.DataBases:
			m.db = v.(repository.DataBases)
			m.uid = m.db.UID
		case repository.PropertyDB:
			m.propertyDB = v.(repository.PropertyDB)

			m.uid = m.db.UID

			now := time.Now()

			m.inputs.startBlock.Placeholder = "Start blocking"
			m.inputs.startBlock.SetValue(fmt.Sprintf("%d-%d-%d %s", now.Year(), now.Month(), now.Day(), m.propertyDB.StartBlock))
			m.inputs.startBlock.Focus()
			m.inputs.startBlock.PromptStyle = styles.FocusedStyleFB
			m.inputs.startBlock.TextStyle = styles.FocusedStyleFB
			m.inputs.startBlock.CharLimit = 19

			m.inputs.finishBlock.Placeholder = "Finish blocking"
			m.inputs.finishBlock.SetValue(fmt.Sprintf("%d-%d-%d %s", now.Year(), now.Month(), now.Day(), m.propertyDB.FinishBlock))
			m.inputs.finishBlock.CharLimit = 19

			m.inputs.keyUnlock.Placeholder = "Unlock key"
			m.inputs.keyUnlock.SetValue(m.propertyDB.KeyUnlock)

			m.inputs.nameUser.Placeholder = "Name user"
			m.inputs.nameUser.SetValue(m.propertyDB.NameUser)

			m.inputs.passwordUser.Placeholder = "Password user"
			m.inputs.passwordUser.SetValue(m.propertyDB.PasswordUser)
			m.inputs.passwordUser.EchoMode = textinput.EchoPassword
			m.inputs.passwordUser.EchoCharacter = '•'

			m.areas.massage.SetWidth(250)
			m.areas.massage.Placeholder = "massage"

			ch := charm.Checkbox{
				Choices:  []string{"Block DB"},
				Selected: make(map[int]struct{}),
			}
			if m.propertyDB.Block {
				ch.Selected[0] = struct{}{}
			}
			m.checkboxes.blocked = ch

		}
	}

	m.areas.massage.SetValue(fmt.Sprintf(constants.MASSAGE, m.propertyDB.StartBlock, m.propertyDB.FinishBlock,
		m.db.Server, m.db.NameOnServer))

	m.Init()
}

func NewFormPropertyDB(db repository.DataBases, propertyDB repository.PropertyDB) *FormBlock {

	newSpinner := spinner.New()
	newSpinner.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))
	newSpinner.Spinner = spinner.Pulse

	m := FormBlock{
		spinnering: false,
		spinner:    newSpinner,
	}

	var t textinput.Model
	var a textarea.Model

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

	a = textarea.New()
	a.CharLimit = 2000
	m.areas.massage = a

	ch := charm.Checkbox{
		Choices:  []string{"Block DB"},
		Selected: make(map[int]struct{}),
	}
	if propertyDB.Block {
		ch.Selected[0] = struct{}{}
	}
	m.checkboxes.blocked = ch

	m.db = db
	m.propertyDB = propertyDB

	return &m

}

func (m *FormBlock) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m *FormBlock) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	cmds := make([]tea.Cmd, 0)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		strMsg := msg.String()
		switch strMsg {
		case "esc", "ctrl+q":
			frm := client.Models[constants.FormMain].(*FormMain)
			frm.initLists(204, 31)

			return frm, nil
		case "ctrl+c", "f10":

			return m, tea.Quit

		case "tab", "shift+tab", "up", "down":

			if m.focusIndex == 5 &&
				(strMsg == "up" || strMsg == "down") {

				if m.focusIndex == 5 {
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
			case 0: //startBlock
				m.checkboxes.blocked.Cursor = 0

				m.inputs.startBlock.Focus()
				m.inputs.startBlock.PromptStyle = styles.FocusedStyleFB
				m.inputs.startBlock.TextStyle = styles.FocusedStyleFB

				m.inputs.finishBlock.Blur()
				m.inputs.finishBlock.PromptStyle = styles.NoStyleFB
				m.inputs.finishBlock.TextStyle = styles.NoStyleFB

			case 1: //finishBlock
				m.inputs.startBlock.Blur()
				m.inputs.startBlock.PromptStyle = styles.NoStyleFB
				m.inputs.startBlock.TextStyle = styles.NoStyleFB

				m.inputs.finishBlock.Focus()
				m.inputs.finishBlock.PromptStyle = styles.FocusedStyleFB
				m.inputs.finishBlock.TextStyle = styles.FocusedStyleFB

				m.inputs.keyUnlock.Blur()
				m.inputs.keyUnlock.PromptStyle = styles.NoStyleFB
				m.inputs.keyUnlock.TextStyle = styles.NoStyleFB
			case 2: //keyUnlock
				m.inputs.finishBlock.Blur()
				m.inputs.finishBlock.PromptStyle = styles.NoStyleFB
				m.inputs.finishBlock.TextStyle = styles.NoStyleFB

				m.inputs.keyUnlock.Focus()
				m.inputs.keyUnlock.PromptStyle = styles.FocusedStyleFB
				m.inputs.keyUnlock.TextStyle = styles.FocusedStyleFB

				m.inputs.nameUser.Blur()
				m.inputs.nameUser.PromptStyle = styles.NoStyleFB
				m.inputs.nameUser.TextStyle = styles.NoStyleFB

			case 3: //nameUser
				m.inputs.keyUnlock.Blur()
				m.inputs.keyUnlock.PromptStyle = styles.NoStyleFB
				m.inputs.keyUnlock.TextStyle = styles.NoStyleFB

				m.inputs.nameUser.Focus()
				m.inputs.nameUser.PromptStyle = styles.FocusedStyleFB
				m.inputs.nameUser.TextStyle = styles.FocusedStyleFB

				m.inputs.passwordUser.Blur()
				m.inputs.passwordUser.PromptStyle = styles.NoStyleFB
				m.inputs.passwordUser.TextStyle = styles.NoStyleFB

			case 4: //passwordUser
				m.inputs.nameUser.Blur()
				m.inputs.nameUser.PromptStyle = styles.NoStyleFB
				m.inputs.nameUser.TextStyle = styles.NoStyleFB

				m.inputs.passwordUser.Focus()
				m.inputs.passwordUser.PromptStyle = styles.FocusedStyleFB
				m.inputs.passwordUser.TextStyle = styles.FocusedStyleFB

				m.areas.massage.Blur()

			case 5: //user
				m.inputs.passwordUser.Blur()
				m.inputs.passwordUser.PromptStyle = styles.NoStyleFB
				m.inputs.passwordUser.TextStyle = styles.NoStyleFB

				m.areas.massage.Focus()

				m.checkboxes.blocked.Cursor = 0
			case 6: //user
				m.areas.massage.Blur()

				m.checkboxes.blocked.Cursor = 1

				m.inputs.startBlock.Blur()
				m.inputs.startBlock.PromptStyle = styles.NoStyleFB
				m.inputs.startBlock.TextStyle = styles.NoStyleFB
			}
			return m, tea.Batch(cmds...)
		case "ctrl+b":
			go m.executeBlocDB()
		case "ctrl+o":
			model, cmd := m.executeOpenDesignerDB()
			if model != nil {
				return model, cmd
			}
		case "ctrl+d":
			go m.executeDropConnectsDB()
		case "ctrl+u":
			go m.executeDropDoubleConnectsDB()
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

	cmd = m.updateInputs(msg)
	return m, cmd
	//return m, nil
}

func (m *FormBlock) View() string {

	var b strings.Builder

	b.WriteRune('\n')
	b.WriteString(fmt.Sprintf("%s DB %s. Patch Srvr=\"%s%s\";Ref=\"%s\"; %s",
		styles.ShortLine, m.db.Name, m.db.Server, m.db.Port, m.db.NameOnServer, styles.ShortLine))
	b.WriteRune('\n')
	b.WriteRune('\n')

	b.WriteString(fmt.Sprintf("%s %s %s %s\n",
		styles.СaptionStyleFB.Render("Time start block:"), m.inputs.startBlock.View(),
		styles.СaptionStyleFB.Render("Time finish block:"), m.inputs.finishBlock.View()))

	b.WriteString(fmt.Sprintf("%s %s\n", styles.СaptionStyleFB.Render("Key unlock:"), m.inputs.keyUnlock.View()))

	b.WriteString(fmt.Sprintf("%s %s %s %s\n",
		styles.СaptionStyleFB.Render("User:"), m.inputs.nameUser.View(),
		styles.СaptionStyleFB.Render("Password:"), m.inputs.passwordUser.View()))

	b.WriteString(fmt.Sprintf("%s\n%s\n", styles.СaptionStyleFB.Render("Massage:"), m.areas.massage.View()))

	cursor := ""
	choice := "Block DB"
	checked := "[ ]"
	if m.propertyDB.Block {
		checked = "[x]"
	}

	if m.checkboxes.blocked.Cursor == 1 {
		cursor = ">"

		checked = styles.SelectColorFg(checked)
		cursor = styles.SelectColorFg(cursor)
		choice = styles.SelectColorFg(choice)
	}

	b.WriteString(fmt.Sprintf("%s %s %s\n", cursor, checked, choice))

	s := " "
	if m.spinnering {
		s = m.spinner.View()
	}
	s = fmt.Sprintf(" %s ", s)

	statusKey := styles.StatusStyleBlock.Render("[ctrl+b] Block DB")
	if m.propertyDB.Block {
		statusKey = styles.StatusStyle.Render("[ctrl+b] Unblock DB")
	}
	encoding := styles.EncodingStyle.Render("[ctrl+o] Open designer DB")
	fishCake := styles.FishCakeStyle.Render("[ctrl+d] ✘ Drop connects DB")

	fishDoubleCake := styles.FishCakeStyle.Render("[ctrl+u] Drop double connects DB")

	bar := lipgloss.JoinHorizontal(lipgloss.Top,
		statusKey,
		encoding,
		fishCake,
		fishDoubleCake,
	)
	b.WriteString("\n\n" + styles.StatusBarStyle.Render(bar) + "\n\n")

	statusVal := styles.StatusText.Copy().
		Width(styles.Width).Render(m.error)
	b.WriteString(s + statusVal)

	b.WriteString(styles.CursorModeHelpStyleFB.Render("\n\n\n" + " ESC exit main menu | F10 exit program"))

	return b.String()

}

func (m *FormBlock) lenForm() int {
	return itemsInputsFormBlock + itemsAreasFormBlock + itemsButtonFormBlock + itemsCheckBoxFormBlock +
		itemsTableFormBlock
}

func (m *FormBlock) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, itemsInputsFormBlock+itemsAreasFormBlock)

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

	m.areas.massage, _ = m.areas.massage.Update(msg)

	return tea.Batch(cmds...)
}

func (m *FormBlock) executeBlocDB() {

	m.spinnering = true

	m.propertyDB.StartBlock = m.inputs.startBlock.Value()
	m.propertyDB.FinishBlock = m.inputs.finishBlock.Value()
	m.propertyDB.KeyUnlock = m.inputs.keyUnlock.Value()
	m.propertyDB.NameUser = m.inputs.nameUser.Value()
	m.propertyDB.PasswordUser = m.inputs.passwordUser.Value()

	m.propertyDB.Massage = m.areas.massage.Value()

	massageJSON := OneCv8.MassageJSON{
		NameServer:     m.db.Server,
		NameDB:         m.db.NameOnServer,
		NameUser:       m.propertyDB.NameUser,
		PasswordUser:   m.propertyDB.PasswordUser,
		Block:          !m.propertyDB.Block,
		PermissionCode: m.propertyDB.KeyUnlock,
		DeniedMessage:  m.propertyDB.Massage,
		DeniedFrom:     m.propertyDB.StartBlock,
		DeniedTo:       m.propertyDB.FinishBlock,
	}

	err := OneCv8.PropertyDB(massageJSON)
	if err == nil {
		m.propertyDB.Block = !m.propertyDB.Block

		if m.propertyDB.Block {
			m.checkboxes.blocked.Selected[0] = struct{}{}
			m.error = "ОК unblock"
		} else {
			delete(m.checkboxes.blocked.Selected, m.checkboxes.blocked.Cursor)
			if client.PerformedActions.UpdateDB.Find(m.db.NameOnServer) == -1 {
				client.PerformedActions.UpdateDB = append(client.PerformedActions.UpdateDB, m.db.NameOnServer)
			}

			m.error = "ОК block"
		}
	} else {
		m.error = err.Error()
	}

	_, k := repository.GetPropertiesDB(client.Storage.PropertyDB, m.uid)
	if k != -1 {
		client.Storage.PropertyDB[k].Block = m.propertyDB.Block
	}

	m.spinnering = false
}

func (m *FormBlock) executeOpenDesignerDB() (tea.Model, tea.Cmd) {
	configDB := OneCv8.ConfigDB{
		Command: client.Storage.PathExe1C,
		Server:  m.db.Server,
		Port:    m.db.Port,
		DB:      m.db.NameOnServer,
		Key:     m.propertyDB.KeyUnlock,
		User:    m.propertyDB.NameUser,
		Pwd:     m.propertyDB.PasswordUser,
	}

	err := OneCv8.OpenConfigDB(configDB)
	if err != nil {
		m.error = err.Error()
	} else {
		m.error = "ОК open"
	}
	m.Update(nil)
	return nil, tea.Batch(nil)
}

func (m *FormBlock) executeDropConnectsDB() {

	m.spinnering = true

	massageJSON := OneCv8.MassageJSON{
		NameServer:   m.db.Server,
		NameDB:       m.db.NameOnServer,
		NameUser:     m.propertyDB.NameUser,
		PasswordUser: m.propertyDB.PasswordUser,
	}

	_, err := OneCv8.DropUsersDB(massageJSON)

	m.error = "ОК drop"
	if err != nil {
		m.error = err.Error()
	}

	m.spinnering = false
}

func (m *FormBlock) executeDropDoubleConnectsDB() {

	m.spinnering = true

	massageJSON := OneCv8.MassageJSON{
		NameServer:   m.db.Server,
		NameDB:       m.db.NameOnServer,
		NameUser:     m.propertyDB.NameUser,
		PasswordUser: m.propertyDB.PasswordUser,
	}

	var massagesJSON []OneCv8.MassageJSON
	massagesJSON = append(massagesJSON, massageJSON)

	var p uintptr
	const coin uint32 = 0
	err := ole.CoInitializeEx(p, coin)
	if err != nil {
		m.error = err.Error()
	}

	var closedConnects []repository.ClosedConnect
	scdc := OneCv8.SettingCloseDoubleConnection{
		IntervalWebClient: 1,
		OutputMessages:    false,
	}

	chanOut := make(chan repository.ClosedConnect)
	go OneCv8.DropDoubleUsersDB(massagesJSON, scdc, chanOut)

	for {
		cc, ok := <-chanOut
		if !ok {

			break
		}

		closedConnects = append(closedConnects, cc)
	}

	//_, err := OneCv8.DropDoubleUsersDB(massagesJSON)

	m.error = fmt.Sprintf("ОК drop (%d)", len(closedConnects))
	if err != nil {
		m.error = err.Error()
	}

	m.spinnering = false
}
