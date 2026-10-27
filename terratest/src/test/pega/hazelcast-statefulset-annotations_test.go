package pega

import (
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
)

var hazelcastStatefulSetTemplates = map[string]map[string]string{
	"charts/hazelcast/templates/clustering-service-deployment.yaml": {"hazelcast.clusteringServiceEnabled": "true"},
	"charts/hazelcast/templates/pega-hz-deployment.yaml":            {"hazelcast.enabled": "true"},
}

func TestHazelcastStatefulSetAnnotations(t *testing.T) {
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for template, enableValues := range hazelcastStatefulSetTemplates {
		setValues := map[string]string{
			"global.provider":        "k8s",
			"global.actions.execute": "deploy",
		}
		for k, v := range enableValues {
			setValues[k] = v
		}
		options := &helm.Options{
			SetValues:   setValues,
			ValuesFiles: []string{"data/values_with_hazelcast_statefulset_annotations.yaml"},
		}

		var statefulset appsv1.StatefulSet
		UnmarshalK8SYaml(t, RenderTemplate(t, options, helmChartPath, []string{template}), &statefulset)

		require.Equal(t, map[string]string{
			"example.com/boolean":           "true",
			"example.com/number":            "90",
			"example.com/plain":             "plain-value",
			"example.com/templated":         "from-values",
			"example.com/templated-default": "fallback",
		}, statefulset.ObjectMeta.Annotations, template)
		for key := range statefulset.ObjectMeta.Annotations {
			require.NotContains(t, statefulset.Spec.Template.ObjectMeta.Annotations, key, template)
		}
	}
}

func TestHazelcastStatefulSetWithoutAnnotations(t *testing.T) {
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for template, enableValues := range hazelcastStatefulSetTemplates {
		setValues := map[string]string{
			"global.provider":        "k8s",
			"global.actions.execute": "deploy",
		}
		for k, v := range enableValues {
			setValues[k] = v
		}

		var statefulset appsv1.StatefulSet
		UnmarshalK8SYaml(t, RenderTemplate(t, &helm.Options{SetValues: setValues}, helmChartPath, []string{template}), &statefulset)

		require.Empty(t, statefulset.ObjectMeta.Annotations, template)
	}
}
