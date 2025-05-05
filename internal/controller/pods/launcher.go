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

func CreateLauncher(jobData *utils.JobData, step *kaasv1.StepData) error {

	uid := string(jobData.Job.GetUID())

	name := fmt.Sprintf("launcher-%s%d-%s", uid, utils.JobSet[uid].RunIndex, jobData.Job.Name)

	fmt.Println(name)
	// check if it's already created
	var isCreated corev1.Pod
	err := jobData.Client.Get(jobData.Context, client.ObjectKey{Namespace: jobData.Job.Namespace, Name: name}, &isCreated)

	if err == nil {
		return nil
	}

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: jobData.Job.Namespace,
			Labels: map[string]string{
				"type": "launcher",
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{{
				Name:  "job-pod",
				Image: "faiithbyte/manapy:debug",
				// Command: []string{"sh", "-c", "/usr/sbin/sshd -D"},
				Command: []string{"sh", "-c", "ls -l /hosts && cat /hosts/hostfile && echo $HOSTNAME && sleep 5"},
				Env:     step.Environment,
				// ImagePullPolicy: corev1.PullAlways,
				VolumeMounts: []corev1.VolumeMount{
					{
						Name:      "data",
						MountPath: "/data",
					},
					{
						Name:      "hostfile-volume",
						MountPath: "/hosts",
					},
					{
						Name:      "ssh-volume",
						MountPath: "/ssh",
					},
				},
			}},
			Volumes: []corev1.Volume{
				{
					Name: "data",
					VolumeSource: corev1.VolumeSource{
						PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
							ClaimName: "shared-volume-rw",
						},
					},
				},
				{
					Name: "hostfile-volume",
					VolumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName: "hosts-" + string(jobData.Job.GetUID()),
						},
					},
				},
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
