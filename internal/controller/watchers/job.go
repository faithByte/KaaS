package watchers

import (
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	kaasv1 "github.com/faithByte/kaas/api/v1"

	"github.com/faithByte/kaas/internal/controller/utils"
)

var JobPredicate = predicate.Funcs{

	CreateFunc: func(e event.CreateEvent) bool {
		fmt.Println("Create job " + e.Object.GetName())
		
		job := e.Object.(*kaasv1.Job)
		uid := string(job.GetUID())
		// index name ==============================================
		for i, step := range job.Spec.Step {
			utils.NewStep(uid, step.Name, i)
		}
		for i, loop := range job.Spec.Automata.Loop {
			utils.NewLoop(uid, loop.Name, i)
		}
		// =========================================================

		if job.Spec.Automata.Run == nil {
			utils.JobSet[uid].RunLen = len(job.Spec.Step)
		} else {
			utils.JobSet[uid].RunLen = len(job.Spec.Automata.Run)
		}

		return true
	},

	UpdateFunc: func(e event.UpdateEvent) bool {
		fmt.Println("Update job " + e.ObjectNew.GetName())
		return true
	},

	DeleteFunc: func(e event.DeleteEvent) bool {
		fmt.Println("Delete job " + e.Object.GetName())
		utils.RemoveJob(string(e.Object.GetUID()))
		return false
	},

	GenericFunc: func(e event.GenericEvent) bool {
		return false
	},
}
