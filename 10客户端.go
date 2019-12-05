package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	//dialog, err := net.Dial("tcp", "39.106.169.153:8088")
	dialog, err := net.Dial("tcp", "127.0.0.1:8088")
	if err != nil {
		fmt.Println("net.Dial", err)
		return
	}

	defer func() {
		dialog.Close()
	}()

	go func() {
		fmt.Print("登陆成功！\n")
		for {
			str := make([]byte, 1024)
			n, err1 := os.Stdin.Read(str)
			if err1 != nil {
				fmt.Println(err1)
				return
			}
			if n > 2 {
				dialog.Write(str[:n-2])
			}

		}

	}()
	for {
		str := make([]byte, 1024)
		n, err1 := dialog.Read(str)
		if err1 != nil {
			fmt.Println(err1)
			return
		}
		fmt.Println(string(str[:n]))
	}
}
