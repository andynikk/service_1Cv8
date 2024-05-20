package cli

import (
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/text/encoding/charmap"
	"strings"

	"Service_1Cv8/internal/cli/charm"
	"Service_1Cv8/internal/cli/styles"
	"Service_1Cv8/internal/constants"
	"Service_1Cv8/internal/token"
	"Service_1Cv8/internal/winsys"
)

const (
	itemsInputsFormResTlg   = 2
	itemsAreasFormResTlg    = 1
	itemsButtonFormResTlg   = 1
	itemsCheckBoxFormResTlg = 20
	itemsTableFormResTlg    = 0
)

type buttonFormResTlg struct {
	buttonSend charm.BoolModel
}

type checkFormResTlg struct {
	checkDBUH  charm.BoolModel
	checkDBDO  charm.BoolModel
	checkDBZUP charm.BoolModel
	checkDBBIT charm.BoolModel
	checkDBERP charm.BoolModel

	checkSrvClearTest  charm.BoolModel
	checkSrvClear1C    charm.BoolModel
	checkSrvClear2C    charm.BoolModel
	checkSrvClearEvo01 charm.BoolModel
	checkSrvClearEvo02 charm.BoolModel

	checkSrvSSTest  charm.BoolModel
	checkSrvSS1C    charm.BoolModel
	checkSrvSS2C    charm.BoolModel
	checkSrvSSEvo01 charm.BoolModel
	checkSrvSSEvo02 charm.BoolModel

	checkSrvRestartTest  charm.BoolModel
	checkSrvRestart1C    charm.BoolModel
	checkSrvRestart2C    charm.BoolModel
	checkSrvRestartEvo01 charm.BoolModel
	checkSrvRestartEvo02 charm.BoolModel
}

type inputsFormResTlg struct {
	sender  textinput.Model
	channel textinput.Model
}

type areasFormResTlg struct {
	message textarea.Model
}

type FormResTlg struct {
	focusIndex int
	inputs     inputsFormResTlg
	areas      areasFormResTlg
	buttons    buttonFormResTlg
	checks     checkFormResTlg
	message    string

	claimStore token.ClaimStore

	spinner    spinner.Model
	spinnering bool

	crawlElements charm.CrawlElements
}

func (m *FormResTlg) SetParameters(args []interface{}) {

	m.inputs.sender.Placeholder = "Sender"
	m.inputs.sender.SetValue("6524372399:AAH7QhU9IO8sXmWdC0ofRNyoqrcf_qUUJho")
	m.inputs.sender.Focus()
	m.inputs.sender.PromptStyle = styles.FocusedStyleFB
	m.inputs.sender.TextStyle = styles.FocusedStyleFB

	m.inputs.channel.Placeholder = ""
	m.inputs.channel.SetValue("1986806055")

	m.areas.message.Placeholder = "key"
	m.areas.message.SetValue("")
	m.areas.message.SetWidth(400)
	m.areas.message.SetHeight(5)
}

func (m *FormResTlg) FillFormResTlgElements() {
	ce := make(charm.CrawlElements)

	ce[0] = &m.inputs.sender
	ce[1] = &m.inputs.channel

	ce[2] = &m.checks.checkDBUH
	ce[3] = &m.checks.checkDBDO
	ce[4] = &m.checks.checkDBZUP
	ce[5] = &m.checks.checkDBBIT
	ce[6] = &m.checks.checkDBERP

	ce[7] = &m.checks.checkSrvClear1C
	ce[8] = &m.checks.checkSrvClear2C
	ce[9] = &m.checks.checkSrvClearEvo01
	ce[10] = &m.checks.checkSrvClearEvo02
	ce[11] = &m.checks.checkSrvClearTest

	ce[12] = &m.checks.checkSrvSS1C
	ce[13] = &m.checks.checkSrvSS2C
	ce[14] = &m.checks.checkSrvSSEvo01
	ce[15] = &m.checks.checkSrvSSEvo02
	ce[16] = &m.checks.checkSrvSSTest

	ce[17] = &m.checks.checkSrvRestart1C
	ce[18] = &m.checks.checkSrvRestart2C
	ce[19] = &m.checks.checkSrvRestartEvo01
	ce[20] = &m.checks.checkSrvRestartEvo02
	ce[21] = &m.checks.checkSrvRestartTest

	ce[22] = &m.areas.message

	ce[23] = &m.buttons.buttonSend

	m.crawlElements = ce
}

