package main

import (
	"fmt"
)

func main() {
	slice := make([]byte, 5)
	slice[0] = 123
	slice[1] = 231
	slice[2] = 321
	fmt.Println(slice[0])
	fmt.Println(slice[1])
	fmt.Println(slice[2])

	slice1 := append(slice, 1, 2, 3, 4, 5)
	for
}
