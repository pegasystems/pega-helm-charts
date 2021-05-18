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
	var supportedOperations = []string{"install", "install-deploy", "upgrade", "upgrade-deploy"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for _, vendor := range supportedVendors {
		for _, operation := range supportedOperations {
			var options = &helm.Options{
				SetValues: map[string]string{
					"global.provider":        vendor,
					"global.actions.execute": operation,
					"global.installer.custom.configurations": "{\"installer\":{\"prconfig.xml\":\"<?xml version=\\\"1.0\\\" encoding=\\\"UTF-8\\\"?>\\n<pegarules>\\n    <env name=\\\"custom/Prconfig\\\" value=\\\"prconfig.xml\\\" />\\n</pegarules>\",\"setupDatabase.properties\":\"# Properties file for use with Pega Deployment Utilities.\\n# For more information, see the Pega Platform help.\\n\\n################### COMMON PROPERTIES - DB CONNECTION ##################\\n########################################################################\\n\\n# CONNECTION INFORMATION\\npega.jdbc.custom.jar={{ .Env.Custom }}\",\"prbootstrap.properties\":\"install.{{ .Env.CUSTOM_DB_TYPE }}.schema={{ .Env.CUSTOM_DATA_SCHEMA }}\\ninitialization.settingsource=file\\ncom.pega.pegarules.priv.LogHelper.USE_LOG4JV2=true\"}}",
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

	compareConfigMapData(t, installConfigData["prconfig.xml"], "data/expectedInstallCustomPrconfig.xml")
	compareConfigMapData(t, installConfigData["setupDatabase.properties"], "data/expectedInstallCustomSetupdatabase.properties")
	compareConfigMapData(t, installConfigData["prbootstrap.properties"], "data/expectedInstallCustomPRbootstrap.properties")
}
