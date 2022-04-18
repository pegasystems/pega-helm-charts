package backingservices

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_shouldNotContainC11NMessagingResourcesWhenDisabled(t *testing.T) {
	helmChartParser := C11NHelmConfigParser(
		NewHelmTest(t, helmChartRelativePath, map[string]string{
			"c11n-messaging.enabled": "false",
			"c11n-messaging.deploymentName": "test-c11n-messaging",
		}),
	)

	for _, i := range c11nMessagingResources {
		require.False(t, helmChartParser.Contains(SearchResourceOption{
			Name: i.Name,
			Kind: i.Kind,
		}))
	}
}

func Test_shouldContainC11NMessagingResourcesWhenEnabled(t *testing.T) {
	helmChartParser := C11NHelmConfigParser(
		NewHelmTest(t, helmChartRelativePath, map[string]string{
			"c11n-messaging.deploymentName": "test-c11n-messaging",
		}),
	)

	for _, i := range c11nMessagingResources {
		require.True(t, helmChartParser.Contains(SearchResourceOption{
			Name: i.Name,
			Kind: i.Kind,
		}))
	}
}

var c11nMessagingResources = []SearchResourceOption{
	{
		Name: "test-c11n-messaging",
		Kind: "Deployment",
	},
	{
		Name: "test-c11n-messaging",
		Kind: "Service",
	},
	{
		Name: "test-c11n-messaging-reg-secret",
		Kind: "Secret",
	},
}
