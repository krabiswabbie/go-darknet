package darknet

/*
#cgo CFLAGS: -I /usr/include
#include <darknet.h>
#include "detection.h"
*/
import "C"
import (
	"image"
)

// BoundingBox represents a bounding box.
type BoundingBox struct {
	StartPoint image.Point
	EndPoint   image.Point
}

// Detection represents a detection.
type Detection struct {
	BoundingBox

	ClassIDs      []int
	ClassNames    []string
	Probabilities []float32
}

func makeDetection(w, h int, det *C.detection, threshold float32, classes int, classNames []string) *Detection {
	dClassIDs := make([]int, 0)
	dClassNames := make([]string, 0)
	dProbs := make([]float32, 0)

	for i := 0; i < int(classes); i++ {
		dProb := float32(
			C.get_detection_probability(det, C.int(i), C.int(classes)))
		if dProb > threshold {
			dClassIDs = append(dClassIDs, i)
			cN := classNames[i]
			dClassNames = append(dClassNames, cN)
			dProbs = append(dProbs, dProb*100)
		}
	}

	fImgW := C.float(w)
	fImgH := C.float(h)
	halfRatioW := det.bbox.w / 2.0
	halfRatioH := det.bbox.h / 2.0
	out := Detection{
		BoundingBox: BoundingBox{
			StartPoint: image.Point{
				X: int((det.bbox.x - halfRatioW) * fImgW),
				Y: int((det.bbox.y - halfRatioH) * fImgH),
			},
			EndPoint: image.Point{
				X: int((det.bbox.x + halfRatioW) * fImgW),
				Y: int((det.bbox.y + halfRatioH) * fImgH),
			},
		},

		ClassIDs:      dClassIDs,
		ClassNames:    dClassNames,
		Probabilities: dProbs,
	}

	return &out
}

func makeDetections(w, h int, detections *C.detection, detectionsLength int,
	threshold float32, classes int, classNames []string) []*Detection {
	// Make list of detection objects.
	ds := make([]*Detection, detectionsLength)
	for i := 0; i < int(detectionsLength); i++ {
		det := C.get_detection(detections, C.int(i), C.int(classes))
		d := makeDetection(w, h, det, threshold, classes, classNames)
		ds[i] = d
	}

	return ds
}
