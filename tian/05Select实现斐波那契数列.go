package main

import (
	"fmt"
)

func fbnq(ch chan float64, quit chan bool) (flag bool) {
	var x float64 = 1
	var y float64 = 1
	for {
		select {
		case ch <- x:
			x, y = y, x+y
		case flag = <-quit:
			return
		}
	}
}
func main() {
	ch := make(chan float64)
	quit := make(chan bool)
	go func() {
		for i := 1; i < 7777; i++ {
			num := <-ch
			fmt.Println(num)
		}
		quit <- true
	}()
	fbnq(ch, quit)
}
