package styles

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const (
	Width       = 96
	ColumnWidth = 30
)

var (

	// General.

	Subtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	Highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	Special   = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}

	Divider = lipgloss.NewStyle().
		SetString("•").
		Padding(0, 1).
		Foreground(Subtle).
		String()

	Url = lipgloss.NewStyle().Foreground(Special).Render

	// Tabs.

	ActiveTabBorder = lipgloss.Border{
		Top:         "-",
		Bottom:      " ",
		Left:        "¦",
		Right:       "¦",
		TopLeft:     "?",
		TopRight:    "?",
		BottomLeft:  "-",
		BottomRight: "L",
	}

	TabBorder = lipgloss.Border{
		Top:         "-",
		Bottom:      "-",
		Left:        "¦",
		Right:       "¦",
		TopLeft:     "?",
		TopRight:    "?",
		BottomLeft:  "+",
		BottomRight: "+",
	}

	Tab = lipgloss.NewStyle().
		Border(TabBorder, true).
		BorderForeground(Highlight).
		Padding(0, 1)

	ShortLine = strings.Repeat("─", Max(0, 20))
	Line      = strings.Repeat("─", Max(0, 60))

	ActiveTab = Tab.Copy().Border(ActiveTabBorder, true)

	TabGap = Tab.Copy().
		BorderTop(false).
		BorderLeft(false).
		BorderRight(false)

	// Title.

	TitleStyle = lipgloss.NewStyle().
			MarginLeft(1).
			MarginRight(5).
			Padding(0, 1).
			Italic(true).
			Foreground(lipgloss.Color("#FFF7DB")).
			SetString("Lip Gloss")

	DescStyle = lipgloss.NewStyle().MarginTop(1)

	InfoStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderTop(true).
			BorderForeground(Subtle)

	// Dialog.

	DialogBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#874BFD")).
			Padding(1, 0).
			BorderTop(true).
			BorderLeft(true).
			BorderRight(true).
			BorderBottom(true)

	ButtonStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFF7DB")).
			Background(lipgloss.Color("#888B7E")).
			Padding(0, 3).
			MarginTop(1)

	ActiveButtonStyle = ButtonStyle.Copy().
				Foreground(lipgloss.Color("#FFF7DB")).
				Background(lipgloss.Color("#F25D94")).
				MarginRight(2).
				Underline(true)

	ListHeader = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForeground(Subtle).
			MarginRight(2).
			Render

	ListItem = lipgloss.NewStyle().PaddingLeft(2).Render

	CheckMark = lipgloss.NewStyle().SetString("?").
			Foreground(Special).
			PaddingRight(1).
			String()

	ListDone = func(s string) string {
		return CheckMark + lipgloss.NewStyle().
			Strikethrough(true).
			Foreground(lipgloss.AdaptiveColor{Light: "#969B86", Dark: "#696969"}).
			Render(s)
	}

	// Paragraphs/History.

	HistoryStyle = lipgloss.NewStyle().
			Align(lipgloss.Left).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(Highlight).
			Margin(1, 3, 0, 0).
			Padding(1, 2).
			Height(19).
			Width(ColumnWidth)

	// Status Bar.

	StatusNugget = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Padding(0, 1)

	StatusBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#343433", Dark: "#C1C6B2"}).
			Background(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#353533"})

	StatusStyle = lipgloss.NewStyle().
			Inherit(StatusBarStyle).
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#FF5F87")).
			Padding(0, 1).
			MarginRight(1)

	StatusStyleBlock = lipgloss.NewStyle().
				Inherit(StatusBarStyle).
				Foreground(lipgloss.Color("#FFFDF5")).
				Background(lipgloss.Color("#20B2AA")).
				Padding(0, 1).
				MarginRight(1)

	EncodingStyle = StatusNugget.Copy().
			Background(lipgloss.Color("#A550DF")).
			Align(lipgloss.Right).
			MarginRight(1)

	StatusText = lipgloss.NewStyle().Inherit(StatusBarStyle)

	FishCakeStyle = StatusNugget.Copy().Background(lipgloss.Color("#6124DF")).MarginRight(1)

	// Page.

	DocStyle = lipgloss.NewStyle().Padding(1, 2, 1, 2)
)

var (
	FocusedStyleFB = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	BlurredStyleFB = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	СaptionStyleFB = lipgloss.NewStyle().Foreground(lipgloss.Color("#33CCCC"))

	CursorStyleFB         = FocusedStyleFB.Copy()
	NoStyleFB             = lipgloss.NewStyle()
	NelpStyleFB           = BlurredStyleFB.Copy()
	CursorModeHelpStyleFB = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	FocusedButtonOk = FocusedStyleFB.Copy().Render("[ Submit ]")
	BlurredButtonOk = fmt.Sprintf("[ %s ]", BlurredStyleFB.Render("Submit"))

	//focusedButtonCancel = focusedStyleFB.Copy().Render("[ Cancel ]")
	//blurredButtonCancel = fmt.Sprintf("[ %s ]", blurredStyleFB.Render("Cancel"))
)

var (
	ColumnStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.HiddenBorder())
	FocusedStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62"))
	HelpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))

	Green   = lipgloss.Color("#04B575")
	GreenFg = lipgloss.NewStyle().Foreground(Green).Render

	Red   = lipgloss.Color("205")
	RedFg = lipgloss.NewStyle().Foreground(Red).Render
)

var (
	SelectColor   = lipgloss.Color("205")
	SelectColorFg = lipgloss.NewStyle().Foreground(SelectColor).Render
)

var (
	BlurredStyle             = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	CursorStyle              = FocusedStyle.Copy()
	NoStyle                  = lipgloss.NewStyle()
	CursorModeHelpStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	CursorModeHelpStyleWhite = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))

	FocusedButton = FocusedStyle.Copy().Render("[ Submit ]")
	BlurredButton = fmt.Sprintf("[ %s ]", BlurredStyle.Render("Submit"))
)

var BaseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
