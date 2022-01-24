package service

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
)

// Storager is an interface to interact with the file storage.
type Storager interface {
	// GetFile is used to take file from the storage.
	GetFile(ctx context.Context, path string) ([]byte, error)

	// UploadFile is used to add the file to the storage.
	UploadFile(ctx context.Context, userID int, fileName string, data []byte) (string, error)

	// DeleteFile is used to delete file from the storage.
	DeleteFile(ctx context.Context, path string) error
}

const (
	jpegType = "jpeg"
	pngType  = "png"
)

const (
	jpegQuality = 100
)

type UnsupportedTypeError struct {
	UnType string
}

func (e UnsupportedTypeError) Error() string {
	return fmt.Sprintf("unsupported type: %q", e.UnType)
}

// decodeImage decodes image from the r.
// Decoding supports only jpeg and png types.
func decodeImage(r io.Reader, imgType string) (image.Image, error) {
	switch imgType {
	case pngType:
		return png.Decode(r)
	case jpegType:
		return jpeg.Decode(r)
	default:
		return nil, &UnsupportedTypeError{imgType}
	}
}

// getResolution function returns the resolution of the image.
func getResolution(i image.Image) (width, height int) {
	return i.Bounds().Dx(), i.Bounds().Dy()
}

// encodeImage encode image with the provided image type, returns bytes of the encoded image.
func encodeImage(i image.Image, imgType string) ([]byte, error) {
	bf := new(bytes.Buffer)

	switch imgType {
	case pngType:
		if err := png.Encode(bf, i); err != nil {
			return nil, err
		}

	case jpegType:
		if err := jpeg.Encode(bf, i, &jpeg.Options{Quality: jpegQuality}); err != nil {
			return nil, err
		}

	default:
		return nil, &UnsupportedTypeError{imgType}
	}

	return bf.Bytes(), nil
}

// Resizer is an interface, which provide method to resize image.
type Resizer interface {
	Resize(im image.Image, ratio float32) image.Image
}
