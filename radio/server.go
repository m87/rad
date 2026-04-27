package radio

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"os"
	"path/filepath"
)

func StateDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".local", "state", "rad")
}

func SaveToLibrary(metadata Metadata) error {
	home, _ := os.UserHomeDir()
	libraryPath := filepath.Join(home, ".rad-library")

	if _, err := os.Stat(libraryPath); os.IsNotExist(err) {
		if err := os.WriteFile(libraryPath, []byte("[]"), 0644); err != nil {
			return fmt.Errorf("failed to create library file: %w", err)
		}
	}

	var library []Metadata
	if data, err := os.ReadFile(libraryPath); err == nil {
		if err := json.Unmarshal(data, &library); err != nil {
			return fmt.Errorf("failed to parse library: %w", err)
		}
	}

	for _, m := range library {
		if m.Title == metadata.Title && m.Artist == metadata.Artist {
			slog.Info("Track already in library", "title", metadata.Title, "artist", metadata.Artist)
			return nil
		}
	}

	library = append(library, metadata)

	data, err := json.MarshalIndent(library, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize library: %w", err)
	}

	if err := os.WriteFile(libraryPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write library: %w", err)
	}

	return nil

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
