package multithreading

import (
	"fmt"
	"time"
)

func Task(id int, data chan int) {
	for taskId := range data {
		if id == 0 {
			fmt.Printf("Worker %d: Task %d\n", id, taskId)
		}
		time.Sleep(time.Second)
	}
}

func Test() {
	channel := make(chan int)
	for i := 0; i < 3; i++ {
		go Task(i, channel)
	}
	for i := 0; i < 10; i++ {
		channel <- i
	}
	close(channel)
}
