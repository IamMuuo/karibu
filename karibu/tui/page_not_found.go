package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
)

type PageNotFound struct {
	session *ssh.Session
	pty     *ssh.Pty
}

func (pnf PageNotFound) Init() tea.Cmd {
	return nil
}

func (pnf PageNotFound) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return pnf, nil
}

func (pnf PageNotFound) View() string {

	var textStyle = lipgloss.NewStyle().
		Padding(8).
		Align(lipgloss.Center).
		BorderStyle(lipgloss.DoubleBorder()).
		Bold(true).
		Background(lipgloss.Color("#684eff")).
		Foreground(lipgloss.Color("#ffc44c"))

	var helptext = lipgloss.PlaceVertical(10,
		lipgloss.Bottom, lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ff3d8e")).
			Render("Hit <q> or <ctrl+c> to quit"),
		lipgloss.WithWhitespaceChars(""),
	)

	var content = lipgloss.Place(pnf.pty.Window.Width,
		pnf.pty.Window.Height-10,
		lipgloss.Center,
		lipgloss.Center,
		textStyle.Render(
			fmt.Sprintf(
				"Hey %s, looks like we hit a 404 on the route you requested. Please check your path and try again.",
				(*pnf.session).User(),
			),
		),
		lipgloss.WithWhitespaceChars(""))

	return fmt.Sprintf("%s %s ", content, helptext)
}
