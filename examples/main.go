package main

import (
	"fmt"
	"time"
)

func main() {
	i := 0
	for {
		// Simple counter
		fmt.Println(i)
		time.Sleep(time.Second)
		i++
	}
}
