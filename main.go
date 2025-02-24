package main

import (
	"net/http"
	"os/exec"

	"go.mattglei.ch/timber"
	"go.mattglei.ch/webterm/internal/cmd"
)

func main() {
	command := exec.Command("bash", "-c", "while true; do date; sleep 1; done")

	http.HandleFunc("/ws", cmd.Stream(command))
	timber.Info("starting server")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		timber.Fatal(err, "failed to start server")
	}
}
