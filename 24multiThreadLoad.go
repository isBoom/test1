package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	file       *os.File
	mutex      sync.Mutex
	wg         sync.WaitGroup
	url        string
	filename   string
	ThreadsNum int
	maxSize    int
)

func resStr(m int, c string) string {
	str := ""
	for i := 0; i < m; i++ {
		str += c
	}
	return str
}

func download(s, i, index int) {
	wg.Add(1)
	p, err := http.NewRequest("get", url, nil)
	if err != nil {
		fmt.Println("newrequest err", err)
		wg.Done()
		return
	}
	Range := fmt.Sprintf("bytes=%d-%d", s, i)
	p.Header.Set("Range", Range)
	Client := &http.Client{}
	resp, err := Client.Do(p)
	if resp.StatusCode != 206 {
		fmt.Printf("[%d]-[%d]下载失败", s, i)
		wg.Done()
		return
	}
	f, err := os.Create(fmt.Sprintf(".temp/.%d%s", index, filename))
	if err != nil {
		fmt.Println(" os.Create err", err)
		wg.Done()
		return
	}
	defer f.Close()
	defer resp.Body.Close()
	buf := make([]byte, 4096)
	for {
		n, _ := resp.Body.Read(buf)
		if n == 0 {
			break
		} else {
			f.Write(buf[:n])
		}
	}
	wg.Done()
}

func singleThreadLoad(resp *http.Response) {
	fmt.Println("此文件不支持多线程下载,正在以单线程模式下载")
	var err error
	file, err = os.Create(filename)
	if err != nil {
		fmt.Println("create err", err)
		return
	} else {
		buf := make([]byte, 4096)
		for {
			n, _ := resp.Body.Read(buf)
			if n == 0 {
				break
			} else {
				file.Write(buf[:n])
			}
		}
		file.Close()
	}
	os.RemoveAll(".temp")
}
func multiThreadLoad() {
	fmt.Println("开始多线程下载")
	os.Mkdir(".temp", 0660)
	blockSize := maxSize / ThreadsNum
	time.Sleep(time.Second / 10)
	for i := 0; i < ThreadsNum; i++ {
		if i == ThreadsNum-1 {
			go download(i*blockSize, maxSize, i)

		} else {
			go download(i*blockSize, (i+1)*blockSize-1, i)
		}
	}
	wg.Wait()
	fmt.Println("下载完毕正在合并")
	file, _ = os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0660)
	for i := 0; i < ThreadsNum; i++ {
		f, err := os.Open(fmt.Sprintf(".temp/.%d%s", i, filename))
		if err != nil {
			fmt.Println(err)
			return
		}
		buf := make([]byte, 4096)
		for {
			n, _ := f.Read(buf)
			if n == 0 {
				break
			}
			file.Write(buf[:n])
		}
		f.Close()
	}
	file.Close()
	os.RemoveAll(".temp")
}
func main() {
	fmt.Println("请输入url,线程数")
	fmt.Scanf("%s %d", &url, &ThreadsNum)
	if tempIndex := strings.LastIndex(url, "/"); tempIndex == -1 {
		fmt.Println("文件名不合法")
		return
	} else {
		filename = url[tempIndex+1 : len(url)]
	}

	r, err := http.NewRequest("get", url, nil)
	r.Header.Set("Range", "bytes=0-")
	r.Header.Set("Referer", url)
	if err != nil {
		fmt.Println("request err", err)
		return
	}
	client := &http.Client{}
	resp, err := client.Do(r)
	defer resp.Body.Close()
	fmt.Println(resp.StatusCode)
	if resp.StatusCode == 200 {
		singleThreadLoad(resp)
	} else if resp.StatusCode == 206 {
		//支持多线程下载
		maxSize, _ = strconv.Atoi(resp.Header["Content-Length"][0])
		multiThreadLoad()
	} else {
		fmt.Println("请求失败", resp.StatusCode)
		return
	}
	fmt.Println("下载完成")
}
