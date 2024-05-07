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

const DefaultSRSAuthScope = "pega.search:full"
const DefaultPrivateKeyAlgorithm = "RS256"

func TestPegaConfigWithoutSRS(t *testing.T) {
	var supportedVendors = []string{"k8s", "openshift", "eks", "gke", "aks", "pks"}
	var supportedOperations = []string{"deploy", "install-deploy"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for _, vendor := range supportedVendors {
		for _, operation := range supportedOperations {
			fmt.Println(vendor + "-" + operation)

			var options = &helm.Options{
				SetValues: map[string]string{
					"global.provider":                  vendor,
					"global.actions.execute":           operation,
					"pegasearch.externalSearchService": "false",
					"pegasearch.externalURL":           "https://srs:9200",
				},
			}

			yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-environment-config.yaml"})
			VerifyPegaWithoutExternalSRSEnvironmentConfig(t, yamlContent)
		}
	}
}

func TestPegaConfigWithSRSAndAuthDisabled(t *testing.T) {
	var supportedVendors = []string{"k8s", "openshift", "eks", "gke", "aks", "pks"}
	var supportedOperations = []string{"deploy", "install-deploy"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for _, vendor := range supportedVendors {
		for _, operation := range supportedOperations {
			fmt.Println(vendor + "-" + operation)

			var options = &helm.Options{
				SetValues: map[string]string{
					"global.provider":                  vendor,
					"global.actions.execute":           operation,
					"pegasearch.externalSearchService": "true",
					"pegasearch.externalURL":           "https://srs:9200",
				},
			}

			yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-environment-config.yaml"})
			VerifyPegaWithExternalSRSEnvironmentConfig(t, yamlContent, false, "", "")
		}
	}
}

func TestPegaConfigWithSRSAndAuthEnabledAndAllAuthParametersProvided(t *testing.T) {
	var supportedVendors = []string{"k8s", "openshift", "eks", "gke", "aks", "pks"}
	var supportedOperations = []string{"deploy", "install-deploy"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for _, vendor := range supportedVendors {
		for _, operation := range supportedOperations {
			fmt.Println(vendor + "-" + operation)

			var options = &helm.Options{
				SetValues: map[string]string{
					"global.provider":                        vendor,
					"global.actions.execute":                 operation,
					"pegasearch.externalSearchService":       "true",
					"pegasearch.externalURL":                 "https://srs:9200",
					"pegasearch.srsAuth.enabled":             "true",
					"pegasearch.srsAuth.url":                 "https://auth-service",
					"pegasearch.srsAuth.clientId":            "client-id",
					"pegasearch.srsAuth.scopes":              "srs-scope",
					"pegasearch.srsAuth.privateKey":          SRSAuthPrivateKeyExample,
					"pegasearch.srsAuth.privateKeyAlgorithm": "RS512",
				},
			}

			yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-environment-config.yaml"})
			VerifyPegaWithExternalSRSEnvironmentConfig(t, yamlContent, true, "RS512", "srs-scope")
		}
	}
}

func TestPegaConfigWithSRSAndAuthEnabledAndAlgorithmNotProvided(t *testing.T) {
	var supportedVendors = []string{"k8s", "openshift", "eks", "gke", "aks", "pks"}
	var supportedOperations = []string{"deploy", "install-deploy"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for _, vendor := range supportedVendors {
		for _, operation := range supportedOperations {
			fmt.Println(vendor + "-" + operation)

			var options = &helm.Options{
				SetValues: map[string]string{
					"global.provider":                  vendor,
					"global.actions.execute":           operation,
					"pegasearch.externalSearchService": "true",
					"pegasearch.externalURL":           "https://srs:9200",
					"pegasearch.srsAuth.enabled":       "true",
					"pegasearch.srsAuth.url":           "https://auth-service",
					"pegasearch.srsAuth.clientId":      "client-id",
					"pegasearch.srsAuth.scopes":        "srs-scope",
					"pegasearch.srsAuth.privateKey":    SRSAuthPrivateKeyExample,
				},
			}

			yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-environment-config.yaml"})
			VerifyPegaWithExternalSRSEnvironmentConfig(t, yamlContent, true, DefaultPrivateKeyAlgorithm, "srs-scope")
		}
	}
}

func TestPegaConfigWithSRSAndAuthEnabledAndScopeNotProvided(t *testing.T) {
	var supportedVendors = []string{"k8s", "openshift", "eks", "gke", "aks", "pks"}
	var supportedOperations = []string{"deploy", "install-deploy"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for _, vendor := range supportedVendors {
		for _, operation := range supportedOperations {
			fmt.Println(vendor + "-" + operation)

			var options = &helm.Options{
				SetValues: map[string]string{
					"global.provider":                        vendor,
					"global.actions.execute":                 operation,
					"pegasearch.externalSearchService":       "true",
					"pegasearch.externalURL":                 "https://srs:9200",
					"pegasearch.srsAuth.enabled":             "true",
					"pegasearch.srsAuth.url":                 "https://auth-service",
					"pegasearch.srsAuth.clientId":            "client-id",
					"pegasearch.srsAuth.privateKeyAlgorithm": "RS384",
					"pegasearch.srsAuth.privateKey":          SRSAuthPrivateKeyExample,
				},
			}

			yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-environment-config.yaml"})
			VerifyPegaWithExternalSRSEnvironmentConfig(t, yamlContent, true, "RS384", DefaultSRSAuthScope)
		}
	}
}

func VerifyPegaWithExternalSRSEnvironmentConfig(t *testing.T, yamlContent string, isAuthEnabled bool, expectedAlgorithm string, expectedScope string) {
	var envConfigMap k8score.ConfigMap
	statefulSlice := strings.Split(yamlContent, "---")
	for index, statefulInfo := range statefulSlice {
		if index >= 1 {
			UnmarshalK8SYaml(t, statefulInfo, &envConfigMap)
			envConfigData := envConfigMap.Data
			require.Equal(t, "https://srs:9200", envConfigData["SEARCH_AND_REPORTING_SERVICE_URL"])
			if isAuthEnabled {
				VerifyEnvConfigDataWithAuthVariables(t, envConfigData, expectedAlgorithm, expectedScope)
			} else {
				VerifyEnvConfigDataWithoutAuthVariables(t, envConfigData)
			}
		}
	}
}

func VerifyPegaWithoutExternalSRSEnvironmentConfig(t *testing.T, yamlContent string) {
	var envConfigMap k8score.ConfigMap
	statefulSlice := strings.Split(yamlContent, "---")
	for index, statefulInfo := range statefulSlice {
		if index >= 1 {
			UnmarshalK8SYaml(t, statefulInfo, &envConfigMap)
			envConfigData := envConfigMap.Data
			require.Empty(t, envConfigData["SEARCH_AND_REPORTING_SERVICE_URL"])
			VerifyEnvConfigDataWithoutAuthVariables(t, envConfigData)
		}
	}
}

func VerifyEnvConfigDataWithoutAuthVariables(t *testing.T, envConfigData map[string]string) {
	authEnvironmentVariables := []string{"SERV_AUTH_URL", "SERV_AUTH_CLIENT_ID", "SERV_AUTH_SCOPES", "SERV_AUTH_PRIVATE_KEY_ALGORITHM", "SERV_AUTH_PRIVAYE_KEY", "SERV_AUTH_CLIENT_SECRET"}
	for _, authEnvironmentVariable := range authEnvironmentVariables {
		require.Emptyf(t, envConfigData[authEnvironmentVariable], "Environment variable '%s' should be empty", authEnvironmentVariable)
	}
}

func VerifyEnvConfigDataWithAuthVariables(t *testing.T, envConfigData map[string]string, expectedAlgorithm string, expectedScope string) {
	require.Equal(t, "https://auth-service", envConfigData["SERV_AUTH_URL"])
	require.Equal(t, "client-id", envConfigData["SERV_AUTH_CLIENT_ID"])
	require.Equal(t, expectedScope, envConfigData["SERV_AUTH_SCOPES"])
	require.Equal(t, expectedAlgorithm, envConfigData["SERV_AUTH_PRIVATE_KEY_ALGORITHM"])
	_, hasPrivateKey := envConfigData["SERV_AUTH_PRIVATE_KEY"]
	require.False(t, hasPrivateKey)
	_, hasClientPrivateKey := envConfigData["SERV_AUTH_CLIENT_SECRET"]
	require.False(t, hasClientPrivateKey)
}
