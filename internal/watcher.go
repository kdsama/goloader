package internal

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

type Watcher struct {
	directory  string
	command    string
	w          *fsnotify.Watcher
	cmd        *exec.Cmd
	lastUpdate time.Time
	ctx        context.Context
	cancel     context.CancelFunc
}

func NewWatcher(directory string, cmd string) *Watcher {

	w, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	wr := Watcher{
		directory: directory,
		w:         w,
		command:   cmd,
	}

	return &wr
}

func (wr *Watcher) Init() {

	if err := wr.w.Add(wr.directory); err != nil {
		log.Println("Error adding watch:", err)
	}

}
func (wr *Watcher) Start() {
	wr.startCommand()
	// Start an event loop to handle events
	for {
		select {
		case event, ok := <-wr.w.Events:
			if !ok {
				return
			}
			if event.Op == fsnotify.Write {

				// if event.Op
				f := strings.Split(event.Name, ".")

				if f[len(f)-1] == "go" {

					// if time.Since(wr.lastUpdate) > 1*time.Second {
					wr.startCommand()
					// }

				}

			}

		case err, ok := <-wr.w.Errors:
			if !ok {
				return
			}
			log.Printf("Error: %s\n", err)
		}
	}
}

func (wg *Watcher) startCommand() {
	cmdArgs := strings.Split(wg.command, " ")
	if wg.cmd != nil {
		wg.cmd.Process.Kill()

	}
	cmd := exec.Command("go", "build", "-o", "./ff", cmdArgs[0])
	cmd.Dir = wg.directory
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Start()

	wg.cmd = exec.Command("./ff")
	wg.cmd.Dir = wg.directory
	wg.cmd.Stdout = os.Stdout
	wg.cmd.Stderr = os.Stderr

	wg.lastUpdate = time.Now()

	wg.run()

}

func (wg *Watcher) run() {
	err := wg.cmd.Start()
	if err != nil {
		fmt.Println("Process Killed", err)
	}
}
