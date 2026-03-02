package radio

import (
	"bytes"
	"io"
	"time"

	"github.com/ebitengine/oto/v3"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
)

type AudioPlayer interface {
	Play(reader io.Reader) error
}

type NativeAudioPlayer struct {
	reader io.Reader
}

func NewNativeAudioPlayer() *NativeAudioPlayer {
	return &NativeAudioPlayer{
		reader: bytes.NewReader(nil),
	}
}

func (p *NativeAudioPlayer) Play(reader io.Reader) error {
	pcmStream, err := mp3.DecodeWithoutResampling(reader)
	if err != nil {
		return err
	}

	sr := pcmStream.SampleRate()

	op := oto.NewContextOptions{
		SampleRate:   sr,
		ChannelCount: 2,
		Format:       oto.FormatSignedInt16LE,
	}

	otoCtx, ready, err := oto.NewContext(&op)
	if err != nil {
		return err
	}
	<-ready

	player := otoCtx.NewPlayer(pcmStream)
	player.Play()

	for {
		if err := player.Err(); err != nil {
			return err
		}
		if err := otoCtx.Err(); err != nil {
			return err
		}

		if !player.IsPlaying() && player.BufferedSize() == 0 {
			return nil
		}

		time.Sleep(200 * time.Millisecond)
	}
}
