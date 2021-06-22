package pega

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	k8score "k8s.io/api/core/v1"
)

const trustStorePassword = "trustStore"
const keyStorePassword = "keyStore"

func TestPegaCredentialsSecretWithCassandraEncryptionPresent(t *testing.T) {
	var supportedVendors = []string{"k8s", "openshift", "eks", "gke", "aks", "pks"}
	var supportedOperations = []string{"install", "install-deploy", "upgrade", "upgrade-deploy"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for _, vendor := range supportedVendors {

		for _, operation := range supportedOperations {

			fmt.Println(vendor + "-" + operation)

			var options = &helm.Options{
				SetValues: map[string]string{
					"global.provider":        vendor,
					"global.actions.execute": operation,
					"installer.upgrade.upgradeType": "zero-downtime",
					"dds.trustStorePassword": trustStorePassword,
					"dds.keyStorePassword":   keyStorePassword,
				},
			}

			yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-credentials-secret.yaml"})
			verifyCredentialsSecret(t, yamlContent, operation)
		}
	}

}

func verifyCredentialsSecret(t *testing.T, yamlContent string, operation string) {

	var secretobj k8score.Secret
	UnmarshalK8SYaml(t, yamlContent, &secretobj)
	secretData := secretobj.Data
	require.Equal(t, "YOUR_JDBC_USERNAME", string(secretData["DB_USERNAME"]))
	require.Equal(t, "YOUR_JDBC_PASSWORD", string(secretData["DB_PASSWORD"]))
	if strings.Contains(operation, "deploy") {
		require.Equal(t, trustStorePassword, string(secretData["CASSANDRA_TRUSTSTORE_PASSWORD"]))
		require.Equal(t, keyStorePassword, string(secretData["CASSANDRA_KEYSTORE_PASSWORD"]))
	} else {
		require.Equal(t, "", string(secretData["CASSANDRA_TRUSTSTORE_PASSWORD"]))
		require.Equal(t, "", string(secretData["CASSANDRA_KEYSTORE_PASSWORD"]))
	}

}
