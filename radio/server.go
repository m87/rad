package radio

import (
	"log/slog"
	"net"
	"os"
	"path/filepath"
)

func stateDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".local", "state", "rad")
}
func sockPath() string { return filepath.Join(stateDir(), "rad.sock") }

func RunServer(radio *Radio, handler func(conn net.Conn, radio *Radio)) error {
	_ = os.MkdirAll(stateDir(), 0755)
	_ = os.Remove(sockPath())

	l, err := net.Listen("unix", sockPath())
	if err != nil {
		return err
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			slog.Error("Failed to accept connection", "err", err)
			continue
		}

		go handler(conn, radio)
	}
}
