package resourceapply

import (
	"context"

	"github.com/google/go-cmp/cmp"
	"github.com/openshift/cluster-version-operator/lib/resourcemerge"
	batchv1 "k8s.io/api/batch/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	batchclientv1 "k8s.io/client-go/kubernetes/typed/batch/v1"
	"k8s.io/klog/v2"
	"k8s.io/utils/pointer"
)

// ApplyJobv1 applies the required Job to the cluster.
func ApplyJobv1(ctx context.Context, client batchclientv1.JobsGetter, required *batchv1.Job, reconciling bool) (*batchv1.Job, bool, error) {
	existing, err := client.Jobs(required.Namespace).Get(ctx, required.Name, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		klog.V(2).Infof("Job %s/%s not found, creating", required.Namespace, required.Name)
		actual, err := client.Jobs(required.Namespace).Create(ctx, required, metav1.CreateOptions{})
		return actual, true, err
	}
	if err != nil {
		return nil, false, err
	}
	// if we only create this resource, we have no need to continue further
	// Also skip if inject-cabundle is set, as it means the resource gets updated by service-ca
	if IsCreateOnly(required) || IsInjectCaBundle(required) {
		return nil, false, nil
	}

	modified := pointer.BoolPtr(false)
	resourcemerge.EnsureJob(modified, existing, *required)
	if !*modified {
		return existing, false, nil
	}

	if reconciling {
		klog.V(2).Infof("Updating Job %s/%s due to diff: %v", required.Namespace, required.Name, cmp.Diff(existing, required))
	}

	actual, err := client.Jobs(required.Namespace).Update(ctx, existing, metav1.UpdateOptions{})
	return actual, true, err
}
