package karibu

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
	"github.com/iammuuo/karibu/config"
	"github.com/jackc/pgx/v5"
)

type Karibu struct {
	regPage      RegistrationPage
	term         string
	profile      string
	width        int
	height       int
	dbConnection *pgx.Conn
	renderer     *lipgloss.Renderer
}

// Configures a new karibu application instance
func NewKaribuApp(cfg *config.Config, conn *pgx.Conn, pty ssh.Pty, renderer *lipgloss.Renderer) (*Karibu, error) {
	var karibu Karibu

	karibu.term = pty.Term
	karibu.profile = renderer.ColorProfile().Name()
	karibu.width = pty.Window.Width
	karibu.height = pty.Window.Height
	karibu.dbConnection = conn
	karibu.renderer = renderer
	karibu.regPage = NewRegistrationPage()

	return &karibu, nil
}

func (k Karibu) Init() tea.Cmd {
	return k.regPage.Init()
}

func (k Karibu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return k.regPage.Update(msg)
}

func (k Karibu) View() string {
	return k.regPage.View()
}
