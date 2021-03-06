package conversion

import (
	"image"

	"github.com/disintegration/imaging"
)

// Function that returns resized picture fith the provided ratio.
func Resize(im image.Image, ratio float32) image.Image {
	newX := int(ratio * float32(im.Bounds().Max.X))
	newY := int(ratio * float32(im.Bounds().Max.Y))

	return imaging.Resize(im, newX, newY, imaging.Lanczos)
}
