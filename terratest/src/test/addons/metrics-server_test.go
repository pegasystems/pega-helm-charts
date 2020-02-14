package addons

import (
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/apps/v1"
	"test/common"
	"testing"
)

func Test_shouldNotContainMetricServerIfDisabled(t *testing.T) {
	helmChartParser := common.NewHelmConfigParser(
		common.NewHelmTest(t, helmChartRelativePath, map[string]string{
			"metrics-server.enabled": "false",
		}),
	)

	for _, i := range metricServerResources {
		require.False(t, helmChartParser.Contains(common.SearchResourceOption{
			Name: i.Name,
			Kind: i.Kind,
		}))
	}
}

func Test_shouldContainMetricServerIfEnabled(t *testing.T) {
	helmChartParser := common.NewHelmConfigParser(
		common.NewHelmTest(t, helmChartRelativePath, map[string]string{
			"metrics-server.enabled": "true",
		}),
	)

	for _, i := range metricServerResources {
		require.True(t, helmChartParser.Contains(common.SearchResourceOption{
			Name: i.Name,
			Kind: i.Kind,
		}))
	}
}

func Test_shouldContainCommandArgs(t *testing.T) {
	helmChartParser := common.NewHelmConfigParser(common.NewHelmTest(t, helmChartRelativePath, map[string]string{
		"metrics-server.enabled": "true",
	}))

	var deployment *v1.Deployment
	helmChartParser.Find(common.SearchResourceOption{
		Name: "release-name-metrics-server",
		Kind: "Deployment",
	}, &deployment)

	require.Contains(t, deployment.Spec.Template.Spec.Containers[0].Command, "--logtostderr")
}

var metricServerResources = []common.SearchResourceOption{
	{
		Name: "release-name-metrics-server",
		Kind: "ServiceAccount",
	},
	{
		Name: "system:metrics-server-aggregated-reader",
		Kind: "ClusterRole",
	},
	{
		Name: "system:release-name-metrics-server",
		Kind: "ClusterRole",
	},
	{
		Name: "release-name-metrics-server:system:auth-delegator",
		Kind: "ClusterRoleBinding",
	},
	{
		Name: "system:release-name-metrics-server",
		Kind: "ClusterRoleBinding",
	},
	{
		Name: "release-name-metrics-server-auth-reader",
		Kind: "RoleBinding",
	},
	{
		Name: "release-name-metrics-server",
		Kind: "Service",
	},
	{
		Name: "release-name-metrics-server-test",
		Kind: "Pod",
	},
	{
		Name: "release-name-metrics-server",
		Kind: "Deployment",
	},
	{
		Name: "v1beta1.metrics.k8s.io",
		Kind: "APIService",
	},
}