func NewFormResTlg() *FormResTlg {

	newSpinner := spinner.New()
	newSpinner.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))
	newSpinner.Spinner = spinner.Pulse

	m := FormResTlg{
		buttons:    buttonFormResTlg{},
		spinnering: false,
		spinner:    newSpinner,
	}

	var t textinput.Model
	var a textarea.Model

	t = textinput.New()
	m.inputs.sender = t

	t = textinput.New()
	m.inputs.channel = t

	a = textarea.New()
	a.CharLimit = 4000
	m.areas.message = a

	return &m
}

func (m *FormResTlg) Init() tea.Cmd {
	return nil
}

func (m *FormResTlg) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	cmds := make([]tea.Cmd, 0)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		strMsg := msg.String()
		switch strMsg {
		case "esc", "ctrl+q":

			model := client.Models[constants.FormExchangeListToken].(*FormResTlg)
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
			if m.focusIndex == 23 {
				m.sendTelegram()
			}
		}
	}

	cmd = m.updateInputs(msg)
	return m, cmd
}

func (m *FormResTlg) View() string {
	var b strings.Builder

	b.WriteRune('\n')
	b.WriteString(fmt.Sprintf("%s Telegram setting %s\n", styles.ShortLine, styles.Line))
	b.WriteRune('\n')
	b.WriteString(fmt.Sprintf("%s %s %s %s\n",
		styles.СaptionStyleFB.Render("Sender:"), m.inputs.sender.View(),
		styles.СaptionStyleFB.Render("Channel:"), m.inputs.channel.View()))
	b.WriteRune('\n')

	///////////////////////////////////////////

	b.WriteString(fmt.Sprintf("%s Updated databases %s\n", styles.ShortLine, styles.Line))
	b.WriteString(m.addCheckBox(&m.checks.checkDBUH, "UH") + "  ")
	b.WriteString(m.addCheckBox(&m.checks.checkDBUH, "DO") + "  ")
	b.WriteString(m.addCheckBox(&m.checks.checkDBUH, "ZUP") + "  ")
	b.WriteString(m.addCheckBox(&m.checks.checkDBUH, "BIT") + "  ")
	b.WriteString(m.addCheckBox(&m.checks.checkDBUH, "ERP") + "  ")
	b.WriteString("\n")

	///////////////////////////////////////////

	b.WriteString(fmt.Sprintf("%s Server clear cash %s\n", styles.ShortLine, styles.Line))
	m.checks.checkSrvClear1C.Value = true
	b.WriteString(m.addCheckBox(&m.checks.checkSrvClear1C, "1C") + "  ")

	m.checks.checkSrvClear2C.Value = true
	b.WriteString(m.addCheckBox(&m.checks.checkSrvClear2C, "2C") + "  ")

	m.checks.checkSrvClearEvo01.Value = true
	b.WriteString(m.addCheckBox(&m.checks.checkSrvClearEvo01, "EVO-1C-01") + "  ")

	m.checks.checkSrvClearEvo02.Value = true
	b.WriteString(m.addCheckBox(&m.checks.checkSrvClearEvo02, "EVO-1C-02") + "  ")

	m.checks.checkSrvClearTest.Value = true
	b.WriteString(m.addCheckBox(&m.checks.checkSrvClearTest, "TK-TEST-APP") + "  ")
	b.WriteString("\n")

	///////////////////////////////////////////

	b.WriteString(fmt.Sprintf("%s Server clear cash %s\n", styles.ShortLine, styles.Line))
	m.checks.checkSrvClear1C.Value = true
	b.WriteString(m.addCheckBox(&m.checks.checkSrvClear1C, "1C") + "  ")

	m.checks.checkSrvClear2C.Value = true
	b.WriteString(m.addCheckBox(&m.checks.checkSrvClear2C, "2C") + "  ")

	m.checks.checkSrvClearEvo01.Value = true
	b.WriteString(m.addCheckBox(&m.checks.checkSrvClearEvo01, "EVO-1C-01") + "  ")

	m.checks.checkSrvClearEvo02.Value = true
	b.WriteString(m.addCheckBox(&m.checks.checkSrvClearEvo02, "EVO-1C-02") + "  ")

	m.checks.checkSrvClearTest.Value = true
	b.WriteString(m.addCheckBox(&m.checks.checkSrvClearTest, "TK-TEST-APP") + "  ")
	b.WriteString("\n")

	///////////////////////////////////////////

	b.WriteString(fmt.Sprintf("%s Restert services %s\n", styles.ShortLine, styles.Line))
	m.checks.checkSrvSS1C.Value = true
	b.WriteString(m.addCheckBox(&m.checks.checkSrvSS1C, "1C") + "  ")

	m.checks.checkSrvSS2C.Value = true
	b.WriteString(m.addCheckBox(&m.checks.checkSrvSS2C, "2C") + "  ")

	m.checks.checkSrvSSEvo01.Value = true
	b.WriteString(m.addCheckBox(&m.checks.checkSrvSSEvo01, "EVO-1C-01") + "  ")

	m.checks.checkSrvSSEvo02.Value = true
	b.WriteString(m.addCheckBox(&m.checks.checkSrvSSEvo02, "EVO-1C-02") + "  ")

	m.checks.checkSrvSSTest.Value = true
	b.WriteString(m.addCheckBox(&m.checks.checkSrvSSTest, "TK-TEST-APP") + "  ")

	b.WriteString("\n")

	///////////////////////////////////////////

	b.WriteString(fmt.Sprintf("%s Restert services %s\n", styles.ShortLine, styles.Line))
	m.checks.checkSrvRestart1C.Value = true
	b.WriteString(m.addCheckBox(&m.checks.checkSrvRestart1C, "1C") + "  ")

	m.checks.checkSrvRestart2C.Value = true
	b.WriteString(m.addCheckBox(&m.checks.checkSrvRestart2C, "2C") + "  ")

	m.checks.checkSrvSSEvo01.Value = true
	b.WriteString(m.addCheckBox(&m.checks.checkSrvSSEvo01, "EVO-1C-01") + "  ")

	m.checks.checkSrvSSEvo02.Value = true
	b.WriteString(m.addCheckBox(&m.checks.checkSrvSSEvo02, "EVO-1C-02") + "  ")

	m.checks.checkSrvSSTest.Value = true
	b.WriteString(m.addCheckBox(&m.checks.checkSrvSSTest, "TK-TEST-APP") + "  ")

	b.WriteString("\n")

	/////////

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

func (m *FormResTlg) addCheckBox(check *charm.BoolModel, text string) string {
	cursor := " "
	checked := " "
	if check.Active {
		cursor = styles.FocusedStyleFB.Render(">")
	}
	if check.Value {
		checked = "x"
	}
	return fmt.Sprintf("%s [%s] %s", cursor, checked, text)
}

func (m *FormResTlg) lenForm() int {

	return itemsInputsFormResTlg + itemsAreasFormResTlg + itemsButtonFormResTlg +
		itemsCheckBoxFormResTlg + itemsTableFormResTlg
}

func (m *FormResTlg) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, itemsInputsFormResTlg+itemsAreasFormResTlg)

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

	m.inputs.sender, _ = m.inputs.sender.Update(msg)
	m.inputs.channel, _ = m.inputs.channel.Update(msg)
	m.areas.message, _ = m.areas.message.Update(msg)

	return tea.Batch(cmds...)
}

func (m *FormResTlg) sendTelegram() {

}
