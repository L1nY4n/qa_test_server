package tcpserver

import (
	"fmt"
	"net"
)

var buf [128]byte

func Tcpserver() {
	fmt.Println("tcp server test---lf 2022-4-28")
	listerer, err := net.Listen("tcp", "0.0.0.0:9090")
	if err != nil {
		fmt.Printf("listen fail,err: %v\n", err)
		return
	}
	for {
		//等待连接
		conn, err := listerer.Accept()
		if err != nil {
			fmt.Printf("accept fail, err: %v\n", err)
			continue
		}
		go process(conn)
	}
}

func process(conn net.Conn) {
	defer conn.Close()
	for {

		//接受数据
		n, err := conn.Read(buf[:])
		if err != nil {
			fmt.Printf("read from connect failed, err: %v\n", err)
			break
		}
		fmt.Printf("receive from client, data: %v\n", string(buf[:n]))
		//发送数据
		if _, err = conn.Write(buf[:n]); err != nil {
			fmt.Printf("write to client failed, err: %v\n", err)
			break
		}
	}
}
