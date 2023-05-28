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

// GlobalSecretSpec defines the desired state of GlobalSecret
type GlobalSecretSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// +kubebuilder:default=false
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +kubebuilder:printcolumn:JSONPath=".spec.immutable",name="Immutable",type="boolean"
	Immutable bool `json:"immutable"`

	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Namespaces NamespacesRegex `json:"namespaces"`

	// +kubebuilder:validation:Enum={"Opaque","kubernetes.io/service-account-token","kubernetes.io/dockercfg","kubernetes.io/dockerconfigjson","kubernetes.io/basic-auth","kubernetes.io/ssh-auth","kubernetes.io/tls","bootstrap.kubernetes.io/token"}
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +kubebuilder:printcolumn:JSONPath=".spec.type",name="Type",type="string"
	Type string `json:"type"`

	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Data map[string]string `json:"data"`
}

// GlobalSecretStatus defines the observed state of GlobalSecret
type GlobalSecretStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// +operator-sdk:csv:customresourcedefinitions:type=status
	DeployedSecrets []DeployedSecret `json:"deployedsecrets,omitempty"`

	// +operator-sdk:csv:customresourcedefinitions:type=status
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,1,rep,name=conditions"`
}

type DeployedSecret struct{}

// GlobalSecret is the Schema for the globalsecrets API
// +kubebuilder:subresource:status
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=services,shortName=gs;gss
type GlobalSecret struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GlobalSecretSpec   `json:"spec,omitempty"`
	Status GlobalSecretStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// GlobalSecretList contains a list of GlobalSecret
type GlobalSecretList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GlobalSecret `json:"items"`
}

func init() {
	SchemeBuilder.Register(&GlobalSecret{}, &GlobalSecretList{})
}
