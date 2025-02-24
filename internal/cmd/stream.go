package cmd

import (
	"bufio"
	"errors"
	"net/http"
	"os/exec"
	"syscall"

	"github.com/gorilla/websocket"
	"go.mattglei.ch/timber"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func Stream(command *exec.Cmd) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			timber.Error(err, "failed to upgrade connection")
			return
		}
		defer conn.Close()

		stdout, err := command.StdoutPipe()
		if err != nil {
			timber.Error(err, "failed to get stdout")
		}

		err = command.Start()
		if err != nil {
			timber.Error(err, "failed to start command")
			return
		}

		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			err := conn.WriteMessage(websocket.TextMessage, []byte(line))
			if err != nil && !errors.Is(err, syscall.EPIPE) {
				timber.Error(err, "failed to write to websocket")
				break
			}
		}

		err = scanner.Err()
		if err != nil {
			timber.Error(err, "failed to write to websocket")
		}

		command.Wait()
	}
}
