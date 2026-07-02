# Install

Install s1ctl via `go install`, from source, or with a pre-built binary for Linux, macOS, or Windows (amd64 and arm64).

## Go install

```bash
go install danny.vn/s1/cmd/s1ctl@latest
```

## From source

```bash
git clone https://github.com/dannyota/s1ctl.git
cd s1ctl
go build -o s1ctl ./cmd/s1ctl
```

## Releases

Pre-built binaries for Windows, macOS, and Linux (amd64 and arm64) are
available on the
[GitHub releases page](https://github.com/dannyota/s1ctl/releases).

Download, extract, and move to a directory in your `PATH`:

```bash
tar xzf s1ctl_linux_amd64.tar.gz
sudo mv s1ctl /usr/local/bin/
```

## FAQ

### What Go version do I need?

Go 1.26 or later. The module uses Go 1.26 language features.

### How do I update s1ctl?

Re-run `go install danny.vn/s1/cmd/s1ctl@latest` or download the latest
release binary.

### Does s1ctl support Apple Silicon?

Yes. Pre-built binaries are available for both `darwin/amd64` and
`darwin/arm64`.

### Where is the config file stored?

Default: `~/.s1ctl/config.yaml`. See [Configure](guides/configure.md) for
details.
