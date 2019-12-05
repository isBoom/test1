package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan int)
	go func() {
		for {
			select {
			case x := <-ch:
				fmt.Println("接收了一个X")
				fmt.Println(x)
			case <-time.NewTimer(3 * time.Second).C:
				fmt.Println("超时了")
				return

			}
		}
	}()

	for {
		r := 0
		fmt.Scan(&r)
		fmt.Println("传了一个R")
		ch <- r
	}
}

func main11() {

	ch := make(chan int)
	timer := time.NewTimer(3 * time.Second)
	go func() {
		<-timer.C
		fmt.Println("超时了")
	}()
	go func() {
		for {
			select {
			case x := <-ch:
				fmt.Println("接收了一个X")
				timer.Reset(3 * time.Second)
				fmt.Println(x)
			}
		}
	}()

	for {
		r := 0
		fmt.Scan(&r)
		fmt.Println("传了一个R")
		ch <- r
	}
}
