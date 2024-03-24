package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"time"
)

func main() {
	defer func() {
		removeNode()
	}()

	initNode()
	listen()
}

func checkHealth() {
	type HealthResponse struct {
		Status bool `json:"status"`
	}
	for {
		health := HealthResponse{}
		response, err := makeRESTRequest(fmt.Sprintf("http://localhost:8080/health"), "GET", []byte(""))
		err = json.Unmarshal(response, &health)
		if err != nil {
			slog.Error(fmt.Sprintf("error unmarshalling: %v", err))
			go initNode()
			break
		}
		time.Sleep(10 * time.Second)
	}

}
