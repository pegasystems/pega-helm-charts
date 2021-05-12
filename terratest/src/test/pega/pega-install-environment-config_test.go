package pega

import (
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	k8score "k8s.io/api/core/v1"
)

func TestPegaInstallEnvironmentConfig(t *testing.T) {

	var supportedVendors = []string{"k8s", "openshift", "eks", "gke", "aks", "pks"}
	var supportedOperations = []string{"install", "install-deploy"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for _, vendor := range supportedVendors {
		for _, operation := range supportedOperations {
			var options = &helm.Options{
				SetValues: map[string]string{
					"global.provider":        vendor,
					"global.actions.execute": operation,
				},
			}

			yamlContent := RenderTemplate(t, options, helmChartPath, []string{"charts/installer/templates/pega-install-environment-config.yaml"})
			assertInstallerEnvironmentConfig(t, yamlContent, options)
		}
	}
}

func assertInstallerEnvironmentConfig(t *testing.T, configYaml string, options *helm.Options) {
	var installEnvConfigMap k8score.ConfigMap
	UnmarshalK8SYaml(t, configYaml, &installEnvConfigMap)
	installEnvConfigData := installEnvConfigMap.Data
	require.Equal(t, installEnvConfigData["DB_TYPE"], "YOUR_DATABASE_TYPE")
	require.Equal(t, installEnvConfigData["JDBC_URL"], "YOUR_JDBC_URL")
	require.Equal(t, installEnvConfigData["JDBC_CLASS"], "YOUR_JDBC_DRIVER_CLASS")
	require.Equal(t, installEnvConfigData["JDBC_DRIVER_URI"], "YOUR_JDBC_DRIVER_URI")
	require.Equal(t, installEnvConfigData["RULES_SCHEMA"], "YOUR_RULES_SCHEMA")
	require.Equal(t, installEnvConfigData["DATA_SCHEMA"], "YOUR_DATA_SCHEMA")
	require.Equal(t, installEnvConfigData["CUSTOMERDATA_SCHEMA"], "")
	require.Equal(t, installEnvConfigData["SYSTEM_NAME"], "pega")
	require.Equal(t, installEnvConfigData["PRODUCTION_LEVEL"], "2")
	require.Equal(t, installEnvConfigData["MULTITENANT_SYSTEM"], "false")
	require.Equal(t, "ADMIN_PASSWORD", installEnvConfigData["ADMIN_PASSWORD"])
	require.Equal(t, "", installEnvConfigData["STATIC_ASSEMBLER"])
	require.Equal(t, installEnvConfigData["BYPASS_UDF_GENERATION"], "true")
	require.Equal(t, installEnvConfigData["BYPASS_TRUNCATE_UPDATESCACHE"], "false")
	require.Equal(t, installEnvConfigData["JDBC_CUSTOM_CONNECTION"], "")
	require.Equal(t, installEnvConfigData["MAX_IDLE"], "5")
	require.Equal(t, installEnvConfigData["MAX_WAIT"], "-1")
	require.Equal(t, installEnvConfigData["MAX_ACTIVE"], "10")
	require.Equal(t, installEnvConfigData["ZOS_PROPERTIES"], "/opt/pega/config/DB2SiteDependent.properties")
	require.Equal(t, installEnvConfigData["DB2ZOS_UDF_WLM"], "")
	require.Equal(t, installEnvConfigData["DISTRIBUTION_KIT_URL"], "")
	require.Equal(t, installEnvConfigData["ACTION"], options.SetValues["global.actions.execute"])
	require.Equal(t, "", installEnvConfigData["DISTRIBUTION_KIT_URL"])
}
