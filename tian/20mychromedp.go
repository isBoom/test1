package mychromedp

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/chromedp/chromedp"
)

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
		"User-Agent,Mozilla/5.0 (Macintosh; U; Intel Mac OS X 10_6_8; en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50",
		"User-Agent, Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; Maxthon 2.0)",
		"User-Agent,Mozilla/5.0 (Windows; U; Windows NT 6.1; en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50",
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	len := len(agent)
	return agent[r.Intn(len)]
}
func getIp() string {
	ip := [...]string{
		"http://113.121.23.20:9999",
		"http://61.145.8.25:9999",
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	len := len(ip)
	return ip[r.Intn(len)]
}
func do(url string) (string, err) {
	//获取随机代理
	proxy := func(_ *http.Request) (*url.URL, error) {
		return url.Parse(getIp())
	}
	transport := &http.Transport{Proxy: proxy}
	c := &http.Client{Transport: transport}
	//构造请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	req.Header.Add("User-Agent", getAgent())
	request.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	request.Header.Set("Connection", "keep-alive")
	//执行请求
	res, err := c.Do(req)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer res.Body.Close()
	temp, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return string(temp), nil
}

func Getbody(url string) (res string, err error) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	if err = chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Sleep(time.Second),
		chromedp.OuterHTML("body", &res, chromedp.ByQueryAll),
	); err != nil {
		return
	}
	return
}
func Gethtml(url string) (result string, err error) {
	resp, err := do(url)
	if err0 != nil {
		err = err0
		return
	}
	buf := make([]byte, 4096)
	//返回html文本
	for {
		n, err0 := resp.Body.Read(buf)
		if err0 != nil {
			if err0 == io.EOF {
				result += string(buf[:n])
				return
			} else {
				err = err0
				return
			}
		}
		result += string(buf[:n])
	}
	return
}
func Getfile(url string, name string) (err error) {
	resp, err_r := http.Get(url)
	if err_r != nil {
		err = err_r
		return
	}
	defer resp.Body.Close()
	file, err_f := os.Create(name)
	if err_f != nil {
		err = err_f
		return
	}
	defer file.Close()

	buf := make([]byte, 4096)
	//保存图片
	for {
		n, err_img := resp.Body.Read(buf)
		if err_img != nil {
			if err_img == io.EOF {
				file.Write(buf[:n])
				return
			} else {
				err = err_img
				return
			}
		}
		file.Write(buf[:n])
	}
	return
}
