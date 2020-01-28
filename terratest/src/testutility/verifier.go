package testutility

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	k8score "k8s.io/api/apps/v1"
	"k8s.io/api/core/v1"
	"testing"
)

type VerifierImpl struct {
	Verifier
	pegaDeployment
	k8sInformationExtractor _k8sInformationExtractor
	t                       *testing.T
	_helmOptions            *helm.Options
	_nodeType               string
	_initContainers         []string
	_passivationTimeout     string
}

type Verifier interface {
	verify()
}

type _k8sInformationExtractor interface {
	getPod() *k8score.PodSpec
	getDeploymentMetadata() DeploymentMetadata
}
type _k8sDeploymentExtractor struct {
	_k8sDeployment *k8score.Deployment
}
