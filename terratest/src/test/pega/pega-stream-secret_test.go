package pega

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	k8score "k8s.io/api/core/v1"
)

const streamTrustStorePassword = "trustStore"
const streamKeyStorePassword = "keyStore"
const jaasConfig = "jaasConfig"

func TestPegaCredentialsSecretWithExternalStreamArePresent(t *testing.T) {
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
						"global.deployment.name":        depName,
						"global.provider":               vendor,
						"global.actions.execute":        operation,
						"installer.upgrade.upgradeType": getUpgradeTypeForUpgradeAction(operation),
						"stream.trustStorePassword":     streamTrustStorePassword,
						"stream.keyStorePassword":       streamKeyStorePassword,
						"stream.jaasConfig":             jaasConfig,
					},
				}

				yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-stream-secret.yaml"})
				verifyStreamCredentialsSecret(t, yamlContent, operation)
			}
		}
	}

}

func verifyStreamCredentialsSecret(t *testing.T, yamlContent string, operation string) {

	var secretobj k8score.Secret
	UnmarshalK8SYaml(t, yamlContent, &secretobj)
	secretData := secretobj.Data
	require.Equal(t, streamTrustStorePassword, string(secretData["STREAM_TRUSTSTORE_PASSWORD"]))
	require.Equal(t, streamKeyStorePassword, string(secretData["STREAM_KEYSTORE_PASSWORD"]))
	require.Equal(t, jaasConfig, string(secretData["STREAM_JAAS_CONFIG"]))
}
