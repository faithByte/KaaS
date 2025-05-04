package utils

import (
	"context"
	// "fmt"

	// "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	// "k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	kaasv1 "github.com/faithByte/kaas/api/v1"
	// "github.com/faithByte/kaas/internal/controller"
)

const MY_NAMESPACE = "default"

var Log = log.Log

type JobData struct {
	Job     kaasv1.Job
	Client  client.Client
	Context context.Context
	Scheme  *runtime.Scheme
}
