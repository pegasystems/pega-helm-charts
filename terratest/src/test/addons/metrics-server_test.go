package addons

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/apps/v1"
	"path/filepath"
	"testing"
)

func Test_shouldNotContainMetricServerIfDisabled(t *testing.T) {
	t.Parallel()

	helmChartPath, err := filepath.Abs(helmChartRelativePath)
	require.NoError(t, err)

	options := &helm.Options{
		SetValues: map[string]string{
			"metrics-server.enabled": "false",
		},
	}

	helmChartParser := NewHelmConfigParser(t, options, helmChartPath)

	for _, i := range metricServerResources {
		require.False(t, helmChartParser.contains(SearchResourceOption{
			name: i.name,
			kind: i.kind,
		}))
	}
}

func Test_shouldContainMetricServerIfEnabled(t *testing.T) {
	t.Parallel()

	helmChartPath, err := filepath.Abs(helmChartRelativePath)
	require.NoError(t, err)

	options := &helm.Options{
		SetValues: map[string]string{
			"metrics-server.enabled": "true",
		},
	}

	helmChartParser := NewHelmConfigParser(t, options, helmChartPath)

	for _, i := range metricServerResources {
		require.True(t, helmChartParser.contains(SearchResourceOption{
			name: i.name,
			kind: i.kind,
		}))
	}
}

func Test_shouldContainCommandArgs(t *testing.T) {
	helmChartPath, err := filepath.Abs(helmChartRelativePath)
	require.NoError(t, err)

	helmChartParser := NewHelmConfigParser(t, &helm.Options{}, helmChartPath)

	var deployment *v1.Deployment
	helmChartParser.find(SearchResourceOption{
		name: "release-name-metrics-server",
		kind: "Deployment",
	}, &deployment)

	require.Contains(t, deployment.Spec.Template.Spec.Containers[0].Command, "--logtostderr")
}

var metricServerResources = []SearchResourceOption{
	{
		name: "release-name-metrics-server",
		kind: "ServiceAccount",
	},
	{
		name: "system:metrics-server-aggregated-reader",
		kind: "ClusterRole",
	},
	{
		name: "system:release-name-metrics-server",
		kind: "ClusterRole",
	},
	{
		name: "release-name-metrics-server:system:auth-delegator",
		kind: "ClusterRoleBinding",
	},
	{
		name: "system:release-name-metrics-server",
		kind: "ClusterRoleBinding",
	},
	{
		name: "release-name-metrics-server-auth-reader",
		kind: "RoleBinding",
	},
	{
		name: "release-name-metrics-server",
		kind: "Service",
	},
	{
		name: "release-name-metrics-server-test",
		kind: "Pod",
	},
	{
		name: "release-name-metrics-server",
		kind: "Deployment",
	},
	{
		name: "v1beta1.metrics.k8s.io",
		kind: "APIService",
	},
}
