package server

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/logging"
	"github.com/iammuuo/karibu/config"
)

type Server struct {
	SshServer *ssh.Server
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
					wish.Println(sess, "Karibu")
					next(sess)
				}
			},

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
	log.Info("Starting SSH server at address: ", s.SshServer.Addr)
	log.Info("Server version: ", s.SshServer.Version)
	if err := s.SshServer.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		// We ignore ErrServerClosed because it is expected.
		log.Error("Could not start server", "error", err)
	}

	return nil

}

// Shutdown the server
func (s *Server) Shutdown(ctx context.Context) {
	log.Info("Server shutdown requested, attempting to shutdown")
	if err := s.SshServer.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("Could not stop server", "error", err)
		os.Exit(255)
	}
	log.Info("Server shutdown was a success!")
}
