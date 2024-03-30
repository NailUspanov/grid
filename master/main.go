package main

import (
	"time"
)

var master = NewTaskMaster()

func main() {
	go checkHealth()

	go listen()
	time.Sleep(3 * time.Second)
	splitTask(master)
}
