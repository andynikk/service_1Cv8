package forms_exchange

import (
	"fmt"
	"github.com/atotto/clipboard"
	"golang.org/x/text/encoding/charmap"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"Service_1Cv8/internal/cli/charm"
	"Service_1Cv8/internal/cli/styles"
	"Service_1Cv8/internal/compression"
	"Service_1Cv8/internal/constants"
	"Service_1Cv8/internal/repository"
	"Service_1Cv8/internal/token"
	"Service_1Cv8/internal/winsys"
)

const (
	itemsInputsFormKeyToken   = 2
	itemsAreasFormKeyToken    = 1
	itemsButtonFormKeyToken   = 3
	itemsCheckBoxFormKeyToken = 0
	itemsTableFormKeyToken    = 0
)

type buttonFormKeyToken struct {
	buttonCreate charm.BoolModel
	buttonSave   charm.BoolModel
	buttonDel    charm.BoolModel
	//buttonProlong charm.BoolModel
}

type inputsFormKeyToken struct {
	id       textinput.Model
	timeLive textinput.Model
}

type areasFormKeyToken struct {
	key textarea.Model
	//public textarea.Model
	token textarea.Model
}

type FormKeyToken struct {
	focusIndex int
	inputs     inputsFormKeyToken
	areas      areasFormKeyToken
	buttons    buttonFormKeyToken
	message    string

	claimStore token.ClaimStore

	spinner    spinner.Model
	spinnering bool

	crawlElements charm.CrawlElements
}

func (m *FormKeyToken) SetParameters(args []interface{}) {

	m.inputs.id.Placeholder = "ID"
	m.inputs.id.SetValue("")
	m.inputs.id.Focus()
	m.inputs.id.PromptStyle = styles.FocusedStyleFB
	m.inputs.id.TextStyle = styles.FocusedStyleFB

	m.inputs.timeLive.Placeholder = "Time live (week)"
	m.inputs.timeLive.SetValue("1")

	m.areas.key.Placeholder = "key"
	m.areas.key.SetValue("")
	m.areas.key.SetWidth(400)
	m.areas.key.SetHeight(5)

	for _, v := range args {
		switch v.(type) {
		case token.ClaimStore:
			m.claimStore = v.(token.ClaimStore)
			//if !ok {
			//	continue
			//}

			m.inputs.id.SetValue(m.claimStore.Key)
			m.areas.key.SetValue(string(m.claimStore.Secret))
			m.areas.token.SetValue(string(m.claimStore.Value))
			m.inputs.timeLive.SetValue("1")

		}
	}

	m.areas.token.Placeholder = "token"
	m.areas.token.SetValue("")
	m.areas.token.SetWidth(400)
	m.areas.token.SetHeight(5)

	m.FillFormKeyTokenElements()
}

func (m *FormKeyToken) FillFormKeyTokenElements() {
	ce := make(charm.CrawlElements)

	ce[0] = &m.inputs.id
	ce[1] = &m.inputs.timeLive
	ce[2] = &m.areas.key
	ce[3] = &m.buttons.buttonCreate
	ce[4] = &m.buttons.buttonSave
	ce[5] = &m.buttons.buttonDel
	//ce[6] = &m.buttons.buttonProlong
	//ce[5] = &m.buttons.buttonFile
	//ce[7] = &m.areas.public
	//ce[8] = &m.areas.token

	m.crawlElements = ce
}

func NewFormKeyToken() *FormKeyToken {

	newSpinner := spinner.New()
	newSpinner.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))
	newSpinner.Spinner = spinner.Pulse

	m := FormKeyToken{
		buttons:    buttonFormKeyToken{},
		spinnering: false,
		spinner:    newSpinner,
	}

	var t textinput.Model
	var a textarea.Model

	t = textinput.New()
	m.inputs.id = t

	t = textinput.New()
	m.inputs.timeLive = t

	a = textarea.New()
	a.CharLimit = 2000
	m.areas.key = a

	//a = textarea.New()
	//a.CharLimit = 2000
	//m.areas.public = a

	a = textarea.New()
	a.CharLimit = 2000
	m.areas.token = a

	return &m
}

func (m *FormKeyToken) Init() tea.Cmd {
	return nil
}

func (m *FormKeyToken) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	cmds := make([]tea.Cmd, 0)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		strMsg := msg.String()
		switch strMsg {
		case "esc", "ctrl+q":

			model := exchanger.Models[constants.FormExchangeListToken].(*FormListToken)
			model.SetParameters(nil)

			return model.Update(nil)
		case "ctrl+c", "f10":

			return m, tea.Quit

		case "tab", "shift+tab":

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
		case "enter":
			if m.focusIndex == 3 {
				m.message = "OK save settings"

				err := m.createTokenKey()
				if err != nil {
					m.message = err.Error()
				}

			}

			if m.focusIndex == 4 {
				m.saveToken()
			}

			if m.focusIndex == 5 {
				m.delToken()
			}

		}
	}

	cmd = m.updateInputs(msg)
	return m, cmd
}

