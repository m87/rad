package radio

import (
	"bufio"
	"context"
	"io"
	"log/slog"
	"net"
	"net/http"
	"strconv"
	"time"
)

const version = "0.0.1"

type Radio struct {
	url string
}

func NewRadio(url string) *Radio {
	return &Radio{
		url: url,
	}
}

func (r *Radio) Play() error {
	slog.Info("Playing radio from URL", "url", r.url)
	httpClient := &http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout:   5 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
			IdleConnTimeout:       30 * time.Second,
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", r.url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Icy-MetaData", "1")
	req.Header.Set("User-Agent", "rad/"+version)

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	slog.Info("Connected to radio stream", "status", resp.Status)
	slog.Debug("Radio station metadata", "headers", resp.Header)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		panic("Failed to connect to radio stream: " + resp.Status)
	}

	icyMetaInt := 16000
	if val := resp.Header.Get("Icy-MetaInt"); val != "" {
		var err error
		icyMetaInt, err = strconv.Atoi(val)
		if err != nil {
			slog.Warn("Invalid Icy-MetaInt header, ignoring", "value", val)
			icyMetaInt = 0
		}
	}
	slog.Debug("Icy-MetaInt value", "icy_meta_int", icyMetaInt)

	reader := bufio.NewReaderSize(resp.Body, 64*1024)
	buf := make([]byte, 4*1024)
	var total int64
	lastPrint := time.Now()
	var bytesSince int64
	metaLen := make([]byte, 1)

	audioPlayer := NewAudioPlayer()
	err = audioPlayer.Init()
	if err != nil {
		slog.Error("Error initializing audio player", "error", err)
		return err
	}

	for {

		audioLeft := icyMetaInt
		for audioLeft > 0 {
			chunk := len(buf)
			if audioLeft < chunk {
				chunk = audioLeft
			}

			n, err := io.ReadFull(reader, buf[:chunk])
			err = audioPlayer.Write(buf[:n])

			if err != nil {
				slog.Error("Error writing to mpv stdin", "error", err)
				return err
			}

			if err != nil {
				if err == io.EOF {
					slog.Info("Radio stream ended")
					return nil
				}

				slog.Error("Error reading from radio stream", "error", err)
				return err
			}

			audioLeft -= n
		}

		_, err := io.ReadFull(reader, metaLen)
		if err != nil {
			if err == io.EOF {
				slog.Info("Radio stream ended")
				return nil
			}

			slog.Error("Error reading from radio stream", "error", err)
			return err
		}
		metaData := make([]byte, int(metaLen[0])*16)
		if metaLen[0] > 0 {
			_, err := io.ReadFull(reader, metaData)
			if err != nil {
				if err == io.EOF {
					slog.Info("Radio stream ended")
					return nil
				}

				slog.Error("Error reading from radio stream", "error", err)
				return err
			}
			slog.Debug("Received raw metadata", "metadata", string(metaData))
			metadata := ParseMetadata(string(metaData))
			slog.Info("Received metadata", "metadata", metadata)
		}

		total += int64(icyMetaInt)
		bytesSince += int64(icyMetaInt)
		if time.Since(lastPrint) > 10*time.Second {
			slog.Debug("Received audio data", "total_bytes", total, "bytes_since_last_print", bytesSince)
			lastPrint = time.Now()
			bytesSince = 0
		}
	}

}
