package radio

import (
	"log/slog"
	"net"
	"os"
	"path/filepath"
)

func StateDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".local", "state", "rad")
}

func SockPath() string { return filepath.Join(StateDir(), "rad.sock") }

func RunServer(radio *Radio, handler func(conn net.Conn, radio *Radio)) error {
	_ = os.MkdirAll(StateDir(), 0755)
	_ = os.Remove(SockPath())

	l, err := net.Listen("unix", SockPath())
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
