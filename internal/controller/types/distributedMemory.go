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
	phase    enum.Phase
	hostfile string
	podName  string
}

func (data *distributedMemory) SetPhase(phase enum.Phase) {
	data.phase = phase
}

func (data *distributedMemory) Run(reconcilerData utils.ReconcilerData) error {
	uid := string(reconcilerData.Job.GetUID())
	if data.phase == enum.NotStarted {
		for i := range data.needed {
			if err := pods.Create(reconcilerData, pods.GetComputePod(reconcilerData, data, i)); err != nil {
				return err
			}
		}
	} else if data.phase == enum.ComputesCreated {

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

		launcher := pods.GetLauncherPod(reconcilerData, data)
		if err := pods.Create(reconcilerData, launcher); err != nil {
			return err
		}
		data.podName = launcher.Name
		data.phase = enum.Launched
	} else if data.phase == enum.Completed {
		pods.DeleteComputes(uid, data.step.Name, reconcilerData)
	}
	return nil
}

func (data *distributedMemory) AddRunningPod(ip, resources string) bool {
	data.started++
	data.hostfile += ip + " slots=" + resources + "\n"

	if data.started == data.needed {
		data.phase = enum.ComputesCreated
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

func (data *distributedMemory) GetPhase() enum.Phase {
	return data.phase
}

func (data *distributedMemory) GetPodName() string {
	return data.podName
}

func (data *distributedMemory) GetStepData() *kaasv1.StepData {
	return data.step
}
