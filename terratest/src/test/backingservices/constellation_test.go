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

func Test_shouldContainConstellationMessagingWhenEnabled(t *testing.T) {
	helmChartParser := NewHelmConfigParser(
		NewHelmTest(t, helmChartRelativePath, map[string]string{
			"srs.enabled": "false",
			"constellation-messaging.enabled": "true",
			"constellation-messaging.name": "constellation-messaging",
			"constellation-messaging.image": "messaging-image:1.0.0",
			"constellation-messaging.replicas": "3",
			"constellation-messaging.imagePullSecretNames": "{secret1, secret2}",
			"constellation-messaging.ingress.domain": "test.com",
		}),
	)

	for _, i := range constellationMessagingResources {
		require.True(t, helmChartParser.Contains(SearchResourceOption{
			Name: i.Name,
			Kind: i.Kind,
		}))
	}
}

func Test_shouldNotContainConstellationMessagingWhenDisabled(t *testing.T) {
	helmChartParser := NewHelmConfigParser(
		NewHelmTest(t, helmChartRelativePath, map[string]string{
			"srs.enabled": "false",
			"constellation-messaging.enabled": "false",
			"constellation-messaging.name": "constellation-messaging",
			"constellation-messaging.image": "messaging-image:1.0.0",
			"constellation-messaging.replicas": "3",
			"constellation-messaging.imagePullSecretNames": "{secret1, secret2}",
			"constellation-messaging.ingress.domain": "test.com",
		}),
	)

	for _, i := range constellationMessagingResources {
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

var constellationMessagingResources = []SearchResourceOption{
	{
		Name: "constellation-messaging",
		Kind: "Deployment",
	},
	{
		Name: "constellation-messaging",
		Kind: "Service",
	},
	{
		Name: "constellation-messaging",
		Kind: "Ingress",
	},
}
