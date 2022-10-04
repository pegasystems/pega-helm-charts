package pega

import (
	"path/filepath"
	"testing"
	k8score "k8s.io/api/core/v1"
	appsv1 "k8s.io/api/apps/v1"
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	intstr "k8s.io/apimachinery/pkg/util/intstr"
)

func TestConstellationDeployment(t *testing.T) {

	var supportedVendors = []string{"k8s"}
	var supportedOperations =  []string{"deploy","install-deploy","upgrade-deploy"}

		helmChartPath, err := filepath.Abs(PegaHelmChartPath)
		require.NoError(t, err)

	for _, vendor := range supportedVendors {
		for _, operation := range supportedOperations {
			var options = &helm.Options{
				SetValues: map[string]string{
					"global.provider":        vendor,
					"global.actions.execute": operation,
					"constellation.enabled":  "true",
					"installer.upgrade.upgradeType": "zero-downtime",
				},
			}
	yamlContent := RenderTemplate(t, options, helmChartPath, []string{"charts/constellation/templates/clln-deployment.yaml"})
	assertConstellationDeployment(t, yamlContent, options)
	constellationService := RenderTemplate(t, options, helmChartPath, []string{"charts/constellation/templates/clln-service.yaml"})
	assertConstellationService(t, constellationService, options)
  }
 }
}

func assertConstellationDeployment(t *testing.T, deploymentYaml string, options *helm.Options) {
	var constellationDeploymentObj appsv1.Deployment
	UnmarshalK8SYaml(t, deploymentYaml, &constellationDeploymentObj)
	require.Equal(t, *constellationDeploymentObj.Spec.Replicas, int32(2))
	deploymentSpec := constellationDeploymentObj.Spec.Template.Spec
	require.Equal(t, "/c11n/v860/ping", deploymentSpec.Containers[0].LivenessProbe.HTTPGet.Path)
}

func assertConstellationService(t *testing.T, constellationService string, options *helm.Options) {
	var constellationServiceObj k8score.Service
	helm.UnmarshalK8SYaml(t, constellationService, &constellationServiceObj)
	require.Equal(t, constellationServiceObj.Spec.Selector["app"], "constellation")
	// require.Equal(t, constellationServiceObj.Spec.Ports[0].Protocol, "TCP")
	require.Equal(t, constellationServiceObj.Spec.Ports[0].Port, int32(3000))
	require.Equal(t, constellationServiceObj.Spec.Ports[0].TargetPort, intstr.FromInt(3000))
}
