package addons

import (
	"github.com/stretchr/testify/require"
	"k8s.io/api/apps/v1beta2"
	v1 "k8s.io/api/core/v1"
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

func Test_shouldBeFullnameOverrideForElasticsearch(t *testing.T) {
	helmChartParser := NewHelmConfigParser(
		NewHelmTest(t, helmChartRelativePath, map[string]string{
			"elasticsearch.enabled": "true",
		}),
	)
	require.True(t, helmChartParser.Contains(SearchResourceOption{
		Name: "elastic-search",
		Kind: "ConfigMap",
	}))
}

func Test_shouldBeElasticSearchUrlForKibana(t *testing.T) {
	helmChartParser := NewHelmConfigParser(
		NewHelmTest(t, helmChartRelativePath, map[string]string{
			"kibana.enabled": "true",
			"kibana.files.kibana.yml.elasticsearch.url": "http://elastic-search-client:9200",
		}),
	)

	var configmap *v1.ConfigMap
	helmChartParser.Find(SearchResourceOption{
		Name: "pega-kibana",
		Kind: "ConfigMap",
	}, &configmap)

	require.Contains(t, configmap.Data["kibana.yml"], "http://elastic-search-client:9200")
}
func Test_shouldBeServiceExternalPortForKibana(t *testing.T) {
	helmChartParser := NewHelmConfigParser(
		NewHelmTest(t, helmChartRelativePath, map[string]string{
			"kibana.enabled":              "true",
			"kibana.service.externalPort": "80",
		}),
	)

	var service *v1.Service
	helmChartParser.Find(SearchResourceOption{
		Name: "pega-kibana",
		Kind: "Service",
	}, &service)

	require.Equal(t, int32(80), service.Spec.Ports[0].Port)
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
			"elasticsearch.host":            "elastic-search-client",
		}),
	)

	var daemon *v1beta2.DaemonSet
	helmChartParser.Find(SearchResourceOption{
		Name: "pega-fluentd-elasticsearch",
		Kind: "DaemonSet",
	}, &daemon)

	require.Equal(t, "elastic-search-client", daemon.Spec.Template.Spec.Containers[0].Env[1].Value)
}

func Test_shouldBeBufferChunkLimitForElasticsearch(t *testing.T) {
	helmChartParser := NewHelmConfigParser(
		NewHelmTest(t, helmChartRelativePath, map[string]string{
			"fluentd-elasticsearch.enabled":    "true",
			"elasticsearch.buffer_chunk_limit": "250M",
		}),
	)

	var daemon *v1beta2.DaemonSet
	helmChartParser.Find(SearchResourceOption{
		Name: "pega-fluentd-elasticsearch",
		Kind: "DaemonSet",
	}, &daemon)

	require.Equal(t, "250M", daemon.Spec.Template.Spec.Containers[0].Env[4].Value)
}

func Test_shouldBeBufferQueueLimitForElasticsearch(t *testing.T) {
	helmChartParser := NewHelmConfigParser(
		NewHelmTest(t, helmChartRelativePath, map[string]string{
			"fluentd-elasticsearch.enabled":    "true",
			"elasticsearch.buffer_queue_limit": "30",
		}),
	)

	var daemon *v1beta2.DaemonSet
	helmChartParser.Find(SearchResourceOption{
		Name: "pega-fluentd-elasticsearch",
		Kind: "DaemonSet",
	}, &daemon)

	require.Equal(t, "30", daemon.Spec.Template.Spec.Containers[0].Env[5].Value)
}

var deployEfkResources = []SearchResourceOption{
	{
		Name: "pega-kibana",
		Kind: "ConfigMap",
	},
	{
		Name: "pega-kibana",
		Kind: "Service",
	},
	{
		Name: "pega-kibana",
		Kind: "Deployment",
	},
	{
		Name: "elastic-search",
		Kind: "ConfigMap",
	},
	{
		Name: "elastic-search-data",
		Kind: "StatefulSet",
	},
	{
		Name: "elastic-search-master",
		Kind: "StatefulSet",
	},
	{
		Name: "elastic-search-client",
		Kind: "ServiceAccount",
	},
	{
		Name: "elastic-search-data",
		Kind: "ServiceAccount",
	},
	{
		Name: "elastic-search-master",
		Kind: "ServiceAccount",
	},
	{
		Name: "elastic-search-client",
		Kind: "Service",
	},
	{
		Name: "elastic-search-discovery",
		Kind: "Service",
	},
	{
		Name: "elastic-search-client",
		Kind: "Deployment",
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
		Kind: "ClusterRole",
	},
	{
		Name: "pega-fluentd-elasticsearch",
		Kind: "ClusterRoleBinding",
	},
	{
		Name: "pega-fluentd-elasticsearch",
		Kind: "DaemonSet",
	},
}
