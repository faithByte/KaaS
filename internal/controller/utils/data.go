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

type Data struct {
	Index    int
	Hostfile string
}

type JobSetData struct {
	Step     string
	Needed   int
	Started  int
	launched bool
	Hostfile string
}

var jobSet = make(map[string]*JobSetData)

func AddJob(uid string) {
	jobSet[uid] = new(JobSetData)
}

func RemoveJob(uid string) {
	delete(jobSet, uid)
	println(jobSet)
}

func JobExists(uid string) bool {
	_, job := jobSet[uid]
	return job
}

func JobStartStep(uid string, stepName string, podsNeeded int) {
	jobSet[uid].Step = stepName
	jobSet[uid].Needed = podsNeeded
	jobSet[uid].Started = 0
	jobSet[uid].launched = false
	jobSet[uid].Hostfile = ""
}

func AddRunningPod(uid string, ip string, resources string) {
	jobSet[uid].Started += 1
	jobSet[uid].Hostfile += "\n" + ip + " slots=" + resources
}
