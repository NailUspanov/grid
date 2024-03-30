package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

// Node представляет собой структуру данных для узла
type Node struct {
	IP   string `json:"ip"`
	Port int    `json:"port"`
}

var nodes = make(map[string]Node)
var nodesConcurrent = NewConcurrentMap(nodes)
var offNodes = make(map[string]struct{})
var pool = make(chan Node)

func checkHealth() {
	type HealthResponse struct {
		Status bool `json:"status"`
	}
	for {
		for url, _ := range nodesConcurrent.m {
			health := HealthResponse{}
			response, err := makeRESTRequest(fmt.Sprintf("http://%s/health", url), "GET", []byte(""))
			err = json.Unmarshal(response, &health)
			if !health.Status || err != nil {
				delete(nodes, url)
				offNodes[url] = struct{}{}

				go func(url string) {
					if b, ok := processingNodes[url]; ok {
						n := <-pool
						delete(processingNodes, url)
						processingNodes[fmt.Sprintf("%s:%d", n.IP, n.Port)] = b
						slog.Info(fmt.Sprintf("Node %s down. Resend task to %s:%d", url, n.IP, n.Port))
						_, err = makeRESTRequest(fmt.Sprintf("http://%s:%d/execute", n.IP, n.Port), "POST", b)
						if err != nil {
							slog.Error(err.Error())
							return
						}
					}
				}(url)

			}
		}
		time.Sleep(10 * time.Second)
	}

}

// registerNode обрабатывает запрос POST для регистрации узла
func registerNode(w http.ResponseWriter, r *http.Request) {
	var newNode Node
	err := json.NewDecoder(r.Body).Decode(&newNode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	nodeString := fmt.Sprintf("%s:%d", newNode.IP, newNode.Port)
	nodesConcurrent.Set(nodeString, newNode)
	delete(offNodes, nodeString)
	pool <- newNode
	fmt.Printf("Узел зарегистрирован: %+v\n", newNode)

	w.WriteHeader(http.StatusCreated)
}

func check(newNode Node) {
	pool <- newNode
}

// getNodes обрабатывает запрос GET для получения списка активных узлов
func getNodes(w http.ResponseWriter, r *http.Request) {
	nodesList := make([]Node, 0, len(nodes))
	for _, node := range nodes {
		nodesList = append(nodesList, node)
	}

	jsonNodes, err := json.Marshal(nodesList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonNodes)
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	type HealthResponse struct {
		Status bool `json:"status"`
	}

	resp := HealthResponse{Status: true}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

}

func removeNode(w http.ResponseWriter, r *http.Request) {
	var node Node
	err := json.NewDecoder(r.Body).Decode(&node)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	nodeString := fmt.Sprintf("%s:%d", node.IP, node.Port)
	delete(nodes, nodeString)
	offNodes[nodeString] = struct{}{}
	fmt.Printf("Узел удален: %+v\n", node)

	w.WriteHeader(http.StatusOK)
}
