package knn2

import (
	"github.com/xlvector/hector/core"
	"github.com/xlvector/hector/eval"
	"math"
	//"math/rand"
	"strconv"
	//"fmt"
)

type KNN struct {
	sv     []*core.Vector
	labels []int
	k      int64
	Xset   map[int64]float64
	Yset   map[int64]float64
}

func (self *KNN) SaveModel(path string) {

}

func (self *KNN) LoadModel(path string) {

}

func (c *KNN) Init(params map[string]string) {
	K, _ := strconv.ParseInt(params["k"], 10, 64)
	c.k = K
	c.Xset = make(map[int64]float64)
	c.Yset = make(map[int64]float64)
}

func (c *KNN) Kernel(x, y *core.Vector) float64 {
	z := x.Copy()
	//z.AddVector(y, -1.0)
	for key,_ := range z.Data{
		_, ok := y.Data[key]
		if ok {
			z.Data[key] += y.Data[key]*(-1.0)
		} else {
			z.Data[key] += (-100.0)
		}
	}
	ret := math.Exp(-1.0 * z.NormL2() / 20.0)
	//ret := z.NormL2()
	return ret
}



func (c *KNN) Predict2(sample *MapBaseSample)(x,y float64){
	x,y= c.PredictMultiClass2(sample)
	return 
}


func (c *KNN) Train2(fingers []*Finger2){
	c.sv = []*core.Vector{}
	c.labels = []int{}
	if len(fingers)< 1000{
		for _,fg :=range fingers{
			ret:= core.NewVector()
			for k,v :=range fg.Feature{
				ret.SetValue(k,v)
			}
			c.sv = append(c.sv,ret)
			c.labels = append(c.labels, fg.Label)
			c.Xset[int64(fg.Label)] = fg.X
			c.Yset[int64(fg.Label)] = fg.Y
		} 
	}else{
		
	}
}


func (c *KNN) PredictMultiClass2(sample *MapBaseSample) (float64,float64){
	x := core.NewVector()
	for k, v :=range sample.Features {
		x.SetValue(k,float64(v))
	}
	predictions := []* eval.LabelPrediction{}
	for i, s :=range c.sv{
		predictions = append(predictions,&(eval.LabelPrediction{Label: c.labels[i], Prediction: c.Kernel(s, x)}))
	}
	compare := func(p1, p2 *eval.LabelPrediction) bool {
        return p1.Prediction > p2.Prediction
    }
	eval.By(compare).Sort(predictions)
	
	var tx,ty float64
	tx =0.0
	ty =0.0
	for i, pred:= range predictions{
		if i>= int(c.k){
			break
		}
		tx += c.Xset[int64(pred.Label)]
		ty += c.Yset[int64(pred.Label)]
		//fmt.Println(c.Xset[int64(pred.Label)],c.Yset[int64(pred.Label)],i,pred.Prediction)
	}
	return tx/float64(c.k), ty/float64(c.k)
}

