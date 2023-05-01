package resourceapply

import (
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CreateOnlyAnnotation means that this resource should be created if it does not exist, but should not be updated
// if it already exists.  It is uniformly respected across all resources, but the first known use-cases are for
// empty config.openshift.io and initial low-level operator resources.
// Set .metadata.annotations["release.openshift.io/create-only"]="true" to have a create-only resource.
const CreateOnlyAnnotation = "release.openshift.io/create-only"

// IsCreateOnly takes metadata and returns true if the resource should only be created, not updated.
func IsCreateOnly(metadata metav1.Object) bool {
	return strings.EqualFold(metadata.GetAnnotations()[CreateOnlyAnnotation], "true")
}

// InjectCaBundleAnnotation means that this resource gets updated by service-ca to inject the CA Bundle,
// so it should not get updated.
// Set .metadata.annotations["service.beta.openshift.io/inject-cabundle"]="true" to have a create-only resource.
const InjectCaBundleAnnotation = "service.beta.openshift.io/inject-cabundle"

// IsCreateOnly takes metadata and returns true if the resource should only be created, not updated.
func IsInjectCaBundle(metadata metav1.Object) bool {
	return strings.EqualFold(metadata.GetAnnotations()[InjectCaBundleAnnotation], "true")
}
