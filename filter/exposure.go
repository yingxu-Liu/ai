package filter

import (
	"gocv.io/x/gocv"
	"log"
	"math"
	"time"
)

// Exposure applies an exposure filter to the given image
// The gamma parameter controls image brightness.
// If gamma < 1 then the image will be darker
// If gamma > 1 then the image will be lighter
// gamma = 1 has no effect
func Exposure(img gocv.Mat, gamma float64) gocv.Mat {
	start := time.Now()
	defer func(start time.Time) {
		dur := time.Since(start)
		log.Printf("func gammaCorrection took %0.0fs to execute\n", dur.Seconds())
	}(start)
	invGamma := 1 / gamma

	wbLUT := gocv.NewMatWithSize(1, 256, gocv.MatTypeCV8U)
	for i := 0; i < 256; i++ {
		wbLUT.SetUCharAt(0, i, uint8(math.Pow(float64(i)/255, invGamma)*255))
	}
	dst := gocv.NewMat()
	gocv.LUT(img, wbLUT, &dst)
	return dst
}
