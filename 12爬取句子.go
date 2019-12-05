package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	//"time"
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
	now_len_sentence := 0
	p_page := 1
	url := "https://www.zuijuzi.com/article/" + strconv.Itoa(page) + "#p"
	result, err0 := httpget(url + strconv.Itoa(p_page))
	if err0 != nil {
		fmt.Println("httpget", err0)
		return
	}
	//获取标题并创建相应文件
	res_title := regexp.MustCompile(`<title>(?s:(.*?))</title>`)
	//有多少句子
	res_len := regexp.MustCompile(`class="active">句子\((?s:(.*?))\)</a></li>`)
	//匹配规则
	res_http := regexp.MustCompile(`<div class="content"><a href="(?s:(.*?))">`)
	//获取名称
	text_name := strings.TrimSpace(res_title.FindAllStringSubmatch(result, -1)[0][1]) + ".txt"
	file, _ := os.Create(text_name)
	defer file.Close()

	//获取句长
	len_sentence, _ := strconv.Atoi(res_len.FindAllStringSubmatch(result, -1)[0][1])
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
		file.Write([]byte(res_txt.FindAllStringSubmatch(result_2, -1)[0][0] + "\n"))
		now_len_sentence++
	}
	if now_len_sentence < len_sentence {
		for {
			p_page++
			result, err0 := httpget(url + strconv.Itoa(p_page))
			if err0 != nil {
				fmt.Println("httpget", err0)
				return
			}
			for _, date := range res_http.FindAllStringSubmatch(result, -1) {
				result_2, err2 := httpget(date[1])
				if err2 != nil {
					fmt.Println("tttpget2", err2)
					continue
				}
				res_txt := regexp.MustCompile(`<div class="content " >(?s:(.*?))</div>`)
				file.Write([]byte(res_txt.FindAllStringSubmatch(result_2, -1)[0][1] + "\n\n"))
				now_len_sentence++
			}
			if now_len_sentence >= len_sentence {
				break
			}
		}
	}

	//fmt.Println(text_name, len_sentence)
}
func main() {
	start, end := 1, 1
	for i := start; i <= end; i++ {
		Dowork(i)
	}
}
