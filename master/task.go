package main

import (
	"log/slog"
	"strconv"
)

type Master interface {
	GetNextTask() string
	HandleResponse(out string) bool
}

var (
	asciiMin    = 32  // Минимальное значение ASCII символа
	asciiMax    = 126 // Максимальное значение ASCII символа
	taskInWork  = make(map[string]struct{})
	counterTask = 0
	batchSize   = 1
)

type TaskMaster struct {
}

func NewTaskMaster() Master {
	return &TaskMaster{}
}

func (t TaskMaster) GetNextTask() string {

	//start := asciiMin + counterTask*batchSize
	//end := start + batchSize - 1

	start := 50 + counterTask*batchSize
	end := start + batchSize
	//start := 48
	//end := 57
	if end > asciiMax {
		end = asciiMax
	}
	var task = `
	package main

import (
	"fmt"
	"github.com/yeka/zip"
	"io/ioutil"
	"log"
)

const (
	asciiMin = 32  // Минимальное значение ASCII символа
	asciiMax = 126 // Максимальное значение ASCII символа
)

// Функция для генерации уникального пароля заданной длины
func generateUniquePassword(n, length int, start, end int, resultChan chan string) {
	defer close(resultChan)

	// Функция для получения следующей комбинации пароля
	nextPassword := func(password []byte) {
		for i := len(password) - 1; i >= 0; i-- {
			if password[i] < byte(asciiMax) {
				password[i]++
				return
			}
			password[i] = byte(start)
		}
	}

	// Инициализируем пароль начальным значением
	password := make([]byte, length)
	for i := range password {
		password[i] = byte(start)
	}

	// Перебираем все комбинации пароля и отправляем их в канал результатов
	for {
		resultChan <- fmt.Sprintf("%s", string(password))

		// Проверяем, достигли ли конца диапазона комбинаций
		if password[0] == byte(end) {
			break
		}

		// Генерируем следующую комбинацию пароля
		nextPassword(password)
	}
}

func main() {
	length := 3 // Длина пароля

	resultChan := make(chan string)

	go generateUniquePassword(1, length,` + strconv.Itoa(start) + `, ` + strconv.Itoa(end) + `, resultChan)

	for password := range resultChan {
		if isValidArchivePassword(password) {
			fmt.Printf(password)
			return
		}
	}
	fmt.Printf("false")
}

func isValidArchivePassword(password string) bool {
	r, err := zip.OpenReader("archive.zip")
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	for _, f := range r.File {
		if f.IsEncrypted() {
			f.SetPassword(password)
		}

		r, err := f.Open()
		if err != nil {
			return false
		}

		_, err = ioutil.ReadAll(r)
		if err != nil {
			return false
		}
		defer r.Close()

		return true
	}
	return false
}
`
	counterTask++
	return task
}

func (t TaskMaster) HandleResponse(out string) bool {
	if out != "false" {
		slog.Info("SUCCESS!!!", "password=", out)
		return true
	}
	return false
}
