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
	"fmt"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	kaasv1 "github.com/faithByte/kaas/api/v1"
	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/faithByte/kaas/internal/controller/secrets"
	"github.com/faithByte/kaas/internal/controller/utils"
	"github.com/faithByte/kaas/internal/controller/watchers"
)

// +kubebuilder:rbac:groups=faithbyte.kaas,resources=jobs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=faithbyte.kaas,resources=jobs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=faithbyte.kaas,resources=jobs/finalizers,verbs=update

type JobReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *JobReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	utils.Log = log.FromContext(ctx)
	var job kaasv1.Job

	fmt.Println("==========RECONCILE=========================")

	// Get job
	if err := r.Get(ctx, req.NamespacedName, &job); err != nil {
		if errors.IsNotFound(err) {
			utils.Log.Error(err, "NO SUCH A JOB")
			return ctrl.Result{}, err
		}
	}

	data := utils.JobData{
		Job:     job,
		Client:  r.Client,
		Context: ctx,
		Scheme:  r.Scheme,
	}

	job.Status.Phase = "Creating dependencies"
	r.Status().Update(ctx, &job)

	if !utils.JobExists(string(job.GetUID())) {
		fmt.Println("Creaaaaaaaaaaaaaaaaaaate")
		// Create ssh secrete
		if err := secrets.CreateSshSecret(&data); err != nil {
			utils.Log.Error(err, "COULDN'T CREATE SSH SECRET")
			return ctrl.Result{}, err
		}
		utils.AddJob(string(job.GetUID()))

		if len(job.Spec.Automata.Run) > 0 {
			// index name ==============================================
			StepSet := make(map[string]utils.Data)
			for i, step := range job.Spec.Step {
				StepSet[step.Name] = utils.Data{Index: i}
			}
			LoopSet := make(map[string]int)
			for i, loop := range job.Spec.Automata.Loop {
				LoopSet[loop.Name] = i
			}
			// =========================================================

			for _, j := range job.Spec.Automata.Run {
				stepName, step := j["step"]
				if step {
					RunStep(&data, &job.Spec.Step[StepSet[stepName].Index])
				} else {
					_, loop := j["loop"]
					// loopName, loop := j["loop"]
					if !loop {
						continue
					}
					// RunLoop(&data, &job.Spec.Step[LoopSet[loopName]])
				}
			}
		} else {
			for _, step := range job.Spec.Step {
				RunStep(&data, &step)
			}
		}

	} else {
		fmt.Println("HEEEEEEEEEEEEEEEEEERRRRE")

	}
	utils.Log.Info("Done")
	return ctrl.Result{}, nil
}

func (r *JobReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&kaasv1.Job{}, builder.WithPredicates(watchers.JobPredicate)).
		Owns(&corev1.Pod{}, builder.WithPredicates(watchers.PodPredicate)).
		Complete(r)
}
