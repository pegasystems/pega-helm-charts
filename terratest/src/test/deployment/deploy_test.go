package deployment

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
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

func (p *AksSpecificPodValidator) specificInitContainerValidation(t *testing.T, pod *k8score.PodSpec, options *helm.Options) {
	container := pod.Containers[0]
	if options.SetValues["global.provider"] == "aks" && options.SetValues["global.actions.execute"] == "upgrade-deploy" {
		require.Equal(t, container.Env[0].Name, "KUBERNETES_SERVICE_HOST")
		require.Equal(t, container.Env[0].Value, "API_SERVICE_ADDRESS")
		require.Equal(t, container.Env[1].Name, "KUBERNETES_SERVICE_PORT_HTTPS")
		require.Equal(t, container.Env[1].Value, "SERVICE_PORT_HTTPS")
		require.Equal(t, container.Env[2].Name, "KUBERNETES_SERVICE_PORT")
		require.Equal(t, container.Env[2].Value, "SERVICE_PORT_HTTPS")
	}
}
