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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// STEP - NEEDS
type NeedsData struct {
	Ntasks        int `json:"ntasks,omitempty"`
	CpusPerTask   int `json:"cpus-per-task,omitempty"`
	Nodes         int `json:"nodes,omitempty"`
	NtasksPerNode int `json:"ntasks-per-node,omitempty"`
}

// STEP
type StepData struct {
	Name string `json:"name"`
	// +kubebuilder:validation:Enum=shared_mem;distributed_mem;hybrid_mem
	Type         string            `json:"type"`
	Image        string            `json:"image"`
	Needs        NeedsData         `json:"needs,omitempty"`
	Command      string            `json:"command"`
	Environment  []corev1.EnvVar   `json:"environment,omitempty"`
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
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

// JOB - SPEC
type JobSpec struct {
	Step      []StepData      `json:"step"`
	Condition []ConditionData `json:"condition,omitempty"`
	Automata  AutomataData    `json:"automata,omitempty"`
}

// JOB - STATUS
type JobStatus struct {
	Phase   string `json:"phase,omitempty"`
	PodName string `json:"PodName,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="STATUS",type=string,JSONPath=".status.phase"

// JOB
type Job struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   JobSpec   `json:"spec,omitempty"`
	Status JobStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

type JobList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Job `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Job{}, &JobList{})
}
