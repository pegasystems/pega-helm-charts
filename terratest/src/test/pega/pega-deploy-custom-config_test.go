package pega

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	k8score "k8s.io/api/core/v1"
	"path/filepath"
	"strings"
	"testing"
)

func TestPegaDeployCustomConfig(t *testing.T) {
	var supportedVendors = []string{"k8s", "openshift", "eks", "gke", "aks", "pks"}
	var supportedOperations = []string{"deploy"}
	var custom_config = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<pegarules>\n    <env name=\"custom/Prconfig\" value=\"prconfig.xml\" />\n</pegarules>"
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for _, vendor := range supportedVendors {
		for _, operation := range supportedOperations {
			var options = &helm.Options{
				SetValues: map[string]string{
					"global.provider":        vendor,
					"global.actions.execute": operation,
					"global.configurations.web.prconfig": custom_config,
				},
			}

			yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-config.yaml"})
			assertPegaDeployCustomConfig(t, yamlContent)
		}
	}
}

func assertPegaDeployCustomConfig(t *testing.T, configYaml string) {
	var pegaConfigMap k8score.ConfigMap
	configSlice := strings.Split(configYaml, "---")
	for index, configData := range configSlice {
		if index == 1 {
			UnmarshalK8SYaml(t, configData, &pegaConfigMap)
			deployConfigData := pegaConfigMap.Data
			compareConfigMapData(t, deployConfigData["prconfig"], "data/expectedDeployCustomPrconfig.xml")
		}
	}
}
