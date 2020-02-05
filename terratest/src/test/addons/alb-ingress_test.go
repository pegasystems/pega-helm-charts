package addons

import (
	"github.com/stretchr/testify/require"
	"test/testhelpers"
	"testing"
)

func TestShouldNotContainAlbIngressIfDisabled(t *testing.T) {
	helmChartParser := testhelpers.NewHelmConfigParser(
		testhelpers.NewHelmTest(t, helmChartRelativePath, map[string]string{
			"aws-alb-ingress-controller.enabled": "false",
		}),
	)

	for _, i := range albIngressResources {
		require.False(t, helmChartParser.Contains(testhelpers.SearchResourceOption{
			Name: i.Name,
			Kind: i.Kind,
		}))
	}
}

func TestAlbIngressShouldContainAllResources(t *testing.T) {
	helmChartParser := testhelpers.NewHelmConfigParser(
		testhelpers.NewHelmTest(t, helmChartRelativePath, map[string]string{
			"aws-alb-ingress-controller.enabled": "true",
		}),
	)

	for _, i := range albIngressResources {
		require.True(t, helmChartParser.Contains(testhelpers.SearchResourceOption{
			Name: i.Name,
			Kind: i.Kind,
		}))
	}
}

var albIngressResources = []testhelpers.SearchResourceOption{
	{
		Name: "release-name-aws-alb-ingress-controller",
		Kind: "ServiceAccount",
	},
	{
		Name: "release-name-aws-alb-ingress-controller",
		Kind: "ClusterRole",
	},
	{
		Name: "release-name-aws-alb-ingress-controller",
		Kind: "ClusterRoleBinding",
	},
	{
		Name: "release-name-aws-alb-ingress-controller",
		Kind: "Deployment",
	},
}
