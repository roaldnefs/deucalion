[![Travis CI](https://img.shields.io/travis/roaldnefs/deucalion.svg?style=for-the-badge)](https://travis-ci.org/roaldnefs/deucalion)
[![GoDoc](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=for-the-badge)](https://godoc.org/github.com/roaldnefs/deucalion)
[![Github All Releases](https://img.shields.io/github/downloads/roaldnefs/deucalion/total.svg?style=for-the-badge)](https://github.com/roaldnefs/deucalion/releases)

Named after the son of Prometheus. A tool for querying a Prometheus instance for alerts and run commands based on the result.

* [Installation](README.md#installation)
     * [Binaries](README.md#binaries)
     * [Via Go](README.md#via-go)
* [Usage](README.md#usage)

## Installation

### Binaries

For installation instructions from binaries please visit the [Releases Page](https://github.com/roaldnefs/deucalion/releases).

### Via Go

```console
$ go get github.com/roaldnefs/deucalion
```

## Usage

```console
$ deucalion -h
A tool for querying a Prometheus instance for alerts and run commands based on the result.

Usage:
  deucalion [flags]

Flags:
      --config string     config file (default is $HOME/.deucalion.yaml)
  -d, --debug             enable debug logging
  -f, --firing string     command to execute when alerts are firing
  -h, --help              help for deucalion
      --severity string   the severity label
  -s, --silent string     command to execute when alerts aren't firing
  -u, --url string        Promtheus URL (default "http://localhost:9090/")
  -w, --warning string    command to execute when alerts are firing and below specified severity
```

## Configuration

The configuration file (`--config`) is written as follows:

```yaml
---
url: <URL>
firing: <COMMAND>
warning: <COMMAND>
silent: <COMMAND>
severity: <LABELVALUE>
```

## Scheduled run

We recommend running the service and timer as your own user, and therefore suggest you run them in systemd user mode. To achieve this, please do the following:

### Prerequisites

```console
mkdir -p ~/.local/share/systemd/user
sudo loginctl enable-linger $USER
```

Create the required files:
`~/.local/share/systemd/user/deucalion.service`:

```console
[Unit]
Description=Deucalion run
After=syslog.target network-online.target
Wants=network-online.target

[Service]
Type=simple
ExecStart=/usr/local/bin/deucalion
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

`~/.local/share/systemd/user/deucalion.timer`:

```console
[Unit]
Description=Run Deucalion every minute

[Timer]
Persistent=false
OnBootSec=80
OnCalendar=minutely
Unit=deucalion.service

[Install]
WantedBy=timers.target
```

This requires you have `deucalion` installed in `/usr/local/bin/deucalion`. But obviously, feel free to change the path. The timer will run the command every minute, using the user given in the commands below.

To enable and start the timer, please run the following commands:

```console
systemctl --user daemon-reload
systemctl --user enable --now deucalion.timer
```
