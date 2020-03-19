#include <stdlib.h>
#include <darknet.h>
#include "network.h"

int get_network_layer_classes(network *n, int index) {
	return n->layers[index].classes;
}

struct network_box_result perform_network_detect(network *n, int w, int h, int classes, float thresh, float hier_thresh, float nms) {
    network_predict(n);

    struct network_box_result result = { NULL };
    result.detections = get_network_boxes(n, w, h, thresh, hier_thresh, 0, 1, &result.detections_len);
    if (nms) do_nms_sort(result.detections, result.detections_len, classes, nms);

    return result;
}
