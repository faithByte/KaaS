/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	ctrl "sigs.k8s.io/controller-runtime"
	kaasv1 "github.com/faithByte/KaaS/api/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

type JobReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=kaas.faithbyte.kaas,resources=jobs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=kaas.faithbyte.kaas,resources=jobs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=kaas.faithbyte.kaas,resources=jobs/finalizers,verbs=update

func (r *JobReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = logf.FromContext(ctx)

	// get job data
	var job kaasv1.Job
	if err := r.Get(ctx, req.NamespacedName, &job); err != nil {
		return ctrl.Result{}, err
	}

	// ssh secret
	if err := r.createSshSecret(ctx, job); err != nil {
		return ctrl.Result{}, err
	}

	// hostfile secret

	return ctrl.Result{}, nil
}

func (r *JobReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&kaasv1.Job{}).
		Named("job").
		Complete(r)
}
