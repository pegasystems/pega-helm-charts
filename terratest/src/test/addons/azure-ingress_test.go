package addons

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	"path/filepath"
	"testing"
)

func TestShouldNotContainAzureIngressIfDisabled(t *testing.T) {
	t.Parallel()
	helmChartPath, err := filepath.Abs(helmChartRelativePath)

	require.NoError(t, err)
	options := &helm.Options{
		SetValues: map[string]string{
			"ingress-azure.enabled": "false",
		},
	}

	helmChart := NewHelmConfigParser(t, options, helmChartPath)

	for _, i := range azureIngressResources {
		require.False(t, helmChart.contains(SearchResourceOption{
			name: i.name,
			kind: i.kind,
		}))
	}
}

func TestAzureIngressShouldContainAllResources(t *testing.T) {
	t.Parallel()
	helmChartPath, err := filepath.Abs(helmChartRelativePath)

	require.NoError(t, err)
	options := &helm.Options{
		SetValues: map[string]string{
			"ingress-azure.enabled": "true",
		},
	}

	helmChart := NewHelmConfigParser(t, options, helmChartPath)

	for _, i := range azureIngressResources {
		require.True(t, helmChart.contains(SearchResourceOption{
			name: i.name,
			kind: i.kind,
		}))
	}
}

var azureIngressResources = []SearchResourceOption{
	{
		name: "networking-appgw-k8s-azure-service-principal",
		kind: "Secret",
	},
	{
		name: "release-name-cm-ingress-azure",
		kind: "ConfigMap",
	},
	{
		name: "release-name-sa-ingress-azure",
		kind: "ServiceAccount",
	},
	{
		name: "release-name-ingress-azure",
		kind: "ClusterRole",
	},
	{
		name: "release-name-ingress-azure",
		kind: "ClusterRoleBinding",
	},
	{
		name: "release-name-ingress-azure",
		kind: "Deployment",
	},
}
