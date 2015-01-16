
package main

import(
	"fmt"
)
type Element struct{
	Key  int
	Value interface{}
} 
type MinHeap struct{
	size int
	max  int
	items []*Element
}
func NewMinHeap(max int) *MinHeap{
	return &MinHeap{size:0, max:max, items:make([]*Element,max,max)}
}

func (h *MinHeap)full() bool{
	return h.size >=h.max
}

func (h *MinHeap)empty() bool{
	return h.size<=0
}

func (h *MinHeap)Min()  *Element{
	if h.empty() {
		return nil
	}else{
		return h.items[0]
	}
}

func (h *MinHeap)add(e *Element){
	if e==nil{
		return
	}
	
	if h.full(){
		h.grow()
	}

	s:= h.size
	h.size++
	for{
		if s<= 0 {
			break
		}

