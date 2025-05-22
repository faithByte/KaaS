package scheduler

import (
	corev1 "k8s.io/api/core/v1"
)

type Node struct {
	name   string
	labels map[string]string
	score int8
	node  *corev1.Node
}

var PrimaryNodeList = make(map[string]*Node)
var SecondaryNodeList = make(map[string]*Node)
