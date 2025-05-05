package pega

import (
	"fmt"
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	k8score "k8s.io/api/core/v1"
	"path/filepath"
	"testing"
)

func TestPegaUpgradeRESTSecret(t *testing.T) {
	var supportedVendors = []string{"k8s", "openshift", "eks", "gke", "aks", "pks"}
	var supportedOperations = []string{"install", "deploy", "install-deploy", "upgrade-deploy"}
	var deploymentNames = []string{"pega", "myapp-dev"}

    var expectedValues = map[string]string{"PEGA_REST_SERVER_URL": "http://pega-web:80/prweb/PRRestService", "PEGA_REST_USERNAME": "username", "PEGA_REST_PASSWORD": "password"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for _, vendor := range supportedVendors {

		for _, operation := range supportedOperations {

			for _, depName := range deploymentNames {
				fmt.Println(vendor + "-" + operation)

				var options = &helm.Options{
					SetValues: map[string]string{
						"global.deployment.name":             depName,
                        "global.provider":                    vendor,
                        "global.actions.execute":             operation,
                        "installer.upgrade.upgradeType":      "zero-downtime",
                        "installer.upgrade.pegaRESTUsername": expectedValues["PEGA_REST_USERNAME"],
                        "installer.upgrade.pegaRESTPassword": expectedValues["PEGA_REST_PASSWORD"],
                        "installer.upgrade.pega_rest_external_secret_name": "",
					},
				}

				yamlContent := RenderTemplate(t, options, helmChartPath, []string{"charts/installer/templates/pega-upgrade-rest-secret.yaml"})
				VerifyRESTSecret(t, yamlContent, options, expectedValues)
			}
		}
	}
}

func VerifyRESTSecret(t *testing.T, yamlContent string, options *helm.Options, expectedValues map[string]string) {
	var secretobj k8score.Secret
	UnmarshalK8SYaml(t, yamlContent, &secretobj)

	require.Equal(t, getObjName(options, "-upgrade-rest-secret"), secretobj.ObjectMeta.Name)

	secretData := secretobj.Data
	require.Equal(t, expectedValues["PEGA_REST_USERNAME"], string(secretData["PEGA_REST_USERNAME"]))
	require.Equal(t, expectedValues["PEGA_REST_PASSWORD"], string(secretData["PEGA_REST_PASSWORD"]))
}
