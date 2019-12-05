package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
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
func download(src string) {
	img, _ := os.Create(src)
	defer img.Close()

	re_img, _ := http.Get(src)
	defer re_img.Body.Close()

	buf_img := make([]byte, 4096)
	for {
		n, _ := re_img.Body.Read(buf_img)
		if n == 0 {
			break
		}
		img.Write(buf_img[:n])
	}
	fmt.Println(src + "下载成功")
}
func main() {

	url := "https://yandex.ru/images/search?text=%E5%8A%A8%E6%BC%AB&isize=large"
	result, err0 := httpget(url)
	if err0 != nil {
		fmt.Println("httpget", err0)
		return
	}
	fmt.Println(result)
	re_regexp := regexp.MustCompile(`https?://.*?\.(jpg|png|gif|webp|svg|tif|url)`)
	if re_regexp == nil {
		fmt.Println("regexp.MustCompile err")
		return
	}
	for i := 0; i < 300; i++ {
		data := re_regexp.FindAllSubmatch([]byte(result), -1)
		fmt.Println(string(data[i][0]) + "开始下载")
		go download(string(data[i][0]))
	}

}
