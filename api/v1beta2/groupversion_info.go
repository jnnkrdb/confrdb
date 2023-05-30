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
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
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
const FinalizerGlobal string = "globals.jnnkrdb.de/v1beta2.finalizer"

// get/set the labels, whehter to compare or to set
func MatchingLables(uid types.UID) client.MatchingLabels {
	return client.MatchingLabels{
		"globals.jnnkrdb.de/confrdb.version": GroupVersion.Version,
		"globals.jnnkrdb.de/confrdb.uid":     string(uid),
	}
}
