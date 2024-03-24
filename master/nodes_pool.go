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
		for url, _ := range nodes {
			health := HealthResponse{}
			response, err := makeRESTRequest(fmt.Sprintf("http://%s/health", url), "GET", []byte(""))
			err = json.Unmarshal(response, &health)
			if err != nil {
				slog.Error(fmt.Sprintf("error unmarshalling: %v", err))
				continue
			}
			if !health.Status {
				delete(nodes, url)
				offNodes[url] = struct{}{}
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
