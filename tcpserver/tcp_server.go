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
	"time"
	"unsafe"
)

func Tcpserver() {
	go proxyStart(4003, 7777)

	//
	go packed_test("Virtual Test")

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

func FemtoDeviceTest() {
	//panic("unimplemented")
}

func process(conn *net.TCPConn) {

	var temp model.Femto_msg_packed

	//conn.SetReadDeadline(time.Now() + time.Millisecond*100)

	defer conn.Close()

	var len int = 0
	//var fram_start int = 0
	var time_pre, time_cur int64
	var packet_len int = int(2 * (model.MOD_ADDR_HOLDINGD_END + model.MOD_ADDR_INPUT_END))
	//循环处理数据
	for {
		//time.Sleep(time.Millisecond * 100)
		buffer := make([]byte, 8192)
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
			//fram_start = len
			len = 0
		}
		//fmt.Printf("###time=[cur=%04d, now=%d,diff=%d] \r\n\n", time_cur, time_pre, diff)

		//fmt.Printf("n==%d,len==%d,packet=%d\n", n, len, int(unsafe.Sizeof(temp)))
		len += n
		//if len < int(2*(model.MOD_ADDR_HOLDINGD_END+model.MOD_ADDR_INPUT_END)) {
		if len < packet_len {
			goto retry
		}

		// if len != packet_len {
		// 	len = 0
		// 	continue
		// }

		fmt.Printf("n==%d,len==%d\n", n, len)
		//新建二进制流

		// for i := 0; i < 1000; i++ {
		// 	if buffer[i] == 0x55 && buffer[i-1] == 0xaa {
		// 		fmt.Printf("###head check=[len=%d]\n", i)
		// 	}

		// }

		buf_w := &bytes.Buffer{}
		//读入六
		//var time [4]uint64 = 0

		//buf_w.Write(buffer[:8])
		buf_w.Write(buffer[0:])
		//fmt.Printf("###packet=[len=%04d],[size of temp=%d] \r\n", len, int(unsafe.Sizeof(temp.Femto_holding_reg)))

		fmt.Printf("packet_size check [%d,%d,%d]\r\n\n", unsafe.Sizeof(temp), unsafe.Sizeof(temp.Femto_holding_reg), unsafe.Sizeof(temp.Femto_input_reg))
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

		//fmt.Printf("msg %+v\n", temp.Femto_input_reg.Time)

		//device.DB.Debug().Create(&device.Dev_cap)

		//回写
		//if _, err = conn.Write(buffer[:n]); err != nil {
		//	fmt.Printf("write to client failed, err: %v\n", err)
		//	break
		//}
		//time.Sleep(time.Millisecond * 100)
	}
}

func packed_test(sn string) {

	var temp model.Femto_msg_packed
	buffer := make([]byte, 8192)

	buf_w := &bytes.Buffer{}
	//读入六
	//var time [4]uint64 = 0

	//buf_w.Write(buffer[:8])
	buf_w.Write(buffer[0:])
	//fmt.Printf("###packet=[len=%04d],[size of temp=%d] \r\n", len, int(unsafe.Sizeof(temp.Femto_holding_reg)))

	fmt.Printf("packet_size check [%d,%d,%d]\r\n\n", unsafe.Sizeof(temp), unsafe.Sizeof(temp.Femto_holding_reg), unsafe.Sizeof(temp.Femto_input_reg))
	binary.Read(buf_w, binary.LittleEndian, &temp)
	//fmt.Printf("###packet=[len=%04d] %+v\r\n\n", len, temp)
	//fmt.Printf("###packet111=[len=%04d,head=%d] \r\n\n", len, temp.Holding_reg.Laser_para.Head)

	copy(temp.Femto_holding_reg.Laser_para.Laser_info.SN[:], sn)
	// //数据解析

	for {

		dev := model.Femto_Decode(temp)
		manager.ManagerGlabal.Update(dev)
		go func() {
			if d, err := json.Marshal(dev); err == nil {
				web.WsManager.Groupbroadcast("device_upload", d)
			}
		}()
		temp.Femto_input_reg.Bate.Hardware_bate++
		temp.Femto_input_reg.Time.Uptime[0]++
		temp.Femto_input_reg.Time.Uptime[1] += 10
		time.Sleep(time.Millisecond * 500)

		//para def
	}

}