func (m *FormKeyToken) View() string {
	var b strings.Builder

	b.WriteRune('\n')
	b.WriteString(fmt.Sprintf("%s Setting %s\n", styles.ShortLine, styles.Line))
	b.WriteRune('\n')

	b.WriteString(fmt.Sprintf("%s %s %s %s\n",
		styles.СaptionStyleFB.Render("ID:"), m.inputs.id.View(),
		styles.СaptionStyleFB.Render("Time live (week):"), m.inputs.timeLive.View()))

	b.WriteRune('\n')
	b.WriteString(fmt.Sprintf("%s\n%s\n", styles.СaptionStyleFB.Render("Key:"), m.areas.key.View()))

	b.WriteRune('\n')
	if m.buttons.buttonCreate.Active {
		b.WriteString(fmt.Sprintf("[ %s ]", styles.FocusedStyleFB.Render("edit/created")))
	} else {
		b.WriteString(fmt.Sprintf("[ %s ]", styles.BlurredStyleFB.Render("edit/created")))
	}
	b.WriteString(" ")

	if m.buttons.buttonSave.Active {
		b.WriteString(fmt.Sprintf("[ %s ]", styles.FocusedStyleFB.Render("save")))
	} else {
		b.WriteString(fmt.Sprintf("[ %s ]", styles.BlurredStyleFB.Render("save")))
	}
	b.WriteString(" ")

	if m.buttons.buttonDel.Active {
		b.WriteString(fmt.Sprintf("[ %s ]", styles.FocusedStyleFB.Render("delete")))
	} else {
		b.WriteString(fmt.Sprintf("[ %s ]", styles.BlurredStyleFB.Render("delete")))
	}
	b.WriteString("\n")

	//b.WriteString(fmt.Sprintf("%s\n%s\n", styles.СaptionStyleFB.Render("Public:"), m.areas.public.View()))

	b.WriteString(fmt.Sprintf("%s\n%s\n", styles.СaptionStyleFB.Render("Token:"), m.areas.token.View()))

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

func (m *FormKeyToken) lenForm() int {

	return itemsInputsFormKeyToken + itemsAreasFormKeyToken + itemsButtonFormKeyToken +
		itemsCheckBoxFormKeyToken + itemsTableFormKeyToken
}

func (m *FormKeyToken) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, itemsInputsFormKeyToken+itemsAreasFormKeyToken)

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
	m.inputs.timeLive, _ = m.inputs.timeLive.Update(msg)
	m.areas.key, _ = m.areas.key.Update(msg)
	//m.areas.public, _ = m.areas.public.Update(msg)
	m.areas.token, _ = m.areas.token.Update(msg)

	return tea.Batch(cmds...)
}

func (m *FormKeyToken) createTokenKey() error {

	sk := m.areas.key.Value()
	if sk == "" {
		sk = constants.SecretKey
	}
	byteSK := []byte(sk)

	timeLive := m.inputs.timeLive.Value()
	intTimeLive, err := strconv.Atoi(timeLive)
	if err != nil {
		intTimeLive = 1
	}
	intTimeLive = intTimeLive * constants.TimeLiveToken

	tc := token.NewClaims(m.inputs.id.Value(), time.Duration(intTimeLive))
	tokenString, err := tc.GenerateJWT(byteSK)
	if err != nil {
		return err
	}

	gzipTokenString, err := compression.Compress([]byte(tokenString))
	if err != nil {
		return err
	}
	gziSecretString, err := compression.Compress(byteSK)
	if err != nil {
		return err
	}

	m.claimStore = token.ClaimStore{
		Key:    m.inputs.id.Value(),
		Value:  gzipTokenString,
		Secret: gziSecretString,
	}

	m.areas.token.SetValue(tokenString)

	return nil
}

func (m *FormKeyToken) saveToken() {

	if m.claimStore.Key == "" {
		return
	}

	err := repository.SetToken(&m.claimStore)
	if err != nil {
		m.message = err.Error()
		return
	}

	t, err := compression.Decompress(m.claimStore.Value)
	if err != nil {
		m.message = err.Error()
		return
	}

	err = clipboard.WriteAll(string(t))
	if err != nil {
		m.message = err.Error()
		return
	}

	m.message = "OK! Token save & write to clipboard "
}

func (m *FormKeyToken) delToken() {
	if m.claimStore.Key == "" {
		return
	}

	err := repository.DelToken(&m.claimStore)
	if err != nil {
		m.message = err.Error()
		return
	}

	m.message = "Del OK!"
}
