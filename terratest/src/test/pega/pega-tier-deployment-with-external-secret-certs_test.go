package pega

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	"path/filepath"
	"strings"
	"testing"
)

func TestPegaDeploymentWithExternalSecretCerts(t *testing.T) {

	var supportedVendors = []string{"k8s"}
	var supportedOperations = []string{"deploy", "install-deploy", "upgrade-deploy"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for _, vendor := range supportedVendors {
		for _, operation := range supportedOperations {

			var options = &helm.Options{
				ValuesFiles: []string{"data/values_with_externalcerts.yaml"},
				SetValues: map[string]string{
					"global.deployment.name":        "pega",
					"global.provider":               vendor,
					"global.actions.execute":        operation,
					"installer.upgrade.upgradeType": "zero-downtime",
				},
			}
			deploymentYaml := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-deployment.yaml"})
			yamlSplit := strings.Split(deploymentYaml, "---")
			assertWeb(t, yamlSplit[1], options)
			assertVolumeAndMount(t, yamlSplit[1], options, true)
			assertMountedCertificateSecrets(t, yamlSplit[1], options)
		}
	}
}
func assertMountedCertificateSecrets(t *testing.T, tierYaml string, options *helm.Options) {
	var deploymentObj appsv1.Deployment
	UnmarshalK8SYaml(t, tierYaml, &deploymentObj)
	pod := deploymentObj.Spec.Template.Spec
	var foundCert = false
	for _, vol := range pod.Volumes {
		if vol.Name == "pega-volume-import-certificates" {
			sources := vol.VolumeSource.Projected.Sources
			require.Equal(t, len(sources), 2)
			for _, source := range sources {
				if source.Secret.LocalObjectReference.Name == "pega-import-certificates-secret" {
					foundCert = true
					break
				}
			}
			break
		}
	}
	require.Equal(t, foundCert, false)
}
