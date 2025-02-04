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

func TestPegaHazelcastEnvironmentConfigForClient(t *testing.T) {
	var supportedVendors = []string{"k8s", "openshift", "eks", "gke", "aks", "pks"}
	var supportedOperations = []string{"deploy", "install-deploy"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for _, vendor := range supportedVendors {

		for _, operation := range supportedOperations {

			fmt.Println(vendor + "-" + operation)

			var options = &helm.Options{
				SetValues: map[string]string{
					"global.provider":                           vendor,
					"global.actions.execute":                    operation,
					"hazelcast.enabled":                         "true",
					"hazelcast.migration.embeddedToCSMigration": "false",
				},
			}

			yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-environment-config.yaml"})
			VerifyPegaHazelcastEnvironmentConfigForClient(t, yamlContent, options)

		}
	}

	for _, vendor := range supportedVendors {

		for _, operation := range supportedOperations {

			fmt.Println(vendor + "-" + operation)

			var options = &helm.Options{
				SetValues: map[string]string{
					"global.provider":                           vendor,
					"global.actions.execute":                    operation,
					"hazelcast.enabled":                         "false",
					"hazelcast.clusteringServiceEnabled":        "true",
					"hazelcast.migration.embeddedToCSMigration": "false",
				},
			}

			yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-environment-config.yaml"})
			VerifyClusteringServiceEnvironmentConfigForClient(t, yamlContent, options, false, false)

		}
	}

	for _, vendor := range supportedVendors {

		for _, operation := range supportedOperations {

			fmt.Println(vendor + "-" + operation)

			var options = &helm.Options{
				SetValues: map[string]string{
					"global.provider":                           vendor,
					"global.actions.execute":                    operation,
					"hazelcast.enabled":                         "false",
					"hazelcast.clusteringServiceEnabled":        "true",
					"hazelcast.migration.embeddedToCSMigration": "false",
				},
			}

			yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-environment-config.yaml"})
			VerifyClusteringServiceEnvironmentConfigForClient(t, yamlContent, options, false, false)

		}
	}

}

func TestPegaHazelcastEnvironmentConfigForClientWithSSL(t *testing.T) {
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
					"hazelcast.enabled":                    "false",
					"hazelcast.clusteringServiceEnabled":   "true",
					"hazelcast.encryption.enabled":         "true",
					"global.highlySecureCryptoModeEnabled": "false",
				},
			}

			yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-environment-config.yaml"})
			VerifyClusteringServiceEnvironmentConfigForClient(t, yamlContent, options, true, false)

		}
	}
}

func TestPegaHazelcastEnvironmentConfigForClientWithHighlySecureCryptoModeEnabled(t *testing.T) {
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
					"hazelcast.enabled":                    "false",
					"hazelcast.clusteringServiceEnabled":   "true",
					"hazelcast.encryption.enabled":         "true",
					"global.highlySecureCryptoModeEnabled": "true",
				},
			}

			yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-environment-config.yaml"})
			VerifyClusteringServiceEnvironmentConfigForClient(t, yamlContent, options, true, true)

		}
	}
}

func VerifyPegaHazelcastEnvironmentConfigForClient(t *testing.T, yamlContent string, options *helm.Options) {

	var envConfigMap k8score.ConfigMap
	statefulSlice := strings.Split(yamlContent, "---")
	for index, statefulInfo := range statefulSlice {
		if index >= 1 {
			UnmarshalK8SYaml(t, statefulInfo, &envConfigMap)
			envConfigData := envConfigMap.Data
			require.Equal(t, envConfigData["HZ_DISCOVERY_K8S"], "true")
			require.Equal(t, envConfigData["HZ_CLIENT_MODE"], "true")
			require.Equal(t, envConfigData["HZ_CLUSTER_NAME"], "PRPC")
			require.Equal(t, envConfigData["HZ_SERVER_HOSTNAME"], "pega-hazelcast-service.default.svc.cluster.local")
		}
	}
}

func VerifyClusteringServiceEnvironmentConfigForClient(t *testing.T, yamlContent string, options *helm.Options,
	ssl bool, highlySecureCryptoModeEnabled bool) {

	var envConfigMap k8score.ConfigMap
	statefulSlice := strings.Split(yamlContent, "---")
	for index, statefulInfo := range statefulSlice {
		if index >= 1 {
			UnmarshalK8SYaml(t, statefulInfo, &envConfigMap)
			envConfigData := envConfigMap.Data
			require.Equal(t, envConfigData["HZ_DISCOVERY_K8S"], "true")
			require.Equal(t, envConfigData["HZ_CLIENT_MODE"], "true")
			require.Equal(t, envConfigData["HZ_CLUSTER_NAME"], "prpchz")
			require.Equal(t, envConfigData["HZ_SERVER_HOSTNAME"], "clusteringservice-service.default.svc.cluster.local")
			if ssl {
				require.Equal(t, envConfigData["HZ_SSL_ENABLED"], "true")
				require.Equal(t, envConfigData["HZ_SSL_PROTOCOL"], "TLS")
				require.Equal(t, envConfigData["HZ_SSL_KEY_STORE_NAME"], "cluster-keystore.jks")
				require.Equal(t, envConfigData["HZ_SSL_TRUST_STORE_NAME"], "cluster-truststore.jks")
				if highlySecureCryptoModeEnabled {
					require.Equal(t, envConfigData["HIGHLY_SECURE_CRYPTO_MODE_ENABLED"], "true")
					require.Equal(t, envConfigData["HZ_SSL_ALGO"], "PKIX")
				} else {
					require.Equal(t, envConfigData["HZ_SSL_ALGO"], "SunX509")
				}
			}
		}
	}
}
