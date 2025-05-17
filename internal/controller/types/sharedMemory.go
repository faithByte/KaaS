package types

import (
	"fmt"

	"k8s.io/apimachinery/pkg/api/resource"

	kaasv1 "github.com/faithByte/kaas/api/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/faithByte/kaas/internal/controller/pods"
	"github.com/faithByte/kaas/internal/controller/utils"
)

type sharedMemory struct {
	step       *kaasv1.StepData
	status     utils.Status
	cpusNumber string
}

func (data *sharedMemory) SetStatus(status utils.Status) {
	data.status = status
}

func (data *sharedMemory) Run(reconcilerData utils.ReconcilerData) error {
	if (data.status == utils.Completed) || (data.status == utils.Error) {
		return nil
	}

	name := fmt.Sprintf("%s-%s", reconcilerData.Job.Name, data.step.Name)

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
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
	data.status = utils.Launched
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

func (data *sharedMemory) GetStatus() utils.Status {
	return data.status
}

func (data *sharedMemory) GetStepData() *kaasv1.StepData {
	return data.step
}
