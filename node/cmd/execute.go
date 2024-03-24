package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type Request struct {
	Task string `json:"task"`
}

func execute(w http.ResponseWriter, r *http.Request) {
	var req Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Создание временного файла с кодом функции
	file, err := os.CreateTemp("", "temp*.go")
	if err != nil {
		fmt.Println("Ошибка создания временного файла:", err)
		return
	}
	defer os.Remove(file.Name())

	// Запись кода функции в файл
	_, err = file.WriteString(req.Task)
	if err != nil {
		fmt.Println("Ошибка записи в файл:", err)
		return
	}
	file.Close()

	// Компиляция и исполнение файла
	cmd := exec.Command("go", "run", file.Name())
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("Ошибка выполнения программы:", err)
		return
	}

	// Вывод результата
	fmt.Println("Результат выполнения функции:", string(out))

	resp := struct {
		Out string `json:"out"`
	}{
		Out: string(out),
	}
	jsonResponse, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func unescapeString(s string) string {
	// Создаем map для замены экранированных символов на нормальные
	replacements := map[string]string{
		`\\`: ` `, // Обратный слэш
		`\"`: `"`, // Кавычки
		`\n`: " ", // Символ новой строки
		// Добавьте другие специальные символы, если необходимо
	}

	// Заменяем экранированные символы на нормальные
	for escaped, normal := range replacements {
		s = strings.ReplaceAll(s, escaped, normal)
	}

	return s
}
