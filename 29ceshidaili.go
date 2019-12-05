package main

import (
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	DB    *sqlx.DB
	sw    sync.WaitGroup
	ch    = make(chan int, 10)
	count int
)

type theIp struct {
	Ip        string `json:"ip" db:"ip"`
	IsChecked int    `json:"isChecked" db:"isChecked"`
}

func init() {
	var err error
	DB, err = sqlx.Connect("mysql", "root:root@tcp(39.106.169.153:3306)/ipPond?charset=utf8&&parseTime=true")
	if err != nil {
		panic(err)
	}
}
func testIp(mod theIp) {
	sw.Add(1)
	proxy := func(_ *http.Request) (*url.URL, error) {
		return url.Parse(mod.Ip)
	}
	transport := &http.Transport{Proxy: proxy}
	c := &http.Client{Transport: transport, Timeout: time.Second * 3}
	req, err := http.NewRequest("GET", "https://xxxholic.top", nil)
	if err != nil {
		sw.Done()
		<-ch
		DB.Exec("update ipPond set isChecked=1 where ip=?", mod.Ip)
		return
	}
	res, err := c.Do(req)
	if err != nil {
		sw.Done()
		<-ch
		DB.Exec("update ipPond set isChecked=1 where ip=?", mod.Ip)
		return
	}
	defer res.Body.Close()
	if err != nil {
		sw.Done()
		<-ch
		DB.Exec("update ipPond set isChecked=1 where ip=?", mod.Ip)
		return
	} else if res.StatusCode != 200 {
		sw.Done()
		<-ch
		DB.Exec("update ipPond set isChecked=1 where ip=?", mod.Ip)
		return
	}
	DB.Exec("insert into availIp VALUE(?,0)", mod.Ip)
	fmt.Printf("ip[%s]已存入数据库\n", mod.Ip)
	DB.Exec("update ipPond set isChecked=1 where ip=?", mod.Ip)
	<-ch
	sw.Done()
}
func main() {
	mod := make([]theIp, 0)
	if err := DB.Select(&mod, "select * from ipPond where isChecked = 0 or isChecked = 1"); err != nil {
		fmt.Println("selectErr", err)
		return
	}
	for _, temp := range mod {
		ch <- 0
		go testIp(temp)
	}
	time.Sleep(time.Second / 10)
	sw.Wait()
}
