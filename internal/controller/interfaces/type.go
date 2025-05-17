package interfaces

import (
	kaasv1 "github.com/faithByte/kaas/api/v1"
	corev1 "k8s.io/api/core/v1"

	"github.com/faithByte/kaas/internal/controller/utils"
)

type Type interface {
	Run(reconcilerData utils.ReconcilerData) error
	SetStatus(status utils.Status)
	GetStatus() utils.Status
	GetStepData() *kaasv1.StepData
	GetResources() corev1.ResourceRequirements
	AddRunningPod(ip, resources string) bool
}
