package multithreading

import (
	"fmt"
	"sync"
	"time"
)

var cakes []int
var pancakeUsed = 0
var mutex sync.Mutex

func MakingCakes(id int, data chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for cakeID := range data {
		time.Sleep(3 * time.Second)
		mutex.Lock()
		fmt.Printf("Worker %d: Created pancake %d\n", id+1, cakeID+1)
		cakes = append(cakes, cakeID+1)
		fmt.Println(cakes)
		mutex.Unlock()
	}
}

func ServingCakes(id int, data chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for cakeID := range data {
		mutex.Lock()
		pancakeUsed++
		fmt.Printf("Guest %d: Used pancake %d\n", id+1, cakeID)
		fmt.Println("Pancake used:", pancakeUsed)
		mutex.Unlock()
		time.Sleep(1 * time.Second)
	}
}

func Func(data chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for pancakeUsed < 10 {
		mutex.Lock()
		if len(cakes) == 0 {
			fmt.Printf("Waiting for more pancakes :(\n")
			mutex.Unlock()
			time.Sleep(time.Second)
		} else {
			cakeID := cakes[len(cakes)-1]
			cakes = cakes[:len(cakes)-1]
			mutex.Unlock()
			data <- cakeID
		}
	}
	close(data)
}

func Test() {
	var wg sync.WaitGroup
	channel := make(chan int)
	channel2 := make(chan int)

	wg.Add(1)
	go Func(channel2, &wg)

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go MakingCakes(i, channel, &wg)
	}

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go ServingCakes(i, channel2, &wg)
	}

	for i := 0; i < 10; i++ {
		channel <- i
	}

	close(channel)
	wg.Wait()
}
