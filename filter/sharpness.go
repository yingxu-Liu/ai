package filter

import (
	"gocv.io/x/gocv"
	"log"
	"time"
)

// SharpenWithSobel applies a sharpening filter by using Sobel algorithm to the given image
// The level parameter controls the strength of the filter
// The level parameter must be an odd number
func SharpenWithSobel(img gocv.Mat, level int) gocv.Mat {
	start := time.Now()
	defer func(start time.Time) {
		dur := time.Since(start)
		log.Printf("func sharpenWithSobel took %0.0fs to execute\n", dur.Seconds())
	}(start)
	dst := gocv.NewMat()
	//src：输入图像 dst:输出图像 ddepth:图像数据深度 dx: x方向导数 dy：y方向导数 ksize: Sobel 核的大小；它必须是1、3、5或7 (单数) scale: 可选比例因子,默认情况下不应用缩放 delta:可选的在结果上添加偏移值 borderType:边界类型
	gocv.Sobel(img, &dst, gocv.MatTypeCV8U, 1, 0, level, 1, 0, gocv.BorderDefault)
	gocv.BitwiseOr(img, dst, &img)
	return img
}

// SharpenWithLaplacian applies a sharpening filter by using Laplacian algorithm to the given image
// The level parameter controls the strength of the filter
// The level parameter must be an odd number
func SharpenWithLaplacian(img gocv.Mat, level int) gocv.Mat {
	start := time.Now()
	defer func(start time.Time) {
		dur := time.Since(start)
		log.Printf("func sharpenWithLaplacian took %0.0fs to execute\n", dur.Seconds())
	}(start)
	dst := gocv.NewMat()
	//src：输入图像 dst:输出图像 ddepth:图像数据深度 ksize: 用于计算二阶导数滤波器的大小 scale: 可选比例因子,默认情况下不应用缩放 delta:可选的在结果上添加偏移值 borderType:边界类型
	gocv.Laplacian(img, &dst, gocv.MatTypeCV8U, level, 1, 0, gocv.BorderDefault)
	gocv.BitwiseOr(img, dst, &img)
	return img
}
