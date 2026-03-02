package radio

import (
	"io"
	"log/slog"
	"os"
	"os/exec"
)

type MpvAudioPlayer struct {
}

func NewMpvAudioPlayer() *MpvAudioPlayer {
	return &MpvAudioPlayer{}
}

func (p *MpvAudioPlayer) Play(reader io.Reader) error {
	cmd := exec.Command("mpv", "--no-terminal", "--quiet", "--audio-client-name=Rad Player", "-")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		_ = stdin.Close()
		return err
	}

	copyErrCh := make(chan error, 1)
	go func() {
		_, err := io.Copy(stdin, reader)
		_ = stdin.Close()
		copyErrCh <- err
	}()

	waitErr := cmd.Wait()

	copyErr := <-copyErrCh

	if waitErr != nil {
		return waitErr
	}
	if copyErr != nil {
		slog.Debug("copy to mpv ended", "err", copyErr)
	}
	return nil
}
