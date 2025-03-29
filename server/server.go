package server

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
	"github.com/iammuuo/karibu/config"
	"github.com/iammuuo/karibu/karibu"
	"github.com/jackc/pgx/v5"
)

type Server struct {
	SshServer *ssh.Server
	conn      *pgx.Conn
}

//go:embed banner.txt
var banner string

// NewServer returns a Server instance with a default SSH server with the provided Middleware.
// A new SSH key pair of type ed25519 will be created if one does not exist.
// By default this server will accept all incoming connections, password and public key.
// Incase of an error an error is returned gracefuly for handling
func NewServer(cfg *config.Config) (*Server, error) {
	var server Server
	var err error

	// Connect to the database

	server.conn, err = pgx.Connect(context.Background(), fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.DatabaseConfig.User,
		cfg.DatabaseConfig.Password,
		cfg.DatabaseConfig.Host,
		cfg.DatabaseConfig.Port,
		cfg.DatabaseConfig.Name,
	))
	if err != nil {
		return nil, err
	}

	log.Info("Connected to database sucessfully")

	server.SshServer, err = wish.NewServer(
		wish.WithAddress(net.JoinHostPort(cfg.AppConfig.Addres, strconv.Itoa(cfg.AppConfig.Port))),

		// The SSH server need its own keys, this will create a keypair in the
		// given path if it doesn't exist yet.
		// By default, it will create an ED25519 key.
		wish.WithHostKeyPath(".ssh/id_ed25519"),

		// banner
		wish.WithBannerHandler(func(ctx ssh.Context) string {
			return fmt.Sprintf(banner, ctx.User())
		}),
		wish.WithVersion("1.0.0"),
		//
		// wish.WithPasswordAuth(func(ctx ssh.Context, password string) bool {
		// 	return password == "hello"
		// }),
		//
		// Middlewares do something on a ssh.Session, and then call the next
		// middleware in the stack.
		wish.WithMiddleware(
			func(next ssh.Handler) ssh.Handler {
				return func(sess ssh.Session) {
					wish.Println(sess, "Kwaheri ya kuonana!")
				}
			},
			bubbletea.Middleware(func(s ssh.Session) (tea.Model, []tea.ProgramOption) {
				return teaHandler(s, cfg, server.conn)
			}),
			// The last item in the chain is the first to be called.
			logging.Middleware(),
		),
	)

	if err != nil {
		return nil, err
	}

	return &server, nil
}

// Runs the server
func (s *Server) Run() error {
	log.Infof("Starting SSH server at address: %s", s.SshServer.Addr)
	log.Infof("Server version: %s", s.SshServer.Version)
	if err := s.SshServer.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		// We ignore ErrServerClosed because it is expected.
		log.Error("Could not start server", "error", err)
	}

	return nil

}

// You can wire any Bubble Tea model up to the middleware with a function that
// handles the incoming ssh.Session. Here we just grab the terminal info and
// pass it to the new model. You can also return tea.ProgramOptions (such as
// tea.WithAltScreen) on a session by session basis.
func teaHandler(s ssh.Session, cfg *config.Config, conn *pgx.Conn) (tea.Model, []tea.ProgramOption) {
	// This should never fail, as we are using the activeterm middleware.
	pty, _, _ := s.Pty()

	// When running a Bubble Tea app over SSH, you shouldn't use the default
	// lipgloss.NewStyle function.
	// That function will use the color profile from the os.Stdin, which is the
	// server, not the client.
	// We provide a MakeRenderer function in the bubbletea middleware package,
	// so you can easily get the correct renderer for the current session, and
	// use it to create the styles.
	// The recommended way to use these styles is to then pass them down to
	// your Bubble Tea model.
	renderer := bubbletea.MakeRenderer(s)

	karibu, err := karibu.NewKaribuApp(cfg, conn, pty, renderer, s)
	if err != nil {
		panic(err)
	}

	return karibu.Launch(), []tea.ProgramOption{tea.WithAltScreen()}
}

// Shutdown the server
func (s *Server) Shutdown(ctx context.Context) {
	log.Info("Server shutdown requested, attempting to shutdown")
	if err := s.conn.Close(ctx); err != nil {
		log.Error("Could not disconect from the database", "error", err)
		os.Exit(255)
	}

	log.Info("Database disconnected sucessfully!")
	if err := s.SshServer.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("Could not stop server", "error", err)
		os.Exit(255)
	}
	log.Info("Server shutdown was a success!")
}
