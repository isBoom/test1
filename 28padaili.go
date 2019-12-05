package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type availIp struct {
	Ip       string `json:"ip" db:"ip"`
	ErrCount int    `json:"errCount" db:"errCount"`
}

var (
	DB  *sqlx.DB
	mod = make([]availIp, 0)
	ch  = make(chan int, 10)
	sw  sync.WaitGroup
)

func init() {
	var err error
	DB, err = sqlx.Connect("mysql", "root:root@tcp(39.106.169.153:3306)/ipPond?charset=utf8&&parseTime=true")
	if err != nil {
		panic(err)
	}
}
func getAgent() string {
	agent := [...]string{
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/535.1 (KHTML, like Gecko) Chrome/14.0.835.163 Safari/535.1",
		"Mozilla/5.0 (Windows NT 6.1; WOW64; rv:6.0) Gecko/20100101 Firefox/6.0",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50",
		"Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; Win64; x64; Trident/5.0; .NET CLR 2.0.50727; SLCC2; .NET CLR 3.5.30729; .NET CLR 3.0.30729; Media Center PC 6.0; InfoPath.3; .NET4.0C; Tablet PC 2.0; .NET4.0E)",
		"Mozilla/4.0 (compatible; MSIE 8.0; Windows NT 6.1; WOW64; Trident/4.0; SLCC2; .NET CLR 2.0.50727; .NET CLR 3.5.30729; .NET CLR 3.0.30729; Media Center PC 6.0; .NET4.0C; InfoPath.3)",
		"Mozilla/4.0 (compatible; MSIE 8.0; Windows NT 5.1; Trident/4.0; GTB7.0)",
		"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1)",
		"Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1; SV1)",
		"Mozilla/5.0 (Windows; U; Windows NT 6.1; ) AppleWebKit/534.12 (KHTML, like Gecko) Maxthon/3.0 Safari/534.12",
		"Mozilla/5.0 (Windows; U; Windows NT 6.1; ) AppleWebKit/534.12 (KHTML, like Gecko) Maxthon/3.0 Safari/534.12",
		"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 6.1; WOW64; Trident/5.0; SLCC2; .NET CLR 2.0.50727; .NET CLR 3.5.30729; .NET CLR 3.0.30729; Media Center PC 6.0; InfoPath.3; .NET4.0C; .NET4.0E; SE 2.X MetaSr 1.0)",
		"Mozilla/5.0 (Windows; U; Windows NT 6.1; en-US) AppleWebKit/534.3 (KHTML, like Gecko) Chrome/6.0.472.33 Safari/534.3 SE 2.X MetaSr 1.0",
		"Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; WOW64; Trident/5.0; SLCC2; .NET CLR 2.0.50727; .NET CLR 3.5.30729; .NET CLR 3.0.30729; Media Center PC 6.0; InfoPath.3; .NET4.0C; .NET4.0E)",
		"Mozilla/5.0 (Windows NT 6.1) AppleWebKit/535.1 (KHTML, like Gecko) Chrome/13.0.782.41 Safari/535.1 QQBrowser/6.9.11079.201",
		"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 6.1; WOW64; Trident/5.0; SLCC2; .NET CLR 2.0.50727; .NET CLR 3.5.30729; .NET CLR 3.0.30729; Media Center PC 6.0; InfoPath.3; .NET4.0C; .NET4.0E) QQBrowser/6.9.11079.201",
		"Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; WOW64; Trident/5.0)",

		"Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:50.0) Gecko/20100101 Firefox/50.0",
		"Opera/9.80 (Macintosh; Intel Mac OS X 10.6.8; U; en) Presto/2.8.131 Version/11.11",
		"Opera/9.80 (Windows NT 6.1; U; en) Presto/2.8.131 Version/11.11",
		"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; 360SE)",
		"Mozilla/5.0 (Windows NT 6.1; rv:2.0.1) Gecko/20100101 Firefox/4.0.1",
		"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; The World)",
		"Mozilla/5.0 (Macintosh; U; Intel Mac OS X 10_6_8; en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50",
		"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; Maxthon 2.0)",
		"Mozilla/5.0 (Windows; U; Windows NT 6.1; en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50",
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	len := len(agent)
	return agent[r.Intn(len)]
}
func getIp() availIp {
	if err := DB.Select(&mod, "select * from availIp"); err != nil {
		panic(err)
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	len := len(mod)
	return mod[r.Intn(len)]
}
func getErr(temp availIp) {
	if temp.ErrCount > 5 {
		if _, err := DB.Exec("delete from availIp where ip = ?", temp.Ip); err == nil {
			fmt.Printf("删除ip[%s]成功\n", temp.Ip)
		} else {
			fmt.Printf("删除ip[%s]失败 %v\n", temp.Ip, err)
		}
	} else {
		DB.Exec("update availIp set errCount = ? where ip = ?", temp.ErrCount+1, temp.Ip)
		fmt.Println("访问失败正在切换ip")
	}
}
func do(address string) (string, error) {
	for {
		mod := getIp()
		fmt.Printf("正在用ip[%s]爬取%s\n", mod.Ip, address)
		proxy := func(_ *http.Request) (*url.URL, error) {
			return url.Parse(mod.Ip)
		}
		transport := &http.Transport{Proxy: proxy}
		c := &http.Client{Transport: transport, Timeout: time.Second * 5}
		req, err := http.NewRequest("GET", address, nil)
		if err != nil {
			getErr(mod)
			continue
		}
		req.Header.Add("User-Agent", getAgent())
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
		req.Header.Set("Connection", "keep-alive")
		req.Header.Set("Referer", address)
		//执行请求
		res, err := c.Do(req)
		if err != nil {
			getErr(mod)
			continue
		}
		defer res.Body.Close()
		temp, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println("readall", err)
			continue
		} else if res.StatusCode != 200 {
			getErr(mod)
			continue
		}
		count := mod.ErrCount - 1
		if count < -5 {
			count = -5
		}
		DB.Exec("update availIp set errCount= ? where ip = ?", count, mod.Ip)
		return string(temp), nil
	}

}
func Gethtml(url string) (result string, err error) {
	return do(url)
}
func insertIp(ip string) error {
	if _, err := DB.Exec("insert into ipPond VALUE(?,0)", ip); err != nil {
		return err
	}
	return nil
}
func pa(url string) {
	sw.Add(1)
	str, err := Gethtml(url)
	if err != nil {
		fmt.Println("Gethtmlerr", err)
	} else {
		temp := regexp.MustCompile(`<tr class="odd">(?s:(.*?))<div title=`)
		for _, data := range temp.FindAllStringSubmatch(str, -1) {
			temp = regexp.MustCompile(`<td>(.*?)</td>`)
			tempData := temp.FindAllStringSubmatch(data[1], -1)
			if tempData[2][1] != "HTTPS" {
				tempData[2][1] = "http"
			}
			ip := fmt.Sprintf("%s://%s:%s", strings.ToLower(tempData[2][1]), tempData[0][1], tempData[1][1])
			if err := insertIp(ip); err == nil {
				fmt.Printf("ip[%s]已存入数据库\n", ip)
			}
		}
	}
	sw.Done()
	<-ch
}
func main() {
	for i := 1; i <= 3916; i++ {
		ch <- 0
		go pa(fmt.Sprintf("https://www.xicidaili.com/nn/%d", i))
	}
	sw.Wait()
}
