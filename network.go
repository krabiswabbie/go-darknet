package darknet

/*
#cgo CFLAGS: -I /usr/include
#include <stdlib.h>
#include <darknet.h>
#include "network.h"
*/
import "C"
import (
	"errors"
	"time"
	"unsafe"
)

// YOLONetwork represents a neural network using YOLO.
type YOLONetwork struct {
	GPUDeviceIndex           int
	DataConfigurationFile    string
	NetworkConfigurationFile string
	WeightsFile              string
	Threshold                float32

	ClassNames []string
	Classes    int

	cNet                *C.network
	hierarchalThreshold float32
	nms                 float32

	// GPU input buffer size
	Width, Height, Colors int32
	GpuInput              *float32
}

// DetectionResult represents the inference results from the network.
type DetectionResult struct {
	Detections           []*Detection
	NetworkOnlyTimeTaken time.Duration
	OverallTimeTaken     time.Duration
}

var errNetworkNotInit = errors.New("network not initialised")
var errUnableToInitNetwork = errors.New("unable to initialise")

// Init the network.
func (n *YOLONetwork) Init() error {
	nCfg := C.CString(n.NetworkConfigurationFile)
	defer C.free(unsafe.Pointer(nCfg))
	wFile := C.CString(n.WeightsFile)
	defer C.free(unsafe.Pointer(wFile))

	var ns C.netSize
	// GPU device ID must be set before `load_network()` is invoked.
	// C.cuda_set_device(C.int(n.GPUDeviceIndex))
	// darknet > network.c > load_network()
	n.cNet = C.load_network(nCfg, wFile, C.int(n.GPUDeviceIndex), &ns)
	n.Width = int32(ns.width)
	n.Height = int32(ns.height)
	n.Colors = int32(ns.colors)
	n.GpuInput = (*float32)(ns.input)

	if n.cNet == nil {
		return errUnableToInitNetwork
	}

	C.set_batch_network(n.cNet, 1)
	C.srand(2222222)

	// Currently, hierarchal threshold is always 0.5.
	n.hierarchalThreshold = .5

	// Currently NMS is always 0.45.
	n.nms = .45

	n.Classes = int(C.get_network_layer_classes(n.cNet, n.cNet.n-1))
	cClassNames := loadClassNames(n.DataConfigurationFile)
	defer freeClassNames(cClassNames)
	n.ClassNames = makeClassNames(cClassNames, n.Classes)

	return nil
}

// Close and release resources.
func (n *YOLONetwork) Close() error {
	if n.cNet == nil {
		return errNetworkNotInit
	}

	C.free_network(n.cNet)
	n.cNet = nil
	return nil
}

// Detect specified image.
func (n *YOLONetwork) Detect(srcw, srch int) (*DetectionResult, error) {
	if n.cNet == nil {
		return nil, errNetworkNotInit
	}

	startTime := time.Now()
	result := C.perform_network_detect(n.cNet, C.int(srcw), C.int(srch), C.int(n.Classes),
		C.float(n.Threshold), C.float(n.hierarchalThreshold), C.float(n.nms))
	endTime := time.Now()
	defer C.free_detections(result.detections, result.detections_len)

	ds := makeDetections(srcw, srch, result.detections, int(result.detections_len),
		n.Threshold, n.Classes, n.ClassNames)

	endTimeOverall := time.Now()

	out := DetectionResult{
		Detections:           ds,
		NetworkOnlyTimeTaken: endTime.Sub(startTime),
		OverallTimeTaken:     endTimeOverall.Sub(startTime),
	}

	return &out, nil
}
