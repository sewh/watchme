watchme
=======

watchme is a simple replacement for the [acme watch command](https://github.com/9fans/go/tree/master/acme/Watch) that is cross-platform (thanks to fsnotify), a bit simpler in implementation, and probably much more naive.

This is heavily based on the original watch source code, which can be found at the link above.

## Installation

```bash
go get github.com/sewhs/watchme
go install github.com/sewhs/watchme

# Ensure $(go env GOPATH)/bin is in your path
```

## Usage

With acme running: `watchme command [arguments]`

This will watch the current working directory for any changes and re-run `command [arguments]` in an acme window when a file gets modified or created.