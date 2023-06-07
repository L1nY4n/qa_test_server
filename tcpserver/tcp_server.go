package tcpserver

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"time"
	"unsafe"

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

	//test

}

func FemtoDeviceTest() {
	//panic("unimplemented")
}

func process(conn *net.TCPConn) {

	var temp model.Femto_msg_packed

	//conn.SetReadDeadline(time.Now() + time.Millisecond*100)

	defer conn.Close()

	buffer := make([]byte, 8192)
	var len int = 0
	var fram_start int = 0
	var time_pre, time_cur int64
	//循环处理数据
	for {
		//time.Sleep(time.Millisecond * 100)

	retry:

		n, err := conn.Read(buffer[len:])
		//异常处理
		if err != nil {
			fmt.Printf("read from connect failed, err: %v\n", err)
			break
		}
		//记录时间
		time_cur = time.Now().Local().UnixNano() / 1e6
		diff := time_cur - time_pre
		time_pre = time_cur
		if diff > 10 {
			//帧头识别
			fram_start = len
			len = 0
		}
		//fmt.Printf("###time=[cur=%04d, now=%d,diff=%d] \r\n\n", time_cur, time_pre, diff)

		//	fmt.Printf("n==%d,len==%d\n", n, len)
		len += n
		if len < int(unsafe.Sizeof(temp)) {

			goto retry
		}

		if len != int(unsafe.Sizeof(temp)) {
			len = 0
			continue
		}

		//新建二进制流
		buf_w := &bytes.Buffer{}
		//读入六
		//var time [4]uint64 = 0

		//buf_w.Write(buffer[:8])
		buf_w.Write(buffer[fram_start:])
		//fmt.Printf("###packet=[len=%04d],[size of temp=%d] \r\n", len, int(unsafe.Sizeof(temp.Holding_reg)))

		//fmt.Printf("packet_size check [%d,%d,%d]\r\n\n", unsafe.Sizeof(temp), unsafe.Sizeof(temp.Holding_reg), unsafe.Sizeof(temp.Input_reg))
		binary.Read(buf_w, binary.LittleEndian, &temp)
		//fmt.Printf("###packet=[len=%04d] %+v\r\n\n", len, temp)
		//fmt.Printf("###packet111=[len=%04d,head=%d] \r\n\n", len, temp.Holding_reg.Laser_para.Head)

		if temp.Femto_holding_reg.Laser_para.Head != 0x55aa {
			len = 0
			continue
		}
		// //数据解析
		dev := model.Femto_Decode(temp)
		manager.ManagerGlabal.Update(dev)
		go func() {
			if d, err := json.Marshal(dev); err == nil {
				web.WsManager.Groupbroadcast("device_upload", d)
			}
		}()

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
