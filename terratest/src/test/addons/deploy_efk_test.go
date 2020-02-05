package addons

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	"path/filepath"
	"testing"
)

func TestShouldNotContainDeploy_EFKIfDisabled(t *testing.T) {
	t.Parallel()
	helmChartPath, err := filepath.Abs(helmChartRelativePath)

	require.NoError(t, err)
	options := &helm.Options{
		SetValues: map[string]string{
			"elasticsearch.enabled":         "false",
			"kibana.enabled":                "false",
			"fluentd-elasticsearch.enabled": "false",
		},
	}

	helmChart := NewHelmConfigParser(t, options, helmChartPath)

	for _, i := range deployEfkResources {
		require.False(t, helmChart.contains(SearchResourceOption{
			name: i.name,
			kind: i.kind,
		}))
	}
}
func TestShouldDeploy_EFKContainAllResourcesIfEnabled(t *testing.T) {
	t.Parallel()
	helmChartPath, err := filepath.Abs(helmChartRelativePath)

	require.NoError(t, err)
	options := &helm.Options{
		SetValues: map[string]string{
			"elasticsearch.enabled":         "true",
			"kibana.enabled":                "true",
			"fluentd-elasticsearch.enabled": "true",
		},
	}
	helmChart := NewHelmConfigParser(t, options, helmChartPath)

	for _, i := range deployEfkResources {
		require.True(t, helmChart.contains(SearchResourceOption{
			name: i.name,
			kind: i.kind,
		}))
	}
}

func Test_shouldBeElasticSearchUrlForKibana(t *testing.T) {
	helmChartPath, err := filepath.Abs(helmChartRelativePath)
	require.NoError(t, err)

	helmChartParser := NewHelmConfigParser(t, &helm.Options{
		SetValues: map[string]string{
			"kibana.enabled": "true",
			"kibana.files.kibana.yml.elasticsearch.url": "http://elastic-search-client:9200",
		},
	}, helmChartPath)

	var configmap *v1.ConfigMap
	helmChartParser.find(SearchResourceOption{
		name: "release-name-kibana",
		kind: "ConfigMap",
	}, &configmap)

	require.Contains(t, configmap.Data["kibana.yml"], "http://elastic-search-client:9200")
}
func Test_shouldBeServiceExternalPortForKibana(t *testing.T) {
	helmChartPath, err := filepath.Abs(helmChartRelativePath)
	require.NoError(t, err)

	helmChartParser := NewHelmConfigParser(t, &helm.Options{
		SetValues: map[string]string{
			"kibana.enabled":              "true",
			"kibana.service.externalPort": "80",
		},
	}, helmChartPath)

	var service *v1.Service
	helmChartParser.find(SearchResourceOption{
		name: "release-name-kibana",
		kind: "Service",
	}, &service)

	//require.Contains(t, configmap.Data["kibana.yml"], "http://elastic-search-client:9200")
	require.Equal(t, int32(80), service.Spec.Ports[0].Port)
}

var deployEfkResources = []SearchResourceOption{
	{
		name: "release-name-kibana",
		kind: "ConfigMap",
	},
	{
		name: "release-name-kibana",
		kind: "Service",
	},
	{
		name: "release-name-kibana",
		kind: "Deployment",
	},
	{
		name: "release-name-kibana",
		kind: "Ingress",
	},
	{
		name: "elastic-search",
		kind: "ConfigMap",
	},
	{
		name: "elastic-search-data",
		kind: "StatefulSet",
	},
	{
		name: "elastic-search-master",
		kind: "StatefulSet",
	},
	{
		name: "elastic-search-client",
		kind: "ServiceAccount",
	},
	{
		name: "elastic-search-data",
		kind: "ServiceAccount",
	},
	{
		name: "elastic-search-master",
		kind: "ServiceAccount",
	},
	{
		name: "elastic-search-client",
		kind: "Service",
	},
	{
		name: "elastic-search-discovery",
		kind: "Service",
	},
	{
		name: "elastic-search-client",
		kind: "Deployment",
	},
	{
		name: "release-name-fluentd-elasticsearch",
		kind: "ConfigMap",
	},
	{
		name: "release-name-fluentd-elasticsearch",
		kind: "ServiceAccount",
	},
	{
		name: "release-name-fluentd-elasticsearch",
		kind: "ClusterRole",
	},
	{
		name: "release-name-fluentd-elasticsearch",
		kind: "ClusterRoleBinding",
	},
	{
		name: "release-name-fluentd-elasticsearch",
		kind: "DaemonSet",
	},
}
