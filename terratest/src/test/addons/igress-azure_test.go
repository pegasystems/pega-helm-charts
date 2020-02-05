package addons

import (
	"github.com/stretchr/testify/require"
	"test/testhelpers"
	"testing"
)

func Test_shouldNotContainIngressAzureIfDisabled(t *testing.T) {
	helmChartParser := testhelpers.NewHelmConfigParser(
		testhelpers.NewHelmTest(t, helmChartRelativePath, map[string]string{
			"ingress-azure.enabled": "false",
		}),
	)

	for _, i := range ingressAzureResources {
		require.False(t, helmChartParser.Contains(testhelpers.SearchResourceOption{
			Name: i.Name,
			Kind: i.Kind,
		}))
	}
}

func Test_shouldContainIngressAzureIfEnabled(t *testing.T) {
	helmChartParser := testhelpers.NewHelmConfigParser(
		testhelpers.NewHelmTest(t, helmChartRelativePath, map[string]string{
			"ingress-azure.enabled": "true",
		}),
	)

	for _, i := range ingressAzureResources {
		require.True(t, helmChartParser.Contains(testhelpers.SearchResourceOption{
			Name: i.Name,
			Kind: i.Kind,
		}))
	}
}

func Test_shouldBeAbleToSetUpValues(t *testing.T) {
	//helmChartParser := NewHelmConfigParser(
	//	NewHelmTest(t, helmChartRelativePath, map[string]string{
	//		"ingress-azure.enabled": "false",
	//	}),
	//)

}

var ingressAzureResources = []testhelpers.SearchResourceOption{
	{
		Name: "release-name-ingress-azure",
		Kind: "Deployment",
	},
	{
		Name: "release-name-ingress-azure",
		Kind: "ClusterRoleBinding",
	},
}
