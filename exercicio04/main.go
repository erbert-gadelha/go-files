package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello, World!")

	for i := 0; i < 1024; i++ {
		fmt.Println(byte(i))
	}
}
