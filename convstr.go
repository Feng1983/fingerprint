
package main

import (
	"log"
	"time"
	"fmt"
)

type ByteSize float64

const (
	_  =	iota
	KB ByteSize = 1<<( 10*iota)
	MB
	GB
	TB
	PB
	EB
	ZB
	YB
)

var a, b, c int

func (b ByteSize)String() string{
	switch {
		case b>=YB:
			return fmt.Sprintf("%.2fYB",b/YB)
		case b>=ZB:
			return fmt.Sprintf("%.2fZB",b/ZB)
		case b>=EB:
			return fmt.Sprintf("%.2fEB",b/EB)
		case b>=PB:
			return fmt.Sprintf("%.2fPB",b/PB)
		case b>=TB:
			return fmt.Sprintf("%.2fTB",b/TB)
		case b>=GB:
			return fmt.Sprintf("%.2fGB",b/GB)
		case b>=MB:
			return fmt.Sprintf("%.2fMB",b/MB)
		case b>=KB:
			return fmt.Sprintf("%.2fKB",b/KB)
	}
	return fmt.Sprintf("%.2fB",b)
}
func main(){
	a=1
	b=2
	go func(){
		c= a+2
