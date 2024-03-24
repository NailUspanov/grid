package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
)

type Response struct {
	Resp string `json:"out"`
}

type Task struct {
	Min       int    `json:"min"`
	Max       int    `json:"max"`
	BatchSize int    `json:"batchSize"`
	Task      string `json:"task"`
}

// Функция для разбиения задачи на подзадачи и отправки на выполнение в отдельный сервис
func splitTask(master Master) {
	//defer close(resultChan)

	i := 0
	for node := range pool {
		task := master.GetNextTask()
		go func(task string, node Node) {
			t := Task{Task: task}
			jsonData, err := json.Marshal(t)
			if err != nil {
				fmt.Println("Error marshalling JSON:", err)
				return
			}
			resp, err := makeRESTRequest(fmt.Sprintf("http://%s:%d/execute", node.IP, node.Port), "POST", jsonData)
			if err != nil {
				slog.Error(err.Error())
				return
			}
			response := Response{}
			err = json.Unmarshal(resp, &response)
			if err != nil {
				slog.Error(err.Error())
				return
			}
			slog.Info("Got response", "from", node, "value=", response.Resp)
			if ok := master.HandleResponse(response.Resp); ok {
				close(pool)
			} else {
				pool <- node
			}
		}(task, node)
		i++
	}
}
