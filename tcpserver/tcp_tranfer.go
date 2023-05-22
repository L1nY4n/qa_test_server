package tcpserver

import (
	"fmt"
	"net"
	//"time"
)

// Start a proxy server listen on fromport
// this proxy will then forward all request from fromport to toport
//
// Notice: a service must has been started on toport
func proxyStart(fromport, toport int) {

	fmt.Println("tcp transfer test---lf 2022-5-10")

	//建立tcp服务器
	proxyaddr := fmt.Sprintf(":%d", fromport)
	proxylistener, err := net.Listen("tcp", proxyaddr)
	if err != nil {
		fmt.Printf("Unable to listen on: %s, error: %s\n", proxyaddr, err.Error())
		//	os.Exit(1)
	}
	defer proxylistener.Close()

	targetaddr := fmt.Sprintf(":%d", toport)
	//targetaddr := fmt.Sprintf("localhost:%d", toport)
	targetlisterner, err := net.Listen("tcp", targetaddr)

	if err != nil {
		fmt.Printf("Unable to connect to: %s, error: %s\n", targetaddr, err.Error())
		//proxyconn.Close()
		//continue
	}
	defer targetlisterner.Close()
	fmt.Println("wait for connect..")

	for {
		//建立连接
		proxyconn, err1 := proxylistener.Accept()
		fmt.Println("wait for connect..")
		targetconn, err2 := targetlisterner.Accept()
		fmt.Println("wait for connect..")

		if err1 != nil || err2 != nil {

			fmt.Printf("Unable to accept a request, error: %s,%s", err1.Error(), err2.Error())
			//time.Sleep(time.Millisecond * 500)

			continue

		}
		fmt.Printf("init_succeful..\r\n")
		go proxyRequest(proxyconn, targetconn)
		go proxyRequest(targetconn, proxyconn)

	}
}

// Forward all requests from r to w
func proxyRequest(r net.Conn, w net.Conn) {
	defer r.Close()
	defer w.Close()
	fmt.Printf("data....\r\n")
	var buffer = make([]byte, 40960)
	for {
		n, err := r.Read(buffer)
		if err != nil {
			fmt.Printf("Unable to read from input, error: %s\n", err.Error())
			break
		}
		fmt.Printf("get_data...%d\r\n", n)
		_, err = w.Write(buffer[:n])
		if err != nil {
			fmt.Printf("Unable to write to output, error: %s\n", err.Error())
			break
		}
	}
}
