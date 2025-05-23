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
