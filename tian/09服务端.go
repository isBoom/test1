package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

//日志文件
var log *os.File

//当前时间
var ctime time.Time

//通信管道
var Message = make(chan string)

//map存储用户列表
var ClientMap map[string]Client = make(map[string]Client)

//成员结构体
type Client struct {
	ch   chan string //通道
	name string      //名字
	addr string      //ip
}

//返回当前时间
func currentTime(now time.Time, type_num int) (str string) {
	now = time.Now()
	if type_num == 1 {
		str = fmt.Sprintf("%d-%d-%d %d:%d:%d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
	} else {
		str = fmt.Sprintf("%d年%d月%d日%d点%d分%d秒", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
	}
	return
}

//封装返回值
func return_str(cli Client, msg string) (str string) {
	str = "[" + cli.addr + "]-" + ClientMap[cli.addr].name + ":\t\t" + msg
	return
}

//发送信息
func send(conn net.Conn, cli Client) {
	for {
		conn.Write([]byte(<-cli.ch))
	}
}

//当前在线用户
func who(conn net.Conn, cli Client) {
	var i int
	for _, tmp := range ClientMap {
		cli.ch <- tmp.name
		i++
	}
	cli.ch <- "当前在线人数为-" + fmt.Sprint(i) + "-\n"
}

//重命名
func reanme(conn net.Conn, cli Client, newname string) {
	str := strings.TrimSpace(newname)
	if len(str) == 0 {
		conn.Write([]byte("请不要输入全是空格的名字"))
	} else if len(str) > 22 {
		conn.Write([]byte("昵称过长"))
	} else {
		//相邻字符之间最多为一个空格
		if strings.Contains(str, "  ") == true {
			conn.Write([]byte("请输入合法字符,相邻字符之间最多为一个空格"))
		} else {
			Message <- cli.name + "已成功改名为" + str
			cli.name = str
			ClientMap[cli.addr] = cli
		}
	}
}

//用户退出
func exit(conn net.Conn, cli Client) {
	conn.Close()
}

//help
func help(conn net.Conn, cli Client) {
	str := "\thelp---获取帮助\n" + "\twho---查看当前在线列表\n" + "\texit---退出登录\n" + "\trename xxx---xxx为新名字,相邻字符之间最多为一个空格,一个汉字三个字符长,不得超过21字符\n" + "" + ""
	cli.ch <- str
}

//协程 收到用户信息就转发
func traverse() {
	for {
		msg := <-Message
		for _, tmp := range ClientMap {
			tmp.ch <- msg
		}
	}
}

func HandleConn(conn net.Conn) {
	//获取ip
	address := conn.RemoteAddr().String()
	//姓名默认是ip
	cli := Client{make(chan string), address, address}
	//加入map
	ClientMap[address] = cli
	//关闭时从用户列表删除此人
	defer func() {
		conn.Close()
		Message <- return_str(cli, "已退出聊天室")
		log.Write([]byte(currentTime(ctime, 1) + return_str(cli, "已退出聊天室\n\n")))
		delete(ClientMap, cli.addr)
	}()
	//新建协程监听管道 一有数据传来就发送 写在这里传参数
	go send(conn, cli)
	//提示用户
	cli.ch <- "来自系统的消息:您的昵称为[" + cli.name + "],输入help获取帮助或更改昵称"
	//广播登录信息
	Message <- return_str(cli, "Login")
	//日记记录
	log.Write([]byte(currentTime(ctime, 1) + return_str(cli, "Login\n\n")))
	//接收此用户信息
	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("conn.Read", err)
			return
		}
		temp_buf := string(buf[:n])
		log.Write([]byte(currentTime(ctime, 1) + return_str(cli, temp_buf) + "\n\n"))

		switch {
		case "who" == temp_buf:
			who(conn, cli) //谁在线
		case n >= 8 && temp_buf[:6] == "rename" && temp_buf[6] == 32:
			reanme(conn, cli, temp_buf[7:n])
		case "help" == temp_buf:
			help(conn, cli)
		case "exit" == temp_buf:
			exit(conn, cli)
		default:
			Message <- return_str(cli, temp_buf)
		}
	}
}
func main() {
	listener, err0 := net.Listen("tcp", ":8088")

	if err0 != nil {
		fmt.Println("net.Listen:", err0)
		return
	}
	//记录日志
	log, _ = os.Create(currentTime(ctime, 0) + ".log")
	defer func() {
		listener.Close()
		log.Close()
	}()
	fmt.Println("server start")
	//协程,有用户发送信息就遍历在线用户广播信息
	go traverse()

	//主协程处理用户登入并使其加入map
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener.Accept:", err)
			continue
		}
		//新建协程处理用户接入
		go HandleConn(conn)
	}
}
