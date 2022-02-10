package addons

import (
	"github.com/stretchr/testify/require"
	"k8s.io/api/apps/v1"
	"testing"
)

func Test_shouldNotContainMetricServerIfDisabled(t *testing.T) {
	helmChartParser := NewHelmConfigParser(
		NewHelmTest(t, helmChartRelativePath, map[string]string{
			"metrics-server.enabled": "false",
		}),
	)

	for _, i := range metricServerResources {
		require.False(t, helmChartParser.Contains(SearchResourceOption{
			Name: i.Name,
			Kind: i.Kind,
		}))
	}
}

func Test_shouldContainMetricServerIfEnabled(t *testing.T) {
	helmChartParser := NewHelmConfigParser(
		NewHelmTest(t, helmChartRelativePath, map[string]string{
			"metrics-server.enabled": "true",
		}),
	)

	for _, i := range metricServerResources {
		require.True(t, helmChartParser.Contains(SearchResourceOption{
			Name: i.Name,
			Kind: i.Kind,
		}))
	}
}

func Test_shouldContainCommandArgs(t *testing.T) {
	helmChartParser := NewHelmConfigParser(NewHelmTest(t, helmChartRelativePath, map[string]string{
		"metrics-server.enabled": "true",
	}))

	var deployment *v1.Deployment
	helmChartParser.Find(SearchResourceOption{
		Name: "pega-metrics-server",
		Kind: "Deployment",
	}, &deployment)

	require.Contains(t, deployment.Spec.Template.Spec.Containers[0].Args, "--logtostderr")
}

var metricServerResources = []SearchResourceOption{
	{
		Name: "pega-metrics-server",
		Kind: "ServiceAccount",
	},
	{
		Name: "system:metrics-server-aggregated-reader",
		Kind: "ClusterRole",
	},
	{
		Name: "system:pega-metrics-server",
		Kind: "ClusterRole",
	},
	{
		Name: "pega-metrics-server:system:auth-delegator",
		Kind: "ClusterRoleBinding",
	},
	{
		Name: "system:pega-metrics-server",
		Kind: "ClusterRoleBinding",
	},
	{
		Name: "pega-metrics-server-auth-reader",
		Kind: "RoleBinding",
	},
	{
		Name: "pega-metrics-server",
		Kind: "Service",
	},
	{
		Name: "pega-metrics-server",
		Kind: "Deployment",
	},
	{
		Name: "v1beta1.metrics.k8s.io",
		Kind: "APIService",
	},
}
