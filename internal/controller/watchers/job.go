package watchers

import (
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	"github.com/faithByte/kaas/internal/controller/utils"
)

var JobPredicate = predicate.Funcs{

	CreateFunc: func(e event.CreateEvent) bool {
		fmt.Println("Create job " + e.Object.GetName())
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
