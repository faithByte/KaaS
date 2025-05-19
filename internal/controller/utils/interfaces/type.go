package interfaces

import (
	kaasv1 "github.com/faithByte/kaas/api/v1"
	corev1 "k8s.io/api/core/v1"

	"github.com/faithByte/kaas/internal/controller/utils"
	enum "github.com/faithByte/kaas/internal/controller/utils/enums"
)

type Type interface {
	Run(reconcilerData utils.ReconcilerData) error
	SetPhase(phase enum.Phase)
	GetPhase() enum.Phase
	GetPodName() string
	GetStepData() *kaasv1.StepData
	GetResources() corev1.ResourceRequirements
	AddRunningPod(ip, resources string) bool
}
