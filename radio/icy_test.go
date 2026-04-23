package radio

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func buildICYStream(metaInt int, chunks []struct {
	audio []byte
	meta  string
}) []byte {
	var buf bytes.Buffer
	for _, c := range chunks {
		buf.Write(c.audio)
		if c.meta == "" {
			buf.WriteByte(0)
		} else {
			padded := c.meta
			for len(padded)%16 != 0 {
				padded += "\x00"
			}
			buf.WriteByte(byte(len(padded) / 16))
			buf.WriteString(padded)
		}
	}
	return buf.Bytes()
}

func TestNewReaderInvalidMetaInt(t *testing.T) {
	_, err := NewReader(bytes.NewReader(nil), 0, nil)
	assert.Error(t, err)

	_, err = NewReader(bytes.NewReader(nil), -1, nil)
	assert.Error(t, err)
}

func TestReaderStripsMetadata(t *testing.T) {
	metaInt := 8
	audio1 := []byte("AAAAAAAA")
	audio2 := []byte("BBBBBBBB")
	meta := "StreamTitle='Art - Song';"

	stream := buildICYStream(metaInt, []struct {
		audio []byte
		meta  string
	}{
		{audio: audio1, meta: meta},
		{audio: audio2, meta: ""},
	})

	var received []Metadata
	r, err := NewReader(bytes.NewReader(stream), metaInt, func(m Metadata, raw string) {
		received = append(received, m)
	})
	require.NoError(t, err)

	out, err := io.ReadAll(r)
	require.NoError(t, err)

	assert.Equal(t, append(audio1, audio2...), out)
	require.Len(t, received, 1)
	assert.Equal(t, "Art", received[0].Artist)
	assert.Equal(t, "Song", received[0].Title)
}

func TestReaderSmallReads(t *testing.T) {
	metaInt := 4
	audio := []byte("ABCD")
	stream := buildICYStream(metaInt, []struct {
		audio []byte
		meta  string
	}{
		{audio: audio, meta: ""},
	})

	r, err := NewReader(bytes.NewReader(stream), metaInt, nil)
	require.NoError(t, err)

	buf := make([]byte, 2)
	var result []byte
	for {
		n, err := r.Read(buf)
		result = append(result, buf[:n]...)
		if err == io.EOF {
			break
		}
		require.NoError(t, err)
	}

	assert.Equal(t, audio, result)
}

func TestReaderEmptyRead(t *testing.T) {
	r, err := NewReader(bytes.NewReader([]byte("AAAA\x00")), 4, nil)
	require.NoError(t, err)

	n, err := r.Read(nil)
	assert.Equal(t, 0, n)
	assert.NoError(t, err)
}

func TestReaderMultipleMetadataBlocks(t *testing.T) {
	metaInt := 4
	chunks := []struct {
		audio []byte
		meta  string
	}{
		{audio: []byte("AAAA"), meta: "StreamTitle='X - Y';"},
		{audio: []byte("BBBB"), meta: "StreamTitle='A - B';"},
		{audio: []byte("CCCC"), meta: ""},
	}
	stream := buildICYStream(metaInt, chunks)

	var received []Metadata
	r, err := NewReader(bytes.NewReader(stream), metaInt, func(m Metadata, raw string) {
		received = append(received, m)
	})
	require.NoError(t, err)

	out, err := io.ReadAll(r)
	require.NoError(t, err)

	assert.Equal(t, []byte("AAAABBBBCCCC"), out)
	require.Len(t, received, 2)
	assert.Equal(t, Metadata{Artist: "X", Title: "Y"}, received[0])
	assert.Equal(t, Metadata{Artist: "A", Title: "B"}, received[1])
}
