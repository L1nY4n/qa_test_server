package tcpserver

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"qa_test_server/device"
	"time"
)

func Tcpserver() {
	go proxyStart(4003, 7777)
	fmt.Println("tcp server test---lf 2022-4-28")
	listerer, err := net.Listen("tcp", "0.0.0.0:4001")

	if err != nil {
		fmt.Printf("listen fail,err: %v\n", err)
		return
	}
	defer listerer.Close()
	
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
	//循环处理数据
	for {
		var buffer [2048]byte
		//接受数据
		n, err := conn.Read(buffer[:])
		if err != nil {
			fmt.Printf("read from connect failed, err: %v\n", err)
			break
		}
		//新建二进制流
		buf_w := &bytes.Buffer{}
		//读入六
		//var time [4]uint64 = 0

		buf_w.Write(buffer[:8])
		buf_w.Write(buffer[:n])
		//写入到结构体

		var temp device.Dev_capture_packed
		binary.Read(buf_w, binary.LittleEndian, &temp)
		temp.Cap_info.Time = uint64(time.Now().Unix())
		dev := device.Decode(temp)
		device.ManagerGlabal.Update(dev)
		//fmt.Printf("##############################")
		//fmt.Printf("msg %d\n", device.Dev_cap.Sys_para.Seed_param.Freq)
		//device.DB.Debug().Create(&device.Dev_cap)
		if _, err = conn.Write(buffer[:n]); err != nil {
			fmt.Printf("write to client failed, err: %v\n", err)
			break
		}

	}
}
