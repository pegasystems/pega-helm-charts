package pega

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	k8score "k8s.io/api/core/v1"
)

const trustStorePassword = "trustStore"
const keyStorePassword = "keyStore"

func TestPegaDDSSecretWithEncryptionPresent(t *testing.T) {
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
						"dds.externalNodes":             "123.45.60.00",
						"dds.trustStorePassword":        trustStorePassword,
						"dds.keyStorePassword":          keyStorePassword,
					},
				}

				yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-dds-secret.yaml"})
				verifyDDSSecret(t, yamlContent)
			}
		}
	}

}

func verifyDDSSecret(t *testing.T, yamlContent string) {

	var secretobj k8score.Secret
	UnmarshalK8SYaml(t, yamlContent, &secretobj)
	secretData := secretobj.Data
	require.Equal(t, trustStorePassword, string(secretData["CASSANDRA_TRUSTSTORE_PASSWORD"]))
	require.Equal(t, keyStorePassword, string(secretData["CASSANDRA_KEYSTORE_PASSWORD"]))
}
