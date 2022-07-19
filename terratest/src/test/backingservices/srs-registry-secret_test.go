package backingservices

import (
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"testing"
)
func TestSRSRegistrySecretDefaultName(t *testing.T){
	helmChartParser := NewHelmConfigParser(
		NewHelmTestFromTemplate(t, helmChartRelativePath, map[string]string{
			"srs.enabled": "true",
			"srs.deploymentName": "test-srs",
			"srs.srsStorage.tls.enabled": "false",
		},
			[]string{"charts/srs/templates/registrysecret.yaml"}),
	)

	var secret corev1.Secret
	helmChartParser.getResourceYAML(SearchResourceOption{
		Name: "test-srs-reg-secret",
		Kind: "Secret",
	}, &secret)
	registryData := secret.Data
	require.Contains(t, string(registryData[".dockerconfigjson"]), "YOUR_DOCKER_REGISTRY")
	require.Contains(t, string(registryData[".dockerconfigjson"]), "WU9VUl9ET0NLRVJfUkVHSVNUUllfVVNFUk5BTUU6WU9VUl9ET0NLRVJfUkVHSVNUUllfUEFTU1dPUkQ=")
}

func TestSRSRegistrySecretCustomName(t *testing.T){
	helmChartParser := NewHelmConfigParser(
		NewHelmTestFromTemplate(t, helmChartRelativePath, map[string]string{
			"srs.enabled": "true",
			"srs.deploymentName": "custom-srs",
			"global.imageCredentials.registry": "docker-repo.acme.io",
			"global.imageCredentials.username": "acmeuser",
			"global.imageCredentials.password": "thisisapwd",
			"srs.srsStorage.tls.enabled": "false",
		},
			[]string{"charts/srs/templates/registrysecret.yaml"}),
	)

	var secret corev1.Secret
	helmChartParser.getResourceYAML(SearchResourceOption{
		Name: "custom-srs-reg-secret",
		Kind: "Secret",
	}, &secret)
	registryData := secret.Data
	require.Contains(t, string(registryData[".dockerconfigjson"]), "docker-repo.acme.io")
	require.Contains(t, string(registryData[".dockerconfigjson"]), "YWNtZXVzZXI6dGhpc2lzYXB3ZA==")
}