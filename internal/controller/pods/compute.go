package pods

import (
	"fmt"

	// "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"

	kaasv1 "github.com/faithByte/kaas/api/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/faithByte/kaas/internal/controller/utils"
)

func CreateCompute(jobData *utils.JobData, step *kaasv1.StepData, index int) error {

	uid := string(jobData.Job.GetUID())

	name := fmt.Sprintf("compute-%s%d-%s-%d", uid, utils.JobSet[uid].RunIndex, jobData.Job.Name, index)

	// check if it's already created
	var isCreated corev1.Pod
	err := jobData.Client.Get(jobData.Context, client.ObjectKey{Namespace: utils.MY_NAMESPACE, Name: name}, &isCreated)

	if err == nil {
		return nil
	}

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: utils.MY_NAMESPACE,
			Labels: map[string]string{
				"type":      "compute",
				"resources": "1",
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{{
				Name:  "job-pod",
				Image: "faiithbyte/manapy:debug",
				// Command: []string{"sh", "-c", "/usr/sbin/sshd -D"},
				Command: []string{"sh", "-c", "sleep 100"},
				Env:     step.Environment,
				// ImagePullPolicy: corev1.PullAlways,
				VolumeMounts: []corev1.VolumeMount{
					// {
					// 	Name:      "data",
					// 	MountPath: "/data",
					// },
					{
						Name:      "ssh-volume",
						MountPath: "/mnt/ssh",
						ReadOnly:  true,
					},
				},
			}},
			Volumes: []corev1.Volume{
				// {
				// 	Name: "data",
				// 	VolumeSource: corev1.VolumeSource{
				// 		PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
				// 			ClaimName: "shared-volume-rw",
				// 		},
				// 	},
				// },
				{
					Name: "ssh-volume",
					VolumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName: "ssh-" + string(jobData.Job.GetUID()),
						},
					},
				},
			},
			RestartPolicy: corev1.RestartPolicyNever,
		},
	}

	// Set owner reference
	if err := ctrl.SetControllerReference(&jobData.Job, pod, jobData.Scheme); err != nil {
		return err
	}

	utils.Log.Info("Creating a new Pod", "Pod.Name", pod.Name)

	if err := jobData.Client.Create(jobData.Context, pod); err != nil {
		return err
	}

	utils.Log.Info("Pod.Name", pod.Name, "Created succesfully")

	return nil
}
