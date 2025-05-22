package watchers

import (
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	corev1 "k8s.io/api/core/v1"

	"github.com/faithByte/kaas/internal/controller/jobs"
	enum "github.com/faithByte/kaas/internal/controller/utils/enums"
)

const APIVersion = "faithbyte.kaas/v1"
const KIND = "JobSteps"

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

		return false
	},
	UpdateFunc: func(e event.UpdateEvent) bool {
		var uid = ""
		old := e.ObjectOld.(*corev1.Pod)

		if uid = isOwnedByMe(old); uid == "" {
			return false
		}
		fmt.Println("Update pod " + e.ObjectNew.GetName())

		new := e.ObjectNew.(*corev1.Pod)
		// new.
		if old.Status.Phase != new.Status.Phase {
			if (new.Status.Phase == "Running") && (new.Labels["type"] == "compute") {
				return jobs.GetStepType(uid).AddRunningPod(new.Status.PodIP, "1")
			} else if new.Status.Phase == "Succeeded" {
				if (new.Labels["type"] == "launcher") || (new.Labels["type"] == "main") {
					jobs.UpdateStepPhase(uid, enum.Completed)
				}
				return true
			}
		}
		return false
	},
	DeleteFunc: func(e event.DeleteEvent) bool {
		fmt.Println("Delete pod " + e.Object.GetName())

		return false
	},

	GenericFunc: func(e event.GenericEvent) bool {
		return false
	},
}
