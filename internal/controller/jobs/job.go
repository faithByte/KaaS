package jobs

import (
	"fmt"

	kaasv1 "github.com/faithByte/kaas/api/v1"

	"github.com/faithByte/kaas/internal/controller/utils"
	interfaces "github.com/faithByte/kaas/internal/controller/utils/interfaces"
)

type JobData struct {
	Step    interfaces.Type
	StepSet map[string]int
	LoopSet map[string]int
}

var jobs = make(map[string]*JobData)

// INIT JOB
func New(uid string, data *utils.ReconcilerData) bool {
	job := &data.Job
	isDistributed := false

	jobs[uid] = new(JobData)
	jobs[uid].StepSet = make(map[string]int)
	jobs[uid].LoopSet = make(map[string]int)

	// index name ==============================================
	for i, step := range job.Spec.Step {
		jobs[uid].StepSet[step.Name] = i
		if (!isDistributed) && ((step.Type == "distributed_mem") || (step.Type == "hybrid_mem")) {
			isDistributed = true
		}
	}
	for i, loop := range job.Spec.Automata.Loop {
		jobs[uid].LoopSet[loop.Name] = i
	}
	// =========================================================

	if job.Spec.Automata.Run == nil {
		job.Status.Total = len(job.Spec.Step)
	} else {
		job.Status.Total = len(job.Spec.Automata.Run)
	}

	job.Status.Phase = "Starting"
	job.Status.Progress = 0
	data.Job.Status.ProgressPerTotal = fmt.Sprintf("%d/%d", data.Job.Status.Progress, data.Job.Status.Total)
	data.Client.Status().Update(data.Context, job)
	return isDistributed
}

// JOB SETTERS
func UpdateStatus(status string, data utils.ReconcilerData) {
	data.Job.Status.Phase = status
	data.Client.Status().Update(data.Context, &data.Job)
}

func IncrementProgress(data utils.ReconcilerData) {
	jobs[string(data.Job.GetUID())].Step = nil

	data.Job.Status.Progress++
	data.Job.Status.ProgressPerTotal = fmt.Sprintf("%d/%d", data.Job.Status.Progress, data.Job.Status.Total)
	data.Client.Status().Update(data.Context, &data.Job)
}

// JOB GETTERS
func Exists(uid string) bool {
	_, job := jobs[uid]
	return job
}

func IsDone(jobStatus kaasv1.JobStepsStatus) bool {
	if jobStatus.Total == 0 {
		return false
	}
	return jobStatus.Progress >= jobStatus.Total
}

// DELETE JOB
func Delete(uid string) {
	delete(jobs, uid)
}
