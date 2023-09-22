package pega

import (
	"fmt"
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	k8score "k8s.io/api/core/v1"
	"path/filepath"
	"testing"
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
