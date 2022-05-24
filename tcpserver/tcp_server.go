package tcpserver

import "C"
import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"qa_test_server/device"
)

var Buf [2048]byte

func Tcpserver() {
	go proxyStart(9999, 7777)
	fmt.Println("tcp server test---lf 2022-4-28")
	listerer, err := net.Listen("tcp", "0.0.0.0:4001")
	defer listerer.Close()
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

var cont uint32 = 1

func process(conn net.Conn) {
	defer conn.Close()
	//循环处理数据
	for {

		//接受数据
		n, err := conn.Read(Buf[:])
		if err != nil {
			fmt.Printf("read from connect failed, err: %v\n", err)
			break
		}
		//新建二进制流
		buf_w := &bytes.Buffer{}
		//读入六
		buf_w.Write(Buf[:n])
		//写入到结构体
		binary.Read(buf_w, binary.LittleEndian, &device.Dev_cap)
		fmt.Printf("Dev_cap= %+v\n", device.Dev_cap)
		fmt.Printf("##############################")
		//fmt.Printf("msg %d\n", device.Dev_cap.Sys_para.Seed_param.Freq)
		//device.DB.Debug().Create(&device.Dev_cap)
		if _, err = conn.Write(Buf[:n]); err != nil {
			fmt.Printf("write to client failed, err: %v\n", err)
			break
		}

	}
}
