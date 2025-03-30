package tui

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
)

type SplashScreen struct {
	session *ssh.Session
	pty     *ssh.Pty
}

func (ss SplashScreen) Init() tea.Cmd {
	return tea.Tick(time.Second*5, func(t time.Time) tea.Msg {
		return SwitchPageMsg{TargetPage: "hello"}
	})
}

func (ss SplashScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	return ss, func() tea.Msg {
		return msg
	}
}

func (ss SplashScreen) View() string {
	var textStyle = lipgloss.NewStyle().
		Padding(8).
		Align(lipgloss.Center).
		BorderStyle(lipgloss.DoubleBorder()).
		Bold(true).
		BorderForeground(lipgloss.Color("#684eff")).
		Foreground(lipgloss.Color("#f21dcb"))

	var helptext = lipgloss.PlaceVertical(10,
		lipgloss.Bottom, lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ff3d8e")).
			Render("Hit <q> or <ctrl+q> to quit"),
		lipgloss.WithWhitespaceChars(""),
	)

	var content = lipgloss.Place(ss.pty.Window.Width,
		ss.pty.Window.Height-10,
		lipgloss.Center,
		lipgloss.Center,
		textStyle.Render(fmt.Sprintf("Karibu %s!", (*ss.session).User())),
		lipgloss.WithWhitespaceChars(""),
	)

	return fmt.Sprintf("%s %s ", content, helptext)
}
