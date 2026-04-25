package pega

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	k8score "k8s.io/api/core/v1"
)

const srsMTLSTrustStorePassword = "trustStore"
const srsMTLSKeyStorePassword = "keyStore"

func TestPegaCredentialsSecretWithExternalSRSAndmTLSArePresent(t *testing.T) {
	var supportedVendors = []string{"k8s", "openshift", "eks", "gke", "aks", "pks"}
	var supportedOperations = []string{"install-deploy", "deploy", "upgrade-deploy"}
	var deploymentNames = []string{"pega", "myapp-dev"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for _, vendor := range supportedVendors {

		for _, operation := range supportedOperations {

			for _, depName := range deploymentNames {
				fmt.Println(vendor + "-" + operation)

				var options = &helm.Options{
					SetValues: map[string]string{
						"global.deployment.name":                depName,
						"global.provider":                       vendor,
						"global.actions.execute":                operation,
						"installer.upgrade.upgradeType":         getUpgradeTypeForUpgradeAction(operation),
						"pegasearch.externalSearchService":      "true",
						"pegasearch.srsMTLS.enabled":            "true",
						"pegasearch.srsMTLS.trustStorePassword": srsMTLSTrustStorePassword,
						"pegasearch.srsMTLS.keyStorePassword":   srsMTLSKeyStorePassword,
					},
				}

				yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-srs-mtls-secret.yaml"})
				verifySRSmTLSCredentialsSecret(t, yamlContent, operation)
			}
		}
	}

}

func verifySRSmTLSCredentialsSecret(t *testing.T, yamlContent string, operation string) {

	var secretobj k8score.Secret
	UnmarshalK8SYaml(t, yamlContent, &secretobj)
	secretData := secretobj.Data
	require.Equal(t, srsMTLSTrustStorePassword, string(secretData["SRS_TRUSTSTORE_PASSWORD"]))
	require.Equal(t, srsMTLSKeyStorePassword, string(secretData["SRS_KEYSTORE_PASSWORD"]))
}

func TestPegaSRSmTLSSecretNotCreatedWhenVaultEnabled(t *testing.T) {
	var supportedVendors = []string{"k8s", "openshift", "eks", "gke", "aks", "pks"}
	var supportedOperations = []string{"install-deploy", "deploy", "upgrade-deploy"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for _, vendor := range supportedVendors {
		for _, operation := range supportedOperations {
			fmt.Println(vendor + "-" + operation)

			var options = &helm.Options{
				SetValues: map[string]string{
					"global.provider":                      vendor,
					"global.actions.execute":               operation,
					"installer.upgrade.upgradeType":        getUpgradeTypeForUpgradeAction(operation),
					"pegasearch.externalSearchService":     "true",
					"pegasearch.srsMTLS.enabled":           "true",
					"pegasearch.srsMTLS.vault.enabled":     "true",
					"pegasearch.srsMTLS.vault.url":         "https://vault.example.com",
					"pegasearch.srsMTLS.vault.role":        "pega-role",
					"pegasearch.srsMTLS.vault.secretPath":  "secret/data/pega/mtls",
					"pegasearch.srsMTLS.vault.tokenSecret": "vault-token-secret",
				},
			}

			_, err := helm.RenderTemplateE(t, options, helmChartPath, "pega", []string{"templates/pega-srs-mtls-secret.yaml"})
			require.Error(t, err, "SRS mTLS secret should not be rendered when vault is enabled")
		}
	}
}
