package main

import(
	"encoding/binary"
	"bytes"
)

const (
	ConstHeader =""
	ConstHeaderLength = 15
	ConstSaveDataLenght = 4
)
//packet
func Packet(message []byte) []byte{
	return append(append([]byte(ConstHeader),Int2Bytes(len(message))...), message...)
}

//unpack
func Unpack(buffer[]byte, readerChannel chan[]byte) []byte {
	length:= len(buffer)
	var i int
	for i=0; i< length; i=i+1{
		if length < i+ConstHeaderLength +ConstSaveDataLenght{
			break
		}
		if string(buffer[i:i+ConstHeaderLength])==ConstHeader{
			messageLength:= Bytes2Int(buffer[i+ConstHeaderLength:i+ConstHeaderLength+ConstSaveDataLenght])
			if length < messageLength {
				break
			}
			data := buffer[i+ConstHeaderLength+ConstSaveDataLenght:i+ConstHeaderLength+ConstSaveDataLenght+messageLength]
			readerChannel<-data
			i+= ConstHeaderLengt+ConstSaveDataLenght-1
			
		}
	}
}

func Int2Bytes(n int) []byte{
	x := int32(n)
	bytesBuffer :=bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer,binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

func Bytes2Int(b []byte) int{
	bytesBuffer:= bytes.NewBuffer(b)
	var x int32
	binary.Read(bytesBuffer,binary.BigEndian, &x)
	return int(x)
}

