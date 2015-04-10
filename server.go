package main
import(
	"log"
	"fmt"
	"net"
	"os"
)

func main(){
	netListen, err := net.Listen("tcp",":9988")
	CheckError(err)
	
	defer netListen.Close()
	Log("Waiting for clients")

	for{
		conn, err:= netListen.Accept()
		if err!=nil{
			continue
		}
		Log(conn.RemoteAddr().String(),"tcp connect success")
		go handleConnection(conn)
      }
}

func handleConnection(conn net.Conn){
	tmpBuffer := make([]byte,0)
	//
	readerChannel :=make(chan[]byte,16)
	go reader(readerChannel)
	buffer:=make([]byte, 1024)
	for{
		n, err:= conn.Read(buffer)
		if err!=nil{
			Log(conn.RemoteAddr().String()," connection error")
			return
		}
		tmpBuffer = protocol.Unpack(append(tmpBuffer,buffer[:n]...),readerChannel)
	}
}

func reader(readerChannel chan[]byte){
	for {
		select {
			case data:= <- readerChannel:
				Log(string(data))
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
