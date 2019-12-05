package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	//dialog, err := net.Dial("tcp", "39.106.169.153:8888")
	dialog, err := net.Dial("tcp", "223.91.51.251:8888")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer dialog.Close()
	go func() {
		fmt.Print("登陆成功！\n输入中...")
		for {
			str := make([]byte, 1024)
			n, err1 := os.Stdin.Read(str)
			if err1 != nil {
				fmt.Println(err1)
				return
			}
			dialog.Write(str[:n-2])
		}

	}()
	for {
		str := make([]byte, 1024)
		n, err1 := dialog.Read(str)
		if err1 != nil {
			fmt.Println(err1)
			return
		}
		if "您已登出" == string(str[:n]) {
			return
		}
		fmt.Print("来自服务器的消息:", string(str[:n]), "\n输入中...")
	}
}
