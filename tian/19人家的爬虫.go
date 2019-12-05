package main

/*本程序由 itruirui@outlook.com 瑞哥 开发制作*/
import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

var Size_min int64 = 100 //图片最小值，为100KB，修改这里既可修改限制
func main() {
	/*虽然GO支持百万级协程，但是，网速不快的你，最好不要开太多*/
	ch := make(chan int, 10) //下载图片，最大协程10
	ch2 := make(chan int, 2) //分析页面代码，最大协程2
	var min int              //搜索开始地址
	var max int              //搜索结束地址
	min = 20000              //最小地址
	max = 22000              //最大地址
	os.Mkdir("DownloadImg", os.ModePerm)
	for i := max; i >= min; i-- {
		url := fmt.Sprintf("http://www.netbian.com/desk/%d-1920x1080.htm", i)
		fmt.Printf("获取 %d.html中\n", i)
		ch2 <- i
		go DownloadUrlImg(url, ".jpg", ch, ch2) //后面的.jpg是要下载的文件类型
	}
}

/*
* 作用：下载某个URL里面的资源
* 参数url：URL地址
* 参数filetype：文件类型，.jpg或者.png
* 返回值sum：一共下载多少个资源
 */
func DownloadUrlImg(url string, filetype string, ch chan int, ch2 chan int) (sum int) {
	body, err := GetHtml(url) //获取html代码
	if err != nil {
		fmt.Println("爬取失败，页面编号：", <-ch2)
		return 0
	}
	imgsrc := GetImg(body, filetype)
	var sz []string
	for i := 0; i < len(imgsrc); i++ {
		sz = append(sz, imgsrc[i][1])
	}
	for _, data := range sz {
		ch <- 1
		go DlImg(data, ch)
	}
	fmt.Println("爬取完毕，页面编号：", <-ch2)
	return 0
}

/*
* 获取某个url的html代码
* 参数：url地址
* 返回值：string类型，error类型
 */
func GetHtml(url string) (value string, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err = errors.New(fmt.Sprintln("错误代码：", resp.StatusCode))
		return value, err
	}
	all, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	return string(all), nil
}

/*
* 解析某个string里面的文件
* 参数：要解析的string
* 返回值：文件内容[][]string
 */
func GetImg(body string, filetype string) (result [][]string) {
	str := fmt.Sprintf("src=\"([^\"]*%s)\"", filetype)
	reg := regexp.MustCompile(str) //正则表达式规则
	if reg == nil {
		fmt.Println("正则表达式，编译错误：", reg)
		return
	} else {
		result = reg.FindAllStringSubmatch(body, -1)
		return result
	}
}

func DlImg(url string, ch chan int) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer func() {
		resp.Body.Close()
		<-ch
	}()
	file, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	if resp.ContentLength < (Size_min * 1024) {
		return
	}
	name := strings.Replace(url, "/", "_", -1)
	name = strings.Replace(name, ":", "-", -1)
	name = "DownloadImg/" + name
	_, err = os.Stat(name)
	if err == nil {
		return
	}
	ioutil.WriteFile(name, file, 0666)
	fmt.Println(name)
}
