package secrets

import (
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/faithByte/kaas/internal/controller/utils"
)

func CreateHostfile(reconcilerData *utils.ReconcilerData) error {
	name := "hosts-" + string(reconcilerData.Job.GetUID())

	var isCreated corev1.Secret
	err := reconcilerData.Client.Get(reconcilerData.Context, client.ObjectKey{Namespace: utils.MY_NAMESPACE, Name: name}, &isCreated)
	if err == nil {
		return errors.NewAlreadyExists(schema.GroupResource{}, "")
	} else if !errors.IsNotFound(err) {
		return err
	}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: utils.MY_NAMESPACE,
		},
	}

	if err := ctrl.SetControllerReference(&reconcilerData.Job, secret, reconcilerData.Scheme); err != nil {
		return err
	}

	if err := reconcilerData.Client.Create(reconcilerData.Context, secret); err != nil {
		return err
	}

	return nil
}
