/*
Copyright 2023.

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

package v1beta2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// GlobalConfigSpec defines the desired state of GlobalConfig
type GlobalConfigSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Namespaces NamespacesRegex `json:"namespaces"`

	//
	//
	//
	//

	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Data map[string]string `json:"data"`
}

// GlobalConfigStatus defines the observed state of GlobalConfig
type GlobalConfigStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// +operator-sdk:csv:customresourcedefinitions:type=status
	DeployedConfigMaps []DeployedConfigMap `json:"deployedconfigmaps,omitempty"`

	// +operator-sdk:csv:customresourcedefinitions:type=status
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,1,rep,name=conditions"`
}

type DeployedConfigMap struct{}

// GlobalConfig is the Schema for the globalconfigs API
// +kubebuilder:subresource:status
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=services,shortName=gc;gcs
type GlobalConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GlobalConfigSpec   `json:"spec,omitempty"`
	Status GlobalConfigStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// GlobalConfigList contains a list of GlobalConfig
type GlobalConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GlobalConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&GlobalConfig{}, &GlobalConfigList{})
}
