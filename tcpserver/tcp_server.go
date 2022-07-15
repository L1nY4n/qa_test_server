package tcpserver

import (
	"bytes"
	"encoding/binary"
	"encoding/json"

	"fmt"
	"net"
	"qa_test_server/manager"
	"qa_test_server/model"
	"qa_test_server/web"
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

		var temp model.Nano_Dev_capture_packed
		binary.Read(buf_w, binary.LittleEndian, &temp)
		//fmt.Printf("Dev_cap= %+v\n", temp)
		dev := model.Decode(temp)
		manager.ManagerGlabal.Update(dev)
		go func() {
			if d, err := json.Marshal(dev); err == nil {
				web.WsManager.Groupbroadcast("device_upload", d)
			}
		}()
		//fmt.Printf("##############################")

		//fmt.Printf("msg %d\n", device.Dev_cap.Sys_para.Seed_param.Freq)
		//device.DB.Debug().Create(&device.Dev_cap)
		if _, err = conn.Write(buffer[:n]); err != nil {
			fmt.Printf("write to client failed, err: %v\n", err)
			break
		}

	}
}
