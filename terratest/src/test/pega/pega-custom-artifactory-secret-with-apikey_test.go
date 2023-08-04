package pega

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	k8score "k8s.io/api/core/v1"
)

const apiKeyHeaderName = "apiKeyHeaderName"
const apiKeyHeaderValue = "apiKeyHeaderValue"

func TestPegaCustomArtifactorySecretWithApiKey(t *testing.T) {
	var supportedVendors = []string{"k8s", "openshift", "eks", "gke", "aks", "pks"}
	var deploymentNames = []string{"pega", "myapp-dev"}
	var supportedOperations = []string{"install", "upgrade", "install-deploy", "deploy", "upgrade-deploy"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for _, vendor := range supportedVendors {

		for _, operation := range supportedOperations {

			for _, depName := range deploymentNames {
				fmt.Println(vendor + "-" + operation)

				var options = &helm.Options{
					SetValues: map[string]string{
						"global.deployment.name":        depName,
						"global.provider":               vendor,
						"global.actions.execute":        operation,
						"installer.upgrade.upgradeType": getUpgradeTypeForUpgradeAction(operation),
						"global.customArtifactory.authentication.apiKey.headerName": apiKeyHeaderName,
						"global.customArtifactory.authentication.apiKey.value":      apiKeyHeaderValue,
					},
				}
				yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-custom-artifactory-secret.yaml"})
				verifyCustomArtifactorySecretApiKey(t, yamlContent)
			}
		}
	}

}

func verifyCustomArtifactorySecretApiKey(t *testing.T, yamlContent string) {

	var secretobj k8score.Secret
	UnmarshalK8SYaml(t, yamlContent, &secretobj)
	secretData := secretobj.Data
	require.Equal(t, apiKeyHeaderName, string(secretData["CUSTOM_ARTIFACTORY_APIKEY_HEADER"]))
	require.Equal(t, apiKeyHeaderValue, string(secretData["CUSTOM_ARTIFACTORY_APIKEY"]))
}
