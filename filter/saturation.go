package filter

import (
	"gocv.io/x/gocv"
	"log"
	"time"
)

// SaturationWithBRG using BGR colorspace to increase or decrease the saturation of the image.
// The percent param is the percentage of saturation to increase or decrease.
// if percent is positive, the saturation will be increased.
// if percent is negative, the saturation will be decreased.
func SaturationWithBRG(img gocv.Mat, percent int64) gocv.Mat {
	start := time.Now()
	defer func(start time.Time) {
		dur := time.Since(start)
		log.Printf("func saturationWithBGR took %0.0fs to execute\n", dur.Seconds())
	}(start)
	channels := gocv.Split(img)
	var (
		row       = img.Rows()
		col       = img.Cols()
		increment = float64(percent) / float64(100)
		channelB  = channels[0]
		channelG  = channels[1]
		channelR  = channels[2]
	)
	// 默认BGR空间
	for i := 0; i < row; i++ {
		for j := 0; j < col; j++ {
			var (
				b, g, r      uint8
				nb, ng, nr   uint8
				delta, value float64
				L, S, alpha  float64
			)
			b = channelB.GetUCharAt(i, j)
			g = channelG.GetUCharAt(i, j)
			r = channelR.GetUCharAt(i, j)

			max := Max3(r, g, b)
			min := Min3(r, g, b)
			delta = (max - min) / 255
			//灰点不做处理
			if delta == 0 {
				continue
			}
			value = (max + min) / 255
			L = value / 2
			if L < 0.5 {
				S = delta / value
			} else {
				S = delta / (2 - value)
			}
			// 饱和度增加
			if increment >= 0 {
				if increment+S >= 1 {
					alpha = S
				} else {
					alpha = 1 - increment
				}
				alpha = 1/alpha - 1
				nr = r + uint8((float64(r)-L*255)*alpha)
				ng = g + uint8((float64(g)-L*255)*alpha)
				nb = b + uint8((float64(b)-L*255)*alpha)
			} else {
				// 饱和度降低
				alpha = increment
				nr = uint8(L*255 + (float64(r)-L*255)*(1+alpha))
				ng = uint8(L*255 + (float64(g)-L*255)*(1+alpha))
				nb = uint8(L*255 + (float64(b)-L*255)*(1+alpha))
			}
			channelB.SetUCharAt(i, j, nb)
			channelG.SetUCharAt(i, j, ng)
			channelR.SetUCharAt(i, j, nr)
		}
	}
	gocv.Merge([]gocv.Mat{channelB, channelG, channelR}, &img)
	return img
}

// SaturationWithHSV using HSV colorspace to increase or decrease the saturation of the image.
// The percent param is the percentage of saturation to increase or decrease.
// if percent is positive, the saturation will be increased.
// if percent is negative, the saturation will be decreased.
func SaturationWithHSV(img gocv.Mat, percent int64) gocv.Mat {
	start := time.Now()
	defer func(start time.Time) {
		dur := time.Since(start)
		log.Printf("func saturationWithHSV took %0.0fs to execute\n", dur.Seconds())
	}(start)
	hsv := gocv.NewMat()
	gocv.CvtColor(img, &hsv, gocv.ColorBGRToHSV)
	// split() the hsv image into 3 seperate (uchar) channels
	channels := gocv.Split(hsv)
	S := channels[1]
	for row := 0; row < S.Rows(); row++ {
		for col := 0; col < S.Cols(); col++ {
			originVal := S.GetUCharAt(row, col)
			val := float64(originVal) * (1 + float64(percent)/100)
			// 处理极值
			if percent > 0 && val > 255 {
				val = 255
			}
			if percent < 0 && val < 0 {
				val = 0
			}
			S.SetUCharAt(row, col, uint8(val))
		}
	}
	gocv.Merge([]gocv.Mat{channels[0], S, channels[2]}, &hsv)
	gocv.CvtColor(hsv, &img, gocv.ColorHSVToBGR)
	return img
}
