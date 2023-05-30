package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// DummySpec defines the desired state of Dummy
type DummySpec struct {
	Message string `json:"message,omitempty"`
}

// DummyStatus defines the observed state of Dummy
type DummyStatus struct {
	SpecEcho  string `json:"specEcho,omitempty"`
	PodStatus string `json:"podStatus,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Dummy is the Schema for the dummies API
type Dummy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DummySpec   `json:"spec,omitempty"`
	Status DummyStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// DummyList contains a list of Dummy
type DummyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Dummy `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Dummy{}, &DummyList{})
}
