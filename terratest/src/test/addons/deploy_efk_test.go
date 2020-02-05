package addons

import (
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	"test/testhelpers"
	"testing"
)

func TestShouldNotContainDeploy_EFKIfDisabled(t *testing.T) {
	helmChartParser := testhelpers.NewHelmConfigParser(
		testhelpers.NewHelmTest(t, helmChartRelativePath, map[string]string{
			"elasticsearch.enabled":         "false",
			"kibana.enabled":                "false",
			"fluentd-elasticsearch.enabled": "false",
		}),
	)

	for _, i := range deployEfkResources {
		require.False(t, helmChartParser.Contains(testhelpers.SearchResourceOption{
			Name: i.Name,
			Kind: i.Kind,
		}))
	}
}
func TestShouldDeploy_EFKContainAllResourcesIfEnabled(t *testing.T) {
	helmChartParser := testhelpers.NewHelmConfigParser(
		testhelpers.NewHelmTest(t, helmChartRelativePath, map[string]string{
			"elasticsearch.enabled":         "true",
			"kibana.enabled":                "true",
			"fluentd-elasticsearch.enabled": "true",
		}),
	)

	for _, i := range deployEfkResources {
		require.True(t, helmChartParser.Contains(testhelpers.SearchResourceOption{
			Name: i.Name,
			Kind: i.Kind,
		}))
	}
}

func Test_shouldBeElasticSearchUrlForKibana(t *testing.T) {
	helmChartParser := testhelpers.NewHelmConfigParser(
		testhelpers.NewHelmTest(t, helmChartRelativePath, map[string]string{
			"kibana.enabled": "true",
			"kibana.files.kibana.yml.elasticsearch.url": "http://elastic-search-client:9200",
		}),
	)

	var configmap *v1.ConfigMap
	helmChartParser.Find(testhelpers.SearchResourceOption{
		Name: "release-name-kibana",
		Kind: "ConfigMap",
	}, &configmap)

	require.Contains(t, configmap.Data["kibana.yml"], "http://elastic-search-client:9200")
}
func Test_shouldBeServiceExternalPortForKibana(t *testing.T) {
	helmChartParser := testhelpers.NewHelmConfigParser(
		testhelpers.NewHelmTest(t, helmChartRelativePath, map[string]string{
			"kibana.enabled":              "true",
			"kibana.service.externalPort": "80",
		}),
	)

	var service *v1.Service
	helmChartParser.Find(testhelpers.SearchResourceOption{
		Name: "release-name-kibana",
		Kind: "Service",
	}, &service)

	//require.Contains(t, configmap.Data["kibana.yml"], "http://elastic-search-client:9200")
	require.Equal(t, int32(80), service.Spec.Ports[0].Port)
}

var deployEfkResources = []testhelpers.SearchResourceOption{
	{
		Name: "release-name-kibana",
		Kind: "ConfigMap",
	},
	{
		Name: "release-name-kibana",
		Kind: "Service",
	},
	{
		Name: "release-name-kibana",
		Kind: "Deployment",
	},
	{
		Name: "release-name-kibana",
		Kind: "Ingress",
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
		Name: "release-name-fluentd-elasticsearch",
		Kind: "ConfigMap",
	},
	{
		Name: "release-name-fluentd-elasticsearch",
		Kind: "ServiceAccount",
	},
	{
		Name: "release-name-fluentd-elasticsearch",
		Kind: "ClusterRole",
	},
	{
		Name: "release-name-fluentd-elasticsearch",
		Kind: "ClusterRoleBinding",
	},
	{
		Name: "release-name-fluentd-elasticsearch",
		Kind: "DaemonSet",
	},
}
