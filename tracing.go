package main

import (
	"fmt"
)

func TracingMeasure() error {
	ch := make(chan int)

	go func() {
		ch <- 1
	}()

	<-ch
	fmt.Println("Measure function execution time")
	return nil
}
