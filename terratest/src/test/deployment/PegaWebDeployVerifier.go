package deployment

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	k8sapps "k8s.io/api/apps/v1"
	k8score "k8s.io/api/core/v1"
	"testing"
)

type PegaWebDeployVerifier struct {
	DeploymentVerifier
}

func (v *PegaWebDeployVerifier) getPod() *k8score.PodSpec {
	return &v.k8sDeployment.Spec.Template.Spec
}
func NewPegaWebDeployVerifier(t *testing.T, helmOptions *helm.Options, initContainers []string, deployment *k8sapps.Deployment) *PegaWebDeployVerifier {
	v := &PegaWebDeployVerifier{
		DeploymentVerifier: *NewDeploymentVerifier(t, helmOptions, initContainers, deployment),
	}
	v._nodeType = "WebUser"
	v._passivationTimeout = "900"
	return v
}
