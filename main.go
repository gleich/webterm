package main

import (
	"bufio"
	"net/http"
	"os/exec"

	"github.com/gorilla/websocket"
	"go.mattglei.ch/timber"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		timber.Error(err, "failed to upgrade connection")
		return
	}
	defer conn.Close()

	cmd := exec.Command("bash", "-c", "while true; do date; sleep 1; done")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		timber.Error(err, "failed to get stdout")
	}

	err = cmd.Start()
	if err != nil {
		timber.Error(err, "failed to start command")
		return
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		if err := conn.WriteMessage(websocket.TextMessage, []byte(line)); err != nil {
			timber.Error(err, "failed to write to websocket")
			break
		}
	}

	err = scanner.Err()
	if err != nil {
		timber.Error(err, "failed to write to websocket")
	}

	cmd.Wait()
}

func main() {
	http.HandleFunc("/ws", handler)
	timber.Info("starting server")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		timber.Fatal(err, "failed to start server")
	}
}
