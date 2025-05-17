package pods

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/faithByte/kaas/internal/controller/utils"
	"github.com/faithByte/kaas/internal/controller/utils/interfaces"
)

func GetLauncherPod(reconcilerData utils.ReconcilerData, data interfaces.Type) *corev1.Pod {

	// uid := string(reconcilerData.Job.GetUID())
	step := data.GetStepData()
	name := fmt.Sprintf("%s-%s-launcher", reconcilerData.Job.Name, step.Name)

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
			{
				Name: "hostfile-volume",
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: "hosts-" + string(reconcilerData.Job.GetUID()),
					},
				},
			},
		}...)

	volumes_mount := append(step.VolumeMounts,
		[]corev1.VolumeMount{
			{
				Name:      "ssh-volume",
				MountPath: "/mnt/ssh",
			},
			{
				Name:      "hostfile-volume",
				MountPath: "/hosts",
			},
		}...)

	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: reconcilerData.Job.Namespace,
			Labels: map[string]string{
				"type": "launcher",
				"step": step.Name,
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{{
				Name:  step.Name,
				Image: step.Image,
				Command: []string{"sh", "-c",
					"cp /mnt/ssh/ssh-pubkey ~/.ssh/id_rsa.pub && cp /mnt/ssh/ssh-prvkey ~/.ssh/id_rsa && chmod 600 ~/.ssh/id_rsa && chmod 644 ~/.ssh/id_rsa.pub && service ssh start && mpirun --hostfile /hosts/hostfile " + step.Command},
				Env: step.Environment,
				// ImagePullPolicy: corev1.PullAlways,
				Resources:    data.GetResources(),
				VolumeMounts: volumes_mount,
			}},
			Volumes:       volumes,
			NodeSelector:  step.NodeSelector,
			RestartPolicy: corev1.RestartPolicyNever,
		},
	}
}
