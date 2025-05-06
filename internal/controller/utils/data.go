package utils

import (
	"context"
	// "fmt"

	// "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	// "k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	kaasv1 "github.com/faithByte/kaas/api/v1"
	// "github.com/faithByte/kaas/internal/controller"
)

const MY_NAMESPACE = "default"

var Log = log.Log

type JobData struct {
	Job     kaasv1.Job
	Client  client.Client
	Context context.Context
	Scheme  *runtime.Scheme
}

type Status int8

const (
	NotStarted       Status = 0
	ComputesCreating Status = 1
	ComputesCreated  Status = 2
	Launched         Status = 3
	Completed        Status = 4
)

type JobSetData struct {
	RunIndex int
	RunLen   int

	Step     string
	Needed   int
	Started  int
	Status   Status
	Hostfile string

	StepSet map[string]int
	LoopSet map[string]int
}

var JobSet = make(map[string]*JobSetData)

func AddJob(uid string) {
	JobSet[uid] = new(JobSetData)
	JobSet[uid].StepSet = make(map[string]int)
	JobSet[uid].LoopSet = make(map[string]int)
}

func JobExists(uid string) bool {
	_, job := JobSet[uid]
	return job
}

func JobStartStep(uid string, stepName string, podsNeeded int) {
	JobSet[uid].Step = stepName
	JobSet[uid].Needed = podsNeeded
	JobSet[uid].Started = 0
	JobSet[uid].Status = ComputesCreating
	JobSet[uid].Hostfile = ""
}

func AddRunningPod(uid string, ip string, resources string) {
	JobSet[uid].Started += 1
	JobSet[uid].Hostfile += ip + " slots=" + resources + "\n"

	if JobSet[uid].Started == JobSet[uid].Needed {
		JobSet[uid].Status = ComputesCreated
	}
}

func RemoveJob(uid string) {
	delete(JobSet, uid)
}

// func UpdateJobStatus(){

// }
