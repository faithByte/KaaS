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

package v1

import (
	enum "github.com/faithByte/kaas/internal/controller/utils/enums"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// STEP - NEEDS
type NeedsData struct {
	// +kubebuilder:default=1
	CpusPerTask int `json:"cpus-per-task,omitempty"`
	// +kubebuilder:default=1
	Nodes int `json:"nodes,omitempty"`
	// +kubebuilder:default=1
	NtasksPerNode int `json:"ntasks-per-node,omitempty"`
}

// STEP
type StepData struct {
	// +kubebuilder:validation:Pattern=`^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$`
	Name string `json:"name"`
	// +kubebuilder:validation:Enum=shared_mem;distributed_mem;hybrid_mem
	Type         string               `json:"type"`
	Image        string               `json:"image"`
	Needs        NeedsData            `json:"needs,omitempty"`
	Command      string               `json:"command"`
	Environment  []corev1.EnvVar      `json:"environment,omitempty"`
	NodeSelector map[string]string    `json:"nodeSelector,omitempty"`
	VolumeMounts []corev1.VolumeMount `json:"volumeMounts,omitempty"`
	// +kubebuilder:default=false
	IgnoreError bool `json:"ignore_errors,omitempty"`
}

// AUTOMATA - LOOP
type LoopData struct {
	Name      string `json:"name"`
	Condition string `json:"condition"`
	Step      string `json:"step"`
}

// AUTOMATA
type AutomataData struct {
	Run  []map[string]string `json:"run,omitempty"`
	Loop []LoopData          `json:"loop,omitempty"`
}

// CONDITION
type ConditionData struct {
	Name    string `json:"name"`
	Image   string `json:"image"`
	Command string `json:"command"`
}

// Email
type EmailData struct {
	Email string `json:"email"`
	// +kubebuilder:validation:Enum=all;success;error
	// +kubebuilder:default=all
	Status string `json:"type,omitempty"`
	// +kubebuilder:validation:Enum=all;job;step
	// +kubebuilder:default=all
	For enum.EmailFor `json:"for,omitempty"`
}

// JOBSTEPS - SPEC
type JobStepsSpec struct {
	Step      []StepData      `json:"step"`
	Condition []ConditionData `json:"condition,omitempty"`
	Automata  AutomataData    `json:"automata,omitempty"`
	Email     EmailData       `json:"email,omitempty"`
	Volumes   []corev1.Volume `json:"volumes,omitempty"`
}

// JOBSTEPS - STATUS
type JobStepsStatus struct {
	Phase   string `json:"phase,omitempty"`
	PodName string `json:"PodName,omitempty"`

	// +kubebuilder:default=0
	Total int `json:"total,omitempty"`
	// +kubebuilder:default=0
	Progress int `json:"progress,omitempty"`

	ProgressPerTotal string `json:"progresspertotal,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="STATUS",type=string,JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Steps Progress",type=string,JSONPath=".status.progresspertotal"

// JOBSTEPS
type JobSteps struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   JobStepsSpec   `json:"spec,omitempty"`
	Status JobStepsStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

type JobStepsList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []JobSteps `json:"items"`
}

func init() {
	SchemeBuilder.Register(&JobSteps{}, &JobStepsList{})
}
