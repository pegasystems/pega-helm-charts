package pega

import (
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	k8score "k8s.io/api/core/v1"
)

func TestPegaInstallerCustomConfig(t *testing.T) {
	var supportedVendors = []string{"k8s", "openshift", "eks", "gke", "aks", "pks"}
	var supportedOperations = []string{"install", "install-deploy"}
	var custom_config = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<pegarules>\n    <env name=\"custom/Prconfig\" value=\"prconfig.xml\" />\n</pegarules>"
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for _, vendor := range supportedVendors {
		for _, operation := range supportedOperations {
			var options = &helm.Options{
				SetValues: map[string]string{
					"global.provider":        vendor,
					"global.actions.execute": operation,
					"installer.custom.configurations.prconfig": custom_config,
				},
			}

			yamlContent := RenderTemplate(t, options, helmChartPath, []string{"charts/installer/templates/pega-installer-config.yaml"})
			assertInstallerCustomConfig(t, yamlContent)
		}
	}
}

func assertInstallerCustomConfig(t *testing.T, configYaml string) {
	var installConfigMap k8score.ConfigMap
	UnmarshalK8SYaml(t, configYaml, &installConfigMap)
	installConfigData := installConfigMap.Data
	compareConfigMapData(t, installConfigData["prconfig"], "data/expectedInstallCustomPrconfig.xml")
}
