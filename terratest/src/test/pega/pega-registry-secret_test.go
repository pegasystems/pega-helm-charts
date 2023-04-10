package pega

import (
	"fmt"
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	k8score "k8s.io/api/core/v1"
	"path/filepath"
	"testing"
)

func TestPegaRegistrySecret(t *testing.T) {
	var supportedVendors = []string{"k8s", "openshift", "eks", "gke", "aks", "pks"}
	var supportedOperations = []string{"install", "install-deploy", "upgrade-deploy"}
	var deploymentNames = []string{"pega", "myapp-dev"}

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
						"installer.upgrade.upgradeType": "zero-downtime",
					},
				}

				yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-registry-secret.yaml"})
				VerfiyRegistrySecret(t, yamlContent, options)
			}
		}
	}

}

// VerfiyRegistrySecret - Verifies the registry secret deployed with the values as provided in default values.yaml
func VerfiyRegistrySecret(t *testing.T, yamlContent string, options *helm.Options) {
	var registrySecretObj k8score.Secret
	UnmarshalK8SYaml(t, yamlContent, &registrySecretObj)
	require.Equal(t, registrySecretObj.ObjectMeta.Name, getObjName(options, "-registry-secret"))
	reqgistrySecretData := registrySecretObj.Data
	require.Contains(t, string(reqgistrySecretData[".dockerconfigjson"]), "YOUR_DOCKER_REGISTRY")
	require.Contains(t, string(reqgistrySecretData[".dockerconfigjson"]), "WU9VUl9ET0NLRVJfUkVHSVNUUllfVVNFUk5BTUU6WU9VUl9ET0NLRVJfUkVHSVNUUllfUEFTU1dPUkQ=")
}
