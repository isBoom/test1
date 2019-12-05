package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

//爬取网页主体
var file *os.File
var file1 *os.File

func httpget(url string) (result string, err error) {
	resp, err0 := http.Get(url)
	if err0 != nil {
		err = err0
	}
	defer resp.Body.Close()
	buf := make([]byte, 4096)
	for {
		n, _ := resp.Body.Read(buf)
		if n == 0 {
			break
		}
		result += string(buf[:n])
	}
	return
}

//核心工作函数
func Dowork(page int) {
	url := "https://www.zuijuzi.com/ju/" + strconv.Itoa(page)
	result, err0 := httpget(url)
	if err0 != nil {
		fmt.Println("httpget", err0)
		return
	}
	//获取标题并创建相应文件
	res_txt := regexp.MustCompile(`<div class="content " >(?s:(.*?))</div>`)
	str := strings.Replace(res_txt.FindAllStringSubmatch(result, -1)[0][1], "<br/>", "", -1)
	str = strings.Replace(str, "</B>", "", -1)
	file.Write([]byte(str + "\n\n"))
	os.Rename(strconv.Itoa(page-1), strconv.Itoa(page))
}
func main() {
	start, end := 306353, 857777

	file, _ = os.Create("file2.txt")
	file1, _ = os.Create("306352")

	file1.Write([]byte("从第_" + strconv.Itoa(start) + "_开始写入"))
	defer file.Close()
	file1.Close()

	for i := start; i <= end; i++ {
		Dowork(i)
	}
}
