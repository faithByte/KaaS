package utils

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	kaasv1 "github.com/faithByte/kaas/api/v1"
	// corev1 "k8s.io/api/core/v1"
	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const MY_NAMESPACE = "default"
const APIVersion = "faithbyte.kaas/v1"
const KIND = "JobSteps"

var Log = log.Log

type ReconcilerData struct {
	Job     kaasv1.JobSteps
	Client  client.Client
	Context context.Context
	Scheme  *runtime.Scheme
}

func KeyExistInMap(key string, m map[string]interface{}) bool {
	_, exist := m[key]
	return exist
}

// func ValueExistInMap() {

// }
