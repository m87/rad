package radio

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
)

type OnMeta func(metadata Metadata, raw string)

type Reader struct {
	r         *bufio.Reader
	metaInt   int
	audioLeft int

	onMeta OnMeta
}

func NewReader(src io.Reader, metaInt int, onMeta OnMeta) (*Reader, error) {
	if metaInt <= 0 {
		return nil, fmt.Errorf("Invalid metadata interval: %d", metaInt)
	}
	return &Reader{
		r:         bufio.NewReaderSize(src, 64*1024),
		metaInt:   metaInt,
		audioLeft: metaInt,
		onMeta:    onMeta,
	}, nil
}

func (r *Reader) Read(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}

	if r.audioLeft == 0 {
		if err := r.readMeta(); err != nil {
			return 0, err
		}
		r.audioLeft = r.metaInt
	}

	toRead := len(p)
	if toRead > r.audioLeft {
		toRead = r.audioLeft
	}

	_, err = io.ReadFull(r.r, p[:toRead])
	if err != nil {
		return 0, err
	}

	r.audioLeft -= toRead
	return toRead, nil
}

func (r *Reader) readMeta() error {
	var len1 [1]byte
	if _, err := io.ReadFull(r.r, len1[:]); err != nil {
		return err
	}
	metaLen := int(len1[0]) * 16
	if metaLen == 0 {
		return nil
	}

	metaData := make([]byte, metaLen)
	if _, err := io.ReadFull(r.r, metaData); err != nil {
		return err
	}

	raw := string(bytes.TrimRight(metaData, "\x00"))

	metadata := ParseMetadata(raw)
	if r.onMeta != nil {
		r.onMeta(metadata, raw)
	}
	return nil
}
