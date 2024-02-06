package pega

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
)

func TestPegaDeploymentWithAndWithoutCustomCerts(t *testing.T) {

	var supportedVendors = []string{"k8s"}
	var supportedOperations = []string{"deploy", "install-deploy", "upgrade-deploy"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for _, vendor := range supportedVendors {
		for _, operation := range supportedOperations {

			var options = &helm.Options{
				ValuesFiles: []string{"data/values_with_customcerts.yaml"},
				SetValues: map[string]string{
					"global.deployment.name":        "pega",
					"global.provider":               vendor,
					"global.actions.execute":        operation,
					"installer.upgrade.upgradeType": "zero-downtime",
				},
			}
			renderTierForCertTest(t, options, helmChartPath, true)

			options.ValuesFiles = []string{"data/values_without_customcerts.yaml"}

			renderTierForCertTest(t, options, helmChartPath, false)
		}
	}
}

func renderTierForCertTest(t *testing.T, options *helm.Options, helmChartPath string, shouldHaveVol bool) {

	deploymentYaml := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-deployment.yaml"})

	yamlSplit := strings.Split(deploymentYaml, "---")
	assertWeb(t, yamlSplit[1], options)
	assertVolumeAndMount(t, yamlSplit[1], options, shouldHaveVol)

	assertBatch(t, yamlSplit[2], options)
	assertVolumeAndMount(t, yamlSplit[2], options, shouldHaveVol)

	assertStream(t, yamlSplit[3], options)
	assertVolumeAndMount(t, yamlSplit[3], options, shouldHaveVol)

	require.False(t, strings.Contains(deploymentYaml, "Conflict: cannot merge map onto non-map"), "'cannot merge map onto non-map' warning should not be logged")
}

func assertVolumeAndMount(t *testing.T, tierYaml string, options *helm.Options, shouldHaveVol bool) {
	var deploymentObj appsv1.Deployment
	UnmarshalK8SYaml(t, tierYaml, &deploymentObj)
	pod := deploymentObj.Spec.Template.Spec

	var foundVol = false
	for _, vol := range pod.Volumes {
		if vol.Name == "pega-volume-import-certificates" {
			foundVol = true
			break
		}
	}
	require.Equal(t, shouldHaveVol, foundVol)

	var foundVolMount = false
	for _, container := range pod.Containers {
		if container.Name == "pega-web-tomcat" {
			for _, volMount := range container.VolumeMounts {
				if volMount.Name == "pega-volume-import-certificates" {
					require.Equal(t, "/opt/pega/certs", volMount.MountPath)
					foundVolMount = true
					break
				}
			}
			break
		}
	}
	require.Equal(t, shouldHaveVol, foundVolMount)

}
