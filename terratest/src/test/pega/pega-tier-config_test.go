package pega

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	k8score "k8s.io/api/core/v1"
)

func TestPegaTierConfig(t *testing.T) {
	var supportedVendors = []string{"k8s", "openshift", "eks", "gke", "aks", "pks"}
	var supportedOperations = []string{"deploy", "install-deploy", "upgrade-deploy"}
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

				yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-config.yaml"})
				VerifyTierConfig(t, yamlContent, options)

				options.SetValues["global.fips140_3Mode"] = "true"
				yamlContent = RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-config.yaml"})
				VerifyTierConfig(t, yamlContent, options)
			}
		}
	}

}

// VerifyTierConfig - Performs the tier specific configuration assertions with the values as provided in default values.yaml
func VerifyTierConfig(t *testing.T, yamlContent string, options *helm.Options) {
	var pegaConfigMap k8score.ConfigMap
	configSlice := strings.Split(yamlContent, "---")
	for index, configData := range configSlice {
		if index >= 1 && index <= 3 {
			var tierName string
			switch index {
			case 1:
				tierName = "-web"
			case 2:
				tierName = "-batch"
			case 3:
				tierName = "-stream"
			}

			UnmarshalK8SYaml(t, configData, &pegaConfigMap)

			require.Equal(t, pegaConfigMap.ObjectMeta.Name, getObjName(options, tierName))

			pegaConfigMapData := pegaConfigMap.Data
			compareConfigMapData(t, pegaConfigMapData["prconfig.xml"], "data/expectedInstallDeployPrconfig.xml")
			compareConfigMapData(t, pegaConfigMapData["context.xml.tmpl"], "data/expectedInstallDeployContext.xml.tmpl")
			compareConfigMapData(t, pegaConfigMapData["prlog4j2.xml"], "data/expectedInstallDeployPRlog4j2.xml")
			compareConfigMapData(t, pegaConfigMapData["server.xml.tmpl"], "data/expectedInstallDeployServer.xml.tmpl")
			require.Equal(t, "", pegaConfigMapData["web.xml"])
			compareConfigMapData(t, pegaConfigMapData["java.security.overwrite.tmpl"], "data/expectedJava.security.overwrite.tmpl")
		}
	}
}
