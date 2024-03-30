package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
)

type Response struct {
	Resp string `json:"out"`
	Host string `json:"host"`
}

type Task struct {
	MasterNode string `json:"master_node"`
	Task       string `json:"task"`
}

var processingNodes = make(map[string][]byte)

// Функция для разбиения задачи на подзадачи и отправки на выполнение в отдельный сервис
func splitTask(master Master) {
	//defer close(resultChan)

	i := 0
	for node := range pool {
		task := master.GetNextTask()
		go func(task string, node Node) {
			t := Task{Task: task, MasterNode: fmt.Sprintf("localhost:%d", httpPort)}
			jsonData, err := json.Marshal(t)
			if err != nil {
				fmt.Println("Error marshalling JSON:", err)
				return
			}
			processingNodes[fmt.Sprintf("%s:%d", node.IP, node.Port)] = jsonData
			_, err = makeRESTRequest(fmt.Sprintf("http://%s:%d/execute", node.IP, node.Port), "POST", jsonData)
			if err != nil {
				slog.Error(err.Error())
				return
			}
		}(task, node)
		i++
	}
}

func receiveTask(w http.ResponseWriter, r *http.Request) {
	response := Response{}
	err := json.NewDecoder(r.Body).Decode(&response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	parts := strings.Split(response.Host, ":")

	// Ensure we have both host and port
	if len(parts) != 2 {
		fmt.Println("Invalid host:port format")
		return
	}

	// Extract host and port
	host := parts[0]
	port := parts[1]

	delete(processingNodes, fmt.Sprintf("%s:%s", host, port))

	slog.Info("Got response", "from", response, "value=", response.Resp)
	if ok := master.HandleResponse(response.Resp); ok {
		close(pool)
	} else {
		n, _ := strconv.Atoi(port)
		newNode := Node{
			IP:   host,
			Port: n,
		}
		go check(newNode)
	}
}
