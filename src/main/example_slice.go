package main

import (
	"fmt"
)

func main() {
	slice := make([]interface{}, 5)
	slice[0] = 123
	slice[1] = 231
	slice[2] = 321
	fmt.Println(slice[0])
	fmt.Println(slice[1])
	fmt.Println(slice[2])
}
