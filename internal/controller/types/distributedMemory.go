package types

import (
	"k8s.io/apimachinery/pkg/api/resource"
	"sigs.k8s.io/controller-runtime/pkg/client"

	kaasv1 "github.com/faithByte/kaas/api/v1"
	corev1 "k8s.io/api/core/v1"

	"github.com/faithByte/kaas/internal/controller/pods"
	"github.com/faithByte/kaas/internal/controller/utils"
	enum "github.com/faithByte/kaas/internal/controller/utils/enums"
)

type distributedMemory struct {
	step     *kaasv1.StepData
	needed   int
	started  int
	status   enum.Status
	hostfile string
}

func (data *distributedMemory) SetStatus(status enum.Status) {
	data.status = status
}

func (data *distributedMemory) Run(reconcilerData utils.ReconcilerData) error {
	uid := string(reconcilerData.Job.GetUID())
	if data.status == enum.NotStarted {
		for i := range data.needed {
			if err := pods.Create(reconcilerData, pods.GetComputePod(reconcilerData, data, i)); err != nil {
				return err
			}
		}
	} else if data.status == enum.ComputesCreated {

		var hostfile corev1.Secret
		err := reconcilerData.Client.Get(reconcilerData.Context, client.ObjectKey{Namespace: utils.MY_NAMESPACE, Name: "hosts-" + uid}, &hostfile)
		if err != nil {
			return err
		}

		if hostfile.StringData == nil {
			hostfile.StringData = map[string]string{
				"hostfile": data.hostfile,
			}
		}
		if err := reconcilerData.Client.Update(reconcilerData.Context, &hostfile); err != nil {
			return err
		}

		if err := pods.Create(reconcilerData, pods.GetLauncherPod(reconcilerData, data)); err != nil {
			return err
		}
		data.status = enum.Launched
	} else if data.status == enum.Completed {
		pods.DeleteComputes(uid, data.step.Name, reconcilerData)
	}
	return nil
}

func (data *distributedMemory) AddRunningPod(ip, resources string) bool {
	data.started++
	data.hostfile += ip + " slots=" + resources + "\n"

	if data.started == data.needed {
		data.status = enum.ComputesCreated
		return true
	}
	return false
}

// =========================== GETTERS ====================================

func (data *distributedMemory) GetResources() corev1.ResourceRequirements {
	return corev1.ResourceRequirements{
		Requests: corev1.ResourceList{
			corev1.ResourceCPU: resource.MustParse("0.5"),
		},
		Limits: corev1.ResourceList{
			corev1.ResourceCPU: resource.MustParse("0.5"),
		},
	}
}

func (data *distributedMemory) GetStatus() enum.Status {
	return data.status
}

func (data *distributedMemory) GetStepData() *kaasv1.StepData {
	return data.step
}
