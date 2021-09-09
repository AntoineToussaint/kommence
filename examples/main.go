package main

import (
	"fmt"
	"time"
)

func main() {
	i := 0
	for {
		fmt.Println(i)
		time.Sleep(time.Second)
		i++
	}
}
