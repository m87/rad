# rad

`rad` is a lightweight terminal radio player written in Go.

It can:
- play an internet radio stream directly from URL,
- save stations under aliases,
- play by alias,
- expose current track metadata through a local UNIX socket.

## Features

- Direct playback from stream URL
- Alias-based station management in a YAML config file
- Two audio backends:
	- `native` (built-in Go audio path)
	- `mpv` (external `mpv` process)
- ICY metadata parsing (`Artist - Title`)
- `status` command to read current metadata from the running player

## Requirements

- Go `1.25+`
- Linux/macOS environment (UNIX socket support)
- Optional: `mpv` installed if you want to use `--player mpv`

## Installation

### From source

```bash
git clone https://github.com/m87/rad.git
cd rad
go build -o rad .
```

### Install with Go

```bash
go install github.com/m87/rad@latest
```

## Quick start

Play directly from URL:

```bash
rad https://example.com/stream.mp3
```

Play with `mpv` backend:

```bash
rad --player mpv https://example.com/stream.mp3
```

## Commands

### Play station

You can play:
- a full URL,
- an alias prefixed with `@`, e.g. `@jazz`,
- an alias without `@` if it exists in config.

Examples:

```bash
rad @jazz
rad jazz
rad https://example.com/live
```

### Add station alias

```bash
rad add --alias jazz --url https://example.com/jazz.mp3
```

Short flags:

```bash
rad add -a jazz -u https://example.com/jazz.mp3
```

### Show status / metadata

When a station is currently running in `rad`, you can read metadata:

```bash
rad status
```

Expected output is JSON, e.g.:

```json
{"metadata":{"Title":"Song Title","Artist":"Artist Name"}}
```

## Configuration

Default config file:

```text
$HOME/.rad.yaml
```

You can provide custom config path with:

```bash
rad --config /path/to/config.yaml @jazz
```

Example config:

```yaml
stations:
	jazz: "https://example.com/jazz.mp3"
	rock: "https://example.com/rock"
```

## Internal socket path

Runtime status socket location:

```text
$HOME/.local/state/rad/rad.sock
```

`rad status` connects to this socket and requests metadata.

## Troubleshooting

- `Error playing radio`:
	- verify stream URL is reachable,
	- check that the stream serves MP3/ICY-compatible data.
- No output in `rad status`:
	- ensure another `rad` process is currently playing,
	- ensure socket file exists under `$HOME/.local/state/rad/`.
- `--player mpv` fails:
	- install `mpv` and confirm it is available in `PATH`.
