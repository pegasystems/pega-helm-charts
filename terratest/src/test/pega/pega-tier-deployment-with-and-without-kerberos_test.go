package pega

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	"path/filepath"
	"strings"
	"testing"
)

func TestPegaDeploymentWithAndWithoutKerberos(t *testing.T) {

	var supportedVendors = []string{"k8s"}
	var supportedOperations = []string{"deploy", "install-deploy", "upgrade-deploy"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for _, vendor := range supportedVendors {
		for _, operation := range supportedOperations {

			var options = &helm.Options{
				ValuesFiles: []string{"data/values_with_kerberos.yaml"},
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
			assertVolumeAndMountForKerberos(t, yamlSplit[1], true)

			assertBatch(t, yamlSplit[2], options)
			assertVolumeAndMountForKerberos(t, yamlSplit[2], true)

			assertStream(t, yamlSplit[3], options)
			assertVolumeAndMountForKerberos(t, yamlSplit[3], true)

			options.ValuesFiles = []string{"data/values_without_kerberos.yaml"}

			deploymentYaml = RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-deployment.yaml"})
			yamlSplit = strings.Split(deploymentYaml, "---")
			assertWeb(t, yamlSplit[1], options)
			assertVolumeAndMountForKerberos(t, yamlSplit[1], false)

			assertBatch(t, yamlSplit[2], options)
			assertVolumeAndMountForKerberos(t, yamlSplit[2], false)

			assertStream(t, yamlSplit[3], options)
			assertVolumeAndMountForKerberos(t, yamlSplit[3], false)
		}
	}
}

func assertVolumeAndMountForKerberos(t *testing.T, tierYaml string, shouldHaveVol bool) {
	var deploymentObj appsv1.Deployment
	UnmarshalK8SYaml(t, tierYaml, &deploymentObj)
	pod := deploymentObj.Spec.Template.Spec

	var foundVol = false
	for _, vol := range pod.Volumes {
		if vol.Name == "pega-import-kerberos-config" {
			foundVol = true
			break
		}
	}
	require.Equal(t, shouldHaveVol, foundVol)

	var foundVolMount = false
	for _, container := range pod.Containers {
		if container.Name == "pega-web-tomcat" {
			for _, volMount := range container.VolumeMounts {
				if volMount.Name == "pega-import-kerberos-config" {
					require.Equal(t, "/opt/pega/kerberos", volMount.MountPath)
					foundVolMount = true
					break
				}
			}
			break
		}
	}
	require.Equal(t, shouldHaveVol, foundVolMount)

}
