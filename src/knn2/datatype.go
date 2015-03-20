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

type ProcessData struct{
	Timestamp	int64	`zoom:"index"`
	Mac			int64	`zoom:"index"`
	X			float64
	Y			float64
	zoom.DefaultData
}
type FingerOri struct{
	Id		int  `gorm:"column:id"`
	X		float64	`gorm:"column:x"`
	Y		float64	`gorm:"column:y"`
	Ap1		int		`gorm:"column:ap1_no_scan"`
	Ap2		int		`gorm:"column:ap2_no_scan"`
	Ap3		int		`gorm:"column:ap3_no_scan"`
	Freq1	float64	`gorm:"column:ap1_mean_ss"`
	Freq2	float64 `gorm:"column:ap2_mean_ss"`
	Freq3	float64	`gorm:"column:ap3_mean_ss"`
	//zoom.DefaultData
}

type Finger2 struct{
    Id	   int      `zoom:"index"`
	Label  int      `zoom:"index"`
    X      float64
    Y      float64
    Feature map[int64] float64
    zoom.DefaultData
}
