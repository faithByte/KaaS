package watchers

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	"github.com/faithByte/kaas/internal/scheduler"
)

var NodePredicate = predicate.Funcs{

	CreateFunc: func(e event.CreateEvent) bool {
		fmt.Println("Create Node \n" + e.Object.GetName())
		scheduler.NewNode(e.Object.(*corev1.Node))
		return false
	},

	UpdateFunc: func(e event.UpdateEvent) bool {
		fmt.Println("Update Node \n" + e.ObjectNew.GetName())
		scheduler.UpdateNode(e.ObjectOld.(*corev1.Node), e.ObjectNew.(*corev1.Node))
		return false
	},

	DeleteFunc: func(e event.DeleteEvent) bool {
		fmt.Println("Delete Node " + e.Object.GetName())
		scheduler.DeleteNode(string(e.Object.(*corev1.Node).GetUID()))
		return false
	},

	GenericFunc: func(e event.GenericEvent) bool {
		return false
	},
}
