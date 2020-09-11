package addons

import (
	"github.com/stretchr/testify/require"
	"k8s.io/api/apps/v1beta2"
	"k8s.io/api/extensions/v1beta1"
	"testing"
)

func TestShouldNotContainDeploy_EFKIfDisabled(t *testing.T) {
	helmChartParser := NewHelmConfigParser(
		NewHelmTest(t, helmChartRelativePath, map[string]string{
			"elasticsearch.enabled":         "false",
			"kibana.enabled":                "false",
			"fluentd-elasticsearch.enabled": "false",
		}),
	)

	for _, i := range deployEfkResources {
		require.False(t, helmChartParser.Contains(SearchResourceOption{
			Name: i.Name,
			Kind: i.Kind,
		}))
	}
}

func TestShouldDeploy_EFKContainAllResourcesIfEnabled(t *testing.T) {
	helmChartParser := NewHelmConfigParser(
		NewHelmTest(t, helmChartRelativePath, map[string]string{
			"elasticsearch.enabled":         "true",
			"kibana.enabled":                "true",
			"fluentd-elasticsearch.enabled": "true",
		}),
	)

	for _, i := range deployEfkResources {
		require.True(t, helmChartParser.Contains(SearchResourceOption{
			Name: i.Name,
			Kind: i.Kind,
		}))
	}
}

func Test_shouldBeIngressEnabledForKibana(t *testing.T) {
	helmChartParser := NewHelmConfigParser(
		NewHelmTest(t, helmChartRelativePath, map[string]string{
			"kibana.enabled":         "true",
			"kibana.ingress.enabled": "true",
		}),
	)

	require.True(t, helmChartParser.Contains(SearchResourceOption{
		Name: "pega-kibana",
		Kind: "Ingress",
	}))
}

func Test_shouldBeIngressDisabledForKibana(t *testing.T) {
	helmChartParser := NewHelmConfigParser(
		NewHelmTest(t, helmChartRelativePath, map[string]string{
			"kibana.enabled":         "true",
			"kibana.ingress.enabled": "false",
		}),
	)

	require.False(t, helmChartParser.Contains(SearchResourceOption{
		Name: "pega-kibana",
		Kind: "Ingress",
	}))
}

func Test_shouldBeHostForIngressKibana(t *testing.T) {
	helmChartParser := NewHelmConfigParser(
		NewHelmTest(t, helmChartRelativePath, map[string]string{
			"kibana.enabled":         "true",
			"kibana.ingress.enabled": "true",
			"kibana.ingress.hosts":   "{YOUR_WEB.KIBANA.EXAMPLE.COM}",
		}),
	)

	var ingress *v1beta1.Ingress
	helmChartParser.Find(SearchResourceOption{
		Name: "pega-kibana",
		Kind: "Ingress",
	}, &ingress)

	require.Contains(t, ingress.Spec.Rules[0].Host, "YOUR_WEB.KIBANA.EXAMPLE.COM")

}
func Test_shouldBeHostForElasticsearch(t *testing.T) {
	helmChartParser := NewHelmConfigParser(
		NewHelmTest(t, helmChartRelativePath, map[string]string{
			"fluentd-elasticsearch.enabled": "true",
			"fluentd-elasticsearch.elasticsearch.host": "elasticsearch-master:9200",
		}),
	)

	var daemon *v1beta2.DaemonSet
	helmChartParser.Find(SearchResourceOption{
		Name: "pega-fluentd-elasticsearch",
		Kind: "DaemonSet",
	}, &daemon)

	require.Equal(t, "elasticsearch-master:9200", daemon.Spec.Template.Spec.Containers[0].Env[1].Value)
}

var deployEfkResources = []SearchResourceOption{
	{
		Name: "pega-kibana",
		Kind: "Service",
	},
	{
		Name: "pega-kibana",
		Kind: "Deployment",
	},
	{
		Name: "pega-kibana",
		Kind: "Ingress",
	},	
	{
		Name: "elasticsearch-master",
		Kind: "StatefulSet",
	},
	{
		Name: "elasticsearch-master-pdb",
		Kind: "PodDisruptionBudget",
	},
	{
		Name: "elasticsearch-master",
		Kind: "Service",
	},
	{
		Name: "elasticsearch-master-headless",
		Kind: "Service",
	},
	{
		Name: "pega-fluentd-elasticsearch",
		Kind: "ConfigMap",
	},
	{
		Name: "pega-fluentd-elasticsearch",
		Kind: "ServiceAccount",
	},
	{
		Name: "pega-fluentd-elasticsearch",
		Kind: "DaemonSet",
	},
}
