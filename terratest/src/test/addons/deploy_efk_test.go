package addons

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	"path/filepath"
	"testing"
)

func TestShouldNotContainDeploy_EFKIfDisabled(t *testing.T) {
	t.Parallel()
	helmChartPath, err := filepath.Abs(helmChartRelativePath)

	require.NoError(t, err)
	options := &helm.Options{
		SetValues: map[string]string{
			"deploy_efk": "false",
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
			"deploy_efk": "true",
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
