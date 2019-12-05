package main

import (
	"fmt"
	"time"
)

func t1() {
	fmt.Println("编译完成")
	<-time.NewTimer(time.Second).C
	fmt.Println("aaaa")
	t2()
}
func t2() {
	<-time.After(time.Second)
	fmt.Println("aaaa")
	t3()
}
func t3() {
	time.Sleep(time.Second)
	fmt.Println("aaaa")
}
func main() {
	t1()
}
