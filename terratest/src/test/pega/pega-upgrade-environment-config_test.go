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
	var supportedOperations = []string{"upgrade-deploy"}
	var expectedValues = map[string]string{"PEGA_REST_SERVER_URL": "http://pega-web:80/prweb/PRRestService", "PEGA_REST_USERNAME": "username", "PEGA_REST_PASSWORD": "password"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for _, vendor := range supportedVendors {
		for _, operation := range supportedOperations {
			var options = &helm.Options{
				SetValues: map[string]string{
					"global.provider":                     vendor,
					"global.actions.execute":              operation,
					"installer.upgrade.upgradeType":       "zero-downtime",
					"installer.upgrade.pegaRESTServerURL": expectedValues["PEGA_REST_SERVER_URL"],
					"installer.upgrade.pegaRESTUsername":  expectedValues["PEGA_REST_USERNAME"],
					"installer.upgrade.pegaRESTPassword":  expectedValues["PEGA_REST_PASSWORD"],
				},
			}

			yamlContent := RenderTemplate(t, options, helmChartPath, []string{"charts/installer/templates/pega-upgrade-environment-config.yaml"})
			assertUpgradeEnvironmentConfig(t, yamlContent, options, expectedValues)
		}
	}
}

type restURLExpectation struct {
	valuesFile  string
	expectedURL string
}

func TestPegaUpgradeEnvironmentConfig_DetermineRestUrl(t *testing.T) {
	var supportedVendors = []string{"k8s", "openshift", "eks", "gke", "aks", "pks"}
	var supportedOperations = []string{"upgrade-deploy"}
	var expectedValues = map[string]string{"PEGA_REST_SERVER_URL": "http://pega-web:80/prweb/PRRestService", "PEGA_REST_USERNAME": "username", "PEGA_REST_PASSWORD": "password"}

	var expectations = []restURLExpectation{
		{
			valuesFile:  "/data/values_zdt_upgrade_http_non_default_port81.yaml",
			expectedURL: "http://pega-web:81/prweb/PRRestService",
		},
		{
			valuesFile:  "/data/values_zdt_upgrade_https_default_port.yaml",
			expectedURL: "https://pega-web:443/prweb/PRRestService",
		},
		{
			valuesFile:  "/data/values_zdt_upgrade_https_non_default_port.yaml",
			expectedURL: "https://pega-web:5443/prweb/PRRestService",
		},
		{
			valuesFile:  "/data/values_zdt_upgrade_http_default_port_renamed_deployment_renamed_tier.yaml",
			expectedURL: "http://xxx-webcustom:80/somethingotherthanprweb/PRRestService",
		},
	}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	testsPath, err := filepath.Abs(PegaHelmChartTestsPath)
	require.NoError(t, err)

	for _, expectation := range expectations {
		expectedValues["PEGA_REST_SERVER_URL"] = expectation.expectedURL
		for _, vendor := range supportedVendors {
			for _, operation := range supportedOperations {
				var options = &helm.Options{
					SetValues: map[string]string{
						"global.provider":                    vendor,
						"global.actions.execute":             operation,
						"installer.upgrade.upgradeType":      "zero-downtime",
						"installer.upgrade.pegaRESTUsername": expectedValues["PEGA_REST_USERNAME"],
						"installer.upgrade.pegaRESTPassword": expectedValues["PEGA_REST_PASSWORD"],
					},
				}

				yamlContent := RenderTemplate(t, options, helmChartPath, []string{"charts/installer/templates/pega-upgrade-environment-config.yaml"}, "--values", testsPath+expectation.valuesFile)
				assertUpgradeEnvironmentConfig(t, yamlContent, options, expectedValues)
			}
		}
	}
}

func assertUpgradeEnvironmentConfig(t *testing.T, configYaml string, options *helm.Options, expectedValues map[string]string) {
	var upgradeEnvConfigMap k8score.ConfigMap

	UnmarshalK8SYaml(t, configYaml, &upgradeEnvConfigMap)
	upgradeEnvConfigData := upgradeEnvConfigMap.Data

	require.Equal(t, upgradeEnvConfigData["ADMIN_PASSWORD"], "ADMIN_PASSWORD")
	require.Equal(t, upgradeEnvConfigData["DB_TYPE"], "YOUR_DATABASE_TYPE")
	require.Equal(t, upgradeEnvConfigData["JDBC_URL"], "YOUR_JDBC_URL")
	require.Equal(t, upgradeEnvConfigData["JDBC_CLASS"], "YOUR_JDBC_DRIVER_CLASS")
	require.Equal(t, upgradeEnvConfigData["JDBC_DRIVER_URI"], "YOUR_JDBC_DRIVER_URI")
	require.Equal(t, upgradeEnvConfigData["RULES_SCHEMA"], "YOUR_RULES_SCHEMA")
	require.Equal(t, upgradeEnvConfigData["DATA_SCHEMA"], "YOUR_DATA_SCHEMA")
	require.Equal(t, upgradeEnvConfigData["CUSTOMERDATA_SCHEMA"], "")
	require.Equal(t, upgradeEnvConfigData["UPGRADE_TYPE"], "zero-downtime")
	require.Equal(t, upgradeEnvConfigData["MT_SYSTEM"], "")
	require.Equal(t, upgradeEnvConfigData["BYPASS_UDF_GENERATION"], "true")
	require.Equal(t, upgradeEnvConfigData["ZOS_PROPERTIES"], "/opt/pega/config/DB2SiteDependent.properties")
	require.Equal(t, upgradeEnvConfigData["DB2ZOS_UDF_WLM"], "")
	require.Equal(t, upgradeEnvConfigData["JDBC_CUSTOM_CONNECTION"], "")
	require.Equal(t, upgradeEnvConfigData["TARGET_RULES_SCHEMA"], "")
	require.Equal(t, upgradeEnvConfigData["TARGET_ZOS_PROPERTIES"], "/opt/pega/config/DB2SiteDependent.properties")
	require.Equal(t, upgradeEnvConfigData["MIGRATION_DB_LOAD_COMMIT_RATE"], "100")
	require.Equal(t, upgradeEnvConfigData["UPDATE_EXISTING_APPLICATIONS"], "false")
	require.Equal(t, upgradeEnvConfigData["UPDATE_APPLICATIONS_SCHEMA"], "false")
	require.Equal(t, upgradeEnvConfigData["RUN_RULESET_CLEANUP"], "false")
	require.Equal(t, upgradeEnvConfigData["REBUILD_INDEXES"], "false")
	require.Equal(t, upgradeEnvConfigData["PEGA_REST_SERVER_URL"], expectedValues["PEGA_REST_SERVER_URL"])
	require.Equal(t, upgradeEnvConfigData["PEGA_REST_USERNAME"], expectedValues["PEGA_REST_USERNAME"])
	require.Equal(t, upgradeEnvConfigData["PEGA_REST_PASSWORD"], expectedValues["PEGA_REST_PASSWORD"])
	require.Equal(t, upgradeEnvConfigData["DISTRIBUTION_KIT_URL"], "")
	require.Equal(t, upgradeEnvConfigData["ENABLE_CUSTOM_ARTIFACTORY_SSL_VERIFICATION"], "true")
	require.Equal(t, upgradeEnvConfigData["AUTOMATIC_RESUME_ENABLED"], "false")
}
