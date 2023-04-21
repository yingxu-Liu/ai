package filter

import (
	"gocv.io/x/gocv"
	"log"
	"math"
	"time"
)

// Shadow changes the darkest part of the image
// The light param is a value to that increase or decrease the shadow of the image.
// If light is positive, the shadow will be brighter.
// If light is negative, the shadow will be darker.
func Shadow(img gocv.Mat, light int64) gocv.Mat {
	start := time.Now()
	defer func(start time.Time) {
		dur := time.Since(start)
		log.Printf("func shadow took %0.0fs to execute\n", dur.Seconds())
	}(start)
	// 生成灰度图
	gray := gocv.Zeros(img.Rows(), img.Cols(), gocv.MatTypeCV32FC1)
	//gocv.CvtColor(img, &gray, gocv.ColorBGRToGray)
	f := img.Clone()
	f.ConvertTo(&f, gocv.MatTypeCV32FC3)
	pics := gocv.Split(f)
	for row := 0; row < gray.Rows(); row++ {
		for col := 0; col < gray.Cols(); col++ {
			gray.SetFloatAt(row, col, pics[0].GetFloatAt(row, col)*0.299+pics[1].GetFloatAt(row, col)*0.587+pics[2].GetFloatAt(row, col)*0.114)
		}
	}
	gray.DivideFloat(255)

	// 确定阴影区 通过(1-gray)*(1-gray)，得到的图像中原本亮的地方则为亮，取平均值当阈值，进行二值化得到掩膜mask
	thresh := gocv.Zeros(gray.Rows(), gray.Cols(), gray.Type())
	for row := 0; row < gray.Rows(); row++ {
		for col := 0; col < gray.Cols(); col++ {
			val := gray.GetFloatAt(row, col)
			thresh.SetFloatAt(row, col, (1-val)*(1-val))
		}
	}
	t := thresh.Mean()
	mask := gocv.Zeros(gray.Rows(), gray.Cols(), gocv.MatTypeCV8UC1)
	// 二值化 将大于阈值的像素点置为255，小于阈值的置为0
	gocv.Threshold(thresh, &mask, float32(t.Val1), 255, gocv.ThresholdBinary)
	// 对掩膜区边缘进行平滑过渡。假设light为50，那么midrate的掩膜区值为1.5，黑色区为1，过渡区为1~1.5；bright的掩膜区为0.125，黑色区为0，过渡区为0~0.125
	var (
		max        = float32(4.0)
		bright     = float32(light) / 100.0 / max
		mid        = 1 + max*bright
		midrate    = gocv.Zeros(img.Rows(), img.Cols(), gocv.MatTypeCV32FC1)
		brightrate = gocv.Zeros(img.Rows(), img.Cols(), gocv.MatTypeCV32FC1)
		result     = gocv.Zeros(img.Rows(), img.Cols(), img.Type())
	)
	for row := 0; row < img.Rows(); row++ {
		for col := 0; col < img.Cols(); col++ {
			if mask.GetUCharAt(row, col) == 255 {
				midrate.SetFloatAt(row, col, mid)
				brightrate.SetFloatAt(row, col, bright)
				continue
			}
			midrate.SetFloatAt(row, col, (mid-1.0)/float32(t.Val1)*thresh.GetFloatAt(row, col)+1.0)
			brightrate.SetFloatAt(row, col, 1.0/float32(t.Val1)*thresh.GetFloatAt(row, col)*bright)
		}
	}
	// 根据midrate和brightrate，进行高光区提亮
	channels := gocv.Split(img)
	var (
		channelB = channels[0]
		channelG = channels[1]
		channelR = channels[2]
	)
	for row := 0; row < img.Rows(); row++ {
		for col := 0; col < img.Cols(); col++ {
			for c := 0; c < 3; c++ {
				var (
					r = channelR.GetUCharAt(row, col)
					g = channelG.GetUCharAt(row, col)
					b = channelB.GetUCharAt(row, col)
				)
				channelR.SetUCharAt(row, col, calculateTemp(r, row, col, midrate, brightrate))
				channelG.SetUCharAt(row, col, calculateTemp(g, row, col, midrate, brightrate))
				channelB.SetUCharAt(row, col, calculateTemp(b, row, col, midrate, brightrate))
			}
		}
	}
	gocv.Merge([]gocv.Mat{channelB, channelG, channelR}, &result)
	return result
}

func calculateTemp(val uint8, row, col int, midrate gocv.Mat, brightrate gocv.Mat) uint8 {
	temp := math.Pow(float64(val)/255.0, float64(1.0/midrate.GetFloatAt(row, col)*(1.0/(1.0-brightrate.GetFloatAt(row, col)))))
	if temp > 1.0 {
		temp = 1.0
	}
	if temp < 0.0 {
		temp = 0
	}
	return uint8(temp * 255.0)
}
