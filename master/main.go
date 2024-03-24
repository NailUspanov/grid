package main

import (
	"time"
)

func main() {
	go checkHealth()

	go listen()
	time.Sleep(3 * time.Second)

	var task Master
	task = NewTaskMaster()
	splitTask(task)

}
