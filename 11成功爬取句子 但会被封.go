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
	url := "https://www.zuijuzi.com/article/" + strconv.Itoa(page)
	fmt.Println(url)
	result, err0 := httpget(url)
	if err0 != nil {
		fmt.Println("httpget", err0)
		return
	}
	//获取标题并创建相应文件
	res_title := regexp.MustCompile(`<title>(?s:(.*?))</title>`)
	//匹配规则
	res_http := regexp.MustCompile(`<div class="content"><a href="(?s:(.*?))">`)
	//获取名称
	text_name := strings.TrimSpace(res_title.FindAllStringSubmatch(result, -1)[0][1]) + ".txt"
	file, _ := os.Create(text_name)
	defer file.Close()
	if res_http == nil {
		fmt.Println("regexp.MustCompile err")
		return
	}

	for _, date := range res_http.FindAllStringSubmatch(result, -1) {
		result_2, err2 := httpget(date[1])
		if err2 != nil {
			fmt.Println("tttpget2", err2)
			continue
		}
		res_txt := regexp.MustCompile(`<div class="content " >(?s:(.*?))</div>`)
		str := strings.Replace(res_txt.FindAllStringSubmatch(result_2, -1)[0][1], "<br/>", "", -1)
		str = strings.Replace(str, "</B>", "", -1)
		file.Write([]byte(str + "\n\n"))
	}
}
func main() {
	start, end := 1, 69094
	for i := start; i <= end; i++ {
		Dowork(i)
	}
}
