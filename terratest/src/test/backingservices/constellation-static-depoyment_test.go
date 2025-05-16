package backingservices

import (
	"testing"

	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
)

func TestConstellationStaticDeploymentWithAffinity(t *testing.T) {

	var affintiyBasePath = "constellation.affinity.nodeAffinity.requiredDuringSchedulingIgnoredDuringExecution.nodeSelectorTerms[0].matchExpressions[0]."

	helmChartParser := NewHelmConfigParser(
		NewHelmTestFromTemplate(t, helmChartRelativePath, map[string]string{
			"constellation.enabled":        "true",
			affintiyBasePath + "key":       "kubernetes.io/os",
			affintiyBasePath + "operator":  "In",
			affintiyBasePath + "values[0]": "linux",
		},
			[]string{"charts/constellation/templates/clln-deployment.yaml"}),
	)

	var cllnDeploymentObj appsv1.Deployment
	helmChartParser.getResourceYAML(SearchResourceOption{
		Name: "constellation",
		Kind: "Deployment",
	}, &cllnDeploymentObj)

	deploymentAffinity := cllnDeploymentObj.Spec.Template.Spec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution
	require.Equal(t, "kubernetes.io/os", deploymentAffinity.NodeSelectorTerms[0].MatchExpressions[0].Key)
	require.Equal(t, "In", string(deploymentAffinity.NodeSelectorTerms[0].MatchExpressions[0].Operator))
	require.Equal(t, "linux", deploymentAffinity.NodeSelectorTerms[0].MatchExpressions[0].Values[0])
	require.Empty(t, cllnDeploymentObj.Spec.Template.Spec.Tolerations)
}

func TestConstellationStaticDeploymentWithTolerations(t *testing.T) {

	helmChartParser := NewHelmConfigParser(
		NewHelmTestFromTemplate(t, helmChartRelativePath, map[string]string{
			"constellation.enabled":                 "true",
			"constellation.tolerations[0].key":      "key1",
			"constellation.tolerations[0].value":    "value1",
			"constellation.tolerations[0].operator": "Equal",
			"constellation.tolerations[0].effect":   "NotSchedule",
		},
			[]string{"charts/constellation/templates/clln-deployment.yaml"}),
	)

	var cllnDeploymentObj appsv1.Deployment
	helmChartParser.getResourceYAML(SearchResourceOption{
		Name: "constellation",
		Kind: "Deployment",
	}, &cllnDeploymentObj)

	deploymentTolerations := cllnDeploymentObj.Spec.Template.Spec.Tolerations
	require.Equal(t, "key1", deploymentTolerations[0].Key)
	require.Equal(t, "value1", deploymentTolerations[0].Value)
	require.Equal(t, "Equal", string(deploymentTolerations[0].Operator))
	require.Equal(t, "NotSchedule", string(deploymentTolerations[0].Effect))
	require.Empty(t, cllnDeploymentObj.Spec.Template.Spec.Affinity)
}

func TestConstellationStaticDeploymentCustomVolumes(t *testing.T) {
	t.Run("Only customerAssetVolumeClaimName is set", func(t *testing.T) {
		helmChartParser := NewHelmConfigParser(
			NewHelmTestFromTemplate(t, helmChartRelativePath, map[string]string{
				"constellation.enabled": "true",
				"constellation.customerAssetVolumeClaimName": "customer-claim",
			},
				[]string{"charts/constellation/templates/clln-deployment.yaml"}),
		)
	
		var deployment appsv1.Deployment
		helmChartParser.getResourceYAML(SearchResourceOption{
			Name: "constellation",
			Kind: "Deployment",
		}, &deployment)
	
		volumes := deployment.Spec.Template.Spec.Volumes
		var foundCustomer bool
		for _, vol := range volumes {
			if vol.Name == "constellation-appstatic-assets" {
				foundCustomer = true
			}
		}
	
		require.True(t, foundCustomer, "Expected volume with claimName 'customer-claim' not found")
		require.Equal(t, 1, len(volumes), "Expected exactly one volume named 'constellation-appstatic-assets'")

	})

	t.Run("Only custom.volumes is set", func(t *testing.T) {
		helmChartParser := NewHelmConfigParser(
			NewHelmTestFromTemplate(t, helmChartRelativePath, map[string]string{
				"constellation.enabled": "true",
				"constellation.custom.volumes[0].name": "custom-volume",
				"constellation.custom.volumes[0].hostPath.path": "/mnt/custom-path",


			},
				[]string{"charts/constellation/templates/clln-deployment.yaml"}),
		)

		var deployment appsv1.Deployment
		helmChartParser.getResourceYAML(SearchResourceOption{
			Name: "constellation",
			Kind: "Deployment",
		}, &deployment)

		volumes := deployment.Spec.Template.Spec.Volumes
		var foundCustom bool

		for _, vol := range volumes {
			t.Logf("Volume Name: %s", vol.Name)
			if vol.Name == "custom-volume"  {
				foundCustom = true
			}
		}
		require.True(t, foundCustom, "Expected volume with claimName 'custom-volume' not found")

	})

	t.Run("Both customerAssetVolumeClaimName and custom.volumes are set", func(t *testing.T) {
		helmChartParser := NewHelmConfigParser(
			NewHelmTestFromTemplate(t, helmChartRelativePath, map[string]string{
				"constellation.enabled": "true",
				"constellation.customerAssetVolumeClaimName": "customer-claim",
				"constellation.custom.volumes[0].name": "custom-volume",
				"constellation.custom.volumes[0].hostPath.path": "/mnt/custom-path",


			},
				[]string{"charts/constellation/templates/clln-deployment.yaml"}),
		)

		var deployment appsv1.Deployment
		helmChartParser.getResourceYAML(SearchResourceOption{
			Name: "constellation",
			Kind: "Deployment",
		}, &deployment)

		volumes := deployment.Spec.Template.Spec.Volumes
		foundCustomer := false
		foundCustom := false
		for _, vol := range volumes {
			// log the volume names for debugging
			t.Logf("Volume Name: %s", vol.Name)
			if vol.Name == "constellation-appstatic-assets"  {
				foundCustomer = true
			}
			if vol.Name == "custom-volume"  {
				foundCustom = true
			}
		}
		require.True(t, foundCustomer, "Expected volume with claimName 'customer-claim' not found")
		require.True(t, foundCustom, "Expected volume with claimName 'custom-volume' not found")
	})
}
