watchme
=======

watchme is a tool to run a shell command inside of an [acme](https://en.wikipedia.org/wiki/Acme_(text_editor)) window whenever a file has been modified.

It is heavily inspired (and based upon) the original [acme watch command](https://github.com/9fans/go/tree/master/acme/Watch). Unlike the original watch command, watchme is cross-platform and forces the user to choose what files they are interested in monitoring with a [glob](https://en.wikipedia.org/wiki/Glob_(programming)).

## Installation

```bash
go get github.com/sewhs/watchme
go install github.com/sewhs/watchme

# Ensure $(go env GOPATH)/bin is in your path
```

## Usage

With acme running: `watchme glob command [arguments]` e.g. `watchme '*.go' go vet` to run the `go vet` tool whenever a '.go' file changes.

watchme will automatically refresh the list of globbed files, so you don't need to re-run watchme each time you add a new file.