package multithreading

import (
	"fmt"
	"time"
)

var cakes []int
var pancakeUsed = 0

func MakingCakes(id int, data chan int) {
	for cakeID := range data {
		fmt.Printf("Worker %d: Created pancake %d\n", id+1, cakeID+1)
		cakes = append(cakes, cakeID+1)
		time.Sleep(time.Second)
	}
}

func ServingCakes(id int, data chan int) {
	if len(cakes) == 0 {
		fmt.Printf("Guest %d: Waiting for more pancakes :(\n", id+1)
		return
	}
	pancakeUsed++
	fmt.Printf("Guest %d: Used pancake %d\n", id+1, cakes[len(cakes)-1])
	cakes = cakes[:len(cakes)-1]
	time.Sleep(time.Second)
}

func Func(data chan int) {
	for pancakeUsed < 10 {
		for i := 0; i < 3; i++ {
			go ServingCakes(i, data)
			time.Sleep(50 * time.Millisecond)
		}
		time.Sleep(time.Second)
	}

}
func Test() {
	channel := make(chan int)
	go Func(channel)
	for i := 0; i < 3; i++ {
		go MakingCakes(i, channel)
	}
	for i := 0; i < 10; i++ {
		channel <- i
	}
	close(channel)
}
