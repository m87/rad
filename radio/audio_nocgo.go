//go:build !cgo

package radio

import "io"

type AudioPlayer interface {
	Play(reader io.Reader) error
}

type NativeAudioPlayer struct{}

func NewNativeAudioPlayer() *NativeAudioPlayer {
	return &NativeAudioPlayer{}
}

func (p *NativeAudioPlayer) Play(reader io.Reader) error {
	return NewMpvAudioPlayer().Play(reader)
}
