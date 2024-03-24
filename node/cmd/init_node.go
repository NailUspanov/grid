package main

import (
	"fmt"
	"log/slog"
	"time"
)

func initNode() {
	body := fmt.Sprintf(`
		{
  			"ip": "%s",
			"port": %d
		}
	`, "localhost", httpPort)

	for {
		resp, err := makeRESTRequest("http://localhost:8080/register", "POST", []byte(body))
		fmt.Println(string(resp))
		if err != nil {
			slog.Error(err.Error())
			time.Sleep(5 * time.Second)
			continue
		}
		go checkHealth()
		break
	}

}

func removeNode() {
	body := fmt.Sprintf(`
		{
  			"ip": "%s",
			"port": %d
		}
	`, "localhost", httpPort)

	resp, err := makeRESTRequest("http://localhost:8080/remove", "POST", []byte(body))
	fmt.Println(string(resp))
	if err != nil {
		slog.Error(err.Error())
		return
	}
}
