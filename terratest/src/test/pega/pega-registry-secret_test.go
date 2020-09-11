package pega

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	k8score "k8s.io/api/core/v1"
	"path/filepath"
	"testing"
	"fmt"
)



func TestPegaRegistrySecret(t *testing.T){
	var supportedVendors = []string{"k8s","openshift","eks","gke","aks","pks"}
	var supportedOperations =  []string{"install","install-deploy", "upgrade", "upgrade-deploy"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)


	for _,vendor := range supportedVendors{

		for _,operation := range supportedOperations{

			fmt.Println(vendor + "-" + operation)

			var options = &helm.Options{			
				SetValues: map[string]string{
					"global.provider":        vendor,
					"global.actions.execute": operation,
			 	},
		    }

			yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-registry-secret.yaml"})
			VerfiyRegistrySecret(t,yamlContent)
		}
	}


}

// VerfiyRegistrySecret - Verifies the registry secret deployed with the values as provided in default values.yaml
func VerfiyRegistrySecret(t *testing.T, yamlContent string) {
	var registrySecretObj k8score.Secret
	UnmarshalK8SYaml(t, yamlContent, &registrySecretObj)
	reqgistrySecretData := registrySecretObj.Data
	require.Contains(t, string(reqgistrySecretData[".dockerconfigjson"]), "YOUR_DOCKER_REGISTRY")
	require.Contains(t, string(reqgistrySecretData[".dockerconfigjson"]), "WU9VUl9ET0NLRVJfUkVHSVNUUllfVVNFUk5BTUU6WU9VUl9ET0NLRVJfUkVHSVNUUllfUEFTU1dPUkQ=")
}