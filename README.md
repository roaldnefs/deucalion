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
      --config string   config file (default is $HOME/.deucalion.yaml)
  -f, --firing string   Command to execute when alerts are firing
  -h, --help            help for deucalion
  -s, --silent string   Command to execute when alerts aren't firing
  -u, --url string      Promtheus URL (default "http://localhost:9090/")
```
