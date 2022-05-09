package pega

import (
	"fmt"
	"path/filepath"
	"testing"
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	k8score "k8s.io/api/core/v1"
)

func TestCustomArtifactoryCertificatesConfig(t *testing.T) {
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
					//ValuesFiles: []string{"data/values_with_artifactory_cert.yaml"},
					SetValues: map[string]string{
						"global.deployment.name":        depName,
						"global.provider":               vendor,
						"global.actions.execute":        operation,
						"installer.upgrade.upgradeType": "zero-downtime",
					},
				}

				yamlContent, _ := RenderTemplateWithErr(t, options, helmChartPath, []string{"templates/pega-custom-artifactory-certificates-config.yaml"})
				VerifyArtifactoryCertificatesConfig(t, yamlContent, options)
			}
		}
	}
}

func TestCustomArtifactoryCertificatesConfigWhenSSLVerificationIsDisabled(t *testing.T) {
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
						"global.customArtifactory.enableSSLVerification": "false",
						"global.customArtifactory.certificate": "self-signed-certificate.cer: |\n\"-----BEGIN CERTIFICATE-----\\nMIIDdTCCAl2gAwIBAgIENdb1mTANBgkqhkiG9w0BAQsFADBrMQswCQYDVQQGEwJJ\nTjESMBAGA1UECBMJVGVsYW5nYW5hMRIwEAYDVQQHEwlIeWRlcmFiYWQxDTALBgNV\nBAoTBFBlZ2ExDTALBgNVBAsTBFBERFMxFjAUBgNVBAMTDTEwLjIyNS43MS4xNDMw\nHhcNMjIwMzIyMTYwNzQxWhcNMjMwMzE3MTYwNzQxWjBrMQswCQYDVQQGEwJJTjES\nMBAGA1UECBMJVGVsYW5nYW5hMRIwEAYDVQQHEwlIeWRlcmFiYWQxDTALBgNVBAoT\nBFBlZ2ExDTALBgNVBAsTBFBERFMxFjAUBgNVBAMTDTEwLjIyNS43MS4xNDMwggEi\nMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCNYvEPbKRZ1y3u5fmPcvdBaQVK\nXsy3ioY7ToMP6vpNXGmRhI06t6jVzDVU85jtdSb4On2B4uyZwUSxO9cWfOFtI6wW\nnrdhRxmygFvXwinon6LXcoRfK9TjI72C/694UWu/UysaUp8yWyrHf2XfQK2qFqMC\nej57bpRSCME2maAXKC88IGNpeX2XjhICUKPPBrWDpK4Jq7NhcgV6z0hmtCWg8+lj\nmMZRXDoJKMNNhWOKMYhL/djDJZH+PMON7sOGVtE52U0UBLhff6Ee4ERBRummNgCv\nxt4MTmzgewsN1uKQ5MiBjtJduVEmKhhiIV38QetrCPpejAHOJLFe2l5VfKHDAgMB\nAAGjITAfMB0GA1UdDgQWBBTK0eVfaa41Vr4qXTww3RBTgFO78DANBgkqhkiG9w0B\nAQsFAAOCAQEAOAjezNJmMx9j0hnutOspnHC8iOqaFQjW8t6D9cWEQALd2PNPB5S9\nQxlEuaN3x/zbtNI55fxZW6ryP/AJ0DclTs8vwzEk7DJ1Yt7vMfFG6DxbIUlPY677\nDGB23K68BXl8MtSYvOLbDwXYjyMDUzcmojaIjS6RwW8C5yvXW34h2jjwVWQm1yti\n46xANKLHEVTp44LiG+gf/9TxfQjSQXpSdgdMbJB744tMmozyfbtulWE0T5dBvd8w\ncdbPKbgldsv4bc8EojOYRRasYu6nZqP+8Tw/4jHr4IB2kiuJ63gs6IlqnzDyzzQ7\nYc0a+hYe1cTSXQn23aL/c9v/901LUpdAYw==\\n-----END CERTIFICATE-----\\n\"",
					},
				}

				_, err := RenderTemplateWithErr(t, options, helmChartPath, []string{"templates/pega-custom-artifactory-certificates-config.yaml"})
				require.Contains(t, err.Error(),
					"could not find template templates/pega-custom-artifactory-certificates-config.yaml in chart")
			}
		}
	}
}

func VerifyArtifactoryCertificatesConfig(t *testing.T, yamlContent string, options *helm.Options) {
	var artifactoryCertConfigMap k8score.ConfigMap
	UnmarshalK8SYaml(t, yamlContent, &artifactoryCertConfigMap)

	artifactoryCertConfigData := artifactoryCertConfigMap.Data
	if len(artifactoryCertConfigData) != 0 {
		require.Equal(t, artifactoryCertConfigMap.ObjectMeta.Name, getObjName(options, "-custom-artifactory-certificate-config"))
		require.Equal(t, artifactoryCertConfigData["self-signed-certificate.cer"], "-----BEGIN CERTIFICATE-----\ncertificate content\n-----END CERTIFICATE-----")
	} else {
		require.Nil(t, artifactoryCertConfigMap.Data)
	}

}
