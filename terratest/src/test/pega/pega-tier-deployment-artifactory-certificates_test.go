package pega

import (
	"path/filepath"
	"strings"
	"testing"
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
)

func TestPegaDeploymentWithArtifactoryCerts(t *testing.T) {

	var supportedVendors = []string{"k8s"}
	var supportedOperations = []string{"deploy", "install-deploy", "upgrade-deploy"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for _, vendor := range supportedVendors {
		for _, operation := range supportedOperations {

			var options = &helm.Options{
				ValuesFiles: []string{"data/values_with_artifactory_cert.yaml"},
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
			assertArtifactoryCertificatesVolume(t, yamlSplit[1], options, true)

			assertBatch(t, yamlSplit[2], options)
			assertArtifactoryCertificatesVolume(t, yamlSplit[2], options, true)

			//assertStream(t, yamlSplit[3], options)
			//assertArtifactoryCertificatesVolume(t, yamlSplit[3], options, true)

			options.ValuesFiles = []string{"data/values_with_artifactory_sslverification_disabled.yaml"}

			deploymentYaml = RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-deployment.yaml"})
			yamlSplit = strings.Split(deploymentYaml, "---")
			assertWeb(t, yamlSplit[1], options)
			assertArtifactoryCertificatesVolume(t, yamlSplit[1], options, false)

			assertBatch(t, yamlSplit[2], options)
			assertArtifactoryCertificatesVolume(t, yamlSplit[2], options, false)

			//assertStream(t, yamlSplit[3], options)
			//assertArtifactoryCertificatesVolume(t, yamlSplit[3], options, false)
		}
	}
}

func assertArtifactoryCertificatesVolume(t *testing.T, tierYaml string, options *helm.Options, shouldHaveVol bool) {
	var deploymentObj appsv1.Deployment
	UnmarshalK8SYaml(t, tierYaml, &deploymentObj)
	pod := deploymentObj.Spec.Template.Spec

    if (shouldHaveVol) {
        var volumes = pod.Volumes
        var pegaVolumeCustomArtifactoryCertificate = findNamedVolume(volumes, "pega-volume-custom-artifactory-certificate")
        require.NotNil(t, pegaVolumeCustomArtifactoryCertificate)
    }
}