package pega

import (
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	k8score "k8s.io/api/core/v1"
)

func TestPegaUpgradeEnvironmentConfig(t *testing.T) {
	var supportedVendors = []string{"k8s", "openshift", "eks", "gke", "aks", "pks"}
	var supportedOperations = []string{"upgrade", "upgrade-deploy"}

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

			yamlContent := RenderTemplate(t, options, helmChartPath, []string{"charts/installer/templates/pega-upgrade-environment-config.yaml"})
			assertUpgradeEnvironmentConfig(t, yamlContent, options)
		}
	}
}

func assertUpgradeEnvironmentConfig(t *testing.T, configYaml string, options *helm.Options) {
	var upgradeEnvConfigMap k8score.ConfigMap
	UnmarshalK8SYaml(t, configYaml, &upgradeEnvConfigMap)
	upgradeEnvConfigData := upgradeEnvConfigMap.Data

	require.Equal(t, upgradeEnvConfigData["DB_TYPE"], "YOUR_DATABASE_TYPE")
	require.Equal(t, upgradeEnvConfigData["JDBC_URL"], "YOUR_JDBC_URL")
	require.Equal(t, upgradeEnvConfigData["JDBC_CLASS"], "YOUR_JDBC_DRIVER_CLASS")
	require.Equal(t, upgradeEnvConfigData["JDBC_DRIVER_URI"], "YOUR_JDBC_DRIVER_URI")
	require.Equal(t, upgradeEnvConfigData["RULES_SCHEMA"], "YOUR_RULES_SCHEMA")
	require.Equal(t, upgradeEnvConfigData["DATA_SCHEMA"], "YOUR_DATA_SCHEMA")
	require.Equal(t, upgradeEnvConfigData["CUSTOMERDATA_SCHEMA"], "")
	require.Equal(t, upgradeEnvConfigData["UPGRADE_TYPE"], "in-place")
	require.Equal(t, upgradeEnvConfigData["MULTITENANT_SYSTEM"], "false")
	require.Equal(t, upgradeEnvConfigData["BYPASS_UDF_GENERATION"], "true")
	require.Equal(t, upgradeEnvConfigData["ZOS_PROPERTIES"], "/opt/pega/config/DB2SiteDependent.properties")
	require.Equal(t, upgradeEnvConfigData["DB2ZOS_UDF_WLM"], "")
	require.Equal(t, upgradeEnvConfigData["TARGET_RULES_SCHEMA"], "")
	require.Equal(t, upgradeEnvConfigData["TARGET_ZOS_PROPERTIES"], "/opt/pega/config/DB2SiteDependent.properties")
	require.Equal(t, upgradeEnvConfigData["MIGRATION_DB_LOAD_COMMIT_RATE"], "100")
	require.Equal(t, upgradeEnvConfigData["UPDATE_EXISTING_APPLICATIONS"], "false")
	require.Equal(t, upgradeEnvConfigData["UPDATE_APPLICATIONS_SCHEMA"], "false")
	require.Equal(t, upgradeEnvConfigData["RUN_RULESET_CLEANUP"], "false")
	require.Equal(t, upgradeEnvConfigData["REBUILD_INDEXES"], "false")
	require.Equal(t, upgradeEnvConfigData["DISTRIBUTION_KIT_URL"], "")
}
