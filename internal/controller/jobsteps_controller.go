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
	// "fmt"
	"os"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	kaasv1 "github.com/faithByte/kaas/api/v1"
	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/faithByte/kaas/internal/controller/jobs"
	"github.com/faithByte/kaas/internal/controller/secrets"
	"github.com/faithByte/kaas/internal/controller/utils"
	enum "github.com/faithByte/kaas/internal/controller/utils/enums"
	"github.com/faithByte/kaas/internal/controller/utils/mail"
	"github.com/faithByte/kaas/internal/controller/watchers"
)

type JobStepsReconciler struct {
	client.Client
	Scheme             *runtime.Scheme
	SchedulerRecorder  record.EventRecorder
	ControllerRecorder record.EventRecorder
}

// +kubebuilder:rbac:groups=faithbyte.kaas,resources=jobsteps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=faithbyte.kaas,resources=jobsteps/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=faithbyte.kaas,resources=jobsteps/finalizers,verbs=update
// func stepIsDone(uid string, stepType interfaces.Type) {

// }
func (r *JobStepsReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	utils.Log = log.FromContext(ctx)

	var job kaasv1.JobSteps

	// Get job
	if err := r.Get(ctx, req.NamespacedName, &job); err != nil {
		if errors.IsNotFound(err) {
			utils.Log.Error(err, "NO SUCH A JOB")
			return ctrl.Result{}, err
		}
	}

	if (job.Spec.Automata.Run == nil) && (job.Spec.Step == nil) {
		return ctrl.Result{}, nil
	}

	reconcilerData := utils.ReconcilerData{
		Job:     job,
		Client:  r.Client,
		Context: ctx,
		Scheme:  r.Scheme,
	}

	uid := string(job.GetUID())

	if !jobs.Exists(uid) {
		isDistributed := jobs.New(uid, &reconcilerData)
		if isDistributed {
			secrets.CreateSshSecret(&reconcilerData)
			secrets.CreateHostfile(&reconcilerData)
		}
		mail.NewMessage(uid, job.Spec)
	}

	if jobs.IsDone(job.Status) {
		jobs.Delete(uid)
		jobs.UpdateStatus("Completed", reconcilerData)

		go mail.JobEmail(job)
		return ctrl.Result{}, nil
	}

	stepType := jobs.GetStepType(uid)

	if stepType == nil {
		if job.Spec.Automata.Run == nil {
			stepType = jobs.StartStepType(uid, &job.Spec.Step[job.Status.Progress])
		} else {
			stepName, exists := job.Spec.Automata.Run[job.Status.Progress]["step"]
			if exists {
				stepType = jobs.StartStepType(uid, &job.Spec.Step[jobs.GetStepIndex(uid, stepName)])
			} // else {
			// 	loopName, loop := job.Spec.Automata.Run[job.Status.Progress]["loop"]
			// 	if exists {
			// 		RunLoop(&data, &job.Spec.Step[utils.JobSet[uid].StepSet[stepName]])
			// 	}
			// }
		}
	}

	stepType.Run(reconcilerData)

	if stepType.GetPhase() == enum.Completed {
		go mail.Messages[uid].StepMessage(job, stepType, reconcilerData)
		jobs.IncrementProgress(reconcilerData)
	} else if stepType.GetPhase() == enum.Error {
		go mail.Messages[uid].StepMessage(job, stepType, reconcilerData)

		if !stepType.GetStepData().IgnoreError {
			jobs.Delete(uid)
			jobs.UpdateStatus("Error", reconcilerData)
			go mail.JobEmail(job)
		} else {
			jobs.IncrementProgress(reconcilerData)
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *JobStepsReconciler) SetupWithManager(mgr ctrl.Manager) error {
	mail.EmailSender.NewEmail(os.Getenv("MAIL"))

	return ctrl.NewControllerManagedBy(mgr).
		For(&kaasv1.JobSteps{}, builder.WithPredicates(watchers.JobPredicate)).
		Owns(&corev1.Pod{}, builder.WithPredicates(watchers.PodPredicate)).
		Owns(&corev1.Node{}, builder.WithPredicates(watchers.NodePredicate)).
		Named("jobsteps").
		Complete(r)
}
