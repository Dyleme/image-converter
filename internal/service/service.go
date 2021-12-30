package service

import (
	"bytes"
	"context"
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

// decodeImage decodes image from the r.
// Decoding supports only jpeg and png types.
func decodeImage(r io.Reader, imgType string) (image.Image, error) {
	switch imgType {
	case pngType:
		return png.Decode(r)
	case jpegType:
		return jpeg.Decode(r)
	default:
		return nil, ErrUnsupportedType
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
			return bf.Bytes(), err
		}
	case jpegType:
		if err := jpeg.Encode(bf, i, &jpeg.Options{Quality: jpegQuality}); err != nil {
			return bf.Bytes(), err
		}
	}

	return nil, ErrUnsupportedType
}
