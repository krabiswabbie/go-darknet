package darknet

// #include <darknet.h>
import "C"
import (
	"errors"
)

// Image represents the image buffer.
type Image struct {
	Width  int
	Height int

	image C.image
}

var errUnableToLoadImage = errors.New("unable to load image")

// Close and release resources.
func (img *Image) Close() error {
	C.free_image(img.image)
	return nil
}
