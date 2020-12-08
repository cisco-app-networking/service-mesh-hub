package utils

import(
	"fmt"

	v1 "github.com/solo-io/skv2/pkg/api/core.skv2.solo.io/v1"
)

// CreateCredentialsName creates the credentials name for Limited Trust
func CreateCredentialsName(virtualMesh *v1.ObjectRef) string {
	return fmt.Sprintf("%s-mtls-credential", virtualMesh.Name)
}