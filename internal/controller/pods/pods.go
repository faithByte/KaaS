package pods

import (
	"sigs.k8s.io/controller-runtime/pkg/client"

	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/faithByte/kaas/internal/controller/utils"
)

func GetOwnersUID(pod *corev1.Pod) string {
	ref := pod.OwnerReferences[0]

	if ref.APIVersion == utils.APIVersion && ref.Kind == utils.KIND {
		return string(ref.UID)
	}

	return ""
}

func Create(reconcilerData utils.ReconcilerData, pod *corev1.Pod) error {
	// check if it's already created
	var isCreated corev1.Pod
	err := reconcilerData.Client.Get(reconcilerData.Context, client.ObjectKey{Namespace: pod.Namespace, Name: pod.Name}, &isCreated)

	if err == nil {
		return nil
	}

	// Set owner reference
	if err := ctrl.SetControllerReference(&reconcilerData.Job, pod, reconcilerData.Scheme); err != nil {
		return err
	}

	// Create
	if err := reconcilerData.Client.Create(reconcilerData.Context, pod); err != nil {
		return err
	}

	utils.Log.Info("Creating a new launcher Pod", "Pod.Name", pod.Name)

	return nil
}
