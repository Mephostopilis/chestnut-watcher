package main

import (
	"fmt"
)

func main() {
	ch1 := make(chan int, 1)
	ch2 := make(chan int, 1)

	// ch1 <- 1
	// select {
	// case e1 := <-ch1:
	// 	fmt.Println("1th case is selected. e1=%v", e1)
	// case e2 := <-ch2:
	// 	fmt.Println("2th case is selected. e2=%v", e2)
	// default:
	// 	fmt.Println("oh! no.")
	// }

	select {
	case ch1 <- 1:
		e1 := <-ch1
		fmt.Println("1th case is selected. %v", e1)
	case ch2 <- 2:
		e2 := <-ch2
		fmt.Println("2th case is selected. %v", e2)
	default:
		fmt.Println("hello.")
	}
}
