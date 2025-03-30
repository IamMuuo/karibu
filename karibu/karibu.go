package karibu

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
	"github.com/iammuuo/karibu/config"
	"github.com/iammuuo/karibu/karibu/tui"
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
	kum          tui.KaribuUiModel
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

	karibu.kum = tui.NewKaribuUiModel(&session, &pty)

	return &karibu, nil
}

func (k *Karibu) Launch() tea.Model {
	return &k.kum

}
