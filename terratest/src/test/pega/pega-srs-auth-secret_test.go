package pega

import (
	"fmt"
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	k8score "k8s.io/api/core/v1"
	"path/filepath"
	"testing"
)

func TestPegaSRSAuthSecretNotCreatedForDeploymentWithoutSRS(t *testing.T) {
	var supportedVendors = []string{"k8s", "openshift", "eks", "gke", "aks", "pks"}
	var supportedOperations = []string{"install", "install-deploy"}
	var deploymentNames = []string{"pega", "myapp-dev"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for _, vendor := range supportedVendors {
		for _, operation := range supportedOperations {
			for _, depName := range deploymentNames {
				fmt.Println(vendor + "-" + operation)

				var options = &helm.Options{
					SetValues: map[string]string{
						"global.deployment.name":           depName,
						"global.provider":                  vendor,
						"pegasearch.externalSearchService": "false",
					},
				}

				yamlContent, err := RenderTemplateE(t, options, helmChartPath, []string{"templates/pega-srs-auth-secret.yaml"})
				VerifySRSAuthSecretIsNotCreated(t, yamlContent, err)
			}
		}
	}
}

func TestPegaSRSAuthSecretNotCreatedForDeploymentWithDisabledSRSAuth(t *testing.T) {
	var supportedVendors = []string{"k8s", "openshift", "eks", "gke", "aks", "pks"}
	var supportedOperations = []string{"install", "install-deploy"}
	var deploymentNames = []string{"pega", "myapp-dev"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for _, vendor := range supportedVendors {
		for _, operation := range supportedOperations {
			for _, depName := range deploymentNames {
				fmt.Println(vendor + "-" + operation)

				var options = &helm.Options{
					SetValues: map[string]string{
						"global.deployment.name":           depName,
						"global.provider":                  vendor,
						"pegasearch.externalSearchService": "true",
						"pegasearch.srsAuth.enabled":       "false",
					},
				}

				yamlContent, err := RenderTemplateE(t, options, helmChartPath, []string{"templates/pega-srs-auth-secret.yaml"})
				VerifySRSAuthSecretIsNotCreated(t, yamlContent, err)
			}
		}
	}
}

func TestPegaSRSAuthSecretNotCreatedForMissingPrivateKey(t *testing.T) {
	var supportedVendors = []string{"k8s", "openshift", "eks", "gke", "aks", "pks"}
	var supportedOperations = []string{"install", "install-deploy"}
	var deploymentNames = []string{"pega", "myapp-dev"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for _, vendor := range supportedVendors {
		for _, operation := range supportedOperations {
			for _, depName := range deploymentNames {
				fmt.Println(vendor + "-" + operation)

				var options = &helm.Options{
					SetValues: map[string]string{
						"global.deployment.name":           depName,
						"global.provider":                  vendor,
						"pegasearch.externalSearchService": "true",
						"pegasearch.srsAuth.enabled":       "true",
					},
				}

				yamlContent, err := RenderTemplateE(t, options, helmChartPath, []string{"templates/pega-srs-auth-secret.yaml"})
				require.Contains(t, yamlContent, "A valid entry is required for pegasearch.srsAuth.privateKey or pegasearch.srsAuth.external_secret_name, when request authentication mechanism(IDP) is enabled between SRS and Pega Infinity i.e. pegasearch.srsAuth.enabled is true.")
				require.Contains(t, err.Error(), "A valid entry is required for pegasearch.srsAuth.privateKey or pegasearch.srsAuth.external_secret_name, when request authentication mechanism(IDP) is enabled between SRS and Pega Infinity i.e. pegasearch.srsAuth.enabled is true.")
			}
		}
	}
}

func TestPegaSRSAuthSecretNotCreatedForDeploymentWithEnabledSRSAuthAndExternalSecret(t *testing.T) {
	var supportedVendors = []string{"k8s", "openshift", "eks", "gke", "aks", "pks"}
	var supportedOperations = []string{"install", "install-deploy"}
	var deploymentNames = []string{"pega", "myapp-dev"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for _, vendor := range supportedVendors {
		for _, operation := range supportedOperations {
			for _, depName := range deploymentNames {
				fmt.Println(vendor + "-" + operation + "-" + depName)

				var options = &helm.Options{
					SetValues: map[string]string{
						"global.deployment.name":                     depName,
						"global.provider":                            vendor,
						"pegasearch.externalSearchService":           "true",
						"pegasearch.srsAuth.enabled":                 "true",
						"pegasearch.srsAuth.external_secret_name":    "test-external-secret",
					},
				}

				yamlContent, err := RenderTemplateE(t, options, helmChartPath, []string{"templates/pega-srs-auth-secret.yaml"})
				VerifySRSAuthSecretIsNotCreated(t, yamlContent, err)
			}
		}
	}
}

func TestPegaSRSAuthSecretCreatedForDeploymentWithEnabledSRSAuth(t *testing.T) {
	var supportedVendors = []string{"k8s", "openshift", "eks", "gke", "aks", "pks"}
	var supportedOperations = []string{"install", "install-deploy"}
	var deploymentNames = []string{"pega", "myapp-dev"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for _, vendor := range supportedVendors {
		for _, operation := range supportedOperations {
			for _, depName := range deploymentNames {
				fmt.Println(vendor + "-" + operation)

				var options = &helm.Options{
					SetValues: map[string]string{
						"global.deployment.name":           depName,
						"global.provider":                  vendor,
						"pegasearch.externalSearchService": "true",
						"pegasearch.srsAuth.enabled":       "true",
						"pegasearch.srsAuth.privateKey":    SRSAuthPrivateKeyExample,
					},
				}

				yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-srs-auth-secret.yaml"})
				VerifySRSAuthSecretIsCreated(t, yamlContent)
			}
		}
	}
}

func VerifySRSAuthSecretIsNotCreated(t *testing.T, yamlContent string, err error) {
	require.Contains(t, yamlContent, "could not find template templates/pega-srs-auth-secret.yaml")
	require.Contains(t, err.Error(), "could not find template templates/pega-srs-auth-secret.yaml")
}

func VerifySRSAuthSecretIsCreated(t *testing.T, yamlContent string) {
	var secret k8score.Secret
	UnmarshalK8SYaml(t, yamlContent, &secret)

	require.Equal(t, "pega-srs-auth-secret", secret.ObjectMeta.Name)
	secretData := secret.Data
	require.Len(t, secretData, 1)
	require.Equal(t, SRSAuthPrivateKeyExample, string(secretData["privateKey"]))
}
