// Definitions for the Kubernetes types
package v1alpha1

import (
	. "github.com/solo-io/mesh-projects/pkg/api/security.zephyr.solo.io/v1alpha1/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status

// VirtualMeshCertificateSigningRequest is the Schema for the virtualMeshCertificateSigningRequest API
type VirtualMeshCertificateSigningRequest struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VirtualMeshCertificateSigningRequestSpec   `json:"spec,omitempty"`
	Status VirtualMeshCertificateSigningRequestStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// VirtualMeshCertificateSigningRequestList contains a list of VirtualMeshCertificateSigningRequest
type VirtualMeshCertificateSigningRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VirtualMeshCertificateSigningRequest `json:"items"`
}

func init() {
	SchemeBuilder.Register(&VirtualMeshCertificateSigningRequest{}, &VirtualMeshCertificateSigningRequestList{})
}
