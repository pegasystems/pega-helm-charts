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

func Test_shouldContainConstellationResourcesWithOptionalValues(t *testing.T) {
	helmChartParser := NewHelmConfigParser(
		NewHelmTest(t, helmChartRelativePath, map[string]string{
			"constellation.enabled": "true",
			"constellation.customerAssetVolumeClaimName": "claim-name",
			"constellation.ingressAnnotations": "annotation",
			"constellation.imagePullSecretNames": "{secret1, secret2}",
		}),
	)

	for _, i := range constellationResources {
		require.True(t, helmChartParser.Contains(SearchResourceOption{
			Name: i.Name,
			Kind: i.Kind,
		}))
	}
}

func Test_shouldContainConstellationCdnWhenEnabled(t *testing.T) {
	helmChartParser := NewHelmConfigParser(
		NewHelmTest(t, helmChartRelativePath, map[string]string{
			"constellation.constellation-cdn.enabled": "true",
			"constellation.constellation-cdn.name": "constellation-cdn",
			"constellation.constellation-cdn.imagePullSecretNames": "{secret1, secret2}",
			"constellation.constellation-cdn.ingress.domain": "test.com",
		}),
	)

	for _, i := range constellationCdnResources {
		require.True(t, helmChartParser.Contains(SearchResourceOption{
			Name: i.Name,
			Kind: i.Kind,
		}))
	}
}

func Test_shouldNotContainConstellationCdnWhenDisabled(t *testing.T) {
	helmChartParser := NewHelmConfigParser(
		NewHelmTest(t, helmChartRelativePath, map[string]string{
			"constellation.constellation-cdn.enabled": "false",
			"constellation.constellation-cdn.name": "constellation-cdn",
			"constellation.constellation-cdn.imagePullSecretNames": "{secret1, secret2}",
			"constellation.constellation-cdn.ingress.domain": "test.com",
		}),
	)

	for _, i := range constellationCdnResources {
		require.False(t, helmChartParser.Contains(SearchResourceOption{
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

var constellationCdnResources = []SearchResourceOption{
	{
		Name: "constellation-cdn-881",
		Kind: "Deployment",
	},
	{
		Name: "constellation-cdn-881",
		Kind: "Service",
	},
	{
		Name: "constellation-cdn",
		Kind: "Ingress",
	},
}
