package utils

import (
	"bytes"
	"github.com/liujiawm/graphics-go/graphics"
	"github.com/rwcarlsen/goexif/exif"
	"image"
	"io"
	"log"
	"math"
)

// ReadOrientation - return the image orientation from reader.
func ReadOrientation(reader io.Reader) (int, error) {
	x, err := exif.Decode(reader)
	if err != nil && exif.IsCriticalError(err) {
		log.Println("[ReadOrientation] Error Decode")
		return 0, nil
	}

	orientation, err := x.Get(exif.Orientation)
	if err != nil {
		log.Println("[ReadOrientation] Error Get")
		return 0, nil
	}
	orientVal, err := orientation.Int(0)
	if err != nil {
		log.Println("[ReadOrientation] Error Int")
		return 0, nil
	}

	return orientVal, nil
}

// RotateImage - Rotate an image by the giving angle.
func RotateImage(src []byte, angle int) (image.Image, error) {
	var img, _, err = image.Decode(bytes.NewReader(src))
	if err != nil {
		log.Println("[RotateImage] Error Decode")
		return nil, err
	}
	angle = angle % 360

	// Radian conversion
	radian := float64(angle) * math.Pi / 180.0
	cos := math.Cos(radian)
	sin := math.Sin(radian)
	// The width and height of the original image
	w := float64(img.Bounds().Dx())
	h := float64(img.Bounds().Dy())

	// New image height and width
	W := int(math.Max(math.Abs(w*cos-h*sin), math.Abs(w*cos+h*sin)))
	H := int(math.Max(math.Abs(w*sin-h*cos), math.Abs(w*sin+h*cos)))

	dst := image.NewNRGBA(image.Rect(0, 0, W, H))

	err = graphics.Rotate(dst, img, &graphics.RotateOptions{Angle: radian})
	if err != nil {
		log.Println("[RotateImage] Error Rotate")
		return nil, err
	}

	return dst, nil
}
