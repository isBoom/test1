package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

func prc(count int, c string) string {
	str := ""
	for i := 0; i < count; i++ {
		str += c
	}
	return str
}
func main() {
	//要下载的文件名
	filename := ""
	//从哪下载文件
	start := 0
	var sizeFile *os.File
	fmt.Println("请输入要下载的文件名:")
	fmt.Scanf("%s", &filename)
	//看是否已经下载到中途的过文件
	_, err := os.Stat("." + filename)
	if err == nil {
		//文件下载了部分
		sizeFile, err = os.OpenFile("."+filename, os.O_RDWR, 0660)
		temp, err := ioutil.ReadAll(sizeFile)
		if err != nil {
			fmt.Println("ioutil.ReadAll err", err)
			return
		}
		//文件无法表示进度 置0重新下
		start, err = strconv.Atoi(string(temp))
		if err != nil {
			start = 0
		}
		defer sizeFile.Close()
		sizeFile.Seek(0, 0)
	} else {
		//没下过从新下
		sizeFile, err = os.OpenFile("."+filename, os.O_RDWR|os.O_CREATE, 0660)
	}

	req, err := http.NewRequest("get", "http://xxxholic.top:8077/d", nil)
	if err != nil {
		fmt.Println("http.NewRequest err", err)
		return
	}
	req.Header.Set("Download", filename)
	req.Header.Set("Start", fmt.Sprint(start))
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("client.Do err", err)
		return
	}
	defer resp.Body.Close()
	fmt.Println(resp.Header)
	if resp.Header["Status"] != nil && resp.Header["Status"][0] == "404" {
		fmt.Println("请求的文件不存在")
		sizeFile.Close()
		os.Remove("." + filename)
		return
	}

	//文件大小
	max, _ := strconv.Atoi(resp.Header["Filesize"][0])

	//校验文件 此处应该用md5校验完整信息而不是文件大小 没写
	nowFile, err := os.Stat(filename)
	if err == nil && int(nowFile.Size()) == max {
		fmt.Println("已存在完整文件")
		sizeFile.Close()
		os.Remove("." + filename)
		return
	}

	if resp.StatusCode != 200 {
		fmt.Println("失败")
	} else {
		buf := make([]byte, 4096)
		file, err := os.OpenFile("clannad.mp4", os.O_APPEND|os.O_RDWR|os.O_CREATE, 0660)
		if err != nil {
			fmt.Println("os.OpenFile err", err)
			return
		}
		defer file.Close()
		for {
			n, _ := resp.Body.Read(buf)
			if n != 0 {
				start += n
				//存进度信息
				sizeFile.Write([]byte(strconv.Itoa(start)))
				sizeFile.Seek(0, 0)
				//进度条
				fen := float64(start) * 100.0 / float64(max)
				fmt.Printf("\r[%s][%.1f%%][%.1f/%.1fMb]", prc(int(fen+0.5)/2, "=")+prc(50-int(fen+0.5)/2, " "), fen, float64(start)/(1024*1024), float64(max)/(1024*1024))
				file.Write(buf[:n])
			} else {
				break
			}
		}
		fmt.Printf("\nover\n")
		sizeFile.Close()
		file.Close()
		os.Remove(filename)
		os.Remove("." + filename)
	}

}
