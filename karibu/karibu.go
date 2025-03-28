package karibu

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/iammuuo/karibu/config"
	"github.com/jackc/pgx/v5"
)

type Karibu struct {
	term         string
	profile      string
	dbConnection *pgx.Conn
	renderer     *lipgloss.Renderer
	pty          *ssh.Pty
	width        int
	height       int
	session      ssh.Session
}

// Configures a new karibu application instance
func NewKaribuApp(cfg *config.Config, conn *pgx.Conn, pty ssh.Pty,
	renderer *lipgloss.Renderer,
	session ssh.Session,
) (*Karibu, error) {
	var karibu Karibu

	karibu.term = pty.Term
	karibu.width = pty.Window.Width
	karibu.height = pty.Window.Height
	karibu.profile = renderer.ColorProfile().Name()
	karibu.dbConnection = conn
	karibu.renderer = renderer
	karibu.pty = &pty
	karibu.session = session

	return &karibu, nil
}

func (k Karibu) Init() tea.Cmd {
	return nil
}

func (k Karibu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch event := msg.(type) {
	case tea.KeyMsg:
		switch event.String() {

		case "ctrl+c", "esc", "q":
			return k, tea.Quit
		}
	case tea.WindowSizeMsg:
		log.Info("Window resize event!")
		k.width = k.pty.Window.Width
		k.height = k.pty.Window.Height
		return k, nil

	}
	return k, nil

}

func (k Karibu) View() string {
	var style = lipgloss.NewStyle().
		Bold(true).
		Padding(8).
		Foreground(lipgloss.Color("#ce82ff")).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#ffa979")).
		Align(lipgloss.Center)

	quitInstructions := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ffa979")).
		Render("Press esc or q to quit")

	screen := lipgloss.Place(k.pty.Window.Width, k.pty.Window.Height-2,
		lipgloss.Center, lipgloss.Center,
		style.Render(fmt.Sprintf("Karibu dev %s!", k.session.User())),
		lipgloss.WithWhitespaceChars(" "),
	)
	helpText := lipgloss.Place(0, 2,
		lipgloss.Bottom, lipgloss.Bottom,
		quitInstructions, lipgloss.WithWhitespaceChars(" "),
	)

	return fmt.Sprintf("%s %s\n", screen, helpText)
}
