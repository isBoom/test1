package main

import (
	"fmt"
	"io"
	"os"
)

func copy(path string, new string) {
	new_f, err1 := os.Create(new)
	if err1 != nil {
		fmt.Println(err1)
		return
	}

	old_f, err2 := os.Open(path)
	if err2 != nil {
		fmt.Println(err2)
		return
	}
	defer new_f.Close()
	defer old_f.Close()

	buf_len, _ := old_f.Seek(0, os.SEEK_END)
	old_f.Seek(0, os.SEEK_SET)
	byte := make([]byte, buf_len)

	for {
		n, err3 := old_f.Read(byte)
		if err3 != nil {
			if err3 != io.EOF {
				fmt.Println(err3)
			}
			return
		} else {
			new_f.Write(byte[:n])
		}
	}
}
func main() {
	var temp, new string
	fmt.Scan(&temp, &new)
	path := "./" + temp
	copy(path, new)
}
