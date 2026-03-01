package radio

import (
	"io"
	"log/slog"
	"os/exec"
)

type AudioPlayer struct {
	cmd   *exec.Cmd
	stdin io.WriteCloser
}

func NewAudioPlayer() *AudioPlayer {
	return &AudioPlayer{}
}

func (p *AudioPlayer) Init() error {
	p.cmd = exec.Command("mpv", "--no-terminal", "--", "-")
	stdin, _ := p.cmd.StdinPipe()
	p.stdin = stdin
	p.cmd.Start()
	slog.Info("Started mpv process with PID", "pid", p.cmd.Process.Pid)
	return nil
}

func (p *AudioPlayer) Write(data []byte) error {
	_, err := p.stdin.Write(data)
	return err
}
