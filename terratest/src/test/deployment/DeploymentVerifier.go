package deployment

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	k8sapps "k8s.io/api/apps/v1"
	"testing"
)

type DeploymentVerifier struct {
	Verifier
	k8sDeployment *k8sapps.Deployment
}

func NewDeploymentVerifier(t *testing.T, helmOptions *helm.Options, initContainers []string, deployment *k8sapps.Deployment) *DeploymentVerifier {
	verifierImpl := *NewDeployVerifier(t, helmOptions, initContainers)
	verifierImpl.k8sInformationExtractor = K8sDeploymentExtractor{
		K8sDeployment: deployment,
	}
	return &DeploymentVerifier{
		k8sDeployment: deployment,
		Verifier:      verifierImpl,
	}
}
