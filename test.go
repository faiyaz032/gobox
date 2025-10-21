package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan string)

	go func() {
		time.Sleep(2 * time.Second)
		ch <- "Hello from goroutine!"
	}()

	fmt.Println("Waiting to receive...")
	msg := <-ch
	fmt.Println("Received:", msg)
}
