package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type hz struct { 
    name  rune     
	    count int
}
var hzmap map[rune]*hz = make(map[rune]*hz)

var sum int

func main() {
	file, _ := os.Open("file.txt") 
	defer file.Close()
	f := bufio.NewReader(file)
	for {
		buf, err := f.ReadBytes('\n')
		sum++
		if err != nil {
			if err == io.EOF {
				break
			} else {
				break
			}
		}
		str_rune := []rune(string(buf))
		for _, data := range str_rune {
			if int(data) != 10 && int(data) != 13 {

				if _, ok := hzmap[data]; ok == false {
					hzmap[data] = &hz{data, 1}
				} else {
					hzmap[data].count++
				}
			}

		}
	}
	i := 0
	for _, s := range hzmap {
		fmt.Printf("%c--->%d\n", s.name, s.count)
		i++
	}
	fmt.Printf("%d   %d", sum, i)
}
