package main

import (
	"fmt"
	"runtime"
	"time"
)

var i int
var flag int

func main00() {
	go func() {
		for {
			i += 1000
			fmt.Println("子协程=====", i)
			time.Sleep(time.Second)
			if i > 5000 {
				flag = 1
				return
			}
		}
	}()
	time.Sleep(time.Second / 100)
	for {
		if flag == 1 {
			return
		}
		i -= 1
		fmt.Println("主协程=====", i, "========")
		time.Sleep(time.Second)
	}
}

func main() {
	go func() {
		for i := 1; i < 10; i++ {
			runtime.Gosched()
			fmt.Println("hello")

		}
	}()
	for j := 1; j < 10; j++ {
		runtime.Gosched()
		fmt.Println("go")

	}
}
