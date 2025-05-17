package types

import (
	"k8s.io/apimachinery/pkg/api/resource"

	corev1 "k8s.io/api/core/v1"
)

type hybridMemory struct {
	cpusNumber string
	distributedMemory
}

// =========================== GETTERS ====================================

func (data *hybridMemory) GetResources() corev1.ResourceRequirements {
	return corev1.ResourceRequirements{
		Requests: corev1.ResourceList{
			corev1.ResourceCPU: resource.MustParse(data.cpusNumber),
		},
		Limits: corev1.ResourceList{
			corev1.ResourceCPU: resource.MustParse(data.cpusNumber),
		},
	}
}
