package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
)

func donnload(w http.ResponseWriter, r *http.Request) {
	if r.Header["Download"] == nil {
		fmt.Println("fenzhi1")
		w.Write([]byte("您没输出要下载的文件"))
		return
	} else {
		filename := r.Header["Download"][0]
		var start int64
		if r.Header["Start"] != nil {
			var err error
			start, err = strconv.ParseInt(r.Header["Start"][0], 10, 64)
			if err != nil {
				fmt.Println("strconv.ParseInt err", err)
				return
			}
		}

		file, err := os.Open(filename)
		if err != nil {
			w.Header().Set("Status", "404")
			return
		}
		fileinfo, _ := file.Stat()
		w.Header().Set("FileSize", fmt.Sprintf("%v", fileinfo.Size()))
		defer file.Close()
		file.Seek(start, 0)

		fmt.Printf("正在从%d位置开始下载\n", start)
		buf := make([]byte, 4096)
		for {
			n, _ := file.Read(buf)
			if n != 0 {
				start += int64(n)
				fmt.Printf("进度%d/%dMb\n", start/(1024*1024), fileinfo.Size()/(1024*1024))
				w.Write(buf[:n])
			} else {
				break
			}
		}
	}
	fmt.Println("发送完毕")

}
func main() {
	//服务端
	fmt.Println("start")
	http.HandleFunc("/d", donnload)
	http.ListenAndServe(":8077", nil)
}
