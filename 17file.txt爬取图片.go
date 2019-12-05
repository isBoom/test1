package main

import (
	"fmt"
	"mychromedp"
	"os"
	"regexp"
	"strconv"
	"sync"
	"time"
)

var (
	wg    sync.WaitGroup
	count int
)

func download(url string, count int) {
	wg.Add(1)
	err := mychromedp.Getfile(url, "img/"+strconv.Itoa(count)+url[len(url)-4:])
	if err != nil {
		fmt.Println("Getfile:", err)
		wg.Done()
		return
	}
	fmt.Println(url + "  img" + strconv.Itoa(count) + " over")
	wg.Done()
}

func main() {
	file, err_f := os.Open("file.txt")
	if err_f != nil {
		fmt.Println("os.open:", err_f)
		return
	}
	defer file.Close()

	buf := make([]byte, 4096)
	str := ""
	for {
		n, _ := file.Read(buf)
		if n == 0 {
			break
		}
		str += string(buf[:n])
	}

	res := regexp.MustCompile(`https?://.*?(\;|\.((jpg)|(png)|(gif)))`)
	res_url := res.FindAllStringSubmatch(str, -1)
	os.Mkdir("img", 777)
	for _, data := range res_url {
		if data[0][len(data[0])-1] != byte(';') {
			go download(data[0], count)
			count++
		}
	}
	time.Sleep(2 * time.Second)
	wg.Wait()
}
