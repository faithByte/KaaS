package types

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	kaasv1 "github.com/faithByte/kaas/api/v1"
	"github.com/faithByte/kaas/internal/controller/pods"
	"github.com/faithByte/kaas/internal/controller/utils"
	enum "github.com/faithByte/kaas/internal/controller/utils/enums"
)

type sharedMemory struct {
	step       *kaasv1.StepData
	phase      enum.Phase
	cpusNumber string
	podName    string
}

func (data *sharedMemory) SetPhase(phase enum.Phase) {
	data.phase = phase
}

func (data *sharedMemory) Run(reconcilerData utils.ReconcilerData) error {
	if (data.phase == enum.Completed) || (data.phase == enum.Error) {
		return nil
	}

	var isCreated corev1.Pod
	data.podName = fmt.Sprintf("%s-%s", reconcilerData.Job.Name, data.step.Name)
	err := reconcilerData.Client.Get(reconcilerData.Context, client.ObjectKey{Namespace: reconcilerData.Job.Namespace, Name: data.podName}, &isCreated)

	i := 1
	for err == nil {
		data.podName = fmt.Sprintf("%s-%s-%d", reconcilerData.Job.Name, data.step.Name, i)
		err = reconcilerData.Client.Get(reconcilerData.Context, client.ObjectKey{Namespace: reconcilerData.Job.Namespace, Name: data.podName}, &isCreated)
		i++
	}

	if !errors.IsNotFound(err) {
		return err
	}

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      data.podName,
			Namespace: reconcilerData.Job.Namespace,
			Labels: map[string]string{
				"type": "main",
				"step": data.step.Name,
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{{
				Name:    data.step.Name,
				Image:   data.step.Image,
				Command: []string{"sh", "-c", data.step.Command},
				Env:     data.step.Environment,
				// ImagePullPolicy: corev1.PullAlways,
				Resources:    data.GetResources(),
				VolumeMounts: data.step.VolumeMounts,
			}},
			Volumes:       reconcilerData.Job.Spec.Volumes,
			NodeSelector:  data.step.NodeSelector,
			RestartPolicy: corev1.RestartPolicyNever,
		},
	}

	pods.Create(reconcilerData, pod)
	return nil
}

func (data *sharedMemory) AddRunningPod(ip, resources string) bool {
	data.phase = enum.Launched
	return true
}

// =========================== GETTERS ====================================

func (data *sharedMemory) GetResources() corev1.ResourceRequirements {
	return corev1.ResourceRequirements{
		Requests: corev1.ResourceList{
			corev1.ResourceCPU: resource.MustParse(data.cpusNumber),
		},
		Limits: corev1.ResourceList{
			corev1.ResourceCPU: resource.MustParse(data.cpusNumber),
		},
	}
}

func (data *sharedMemory) GetPhase() enum.Phase {
	return data.phase
}

func (data *sharedMemory) GetPodName() string {
	return data.podName
}

func (data *sharedMemory) GetStepData() *kaasv1.StepData {
	return data.step
}
