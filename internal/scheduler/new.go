package scheduler

import (
	corev1 "k8s.io/api/core/v1"
)

// condition types:
// NodeReady
// NodeMemoryPressure
// NodeDiskPressure
// NodePIDPressure
// NodeNetworkUnavailable

func isSchedulable(taints []corev1.Taint, conditions []corev1.NodeCondition, unschedulable bool) (int8, int8) {
	var score int8 = 0

	if (unschedulable) {
		return 0, 0
	}

	for _, condition := range conditions {
		if (condition.Type == corev1.NodeReady) && (condition.Status != corev1.ConditionTrue) {
			return 0, 0
		}

		if condition.Status == corev1.ConditionFalse {
			score++
		} else if condition.Type == corev1.NodeNetworkUnavailable {
			return 0, 0
		}
	}

	if taints == nil {
		return 1, score
	}

	var res int8 = 1
	for _, taint := range taints {
		if (taint.Effect == corev1.TaintEffectNoSchedule) || (taint.Effect == corev1.TaintEffectNoExecute) {
			return 0, score
		}
		if taint.Effect == corev1.TaintEffectPreferNoSchedule {
			res = 2
		}
	}
	return res, score
}

func NewNode(node *corev1.Node) {
	name := node.Name
	var newNode *Node

	n, score := isSchedulable(node.Spec.Taints, node.Status.Conditions, node.Spec.Unschedulable)

	switch n {
	case 0:
		return
	case 1: // primary node
		PrimaryNodeList[name] = new(Node)
		newNode = PrimaryNodeList[name]
	case 2: // secondary node
		PrimaryNodeList[name] = new(Node)
		newNode = SecondaryNodeList[name]
	}

	newNode.name = node.Name
	newNode.labels = node.Labels
	newNode.node = node
	newNode.score = score
}
