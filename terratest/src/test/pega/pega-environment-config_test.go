package pega

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	k8score "k8s.io/api/core/v1"
)

func TestPegaEnvironmentConfig(t *testing.T) {
	var supportedVendors = []string{"k8s", "openshift", "eks", "gke", "aks", "pks"}
	var supportedOperations = []string{"deploy", "install-deploy", "upgrade-deploy"}
	var deploymentNames = []string{"pega", "myapp-dev"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for _, vendor := range supportedVendors {
		for _, operation := range supportedOperations {
			for _, depName := range deploymentNames {
				fmt.Println(vendor + "-" + operation + "-" + depName)

				var options = &helm.Options{
					SetValues: map[string]string{
						"global.deployment.name":        depName,
						"global.provider":               vendor,
						"global.actions.execute":        operation,
						"installer.upgrade.upgradeType": "zero-downtime",
					},
				}

				yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-environment-config.yaml"})
				VerifyEnvironmentConfig(t, yamlContent, options)
			}
		}
	}
}

func TestPegaEnvironmentConfigPegaVersionCheck(t *testing.T) {
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	var options = &helm.Options{
		SetValues: map[string]string{
			"global.provider":        "k8s",
			"global.actions.execute": "deploy",
		},
	}

	yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-environment-config.yaml"})
	VerifyEnvNotPresent(t, yamlContent, "IS_PEGA_25_OR_LATER")

	options.SetValues["global.pegaVersion"] = ""
	yamlContent = RenderTemplate(t, options, helmChartPath, []string{"templates/pega-environment-config.yaml"})
	VerifyEnvNotPresent(t, yamlContent, "IS_PEGA_25_OR_LATER")

	options.SetValues["global.pegaVersion"] = "8.25.0"
	yamlContent = RenderTemplate(t, options, helmChartPath, []string{"templates/pega-environment-config.yaml"})
	VerifyEnvValue(t, yamlContent, "IS_PEGA_25_OR_LATER", "true")

	options.SetValues["global.pegaVersion"] = "8.25.0-dev-1234"
	yamlContent = RenderTemplate(t, options, helmChartPath, []string{"templates/pega-environment-config.yaml"})
	VerifyEnvValue(t, yamlContent, "IS_PEGA_25_OR_LATER", "true")

	options.SetValues["global.pegaVersion"] = "branch-8.25.0-bugfix-BUG-12345-4"
	yamlContent = RenderTemplate(t, options, helmChartPath, []string{"templates/pega-environment-config.yaml"})
	VerifyEnvValue(t, yamlContent, "IS_PEGA_25_OR_LATER", "true")

	options.SetValues["global.pegaVersion"] = "branch-8.24.3-bugfix-BUG-12345-4"
	yamlContent = RenderTemplate(t, options, helmChartPath, []string{"templates/pega-environment-config.yaml"})
	VerifyEnvValue(t, yamlContent, "IS_PEGA_25_OR_LATER", "false")

	options.SetValues["global.pegaVersion"] = "8.26.1"
	yamlContent = RenderTemplate(t, options, helmChartPath, []string{"templates/pega-environment-config.yaml"})
	VerifyEnvValue(t, yamlContent, "IS_PEGA_25_OR_LATER", "true")

	options.SetValues["global.pegaVersion"] = "25.1.0"
	yamlContent = RenderTemplate(t, options, helmChartPath, []string{"templates/pega-environment-config.yaml"})
	VerifyEnvValue(t, yamlContent, "IS_PEGA_25_OR_LATER", "true")

	options.SetValues["global.pegaVersion"] = "26.1.1"
	yamlContent = RenderTemplate(t, options, helmChartPath, []string{"templates/pega-environment-config.yaml"})
	VerifyEnvValue(t, yamlContent, "IS_PEGA_25_OR_LATER", "true")

	options.SetValues["global.pegaVersion"] = "8.24.50"
	yamlContent = RenderTemplate(t, options, helmChartPath, []string{"templates/pega-environment-config.yaml"})
	VerifyEnvValue(t, yamlContent, "IS_PEGA_25_OR_LATER", "false")

	options.SetValues["global.pegaVersion"] = "8.8.5"
	yamlContent = RenderTemplate(t, options, helmChartPath, []string{"templates/pega-environment-config.yaml"})
	VerifyEnvValue(t, yamlContent, "IS_PEGA_25_OR_LATER", "false")

}

func TestPegaEnvironmentConfigJDBCTimeouts(t *testing.T) {
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	var options = &helm.Options{
		SetValues: map[string]string{
			"global.provider":        "k8s",
			"global.actions.execute": "deploy",
		},
	}

	yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-environment-config.yaml"})

	VerifyEnvValue(t, yamlContent, "JDBC_TIMEOUT_PROPERTIES", "")
	VerifyEnvValue(t, yamlContent, "JDBC_TIMEOUT_PROPERTIES_RW", "")
	VerifyEnvValue(t, yamlContent, "JDBC_TIMEOUT_PROPERTIES_RO", "")

	options.SetValues["global.jdbc.connectionTimeoutProperties"] = "socketTimeout=90;"
	yamlContent = RenderTemplate(t, options, helmChartPath, []string{"templates/pega-environment-config.yaml"})

	VerifyEnvValue(t, yamlContent, "JDBC_TIMEOUT_PROPERTIES", "socketTimeout=90;")
	VerifyEnvValue(t, yamlContent, "JDBC_TIMEOUT_PROPERTIES_RW", "")
	VerifyEnvValue(t, yamlContent, "JDBC_TIMEOUT_PROPERTIES_RO", "")

	options.SetValues["global.jdbc.writerConnectionTimeoutProperties"] = "socketTimeout=120;"
	yamlContent = RenderTemplate(t, options, helmChartPath, []string{"templates/pega-environment-config.yaml"})

	VerifyEnvValue(t, yamlContent, "JDBC_TIMEOUT_PROPERTIES", "socketTimeout=90;")
	VerifyEnvValue(t, yamlContent, "JDBC_TIMEOUT_PROPERTIES_RW", "socketTimeout=120;")
	VerifyEnvValue(t, yamlContent, "JDBC_TIMEOUT_PROPERTIES_RO", "")

	options.SetValues["global.jdbc.readerConnectionTimeoutProperties"] = "socketTimeout=150;"
	yamlContent = RenderTemplate(t, options, helmChartPath, []string{"templates/pega-environment-config.yaml"})

	VerifyEnvValue(t, yamlContent, "JDBC_TIMEOUT_PROPERTIES", "socketTimeout=90;")
	VerifyEnvValue(t, yamlContent, "JDBC_TIMEOUT_PROPERTIES_RW", "socketTimeout=120;")
	VerifyEnvValue(t, yamlContent, "JDBC_TIMEOUT_PROPERTIES_RO", "socketTimeout=150;")
}

func TestPegaHighlySecureCryptoModeEnabledEnvConfigParam(t *testing.T) {
	var supportedVendors = []string{"k8s", "openshift", "eks", "gke", "aks", "pks"}
	var supportedOperations = []string{"deploy", "install-deploy"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for _, vendor := range supportedVendors {

		for _, operation := range supportedOperations {

			fmt.Println(vendor + "-" + operation)

			var options = &helm.Options{
				SetValues: map[string]string{
					"global.provider":                      vendor,
					"global.actions.execute":               operation,
					"global.highlySecureCryptoModeEnabled": "false",
				},
			}

			yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-environment-config.yaml"})
			VerifyEnvNotPresent(t, yamlContent, "HIGHLY_SECURE_CRYPTO_MODE_ENABLED")

			options.SetValues["global.highlySecureCryptoModeEnabled"] = "true"
			yamlContent = RenderTemplate(t, options, helmChartPath, []string{"templates/pega-environment-config.yaml"})
			VerifyEnvValue(t, yamlContent, "HIGHLY_SECURE_CRYPTO_MODE_ENABLED", "true")

		}
	}
}

func TestFipsModeParam(t *testing.T) {
	var supportedVendors = []string{"k8s", "openshift", "eks", "gke", "aks", "pks"}
	var supportedOperations = []string{"deploy", "install-deploy"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for _, vendor := range supportedVendors {

		for _, operation := range supportedOperations {

			fmt.Println(vendor + "-" + operation)

			var options = &helm.Options{
				SetValues: map[string]string{
					"global.provider":        vendor,
					"global.actions.execute": operation,
				},
			}

			yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-environment-config.yaml"})
			VerifyEnvNotPresent(t, yamlContent, "FIPS_140_3_MODE")

			options.SetValues["global.fips140_3Mode"] = "false"
			yamlContent = RenderTemplate(t, options, helmChartPath, []string{"templates/pega-environment-config.yaml"})
			VerifyEnvNotPresent(t, yamlContent, "FIPS_140_3_MODE")

			options.SetValues["global.fips140_3Mode"] = "true"
			yamlContent = RenderTemplate(t, options, helmChartPath, []string{"templates/pega-environment-config.yaml"})
			VerifyEnvValue(t, yamlContent, "FIPS_140_3_MODE", "true")

		}
	}
}

func VerifyEnvNotPresent(t *testing.T, yamlContent string, entry string) {
	var envConfigMap k8score.ConfigMap
	UnmarshalK8SYaml(t, yamlContent, &envConfigMap)
	envConfigData := envConfigMap.Data

	_, previouslySet := envConfigData[entry]
	require.Equal(t, false, previouslySet)
}

func VerifyEnvValue(t *testing.T, yamlContent string, entry string, value string) {
	var envConfigMap k8score.ConfigMap
	UnmarshalK8SYaml(t, yamlContent, &envConfigMap)
	envConfigData := envConfigMap.Data

	val, previouslySet := envConfigData[entry]
	require.Equal(t, true, previouslySet)
	require.Equal(t, value, val)
}

func VerifyEnvironmentConfig(t *testing.T, yamlContent string, options *helm.Options) {

	var envConfigMap k8score.ConfigMap
	UnmarshalK8SYaml(t, yamlContent, &envConfigMap)

	require.Equal(t, envConfigMap.ObjectMeta.Name, getObjName(options, "-environment-config"))
	require.Equal(t, envConfigMap.ObjectMeta.Labels["ops.identifier"], "infinity")
	envConfigData := envConfigMap.Data
	require.Equal(t, envConfigData["DB_TYPE"], "YOUR_DATABASE_TYPE")
	require.Equal(t, envConfigData["JDBC_URL"], "YOUR_JDBC_URL")
	require.Equal(t, envConfigData["JDBC_CLASS"], "YOUR_JDBC_DRIVER_CLASS")
	require.Equal(t, envConfigData["JDBC_DRIVER_URI"], "YOUR_JDBC_DRIVER_URI")
	if options.SetValues["global.actions.execute"] == "upgrade-deploy" {
		require.Equal(t, envConfigData["RULES_SCHEMA"], "")
	} else {
		require.Equal(t, envConfigData["RULES_SCHEMA"], "YOUR_RULES_SCHEMA")
	}
	require.Equal(t, envConfigData["DATA_SCHEMA"], "YOUR_DATA_SCHEMA")
	require.Equal(t, envConfigData["CUSTOMERDATA_SCHEMA"], "")
	require.Equal(t, envConfigData["JDBC_CONNECTION_PROPERTIES"], "")
	require.Equal(t, envConfigData["PEGA_SEARCH_URL"], "http://"+getObjName(options, "-search"))
	require.Equal(t, envConfigData["CASSANDRA_CLUSTER"], "true")
	require.Equal(t, envConfigData["CASSANDRA_NODES"], "pega-cassandra")
	require.Equal(t, envConfigData["CASSANDRA_PORT"], "9042")
	require.Equal(t, envConfigData["CASSANDRA_ASYNC_PROCESSING_ENABLED"], "false")
	require.Equal(t, envConfigData["CASSANDRA_KEYSPACES_PREFIX"], "")
	require.Equal(t, envConfigData["CASSANDRA_EXTENDED_TOKEN_AWARE_POLICY"], "false")
	require.Equal(t, envConfigData["CASSANDRA_LATENCY_AWARE_POLICY"], "false")
	require.Equal(t, envConfigData["CASSANDRA_CUSTOM_RETRY_POLICY"], "false")
	require.Equal(t, envConfigData["CASSANDRA_CUSTOM_RETRY_POLICY_ENABLED"], "false")
	require.Equal(t, envConfigData["CASSANDRA_CUSTOM_RETRY_POLICY_COUNT"], "1")
	require.Equal(t, envConfigData["CASSANDRA_SPECULATIVE_EXECUTION_POLICY"], "false")
	require.Equal(t, envConfigData["CASSANDRA_SPECULATIVE_EXECUTION_POLICY_ENABLED"], "false")
	require.Equal(t, envConfigData["CASSANDRA_SPECULATIVE_EXECUTION_DELAY"], "100")
	require.Equal(t, envConfigData["CASSANDRA_SPECULATIVE_EXECUTION_MAX_EXECUTIONS"], "2")
	require.Equal(t, envConfigData["CASSANDRA_JMX_METRICS_ENABLED"], "true")
	require.Equal(t, envConfigData["CASSANDRA_CSV_METRICS_ENABLED"], "false")
	require.Equal(t, envConfigData["CASSANDRA_LOG_METRICS_ENABLED"], "false")
	require.Equal(t, envConfigData["ENABLE_CUSTOM_ARTIFACTORY_SSL_VERIFICATION"], "true")
}

func TestPegaRASPEnvironmentConfig(t *testing.T) {
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	var options = &helm.Options{
		SetValues: map[string]string{
			"global.provider": "k8s",
			//"global.actions.execute": "deploy",
		},
	}

	options.SetValues["global.rasp.enabled"] = "true"
	options.SetValues["global.rasp.action"] = ""
	yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-environment-config.yaml"})

	VerifyEnvValue(t, yamlContent, "IS_RASP_ENABLED", "true")
	VerifyEnvNotPresent(t, yamlContent, "RASP_ACTION")

	options.SetValues["global.rasp.enabled"] = "true"
	options.SetValues["global.rasp.action"] = "WARN"
	yamlContent = RenderTemplate(t, options, helmChartPath, []string{"templates/pega-environment-config.yaml"})

	VerifyEnvValue(t, yamlContent, "IS_RASP_ENABLED", "true")
	VerifyEnvValue(t, yamlContent, "RASP_ACTION", "WARN")

	options.SetValues["global.rasp.enabled"] = "false"
	options.SetValues["global.rasp.action"] = ""
	yamlContent = RenderTemplate(t, options, helmChartPath, []string{"templates/pega-environment-config.yaml"})

	VerifyEnvValue(t, yamlContent, "IS_RASP_ENABLED", "false")
	VerifyEnvNotPresent(t, yamlContent, "RASP_ACTION")

	options.SetValues["global.rasp.enabled"] = "false"
	options.SetValues["global.rasp.action"] = "WARN"
	yamlContent = RenderTemplate(t, options, helmChartPath, []string{"templates/pega-environment-config.yaml"})

	VerifyEnvValue(t, yamlContent, "IS_RASP_ENABLED", "false")
	VerifyEnvNotPresent(t, yamlContent, "RASP_ACTION")
}
