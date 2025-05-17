package pods

import (
	"fmt"

	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/client"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/faithByte/kaas/internal/controller/utils"
	"github.com/faithByte/kaas/internal/controller/utils/interfaces"
)

func DeleteComputes(uid, stepName string, data utils.ReconcilerData) {
	labelSelector := client.MatchingLabels{
		"job":  uid,
		"type": "compute",
		"step": stepName,
	}
	namespaceSelector := client.InNamespace(utils.MY_NAMESPACE)
	data.Client.DeleteAllOf(data.Context, &corev1.Pod{}, labelSelector, namespaceSelector)
}

func GetComputePod(reconcilerData utils.ReconcilerData, data interfaces.Type, index int) *corev1.Pod {

	step := data.GetStepData()

	volumes := append(reconcilerData.Job.Spec.Volumes,
		[]corev1.Volume{
			{
				Name: "ssh-volume",
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: "ssh-" + string(reconcilerData.Job.GetUID()),
					},
				},
			},
		}...)

	volumes_mount := append(step.VolumeMounts,
		[]corev1.VolumeMount{
			{
				Name:      "ssh-volume",
				MountPath: "/mnt/ssh",
				ReadOnly:  true,
			},
		}...)

	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-%s-compute-%d", reconcilerData.Job.Name, step.Name, index),
			Namespace: utils.MY_NAMESPACE,
			Labels: map[string]string{
				"job":  string(reconcilerData.Job.GetUID()),
				"type": "compute",
				"step": step.Name,
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{{
				Name:  step.Name,
				Image: step.Image,
				Command: []string{"sh", "-c",
					"cp -f /mnt/ssh/ssh-pubkey ~/.ssh/authorized_keys && chmod 644 ~/.ssh/authorized_keys && /usr/sbin/sshd -D"},
				Env: step.Environment,
				// ImagePullPolicy: corev1.PullAlways,
				Resources:    data.GetResources(),
				VolumeMounts: volumes_mount,
			}},
			Volumes:                       volumes,
			NodeSelector:                  step.NodeSelector,
			RestartPolicy:                 corev1.RestartPolicyNever,
			TerminationGracePeriodSeconds: pointer.Int64(0),
		},
	}
}
