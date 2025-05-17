package utils

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	kaasv1 "github.com/faithByte/kaas/api/v1"
	corev1 "k8s.io/api/core/v1"
	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const MY_NAMESPACE = "default"
const APIVersion = "faithbyte.kaas/v1"
const KIND = "JobSteps"

var Log = log.Log

type ReconcilerData struct {
	Job     kaasv1.JobSteps
	Client  client.Client
	Context context.Context
	Scheme  *runtime.Scheme
}

// func isOwnedByMe(obj *metav1.Object, uid string) bool {
// 	controller := metav1.GetControllerOf(*obj)

// 	if string(controller.UID) == uid {
// 		return true
// 	}
// 	return false
// }

func GetPodsOwnerUID(pod *corev1.Pod) string {
	ref := pod.OwnerReferences[0]

	if ref.APIVersion == APIVersion && ref.Kind == KIND {
		return string(ref.UID)
	}
	return ""
}
