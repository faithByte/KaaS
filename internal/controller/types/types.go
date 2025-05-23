package types

import (
	"strconv"

	kaasv1 "github.com/faithByte/kaas/api/v1"

	enum "github.com/faithByte/kaas/internal/controller/utils/enums"
)

func NewDistributedMemoryStep(step *kaasv1.StepData, needs *kaasv1.NeedsData) *distributedMemory {
	return &distributedMemory{
		step:     step,
		needed:   needs.Nodes * needs.NtasksPerNode,
		started:  0,
		phase:    enum.NotStarted,
		hostfile: "",
	}
}

func NewHybridMemoryStep(step *kaasv1.StepData, needs *kaasv1.NeedsData) *hybridMemory {
	return &hybridMemory{
		strconv.Itoa(needs.CpusPerTask),
		*NewDistributedMemoryStep(step, needs),
	}
}

func NewsharedMemoryStep(step *kaasv1.StepData, needs *kaasv1.NeedsData) *sharedMemory {
	return &sharedMemory{
		step:       step,
		cpusNumber: strconv.Itoa(needs.CpusPerTask),
		phase:      enum.NotStarted,
	}
}
