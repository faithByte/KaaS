package watchers

import (
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	kaasv1 "github.com/faithByte/kaas/api/v1"

	"github.com/faithByte/kaas/internal/controller/jobs"
)

var JobPredicate = predicate.Funcs{

	CreateFunc: func(e event.CreateEvent) bool {
		fmt.Println("Create job " + e.Object.GetName())
		return true
	},

	UpdateFunc: func(e event.UpdateEvent) bool {
		new := e.ObjectOld.(*kaasv1.JobSteps)
		old := e.ObjectNew.(*kaasv1.JobSteps)

		if (new.Status.Progress != 0) && (new.Status.Progress != old.Status.Progress) {
			return true
		}

		fmt.Println("Update job " + e.ObjectNew.GetName())
		return true
	},

	DeleteFunc: func(e event.DeleteEvent) bool {
		fmt.Println("Delete job " + e.Object.GetName())
		jobs.Delete(string(e.Object.GetUID()))
		return false
	},

	GenericFunc: func(e event.GenericEvent) bool {
		return false
	},
}
