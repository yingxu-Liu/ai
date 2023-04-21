package image

import (
	"gocv.io/x/gocv"
	"log"
)

// ReadImage reads an image from the given path and returns it as a gocv.Mat
func ReadImage(path string) gocv.Mat {
	img := gocv.IMRead(path, gocv.IMReadColor)
	if img.Empty() {
		log.Fatalf("Failed to open image: %s", path)
	}
	return img
}

// WriteImage writes the given image to the given path
func WriteImage(path string, img gocv.Mat) {
	if !gocv.IMWrite(path, img) {
		log.Fatalf("Failed to save image: %s", path)
	}
}
