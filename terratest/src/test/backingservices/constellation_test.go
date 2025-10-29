package backingservices

import (
	"testing"

	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
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

func Test_shouldContainConstellationResourcesWithIngressWhenEnabled(t *testing.T) {
	helmChartParser := NewHelmConfigParser(
		NewHelmTest(t, helmChartRelativePath, map[string]string{
			"constellation.enabled":         "true",
			"constellation.ingress.enabled": "true",
		}),
	)

	for _, i := range constellationResourcesWithIngress {
		require.True(t, helmChartParser.Contains(SearchResourceOption{
			Name: i.Name,
			Kind: i.Kind,
		}))
	}
}

func Test_shouldContainConstellationMessagingWhenEnabled(t *testing.T) {
	helmChartParser := NewHelmConfigParser(
		NewHelmTest(t, helmChartRelativePath, map[string]string{
			"srs.enabled":                                         "false",
			"constellation-messaging.enabled":                     "true",
			"constellation-messaging.deployment.name":             "constellation-messaging",
			"constellation-messaging.docker.messaging.image":      "messaging-image:1.0.0",
			"constellation-messaging.replicas":                    "3",
			"constellation-messaging.docker.imagePullSecretNames": "{secret1, secret2}",
			"constellation-messaging.ingress.enabled":             "true",
			"constellation-messaging.ingress.domain":              "test.com",
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
			"srs.enabled":                                  "false",
			"constellation-messaging.enabled":              "false",
			"constellation-messaging.name":                 "constellation-messaging",
			"constellation-messaging.image":                "messaging-image:1.0.0",
			"constellation-messaging.replicas":             "3",
			"constellation-messaging.imagePullSecretNames": "{secret1, secret2}",
			"constellation-messaging.ingress.domain":       "test.com",
		}),
	)

	for _, i := range constellationMessagingResources {
		require.False(t, helmChartParser.Contains(SearchResourceOption{
			Name: i.Name,
			Kind: i.Kind,
		}))
	}
}

func Test_ConstellationMessagingWithLabels(t *testing.T) {

	var deploymentName string = "constellation-msg"

	helmChartParser := NewHelmConfigParser(
		NewHelmTest(t, helmChartRelativePath, map[string]string{
			"constellation-messaging.enabled":                "true",
			"constellation-messaging.deployment.name":        deploymentName,
			"constellation-messaging.deployment.labels.key1": "value1",
			"constellation-messaging.podLabels.podKey1":      "podValue1",
		}),
	)

	var cllnMsgDeployment appsv1.Deployment
	helmChartParser.getResourceYAML(SearchResourceOption{
		Name: deploymentName,
		Kind: "Deployment",
	}, &cllnMsgDeployment)

	require.Equal(t, cllnMsgDeployment.Name, deploymentName)
	require.Equal(t, cllnMsgDeployment.Labels["key1"], "value1")
	require.Equal(t, cllnMsgDeployment.Spec.Template.Labels["podKey1"], "podValue1")
	require.Equal(t, cllnMsgDeployment.Labels["app"], deploymentName)
	require.Equal(t, cllnMsgDeployment.Spec.Template.Labels["app"], deploymentName)
}

func Test_ConstellationWithLabels(t *testing.T) {

	var deploymentName string = "constellation-static"

	helmChartParser := NewHelmConfigParser(
		NewHelmTest(t, helmChartRelativePath, map[string]string{
			"constellation.enabled":                "true",
			"constellation.deployment.name":        deploymentName,
			"constellation.deployment.labels.key1": "value1",
			"constellation.podLabels.podKey1":      "podValue1",
		}),
	)

	var cllnDeployment appsv1.Deployment
	helmChartParser.getResourceYAML(SearchResourceOption{
		Name: deploymentName,
		Kind: "Deployment",
	}, &cllnDeployment)

	require.Equal(t, cllnDeployment.Name, deploymentName)
	require.Equal(t, cllnDeployment.Labels["key1"], "value1")
	require.Equal(t, cllnDeployment.Spec.Template.Labels["podKey1"], "podValue1")
	require.Equal(t, cllnDeployment.Labels["app"], deploymentName)
	require.Equal(t, cllnDeployment.Spec.Template.Labels["app"], deploymentName)
}

// Test case to verify pod annotations
func Test_ConstellationWithPodAnnotations(t *testing.T) {

	var deploymentName string = "constellation-static"

	helmChartParser := NewHelmConfigParser(
		NewHelmTest(t, helmChartRelativePath, map[string]string{
			"constellation.enabled":                "true",
			"constellation.deployment.name":        deploymentName,
			"constellation.podAnnotations.podAnnotationKey1":      "podAnnotationValue1",
		}),
	)

	var cllnDeployment appsv1.Deployment
	helmChartParser.getResourceYAML(SearchResourceOption{
		Name: deploymentName,
		Kind: "Deployment",
	}, &cllnDeployment)

	require.Equal(t, cllnDeployment.Name, deploymentName)
	require.Equal(t, cllnDeployment.Spec.Template.Annotations["podAnnotationKey1"], "podAnnotationValue1")
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
		Name: "constellation-registry-secret",
		Kind: "Secret",
	},
}

var constellationResourcesWithIngress = []SearchResourceOption{
	{
		Name: "constellation",
		Kind: "Deployment",
	},
	{
		Name: "constellation",
		Kind: "Service",
	},
	{
		Name: "constellation",
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
