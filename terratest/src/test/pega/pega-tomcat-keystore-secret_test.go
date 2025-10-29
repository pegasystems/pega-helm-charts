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
						"global.deployment.name":             depName,
						"global.provider":                    vendor,
						"global.actions.execute":             operation,
						"global.tier[0].service.tls.enabled": "true",
						"installer.upgrade.upgradeType":      "zero-downtime",
					},
				}
				yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tomcat-keystore-secret.yaml"})
				VerifyTomcatCertificatesSecret(t, yamlContent, options)
			}
		}
	}
}

func VerifyTomcatCertificatesSecret(t *testing.T, yamlContent string, options *helm.Options) {

	var importCertSecret k8score.Secret
	UnmarshalK8SYaml(t, yamlContent, &importCertSecret)

	importCertSecretData := importCertSecret.Data
	require.Equal(t, importCertSecret.ObjectMeta.Name, getObjName(options, "-tomcat-keystore-secret"))
	require.Equal(t, "123456", string(importCertSecretData["TOMCAT_KEYSTORE_PASSWORD"]))
}
