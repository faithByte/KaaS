package scheduler

import (
	"reflect"

	corev1 "k8s.io/api/core/v1"
)

func UpdateNode(old, update *corev1.Node) {
	name := update.GetName()

	if !(reflect.DeepEqual(old.Spec.Taints, update.Spec.Taints)) || !(reflect.DeepEqual(old.Status.Conditions, update.Status.Conditions)) {
		var newNode *Node
		var score, n int8

		if !NodeExist(name) {
			NewNode(update)
			return
		}

		// lock scheduler
		n, score = isSchedulable(update.Spec.Taints, update.Status.Conditions, update.Spec.Unschedulable)

		switch n {
		case 0:
			DeleteNode(name)
			return
		case 1: // primary node
			if !NodeInPrimaryList(name) {
				PrimaryNodeList[name] = SecondaryNodeList[name]
				newNode = PrimaryNodeList[name]
				delete(SecondaryNodeList, name)
			}
		case 2: // secondary node
			if !NodeInSecondaryList(name) {
				SecondaryNodeList[name] = PrimaryNodeList[name]
				newNode = SecondaryNodeList[name]
				delete(PrimaryNodeList, name)
			}
		}
		newNode.labels = update.Labels
		newNode.node = update
		newNode.score = score
		// unlock scheduler

	}
}
