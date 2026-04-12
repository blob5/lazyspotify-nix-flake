# lazyspotify

`lazyspotify` is a terminal Spotify client that uses a patched `go-librespot`
daemon for playback.

## Runtime requirements

- A Spotify Premium account.
- A patched `go-librespot` daemon build from
  `https://github.com/dubeyKartikay/go-librespot/`.
- A working system keyring.
  On Linux this is a hard requirement: `lazyspotify` will not fall back to
  plaintext token storage if the keyring is unavailable.
- Linux clipboard integration is optional.
  For the auth screen's copy shortcut, install one of `wl-clipboard`, `xclip`,
  or `xsel`.

## Daemon discovery

At startup, `lazyspotify` resolves the playback daemon in this order:

1. `librespot.daemon.cmd` from `config.yaml`
2. A packaged default daemon path compiled into the binary at build time

Not supported:

- `LAZYSPOTIFY_LIBRESPOT_DAEMON`
- `PATH` lookup
- probing next to the executable
- relocatable bundled daemon discovery

If you install `lazyspotify` without packaging and do not compile in a default
daemon path, you must set `librespot.daemon.cmd`.

## Installation

### Packaged install

When you package `lazyspotify`, ship both binaries together and compile the
daemon path into the `lazyspotify` binary:

- `lazyspotify`
- `lazyspotify-librespot`

Linux packages should install the daemon at:

```text
/usr/libexec/lazyspotify/lazyspotify-librespot
```

Build `lazyspotify` for Linux packages with:

```bash
go build -ldflags "-X github.com/dubeyKartikay/lazyspotify/core/utils.defaultLibrespotDaemonPath=/usr/libexec/lazyspotify/lazyspotify-librespot" -o target/lazyspotify ./cmd/lazyspotify/main.go
```

For Homebrew on macOS, install the daemon under the formula's `opt_libexec`
path and inject that stable absolute path with `-ldflags -X` during the build.

### Source build

Build `lazyspotify`:

```bash
go build -o target/lazyspotify ./cmd/lazyspotify/main.go
```

Build the patched daemon from the forked `go-librespot` repository and set an
explicit config override:

```yaml
librespot:
  daemon:
    cmd:
      - /absolute/path/to/lazyspotify-librespot
```

## Configuration

Config lives under the OS config directory:

- macOS: `~/Library/Application Support/lazyspotify/config.yaml`
- Linux: `~/.config/lazyspotify/config.yaml`

Only overrides are required. Package builds may provide a compiled default
daemon path.

Example:

```yaml
librespot:
  daemon:
    cmd:
      - /absolute/path/to/lazyspotify-librespot
```

## Development

Run the app:

```bash
make run
```

Build the app:

```bash
make build
```

Run tests:

```bash
go test ./...
```
