package backingservices

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// func Test_shouldNotContainC11NMessagingResourcesWhenDisabled(t *testing.T) {
// 	helmChartParser := C11NHelmConfigParser(
// 		NewHelmTest(t, helmChartRelativePath, map[string]string{
// 			"enabled":        "false",
// 			"deploymentName": "c11n-messaging",
// 		}),
// 	)

// 	for _, i := range c11nMessagingResources {
// 		require.False(t, helmChartParser.Contains(SearchResourceOption{
// 			Name: i.Name,
// 			Kind: i.Kind,
// 		}))
// 	}
// }

func Test_shouldContainC11NMessagingResourcesWhenEnabled(t *testing.T) {
	helmChartParser := C11NHelmConfigParser(
		NewHelmTest(t, helmChartRelativePath, map[string]string{
			"c11n-messaging.enabled":        "true",
			"c11n-messaging.deploymentName": "c11n-messaging",
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
		Name: "c11n-messaging",
		Kind: "Deployment",
	},
	{
		Name: "c11n-messaging",
		Kind: "Service",
	},
	{
		Name: "c11n-messaging-reg-secret",
		Kind: "Secret",
	},
}
