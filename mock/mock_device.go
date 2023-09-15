package main

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

func processAlarm(conn net.Conn) {
	// defer conn.Close() // 关闭连接
	for {
		reader := bufio.NewReader(conn)
		var buf [128]byte
		_, err := reader.Read(buf[:]) // 读取数据
		if err != nil {
			fmt.Println("read from client failed, err:", err)
			break
		}

		fmt.Println("收到client端发来的数据：", buf)
		if buf[1] == 0x03 && buf[5] == 0x06 {
			//电子水尺
			// data[3:4]
			r := []byte{0x03, 0x00, 0x00, 0xa0, 0x06}
			conn.Write(r) // 发送数据
		}
		if buf[1] == 0x03 && buf[5] == 0x02 {
			//超声波液位器
			// data[3:6]
			r := []byte{0x03, 0x00, 0x00, 0xab, 0x1f, 0x10, 0x0a}
			conn.Write(r) // 发送数据
		}
		if buf[1] == 0x03 && buf[3] == 0x00 {
			//倾角变送器问询指令x
			// data[3:4]
			r := []byte{0x03, 0x00, 0x00, 0xa3, 0x1a}
			conn.Write(r) // 发送数据
		}
		if buf[1] == 0x03 && buf[3] == 0x01 {
			//倾角变送器问询指令y
			// data[3:4]
			r := []byte{0x03, 0x00, 0x00, 0xa1, 0x16}
			conn.Write(r) // 发送数据
		}

	}
}

func MockDeivce(port string) {

	listen, err := net.Listen("tcp", "127.0.0.1:"+port)
	if err != nil {
		fmt.Println("listen failed, err:", err)
		return
	}
	for {
		conn, err := listen.Accept() // 建立连接
		if err != nil {
			fmt.Println("accept failed, err:", err)
			continue
		}
		go processAlarm(conn) // 启动一个goroutine处理连接
	}
}

func main() {
	// var port string
	// flag.StringVar(&port, "p", "20000", "port")
	// flag.Parse()
	// fmt.Printf("arg:%v", port)

	//
	go MockDeivce("9001")
	go MockDeivce("9002")
	go MockDeivce("9003")
	go MockDeivce("9004")
	go MockDeivce("9005")
	go MockDeivce("9006")

	for {
		time.Sleep(1 * time.Hour)
	}
}
