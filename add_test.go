package intadd                                                                                                                                           

import (
    "testing"
)

func Test_add2int(t *testing.T){
    if( add2int(3,4)!=7){
        t.Error("add2ing cant work...")
    }else{
        t.Log("one test pass...")
    }
}
