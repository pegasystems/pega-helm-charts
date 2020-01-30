package deployment

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	k8score "k8s.io/api/core/v1"
	"testing"
)

func TestPegaK8sDeployment(t *testing.T) {
	NewPegaDeploymentTest(
		"K8s",
		"deploy",
		[]string{"wait-for-pegasearch", "wait-for-cassandra"},
		t,
	).Run()
}

func TestPegaAskDeployment(t *testing.T) {
	test := NewPegaDeploymentTest(
		"aks",
		"deploy",
		[]string{"wait-for-pegasearch", "wait-for-cassandra"},
		t,
	)
	test.setPodValidator(AksSpecificPodValidator{})
	test.Run()
}

type AksSpecificPodValidator struct {
	PodValidator
}

func (p *AksSpecificPodValidator) specificInitContainerValidation(pod *k8score.PodSpec, options *helm.Options) {

}
