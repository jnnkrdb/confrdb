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

// Package v1beta2 contains API Schema definitions for the globals v1beta2 API group
// +kubebuilder:object:generate=true
// +groupName=globals.jnnkrdb.de
package v1beta2

import (
	"regexp"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

var (
	// GroupVersion is group version used to register these objects
	GroupVersion = schema.GroupVersion{Group: "globals.jnnkrdb.de", Version: "v1beta2"}

	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	SchemeBuilder = &scheme.Builder{GroupVersion: GroupVersion}

	// AddToScheme adds the types in this group-version to the given scheme.
	AddToScheme = SchemeBuilder.AddToScheme
)

// set the finalizer objects
const FinalizerGlobal string = "v1beta2.globals.jnnkrdb.de/finalizer"

// struct which contains the information about the namespace regex
type NamespacesRegex struct {

	// +kubebuilder:default={default}
	// +kubebuilder:validation:UniqueItems=true
	// +kubebuilder:validation:MinItems=1
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	AvoidRegex []string `json:"avoidregex"`

	// +kubebuilder:default={default}
	// +kubebuilder:validation:UniqueItems=true
	// +kubebuilder:validation:MinItems=1
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	MatchRegex []string `json:"matchregex"`
}

// check whether a string exists in a list of regexpressions or not
func StringMatchesRegExpList(comp string, regexpList []string) (bool, error) {
	for i := range regexpList {
		if matched, err := regexp.MatchString(regexpList[i], comp); err != nil {
			return false, nil
		} else {
			if matched {
				return true, nil
			}
		}
	}
	return false, nil
}

// find a string in a list of string
func StringInList(comp string, list []string) bool {
	for i := range list {
		if list[i] == comp {
			return true
		}
	}
	return false
}
