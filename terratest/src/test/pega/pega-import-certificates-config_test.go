package pega

import (
	"fmt"
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	k8score "k8s.io/api/core/v1"
	"path/filepath"
	"testing"
)

func TestPegaImportCertificatesConfig(t *testing.T) {
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

				yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-certificates-config.yaml"})
				VerifyImportCertificatesConfig(t, yamlContent, options)
			}
		}
	}
}

func VerifyImportCertificatesConfig(t *testing.T, yamlContent string, options *helm.Options) {

	var importCertConfigMap k8score.ConfigMap
	UnmarshalK8SYaml(t, yamlContent, &importCertConfigMap)

	importCertConfigData := importCertConfigMap.Data
	if len(importCertConfigData) != 0 {
		require.Equal(t, importCertConfigMap.ObjectMeta.Name, getObjName(options, "-import-certificates-config"))
		require.Equal(t, importCertConfigData["dbcacertificate1.cer"], "-----BEGIN CERTIFICATE-----\nMIIDeTCCAmGgAwIBAgIJAMnA8BB8xT6wMA0GCSqGSIb3DQEBCwUAMGIxCzAJBgNV\nBAYTAlVTMRMwEQYDVQQIDApDYWxpZm9ybmlhMRYwFAYDVQQHDA1TYW4gRnJhbmNp\nc2NvMQ8wDQYDVQQKDAZCYWRTU0wxFTATBgNVBAMMDCouYmFkc3NsLmNvbTAeFw0y\nMTEwMTEyMDAzNTRaFw0yMzEwMTEyMDAzNTRaMGIxCzAJBgNVBAYTAlVTMRMwEQYD\nVQQIDApDYWxpZm9ybmlhMRYwFAYDVQQHDA1TYW4gRnJhbmNpc2NvMQ8wDQYDVQQK\nDAZCYWRTU0wxFTATBgNVBAMMDCouYmFkc3NsLmNvbTCCASIwDQYJKoZIhvcNAQEB\nBQADggEPADCCAQoCggEBAMIE7PiM7gTCs9hQ1XBYzJMY61yoaEmwIrX5lZ6xKyx2\nPmzAS2BMTOqytMAPgLaw+XLJhgL5XEFdEyt/ccRLvOmULlA3pmccYYz2QULFRtMW\nhyefdOsKnRFSJiFzbIRMeVXk0WvoBj1IFVKtsyjbqv9u/2CVSndrOfEk0TG23U3A\nxPxTuW1CrbV8/q71FdIzSOciccfCFHpsKOo3St/qbLVytH5aohbcabFXRNsKEqve\nww9HdFxBIuGa+RuT5q0iBikusbpJHAwnnqP7i/dAcgCskgjZjFeEU4EFy+b+a1SY\nQCeFxxC7c3DvaRhBB0VVfPlkPz0sw6l865MaTIbRyoUCAwEAAaMyMDAwCQYDVR0T\nBAIwADAjBgNVHREEHDAaggwqLmJhZHNzbC5jb22CCmJhZHNzbC5jb20wDQYJKoZI\nhvcNAQELBQADggEBAC4DensZ5tCTeCNJbHABYPwwqLUFOMITKOOgF3t8EqOan0CH\nST1NNi4jPslWrVhQ4Y3UbAhRBdqXl5N/NFfMzDosPpOjFgtifh8Z2s3w8vdlEZzf\nA4mYTC8APgdpWyNgMsp8cdXQF7QOfdnqOfdnY+pfc8a8joObR7HEaeVxhJs+XL4E\nCLByw5FR+svkYgCbQGWIgrM1cRpmXemt6Gf/XgFNP2PdubxqDEcnWlTMk8FCBVb1\nnVDSiPjYShwnWsOOshshCRCAiIBPCKPX0QwKDComQlRrgMIvddaSzFFTKPoNZjC+\nCUspSNnL7V9IIHvqKlRSmu+zIpm2VJCp1xLulk8=\n-----END CERTIFICATE-----\n")
	} else {
		require.Nil(t, importCertConfigMap.Data)
	}

}
