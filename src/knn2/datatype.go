package knn2
import (
	"github.com/albrow/zoom"
)

type Finger struct{
    Label  int      `zoom:"index"`
    X      float64
    Y      float64
    Feature map[int64] float64
    zoom.DefaultData
}

type MapBaseSample struct{
    Features    map[int64]int
    Label       int
    Prediction  float64
    Timestamp   int64
    Mac         int64
}
