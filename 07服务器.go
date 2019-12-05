package main

import (
	"fmt"
	"net"
)

func fw(conn net.Conn) {
	defer func() {
		fmt.Println("[", conn.RemoteAddr(), "]", "下线了")
		conn.Write([]byte("下线了"))
		conn.Close()
	}()
	fmt.Println("[", conn.RemoteAddr(), "]", "上线了")
	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println(err)
			return
		}
		if "exit" == string(buf[:n]) {
			conn.Write([]byte("您已登出"))
			conn.Close()
		}
		fmt.Println("[", conn.RemoteAddr(), "]:", string(buf[:n]))
		conn.Write(buf[:n])
	}
}
func main() {
	listener, err := net.Listen("tcp", ":8088")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer listener.Close()
	fmt.Println("服务器上线了")
	for {
		conn, err1 := listener.Accept()
		if err1 != nil {
			fmt.Println(err)
			return
		}
		go fw(conn)
	}

}
