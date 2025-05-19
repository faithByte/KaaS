package pods

import (
	"fmt"
	"io"
	"strings"

	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"

	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/faithByte/kaas/internal/controller/utils"
)

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

func GetLogs(name, namespace string, r *utils.ReconcilerData) string {
	config := ctrl.GetConfigOrDie()

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return ""
	}

	req := clientset.CoreV1().Pods(namespace).GetLogs(name, &corev1.PodLogOptions{})
	stream, err := req.Stream(r.Context)
	if err != nil {
		fmt.Println(fmt.Errorf("%w", err))

		return ""
	}
	defer stream.Close()

	var builder strings.Builder
	_, err = io.Copy(&builder, stream)
	if err != nil {
		fmt.Println(fmt.Errorf("%w", err))
		return ""
	}

	return builder.String()
}
