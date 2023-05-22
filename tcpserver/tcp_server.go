package tcpserver

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"time"

	"fmt"
	"net"
	"qa_test_server/manager"
	"qa_test_server/model"
	"qa_test_server/web"
)

func Tcpserver() {
	go proxyStart(4003, 7777)

	fmt.Println("tcp server test---lf 2022-4-28")

	service := ":4001"
	tcpAddr, _ := net.ResolveTCPAddr("tcp", service)
	listerer, err := net.ListenTCP("tcp", tcpAddr)

	if err != nil {
		fmt.Printf("listen fail,err: %v\n", err)
		return
	}
	defer listerer.Close()
	for {
		//等待连接
		conn, err := listerer.AcceptTCP()
		if err != nil {
			fmt.Printf("accept fail, err: %v\n", err)
			continue
		}
		go process(conn)

	}

}

func process(conn *net.TCPConn) {

	//conn.SetReadDeadline()

	//conn.SetReadDeadline(time.Now() + time.Millisecond*100)
	defer conn.Close()
	//循环处理数据
	for {
		time.Sleep(time.Millisecond * 20)
		var buffer [4096]byte
		var len = 0
		//接受数据
	retry:
		n, err := conn.Read(buffer[len:])

		if err != nil {
			fmt.Printf("read from connect failed, err: %v\n", err)
			break
		}
		len += n
		if len < (model.MOD_ADDR_HOLDINGD_END * 2) {

			goto retry
		}

		//新建二进制流
		buf_w := &bytes.Buffer{}
		//读入六
		//var time [4]uint64 = 0

		//buf_w.Write(buffer[:8])
		buf_w.Write(buffer[:len])

		var temp model.Femto_msg_packed
		//fmt.Printf("packet_size check [%d,%d,%d]\r\n\n", unsafe.Sizeof(temp), unsafe.Sizeof(temp.Holding_reg), unsafe.Sizeof(temp.Input_reg))
		binary.Read(buf_w, binary.LittleEndian, &temp)
		//fmt.Printf("###packet=[len=%04d] %+v\r\n\n", len, temp)
		fmt.Printf("###packet=[len=%04d] \r\n\n", len)

		// //数据解析
		dev := model.Femto_Decode(temp)
		manager.ManagerGlabal.Update(dev)
		go func() {
			if d, err := json.Marshal(dev); err == nil {
				web.WsManager.Groupbroadcast("device_upload", d)
			}
		}()

		//fmt.Printf("##############################")

		//fmt.Printf("msg %d\n", device.Dev_cap.Sys_para.Seed_param.Freq)
		//device.DB.Debug().Create(&device.Dev_cap)

		//回写
		//if _, err = conn.Write(buffer[:n]); err != nil {
		//	fmt.Printf("write to client failed, err: %v\n", err)
		//	break
		//}
		//time.Sleep(time.Millisecond * 100)
	}
}
