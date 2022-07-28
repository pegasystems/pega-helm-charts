package pega

import (
	"fmt"
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	k8score "k8s.io/api/core/v1"
	"path/filepath"
	"testing"
)

func TestPegaTomcatKeystoreSecret(t *testing.T) {
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
					SetValues: map[string]string{
						"global.deployment.name":        depName,
						"global.provider":               vendor,
						"global.actions.execute":        operation,
						"installer.upgrade.upgradeType": "zero-downtime",
					},
				}

				yamlContent, error := RenderTemplateE(t, options, helmChartPath, []string{"templates/pega-tomcat-keystore-secret.yaml"})
				if error != nil {
					VerifyTomcatCertificatesSecret(t, yamlContent, options)
				}

			}
		}
	}
}

func VerifyTomcatCertificatesSecret(t *testing.T, yamlContent string, options *helm.Options) {

	var importCertSecret k8score.Secret
	UnmarshalK8SYaml(t, yamlContent, &importCertSecret)

	importCertSecretData := importCertSecret.StringData
	if len(importCertSecretData) != 0 {
		require.Equal(t, importCertSecret.ObjectMeta.Name, getObjName(options, "-tomcat-keystore-secret"))
		require.Equal(t, importCertSecretData["CERT_PASSWORD"], "123456")
	} else {
		require.Nil(t, importCertSecret.Data)
	}

}
