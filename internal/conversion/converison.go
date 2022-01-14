package conversion

import (
	"image"

	"github.com/disintegration/imaging"
)

type Convert struct{}

// Method that returns resized picture fith the provided ratio.
func (c *Convert) Resize(im image.Image, ratio float32) image.Image {
	newX := int(ratio * float32(im.Bounds().Max.X))
	newY := int(ratio * float32(im.Bounds().Max.Y))

	return imaging.Resize(im, newX, newY, imaging.Lanczos)
}
