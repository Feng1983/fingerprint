//通讯协议处理，主要处理封包和解包的过程
package main

import (
	"bytes"
	"encoding/binary"	
	//"net"
	//"unicode/utf8"
	//"fmt"
	//"strconv"
)

const (
	ConstHeader         = "wl"
	ConstHeaderLength   = 2
	//ConstSaveDataLength = 4
	ConstMacLength      = 6
	ConstRecordLength   = 2
)

//封包
func Packet(message []byte) []byte {
	return append(append([]byte(ConstHeader), IntToBytes(len(message))...), message...)
}

//解包

//unpack
func Unpack(buffer []byte, readerChannel chan []byte) []byte {
	length := len(buffer)
	//var i int
	if length < ConstHeaderLength+ConstMacLength+ConstRecordLength {
		//not  full package
		return buffer
	}
	head := string(buffer[:ConstHeaderLength])
	Log("head... ",head)
	if head == ConstHeader {
		//Log(buffer)
		Macaddr := BytesToMac(buffer[ConstHeaderLength : ConstHeaderLength+ConstMacLength])
		//Log(head,Macaddr,length)
		numRecord := BytesToInt16(buffer[ConstHeaderLength+ConstMacLength : ConstHeaderLength+ConstMacLength+ConstRecordLength])
		Log(Macaddr, numRecord, buffer[ConstHeaderLength+ConstMacLength : ConstHeaderLength+ConstMacLength+ConstRecordLength])
		Log(Macaddr, buffer[:10],length," num recor: ",numRecord)
		totLength := ConstHeaderLength + ConstMacLength + ConstRecordLength + 12 * numRecord + 1
		
		if totLength != length {
			return buffer
		} else {
			for i := ConstHeaderLength + ConstMacLength + ConstRecordLength; i < totLength-1;  {
				data := buffer[i : i+12]
				i = i + 12
				readerChannel <- data
			}
			readerChannel <-buffer[totLength-1:totLength]
			return make([]byte, 0)
		}
	}
	Log("error ,package....", buffer)
	return make([]byte, 0)
}

//整形转换成字节
func IntToBytes(n int) []byte {
	x := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.LittleEndian, x)
	return bytesBuffer.Bytes()
}

//字节转换成整形
func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)
	var x int32
	binary.Read(bytesBuffer, binary.LittleEndian, &x)
	return int(x)
}

func BytesToInt_Big(b []byte) int {
        bytesBuffer := bytes.NewBuffer(b)
        var x int32
        binary.Read(bytesBuffer, binary.BigEndian, &x)
        return int(x)
}


func BytesToInt16(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)
	var x int16
	binary.Read(bytesBuffer, binary.LittleEndian, &x)
	return int(x)
}

func BytesToInt16_Big(b []byte) int {
        bytesBuffer := bytes.NewBuffer(b)
        var x int16
        binary.Read(bytesBuffer, binary.BigEndian, &x)
        return int(x)
}



func BytesToString(b []byte) string {
	bytesBuffer := bytes.NewBuffer(b)
	var x string
	binary.Read(bytesBuffer, binary.LittleEndian, &x)
	return x
}

func BytesToInt8(b []byte)int8{
	bytesBuffer := bytes.NewBuffer(b)
        var x int8
        binary.Read(bytesBuffer, binary.LittleEndian, &x)
        return x
}
func BytesToInt8_Big(b []byte)int8{
        bytesBuffer := bytes.NewBuffer(b)
        var x int8
        binary.Read(bytesBuffer, binary.BigEndian, &x)
        return x
}
func BytesToUint8(b []byte) uint8{
	bytesBuffer := bytes.NewBuffer(b)
	var x uint8
	binary.Read(bytesBuffer, binary.LittleEndian, &x)
	return x
}
func BytesToInt64(b []byte)int64{
	bytesBuffer := bytes.NewBuffer(b)
        var x int64
        binary.Read(bytesBuffer, binary.LittleEndian, &x)
        return x
}

func BytesToMac(b[] byte) uint64{
	var r uint64
	var x uint8
	
	for i:=0;i<6;i++{
		x = BytesToUint8(b[i:i+1])
		r = r<<8 + uint64(x)	
	}
	return r 
}
func Int64toBytes(b int64) []byte{
	var buf =make([]byte, 8)
	binary.LittleEndian.PutUint64(buf,uint64(b))
	return buf 
}

//
func IsBigEndian() bool {
	var i int32 = 0x12345678
	var b byte = byte(i)
	if b == 0x12 {
		return true
	}
	return false
}

/*
func main(){
	//barray:= []byte{0 ,24, 248, 85, 54, 247}
	barray2:=[]byte{228,184,173,229,155,189}
	Log(string(barray2[:]))
	Log(string([]byte{248, 85, 54, 247}))
	Log(BytesToInt([]byte{63,0}))
	Log(BytesToInt16([]byte{63,0}))
	Log(BytesToInt64([]byte{0x00 ,0x18, 0xF8, 0x55, 0x36, 0xF7}))
	fmt.Printf("byte is ... %x",BytesToInt64([]byte{228,184,173,229,155,189}))
	Log(BytesToMac([]byte{0x00 ,0x18, 0xF8, 0x55, 0x36, 0xF7}))
	fmt.Printf("mac byte is ... %x\n",BytesToMac([]byte{0x00 ,0x18, 0xF8, 0x55, 0x36, 0xF7}))
	fmt.Printf("over is ... %d\n", BytesToMac([]byte{0xc2,0xae,0x03,0xae,0xba,0x28}))
	Log(BytesToInt([]byte{0x0F ,0x9C, 0x0E ,0x55}))
	//Log(Int64toBytes(214052641880616))
	fmt.Printf("mac is ... %x\n",uint64(214052641880616))
	vvvv:= IsBigEndian()
	fmt.Println(vvvv)
	fmt.Println([]byte("wl"))	
}*/
