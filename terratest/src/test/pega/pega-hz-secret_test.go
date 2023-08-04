package pega

import (
	"fmt"
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	k8score "k8s.io/api/core/v1"
	"path/filepath"
	"testing"
)

func TestPegaHazelcastSecretWhenHazelcastIsEnabled(t *testing.T) {
	var supportedVendors = []string{"k8s", "openshift", "eks", "gke", "aks", "pks"}
	const HzCsAuthUsername = "HZClusterUser"
	const HzCsAuthPassword = "HZclusterPassword"
	var supportedOperations = []string{"deploy", "install-deploy", "upgrade-deploy"}
	var deploymentNames = []string{"pega", "myapp-dev"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for _, vendor := range supportedVendors {

		for _, operation := range supportedOperations {

			for _, depName := range deploymentNames {
				fmt.Println(vendor + "-" + operation)

				var options = &helm.Options{
					SetValues: map[string]string{
						"global.deployment.name":                    depName,
						"global.provider":                           vendor,
						"global.actions.execute":                    operation,
						"installer.upgrade.upgradeType":             "zero-downtime",
						"hazelcast.enabled":                         "true",
						"hazelcast.migration.embeddedToCSMigration": "false",
						"hazelcast.username":                        HzCsAuthUsername,
						"hazelcast.password":                        HzCsAuthPassword,
					},
				}

				yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-hz-secret.yaml"})
				VerifyHazelcastSecretWhenHazelcastIsEnabled(t, yamlContent, HzCsAuthUsername, HzCsAuthPassword)
			}
		}
	}

	for _, vendor := range supportedVendors {

		for _, operation := range supportedOperations {

			fmt.Println(vendor + "-" + operation)

			var options = &helm.Options{
				SetValues: map[string]string{
					"global.provider":                           vendor,
					"global.actions.execute":                    operation,
					"installer.upgrade.upgradeType":             "zero-downtime",
					"hazelcast.enabled":                         "false",
					"hazelcast.clusteringServiceEnabled":        "true",
					"hazelcast.migration.embeddedToCSMigration": "false",
					"hazelcast.username":                        HzCsAuthUsername,
					"hazelcast.password":                        HzCsAuthPassword,
				},
			}

			yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-hz-secret.yaml"})
			VerifyHazelcastSecretWhenHazelcastIsEnabled(t, yamlContent, HzCsAuthUsername, HzCsAuthPassword)
		}
	}

}

func VerifyHazelcastSecretWhenHazelcastIsEnabled(t *testing.T, yamlContent string, username string, password string) {

	var secretobj k8score.Secret
	UnmarshalK8SYaml(t, yamlContent, &secretobj)
	secretData := secretobj.Data
	require.Equal(t, string(secretData["HZ_CS_AUTH_USERNAME"]), username)
	require.Equal(t, string(secretData["HZ_CS_AUTH_PASSWORD"]), password)
}
