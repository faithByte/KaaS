package watchers

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

const APIVersion = "faithbyte.kaas/v1"
const KIND = "Job"

func isOwnedByMe(pod *corev1.Pod) string {
	ref := pod.OwnerReferences[0]

	if ref.APIVersion == APIVersion && ref.Kind == KIND {
		return string(ref.UID)
	}
	return ""
}

var PodPredicate = predicate.Funcs{

	CreateFunc: func(e event.CreateEvent) bool {
		var uid = ""
		if uid = isOwnedByMe(e.Object.(*corev1.Pod)); uid == "" {
			return false
		}
		fmt.Println("Create pod " + e.Object.GetName())

		return true
	},
	UpdateFunc: func(e event.UpdateEvent) bool {
		var uid = ""
		old := e.ObjectOld.(*corev1.Pod)

		if uid = isOwnedByMe(old); uid == "" {
			return false
		}
		fmt.Println("Update pod " + e.ObjectNew.GetName())

		new := e.ObjectNew.(*corev1.Pod)

		if old.Status.Phase != new.Status.Phase {
			fmt.Println(new.Status.Phase)
		}
		return false
	},
	DeleteFunc: func(e event.DeleteEvent) bool {
		var uid = ""
		if uid = isOwnedByMe(e.Object.(*corev1.Pod)); uid == "" {
			return false
		}

		fmt.Println("Delete pod " + e.Object.GetName())

		return false
	},
}
