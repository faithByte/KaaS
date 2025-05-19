package jobs

import (
	kaasv1 "github.com/faithByte/kaas/api/v1"

	"github.com/faithByte/kaas/internal/controller/types"
	enum "github.com/faithByte/kaas/internal/controller/utils/enums"
	"github.com/faithByte/kaas/internal/controller/utils/interfaces"
)

func StartStepType(uid string, step *kaasv1.StepData) interfaces.Type {
	switch step.Type {
	case "shared_mem":
		jobs[uid].Step = types.NewsharedMemoryStep(step, &step.Needs)
	case "distributed_mem":
		jobs[uid].Step = types.NewDistributedMemoryStep(step, &step.Needs)
	case "hybrid_mem":
		jobs[uid].Step = types.NewHybridMemoryStep(step, &step.Needs)
	}
	return jobs[uid].Step
}

func GetStepType(uid string) interfaces.Type {
	return jobs[uid].Step
}

func UpdateStepPhase(uid string, phase enum.Phase) {
	if Exists(uid) {
		jobs[uid].Step.SetPhase(phase)
	}
}

func GetStepIndex(uid string, name string) int {
	return jobs[uid].StepSet[name]
}
