package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	i := 0
	for {
		// Simple counter
		fmt.Println(i)
		fmt.Fprintf(os.Stderr, "on stderr: %d", i)
		time.Sleep(time.Second)
		i++
	}
}
