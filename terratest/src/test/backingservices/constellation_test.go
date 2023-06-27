package backingservices

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_shouldNotContainConstellationResourcesWhenDisabled(t *testing.T) {
	helmChartParser := NewHelmConfigParser(
		NewHelmTest(t, helmChartRelativePath, map[string]string{
			"constellation.enabled": "false",
		}),
	)

	for _, i := range constellationResources {
		require.False(t, helmChartParser.Contains(SearchResourceOption{
			Name: i.Name,
			Kind: i.Kind,
		}))
	}
}

func Test_shouldContainConstellationResourcesWhenEnabled(t *testing.T) {
	helmChartParser := NewHelmConfigParser(
		NewHelmTest(t, helmChartRelativePath, map[string]string{
			"constellation.enabled": "true",
		}),
	)

	for _, i := range constellationResources {
		require.True(t, helmChartParser.Contains(SearchResourceOption{
			Name: i.Name,
			Kind: i.Kind,
		}))
	}
}

var constellationResources = []SearchResourceOption{
	{
		Name: "constellation",
		Kind: "Deployment",
	},
	{
		Name: "constellation",
		Kind: "Service",
	},
	{
		Name: "constellationingress",
		Kind: "Ingress",
	},
	{
		Name: "constellation-registry-secret",
		Kind: "Secret",
	},
}
