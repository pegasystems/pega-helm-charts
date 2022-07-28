package backingservices

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_shouldNotContainSRSResourcesWhenDisabled(t *testing.T) {
	helmChartParser := NewHelmConfigParser(
		NewHelmTest(t, helmChartRelativePath, map[string]string{
			"srs.enabled": "false",
			"srs.srsStorage.provisionInternalESCluster": "false",
			"srs.srsStorage.tls.enabled": "false",
			"srs.deploymentName": "test-srs",
			"srs.srsStorage.basicAuthentication.enabled": "false",
		}),
	)

	for _, i := range srsResources {
		require.False(t, helmChartParser.Contains(SearchResourceOption{
			Name: i.Name,
			Kind: i.Kind,
		}))
	}

	for _, i := range esResources {
		require.False(t, helmChartParser.Contains(SearchResourceOption{
			Name: i.Name,
			Kind: i.Kind,
		}))
	}
}

func Test_shouldContainSRSResourcesWhenEnabled(t *testing.T) {
	helmChartParser := NewHelmConfigParser(
		NewHelmTest(t, helmChartRelativePath, map[string]string{
			"srs.deploymentName": "test-srs",
			"srs.srsStorage.provisionInternalESCluster": "true",
			"srs.srsStorage.tls.enabled": "true",
			"srs.srsStorage.basicAuthentication.enabled": "false",
		}),
	)

	for _, i := range srsResources {
		require.True(t, helmChartParser.Contains(SearchResourceOption{
			Name: i.Name,
			Kind: i.Kind,
		}))
	}

	for _, i := range esResources {
		require.True(t, helmChartParser.Contains(SearchResourceOption{
			Name: i.Name,
			Kind: i.Kind,
		}))
	}
}

func Test_shouldContainSRSandESResourcesWhenEnabled(t *testing.T) {
	helmChartParser := NewHelmConfigParser(
		NewHelmTest(t, helmChartRelativePath, map[string]string{
			"srs.deploymentName": "test-srs",
			"srs.srsStorage.provisionInternalESCluster": "true",
			"srs.srsStorage.tls.enabled": "false",
		}),
	)

	for _, i := range srsResources {
		require.True(t, helmChartParser.Contains(SearchResourceOption{
			Name: i.Name,
			Kind: i.Kind,
		}))
	}

	for _, i := range esResources {
		require.True(t, helmChartParser.Contains(SearchResourceOption{
			Name: i.Name,
			Kind: i.Kind,
		}))
	}
}

func Test_shouldContainSRSWhenEnabledandNotESResourcesWhenDisabled(t *testing.T) {
	helmChartParser := NewHelmConfigParser(
		NewHelmTest(t, helmChartRelativePath, map[string]string{
			"srs.deploymentName": "test-srs",
			"srs.srsStorage.provisionInternalESCluster": "false",
			"srs.srsStorage.domain": "es.managed.io",
			"srs.srsStorage.port": "9200",
			"srs.srsStorage.protocol": "https",
			"srs.srsStorage.tls.enabled": "false",
			"srs.srsStorage.basicAuthentication.enabled": "false",
		}),
	)

	for _, i := range srsResources {
		require.True(t, helmChartParser.Contains(SearchResourceOption{
			Name: i.Name,
			Kind: i.Kind,
		}))
	}

	for _, i := range esResources {
		require.False(t, helmChartParser.Contains(SearchResourceOption{
			Name: i.Name,
			Kind: i.Kind,
		}))
	}
}

var srsResources = []SearchResourceOption{
	{
		Name: "test-srs",
		Kind: "Deployment",
	},
	{
		Name: "test-srs",
		Kind: "Service",
	},
	{
		Name: "test-srs",
		Kind: "PodDisruptionBudget",
	},
	{
		Name: "test-srs-networkpolicy",
		Kind: "NetworkPolicy",
	},
	{
		Name: "test-srs-reg-secret",
		Kind: "Secret",
	},
}

var esResources = []SearchResourceOption{
	{
		Name: "elasticsearch-master",
		Kind: "Service",
	},
	{
		Name: "elasticsearch-master-headless",
		Kind: "Service",
	},
	{
		Name: "elasticsearch-master",
		Kind: "StatefulSet",
	},
	{
		Name: "elasticsearch-master-pdb",
		Kind: "PodDisruptionBudget",
	},
}
