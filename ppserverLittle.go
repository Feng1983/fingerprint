//服务端解包过程
package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	service := ":29988"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	CheckError(err)

	netListen, err := net.ListenTCP("tcp", tcpAddr)
	CheckError(err)

	defer netListen.Close()

	Log("Waiting for clients")
	for {
		conn, err := netListen.Accept()
		if err != nil {
			continue
		}

		Log(conn.RemoteAddr().String(), " tcp connect success")
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	//声明一个临时缓冲区，用来存储被截断的数据
	tmpBuffer := make([]byte, 0)

	//声明一个管道用于接收解包的数据
	readerChannel := make(chan []byte, 1024)
	go reader(readerChannel, conn)

	buffer := make([]byte, 10240)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			Log(conn.RemoteAddr().String(), " connection error: ", err)
			return
		}
		tmpBuffer = Unpack(append(tmpBuffer, buffer[:n]...), readerChannel)
	}
}

func reader(readerChannel chan []byte, conn net.Conn) {
	for {
		select {
		case data := <-readerChannel:
			if len(data)==1{
				Log("send wl to server")
				conn.Write([]byte("wl"))
			}else{
				//Log(data)
				Log(BytesToInt(data[0:4]), "; ",BytesToMac(data[4:10]),"; ", BytesToInt8(data[10:11]),"; ", BytesToInt8(data[11:12]),data)
				
			}
		}
	}
}

func Log(v ...interface{}) {
	fmt.Println(v...)
}

func CheckError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

