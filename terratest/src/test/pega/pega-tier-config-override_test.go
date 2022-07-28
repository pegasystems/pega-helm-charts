package pega

import (
	"fmt"
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	k8score "k8s.io/api/core/v1"
	"path/filepath"
	"strings"
	"testing"
)

func TestPegaTierConfigOverride(t *testing.T) {
	var supportedVendors = []string{"k8s", "openshift", "eks", "gke", "aks", "pks"}
	var supportedOperations = []string{"deploy", "install-deploy", "upgrade-deploy"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for _, vendor := range supportedVendors {

		for _, operation := range supportedOperations {

			fmt.Println(vendor + "-" + operation)

			var options = &helm.Options{
				ValuesFiles: []string{"data/pega-tier-config-override_values.yaml"},
				SetValues: map[string]string{
					"global.provider":               vendor,
					"global.actions.execute":        operation,
					"installer.upgrade.upgradeType": "zero-downtime",
				},
			}

			yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-config.yaml"})
			VerifyTierConfgOverrides(t, yamlContent, options)

		}
	}

}

// VerifyTierConfgOverrides - Performs the tier specific configuration assertions with the values
func VerifyTierConfgOverrides(t *testing.T, yamlContent string, options *helm.Options) {
	var pegaConfigMap k8score.ConfigMap
	configSlice := strings.Split(yamlContent, "---")
	for index, configData := range configSlice {
		if index >= 1 && index <= 3 {
			UnmarshalK8SYaml(t, configData, &pegaConfigMap)
			pegaConfigMapData := pegaConfigMap.Data
			require.Equal(t, "prconfig override", pegaConfigMapData["prconfig.xml"])
			require.Equal(t, "context.xml override", pegaConfigMapData["context.xml.tmpl"])
			require.Equal(t, "prlog4j2 override", pegaConfigMapData["prlog4j2.xml"])
			require.Equal(t, "server.xml override", pegaConfigMapData["server.xml.tmpl"])
			require.Equal(t, "web.xml override", pegaConfigMapData["web.xml"])
		}
	}
}
