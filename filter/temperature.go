package filter

import (
	"gocv.io/x/gocv"
	"log"
	"math"
	"time"
)

// Temperature increases or decreases the temperature of the given image.
// The level parameter is the amount of change to apply to the image.
// if level is positive, the temperature is increased.
// if level is negative, the temperature is decreased.
func Temperature(img gocv.Mat, level int32) gocv.Mat {
	start := time.Now()
	defer func(start time.Time) {
		dur := time.Since(start)
		log.Printf("func temperature took %0.0fs to execute\n", dur.Seconds())
	}(start)
	// 提高色温即对红绿通道增强、蓝色通道减弱
	// 降低色温即对蓝色通道增强、红绿通道减弱
	channels := gocv.Split(img)
	var (
		rows     = img.Rows()
		cols     = img.Cols()
		channelB = channels[0]
		channelG = channels[1]
		channelR = channels[2]
	)
	level = level / 2
	// 默认BGR空间
	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			var (
				r          = channelR.GetUCharAt(row, col)
				g          = channelG.GetUCharAt(row, col)
				b          = channelB.GetUCharAt(row, col)
				nb, ng, nr int32
			)

			nb = int32(b) - level
			ng = int32(g) + level
			nr = int32(r) + level

			nr = int32(math.Min(float64(nr), float64(255)))
			nr = int32(math.Max(float64(nr), float64(0)))

			ng = int32(math.Min(float64(ng), float64(255)))
			ng = int32(math.Max(float64(ng), float64(0)))

			nb = int32(math.Min(float64(nb), float64(255)))
			nb = int32(math.Max(float64(nb), float64(0)))

			if row == 1024 && col > 1024 && col < 1034 {
				log.Println(b, g, r, level, nb, ng, nr)
			}
			channelB.SetUCharAt(row, col, uint8(nb))
			channelG.SetUCharAt(row, col, uint8(ng))
			channelR.SetUCharAt(row, col, uint8(nr))
		}
	}
	gocv.Merge([]gocv.Mat{channelB, channelG, channelR}, &img)
	return img
}
