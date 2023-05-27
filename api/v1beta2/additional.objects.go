package v1beta2

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
