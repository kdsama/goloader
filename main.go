package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

var cmd *exec.Cmd

func watchSubdirectories(directory string, watcher *fsnotify.Watcher) {
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Println("Error walking directory:", err)
			return err
		}

		// Check if it's a subdirectory
		if info.IsDir() && path != directory {
			// Add the subdirectory to the watcher
			if err := watcher.Add(path); err != nil {
				log.Println("Error adding watch:", err)
			}
			// log.Printf("Monitoring directory: %s\n", path)
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}

func startCommand(path string, c string) {
	cmdArgs := strings.Split(c, " ")
	fmt.Println(cmdArgs[0], cmdArgs[1:])
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	cmd.Dir = path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
}

func restartCommand(path string, c string) {
	fmt.Println("Restarting", cmd, cmd)

	if cmd != nil && cmd.Process != nil {
		err := cmd.Process.Kill()
		if err != nil {
			log.Fatal(err)
		}

		// log.Println("Command stopped")

		// Start the command again
		fmt.Println(cmd.Args)
		startCommand(path, c)
	}
}

func main() {
	// if len(os.Args) < 4 {
	// 	// fmt.Println("Usage: go run main.go /path/to/directory command arg1 arg2 ...")
	// 	return
	// }

	directory := flag.String("d", "./", "Directory on which hot loader should run ")
	commandArgs := flag.String("cmd", "", "Go execution command")
	flag.Parse()

	// Create a new watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// Watch the initial directory
	if err := watcher.Add(*directory); err != nil {
		log.Fatal(err)
	}
	// log.Printf("Monitoring directory: %s\n", directory)

	// Watch subdirectories
	watchSubdirectories(*directory, watcher)

	// Start the initial command
	go startCommand(*directory, *commandArgs)

	// Start an event loop to handle events
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			// Restart the command when an event is caught
			// log.Printf("Event: %s\n", event)
			if event.Op == fsnotify.Write {
				// if event.Op
				f := strings.Split(event.Name, ".")

				if f[len(f)-1] == "go" {

					restartCommand(*directory, *commandArgs)
				}

			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Printf("Error: %s\n", err)
		}
	}
}
