package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"sync"
	"time"
)

var (
	root_url string
	sw       sync.WaitGroup
)

func Gethtml(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	str := ""
	buf := make([]byte, 4096)
	for {
		n, _ := resp.Body.Read(buf)
		if n == 0 {
			break
		}
		str += string(buf[:n])
	}
	return str, nil
}
func Getfile(url string, path string) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	file, _ := os.Create(path)
	defer file.Close()
	buf := make([]byte, 4096)
	for {
		n, _ := resp.Body.Read(buf)
		if n == 0 {
			break
		}
		file.Write(buf[:n])
	}
}
func download(url string, path string) {
	//获取最终html
	html, err := Gethtml(url)
	if err != nil {
		fmt.Println("mychromedp.Gethtml", err)
		return
	}
	//获取当前	页面图片url，name
	temp_item_url := regexp.MustCompile(`<a href="/(.*?)" class="bqba" title="`)
	temp_itrm_name := regexp.MustCompile(`class="bqba" title="(.*?)">`)

	item_url := temp_item_url.FindAllStringSubmatch(html, -1)
	item_name := temp_itrm_name.FindAllStringSubmatch(html, -1)

	for i := 0; i < len(item_url); i++ {
		new_url := root_url + "/" + item_url[i][1]
		new_path := path + "/" + item_name[i][1]
		os.Mkdir(new_path, 777)

		detail_url, err := Gethtml(new_url)
		if err != nil {
			fmt.Println("mychromedp.Gethtml", err)
			continue
		}

		temp_img := regexp.MustCompile(`<div class="swiper-wrapper">(?s:(.*?))<script async src`)
		img_html := temp_img.FindAllStringSubmatch(detail_url, -1)[0][1]

		temp_img_url := regexp.MustCompile(`data-original="(.*?)" title="`)
		temp_img_name := regexp.MustCompile(`" alt="(?s:(.*?))" style="max-height: 100%`)

		img_url := temp_img_url.FindAllStringSubmatch(img_html, -1)
		img_name := temp_img_name.FindAllStringSubmatch(img_html, -1)

		for i := 0; i < len(img_url); i++ {
			go Getfile(img_url[i][1], new_path+"/"+img_name[i][1]+img_url[i][1][len(img_url[i][1])-4:])
			fmt.Println(img_name[i][1] + "   over")
		}

	}
}
func subpage(temp_url string, dir_name string) {
	//添加一个协程指示器
	sw.Add(1)
	//获取各个大分类表情包下的html
	getpagecount, err := Gethtml(root_url + temp_url)
	if err != nil {
		fmt.Println("getpagecount", err)
		sw.Done()
		return
	}

	//当前大分类的页数
	temp_pagecount := regexp.MustCompile(`<span id="mobilepage"><span>1</span> / (\d+?)</span>`)
	page_count := temp_pagecount.FindAllStringSubmatch(getpagecount, -1)[0][1]
	count, err := strconv.Atoi(page_count)
	if err != nil {
		fmt.Println("atoi", err)
		sw.Done()
		return
	}
	//page从1开始 不是0
	for i := 1; i <= count; i++ {
		//大分类下小分类的url
		url := root_url + temp_url[:len(temp_url)-5] + "/page/" + strconv.Itoa(i) + ".html"
		os.Mkdir(dir_name, 777)
		//进入小分类
		download(url, dir_name)
	}
	//一个协程结束
	sw.Done()
}
func main() {
	//主页
	root_url = "https://fabiaoqing.com"
	//首页url 用来获取分类url
	first_index_url := "https://fabiaoqing.com/bqb/lists/type/hot.html"
	res_first, err := Gethtml(first_index_url)
	if err != nil {
		fmt.Println("mychromedp.Getbody", err)
		return
	}
	//获取各个分类的url
	temp_item := regexp.MustCompile(`id="bqbcategory"(?s:(.*?))<a href=`)
	res_str := temp_item.FindAllStringSubmatch(res_first, -1)[0][1]

	temp_item_url := regexp.MustCompile(`href="(.*?)" title="`)
	temp_item_name := regexp.MustCompile(`title="(.*?)">`)

	//得到列表项url切片
	item_url := temp_item_url.FindAllStringSubmatch(res_str, -1)
	//列表项名字切片
	item_name := temp_item_name.FindAllStringSubmatch(res_str, -1)

	//多协程同时下载十几个大分类的表情包
	for i := 0; i < len(item_url); i++ {
		go subpage(item_url[i][1], item_name[i][1][:len(item_name[i][1])-10])
	}
	time.Sleep(time.Second)
	//所有协程结束后进程随之结束
	sw.Wait()
	fmt.Println("over")
}
