package addons

import (
	"github.com/stretchr/testify/require"
	"test/testhelpers"
	"testing"
)

func TestShouldNotContainAzureIngressIfDisabled(t *testing.T) {
	helmChartParser := testhelpers.NewHelmConfigParser(
		testhelpers.NewHelmTest(t, helmChartRelativePath, map[string]string{
			"ingress-azure.enabled": "false",
		}),
	)

	for _, i := range azureIngressResources {
		require.False(t, helmChartParser.Contains(testhelpers.SearchResourceOption{
			Name: i.Name,
			Kind: i.Kind,
		}))
	}
}

func TestAzureIngressShouldContainAllResources(t *testing.T) {
	helmChartParser := testhelpers.NewHelmConfigParser(
		testhelpers.NewHelmTest(t, helmChartRelativePath, map[string]string{
			"ingress-azure.enabled": "true",
		}),
	)

	for _, i := range azureIngressResources {
		require.True(t, helmChartParser.Contains(testhelpers.SearchResourceOption{
			Name: i.Name,
			Kind: i.Kind,
		}))
	}
}

var azureIngressResources = []testhelpers.SearchResourceOption{
	{
		Name: "networking-appgw-k8s-azure-service-principal",
		Kind: "Secret",
	},
	{
		Name: "release-name-cm-ingress-azure",
		Kind: "ConfigMap",
	},
	{
		Name: "release-name-sa-ingress-azure",
		Kind: "ServiceAccount",
	},
	{
		Name: "release-name-ingress-azure",
		Kind: "ClusterRole",
	},
	{
		Name: "release-name-ingress-azure",
		Kind: "ClusterRoleBinding",
	},
	{
		Name: "release-name-ingress-azure",
		Kind: "Deployment",
	},
}
