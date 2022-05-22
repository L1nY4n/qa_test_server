package tcpserver

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
		buf_w := &bytes.Buffer{}
		buf_w.Write(Buf[:])

		//fmt.Printf("buf=: %v\n", buf_w.Bytes())
		binary.Read(buf_w, binary.BigEndian, &device.Dev_cap)
		cont++
		device.Dev_cap.Cap_if.ID = cont

		device.DB.Debug().Create(&device.Dev_cap)
		//fmt.Printf("receive from client, data: %v\n", (Buf[:n]))
		//fmt.Printf("cap_in [name]=: %v\n", (device.Dev_cap))
		//发送数据
		if _, err = conn.Write(Buf[:n]); err != nil {
			fmt.Printf("write to client failed, err: %v\n", err)
			break
		}
	}
}
