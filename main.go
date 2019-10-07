package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"9fans.net/go/acme"
	"github.com/fsnotify/fsnotify"
)

var args []string
var win *acme.Win

func main() {
	// Just pull out all arguments
	flag.Parse()
	args = flag.Args()
	if len(args) == 0 {
		fmt.Println("watchme cmd args...")
		os.Exit(0)
	}

	// Create an acme context
	localWin, err := acme.New()
	win = localWin
	if err != nil {
		fmt.Println("Couldn't open acme!")
		os.Exit(1)
	}

	// Get working directory and open a new acme window
	pwd, _ := os.Getwd()
	win.Name(pwd + "/+watchme")
	win.Ctl("clean")
	win.Fprintf("tag", "Get ")

	// Launch our acme event handler so we don't block events from being processed
	go acmeEventHandler()

	// Open up a new filesystem watcher and add the current directory to it
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("Could not open file system watcher: " + err.Error())
		os.Exit(1)
	}
	watcher.Add(pwd)

	// Start the main execution loop
	var prevCmd *exec.Cmd
	firstGo := true
	for {
		// Wait until we receive an event
		if firstGo {
			// We don't want to block waiting for file system events when running the first command
			firstGo = false
		} else {
			// If we receive a file modified or file created event, kick off running the command
			event, ok := <-watcher.Events
			interested := false
			if ok && (event.Op & fsnotify.Write == fsnotify.Write || event.Op & fsnotify.Create == fsnotify.Create) {
				interested = true
			}

			if !interested {
				continue
			}
		}

		// Kill a command hanging from last time
		if prevCmd != nil {
			prevCmd.Process.Kill()
		}
		prevCmd = nil

		// Fire off  the command and redirect output to our pipe
		cmd := exec.Command(args[0], args[1:]...)
		r, w, err := os.Pipe()
		if err != nil {
			fmt.Println("Could not open OS pipe to connect to process: " + err.Error())
			os.Exit(1)
		}
		win.Addr(",")
		win.Write("data", nil)
		win.Ctl("clean")
		win.Fprintf("body", "$ %s\n", strings.Join(args, " "))
		cmd.Stdout = w
		cmd.Stderr = w
		if err := cmd.Start(); err != nil {
			r.Close()
			w.Close()
			win.Fprintf("body", "%s: %s\n", strings.Join(args, " "), err.Error())
			continue
		}
		prevCmd = cmd
		w.Close()

		// Write the output from the command into the acme window
		buf := make([]byte, 4096)
		for {
			n, err := r.Read(buf)
			if err != nil {
				break
			}
			win.Write("body", buf[:n])
		}

		// Wait for the command to finish up
		if err := cmd.Wait(); err != nil {
			win.Fprintf("body", "%s: %s\n", strings.Join(args, " "), err.Error())
		}

		// Finish altering the acme window state
		win.Fprintf("body", "$\n")
		win.Fprintf("addr", "#0")
		win.Ctl("dot=addr")
		win.Ctl("show")
		win.Ctl("clean")

	}
}

func acmeEventHandler() {
	for e := range win.EventChan() {
		win.WriteEvent(e)
	}

	os.Exit(0)
}