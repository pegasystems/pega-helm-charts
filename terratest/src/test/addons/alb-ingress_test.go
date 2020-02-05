package addons

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	"path/filepath"
	"testing"
)

func TestShouldNotContainAlbIngressIfDisabled(t *testing.T) {
	t.Parallel()
	helmChartPath, err := filepath.Abs(helmChartRelativePath)

	require.NoError(t, err)
	options := &helm.Options{
		SetValues: map[string]string{
			"aws-alb-ingress-controller.enabled": "false",
		},
	}

	helmChart := NewHelmConfigParser(t, options, helmChartPath)

	for _, i := range albIngressResources {
		require.False(t, helmChart.contains(SearchResourceOption{
			name: i.name,
			kind: i.kind,
		}))
	}
}

func TestAlbIngressShouldContainAllResources(t *testing.T) {
	t.Parallel()
	helmChartPath, err := filepath.Abs(helmChartRelativePath)

	require.NoError(t, err)
	options := &helm.Options{
		SetValues: map[string]string{
			"aws-alb-ingress-controller.enabled": "true",
		},
	}

	helmChart := NewHelmConfigParser(t, options, helmChartPath)

	for _, i := range albIngressResources {
		require.True(t, helmChart.contains(SearchResourceOption{
			name: i.name,
			kind: i.kind,
		}))
	}
}

var albIngressResources = []SearchResourceOption{
	{
		name: "release-name-aws-alb-ingress-controller",
		kind: "ServiceAccount",
	},
	{
		name: "release-name-aws-alb-ingress-controller",
		kind: "ClusterRole",
	},
	{
		name: "release-name-aws-alb-ingress-controller",
		kind: "ClusterRoleBinding",
	},
	{
		name: "release-name-aws-alb-ingress-controller",
		kind: "Deployment",
	},
}
