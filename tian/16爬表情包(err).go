package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
)

//记录当前是第几张
var count int

const load = 1    //httpget函数的参数 表示不返回html而创建图像文件
const res_str = 0 //httpget函数的参数 表示返回html内容供后续使用

func httpget(url string, thetype int) (result string, err error) {
	//发送请求
	resp, err0 := http.Get(url)
	if err0 != nil {
		err = err0
		return
	}
	defer resp.Body.Close()

	buf := make([]byte, 4096)
	//保存图片
	if thetype == load {
		count++
		file, err1 := os.Create(strconv.Itoa(count) + "." + url[len(url)-3:len(url)])
		if err1 != nil {
			err = err1
			return
		}
		defer file.Close()
		for {
			n, _ := resp.Body.Read(buf)
			if n == 0 {
				break
			}
			file.Write(buf[:n])
		}
	} else {
		//返回html文本
		for {
			n, _ := resp.Body.Read(buf)
			if n == 0 {
				break
			}
			result += string(buf[:n])
		}
	}
	return
}

//核心工作函数
func Dowork(page int) {
	url := "https://www.doutula.com/article/list/?page=" + strconv.Itoa(page)
	temp1, err0 := httpget(url, res_str)
	if err0 != nil {
		fmt.Println("httpget1", err0)
		return
	}

	//正则匹配第一级网址
	res1 := regexp.MustCompile(`<a href="(https?://(.*?)[0-9]{3,3})" class="`)
	res1_url := res1.FindAllStringSubmatch(temp1, -1)
	//迭代第一次匹配的网站
	for _, data1 := range res1_url {
		temp2, err2 := httpget(data1[1], res_str)
		if err2 != nil {
			fmt.Println("httpget2", err2)
			continue
		}
		//正则匹配第二级网址
		res2 := regexp.MustCompile(`" src="(https?://(.*?))" alt="`)
		res2_url := res2.FindAllStringSubmatch(temp2, -1)
		//迭代最终表情包网址信息
		for _, data2 := range res2_url {
			//第二个参数为load 表示下载图片
			go httpget(data2[1], load)
			fmt.Println(data2[1] + "     第 " + strconv.Itoa(count) + "张   over")
		}
	}
}
func main() {
	count = 1
	//初始页码和最后一个页码
	start, end := 1, 630
	for i := start; i <= end; i++ {
		Dowork(i)
	}
}
