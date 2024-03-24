package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func listen() {
	nodes = make(map[string]Node)

	http.HandleFunc("/register", registerNode)
	http.HandleFunc("/remove", removeNode)
	http.HandleFunc("/nodes", getNodes)
	http.HandleFunc("/health", healthCheckHandler)

	fmt.Println("Сервер запущен на порте 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
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
