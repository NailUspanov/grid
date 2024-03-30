package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
)

func listen() {
	l, close := createListener()
	defer close()
	log.Println("listening at", l.Addr().(*net.TCPAddr).Port)

	http.HandleFunc("/execute", execute)
	http.HandleFunc("/health", healthCheckHandler)

	//log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", httpPort), nil))
	log.Fatal(http.Serve(l, nil))
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	type HealthResponse struct {
		Status bool `json:"status"`
	}
	if shouldReturnHealthy() {
		resp := HealthResponse{Status: true}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	} else {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// shouldReturnHealthy returns true with a probability of 90%
func shouldReturnHealthy() bool {
	//rand.Seed(time.Now().UnixNano())
	//return rand.Intn(10) < 9
	return true
}

func makeRESTRequest(url, method string, requestBody []byte) ([]byte, error) {
	// Создание HTTP-клиента
	client := &http.Client{}

	// Создание запроса
	req, err := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	// Установка заголовков (если нужно)
	// req.Header.Set("Content-Type", "application/json")

	// Выполнение запроса
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Чтение тела ответа
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}

func createListener() (l net.Listener, close func()) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%s", strconv.Itoa(httpPort)))
	if err != nil {
		panic(err)
	}
	return l, func() {
		_ = l.Close()
	}
}
