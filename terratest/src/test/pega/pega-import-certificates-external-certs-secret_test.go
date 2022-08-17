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

func TestPegaImportCertificatesESSecret(t *testing.T) {
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
					ValuesFiles: []string{"data/values_with_externalcerts.yaml"},
					SetValues: map[string]string{
						"global.deployment.name":        depName,
						"global.provider":               vendor,
						"global.actions.execute":        operation,
						"installer.upgrade.upgradeType": "zero-downtime",
					},
				}

				yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-certificates-secret.yaml"})
				yamlSplit := strings.Split(yamlContent, "---")
				//passing split value here to suppress warnings generated
				VerifyImportCertificatesESSecret(t, yamlSplit[1], options)
			}
		}
	}
}

func VerifyImportCertificatesESSecret(t *testing.T, yamlContent string, options *helm.Options) {

	var importCertSecret k8score.Secret
	UnmarshalK8SYaml(t, yamlContent, &importCertSecret)
	require.Empty(t, importCertSecret.Name)
	require.Nil(t, importCertSecret.StringData)
}
