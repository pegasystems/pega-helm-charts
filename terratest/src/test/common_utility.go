package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	k8score "k8s.io/api/core/v1"
)

func VerifyCredentialsSecret(t *testing.T, helmChartPath string) {

	secretOutput := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/pega-credentials-secret.yaml"})
	var secretobj k8score.Secret
	helm.UnmarshalK8SYaml(t, secretOutput, &secretobj)
	secretData := secretobj.Data
	require.Equal(t, string(secretData["DB_USERNAME"]), "YOUR_JDBC_USERNAME")
	require.Equal(t, string(secretData["DB_PASSWORD"]), "YOUR_JDBC_PASSWORD")
}

func VerfiyRegistrySecret(t *testing.T, helmChartPath string) {

	registrySecret := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/pega-registry-secret.yaml"})
	var registrySecretObj k8score.Secret
	helm.UnmarshalK8SYaml(t, registrySecret, &registrySecretObj)
	reqgistrySecretData := registrySecretObj.Data
	require.Contains(t, string(reqgistrySecretData[".dockerconfigjson"]), "YOUR_DOCKER_REGISTRY")
}

func SplitOutput() {

}
