package controller

import (
	kaasv1 "github.com/faithByte/kaas/api/v1"

	"github.com/faithByte/kaas/internal/controller/utils"
)

func RunStep(jobData *utils.JobData, step *kaasv1.StepData) {
	podsNeeded := 1
	nodesNeeded := step.Needs.Nodes
	utils.JobStartStep(string(jobData.Job.GetUID()), step.Name, podsNeeded)
	if nodesNeeded == 0 {

	}
}
