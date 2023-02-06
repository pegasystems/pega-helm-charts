package pega

import (
	"fmt"
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	"path/filepath"
	"strings"
	"testing"
)

func TestPegaDeploymentWithExternalSecretCredentials(t *testing.T) {

	var supportedVendors = []string{"k8s"}
	var supportedOperations = []string{"deploy", "install-deploy", "upgrade-deploy"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for _, vendor := range supportedVendors {
		for _, operation := range supportedOperations {

			var options = &helm.Options{
				ValuesFiles: []string{"data/values_with_externalsecrets.yaml"},
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
			assertCredentialVolumeAndMount(t, yamlSplit[1], options, true)
		}
	}
}
func assertCredentialVolumeAndMount(t *testing.T, tierYaml string, options *helm.Options, shouldHaveVol bool) {
	var deploymentObj appsv1.Deployment
	UnmarshalK8SYaml(t, tierYaml, &deploymentObj)
	pod := deploymentObj.Spec.Template.Spec

	var foundVol = false
	for _, vol := range pod.Volumes {
		if vol.Name == "pega-volume-credentials" {
			sources := vol.VolumeSource.Projected.Sources
			fmt.Println(sources[0].Secret.LocalObjectReference.Name)
			require.Equal(t, len(sources), 5)
			require.Equal(t, sources[0].Secret.LocalObjectReference.Name, "pega-credentials-secret")
			require.Equal(t, sources[1].Secret.LocalObjectReference.Name, "hazelcast-secret")
			require.Equal(t, sources[2].Secret.LocalObjectReference.Name, "customArtifactory-secret")
			require.Equal(t, sources[3].Secret.LocalObjectReference.Name, "dds-secret")
			require.Equal(t, sources[4].Secret.LocalObjectReference.Name, "kafka-secret")
			foundVol = true
			break
		}
	}
	require.Equal(t, shouldHaveVol, foundVol)

	var foundVolMount = false
	for _, container := range pod.Containers {
		if container.Name == "pega-web-tomcat" {
			for _, volMount := range container.VolumeMounts {
				if volMount.Name == "pega-volume-credentials" {
					require.Equal(t, "/opt/pega/secrets", volMount.MountPath)
					foundVolMount = true
					break
				}
			}
			break
		}
	}
	require.Equal(t, shouldHaveVol, foundVolMount)

}
