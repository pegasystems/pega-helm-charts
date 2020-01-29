package verifier

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	k8sapps "k8s.io/api/apps/v1"
	. "test/deployment"
	"testing"
)

type DeploymentVerifier struct {
	VerifierImpl
	k8sDeployment *k8sapps.Deployment
}

func NewDeploymentVerifier(t *testing.T, helmOptions *helm.Options, initContainers []string, deployment *k8sapps.Deployment) *DeploymentVerifier {
	verifierImpl := *NewDeployVerifier(t, helmOptions, initContainers)
	verifierImpl.k8sInformationExtractor = K8sDeploymentExtractor{
		K8sDeployment: deployment,
	}
	return &DeploymentVerifier{
		k8sDeployment: deployment,
		VerifierImpl:  verifierImpl,
	}
}
