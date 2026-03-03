package radio

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

const version = "0.0.1"

type Status struct {
	Metadata Metadata `json:"metadata"`
	playing  bool     `json:"playing"`
}

type Radio struct {
	mu       sync.RWMutex
	url      string
	metadata Metadata
	playing  bool
}

func (r *Radio) GetMetadata() Metadata {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.metadata
}

func NewRadio(url string) *Radio {
	return &Radio{
		url: url,
	}
}

func (r *Radio) Play(player string) error {
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

	src := bufio.NewReaderSize(resp.Body, 64*1024)

	reader, err := NewReader(src, icyMetaInt, func(metadata Metadata, raw string) {
		slog.Info("Received metadata", "metadata", metadata)
		r.mu.Lock()
		r.metadata = metadata
		r.mu.Unlock()
	})
	if err != nil {
		return err
	}

	if player == "mpv" {
		go NewMpvAudioPlayer().Play(reader)
	} else {
		go NewNativeAudioPlayer().Play(reader)
	}

	if err := RunServer(r, func(conn net.Conn, radio *Radio) {
		defer conn.Close()
		br := bufio.NewReader(conn)
		line, err := br.ReadString('\n')
		if err != nil {
			slog.Error("Failed to read command", "err", err)
			return
		}
		line = strings.TrimSpace(strings.ToUpper(line))

		resp := &Status{
			Metadata: radio.GetMetadata(),
			playing:  radio.playing,
		}
		switch line {
		case "METADATA":
			_ = json.NewEncoder(conn).Encode(resp)
		default:
			fmt.Fprintf(conn, "Unknown command: %s\n", line)
		}

	}); err != nil {
		return err
	}
	return nil
}
