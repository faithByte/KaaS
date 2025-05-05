package controller

import (
	kaasv1 "github.com/faithByte/kaas/api/v1"
	corev1 "k8s.io/api/core/v1"

	"github.com/faithByte/kaas/internal/controller/pods"
	// "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	// "github.com/faithByte/kaas/internal/controller/secrets"
	"github.com/faithByte/kaas/internal/controller/utils"
)

func RunStep(jobData *utils.JobData, step *kaasv1.StepData) error {
	podsNeeded := 3

	uid := string(jobData.Job.GetUID())
	status := utils.JobSet[uid].Status
	if (status == utils.NotStarted) || (status == utils.Completed) {
		utils.JobStartStep(uid, step.Name, podsNeeded)

		for i := range podsNeeded {
			if err := pods.CreateCompute(jobData, step, i); err != nil {
				return err
			}
		}
	} else if status == utils.ComputesCreated {
		utils.JobSet[uid].Status = utils.Launched

		var hostfile corev1.Secret
		err := jobData.Client.Get(jobData.Context, client.ObjectKey{Namespace: utils.MY_NAMESPACE, Name: "hosts-" + uid}, &hostfile)
		if err != nil {
			return err
		}

		if hostfile.StringData == nil {
			hostfile.StringData = map[string]string{
				"hostfile": utils.JobSet[uid].Hostfile,
			}
		}
		if err := jobData.Client.Update(jobData.Context, &hostfile); err != nil {
			return err
		}

		if err := pods.CreateLauncher(jobData, step); err != nil {
			return err
		}
	}
	return nil
}
