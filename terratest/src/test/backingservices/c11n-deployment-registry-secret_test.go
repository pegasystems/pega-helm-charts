package backingservices

import (
	"testing"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
)

func TestC11NdeploymentRegistrySecretDefaultName(t *testing.T) {
	helmChartParser := C11NHelmConfigParser(
		NewHelmTestFromTemplate(t, helmChartRelativePath, map[string]string{
			"c11n-messaging.enabled":        "true",
			"c11n-messaging.deploymentName": "c11n-messaging",
		},
			[]string{"charts/constellation-messaging/templates/registrysecret.yaml"}),
	)

	var secret corev1.Secret
	helmChartParser.getResourceYAML(SearchResourceOption{
		Name: "c11n-messaging-reg-secret",
		Kind: "Secret",
	}, &secret)
	registryData := secret.Data
	require.Contains(t, string(registryData[".dockerconfigjson"]), "YOUR_DOCKER_REGISTRY")
	require.Contains(t, string(registryData[".dockerconfigjson"]), "WU9VUl9ET0NLRVJfUkVHSVNUUllfVVNFUk5BTUU6WU9VUl9ET0NLRVJfUkVHSVNUUllfUEFTU1dPUkQ=")
}
