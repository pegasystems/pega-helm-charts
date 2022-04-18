package backingservices

import (
	"testing"

	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	k8score "k8s.io/api/core/v1"
)

func TestC11NDeployment(t *testing.T) {

	helmChartParser := C11NHelmConfigParser(
		NewHelmTestFromTemplate(t, helmChartRelativePath, map[string]string{
			"c11n-messaging.enabled":                                     "true",
			"c11n-messaging.deploymentName":                              "test-c11n-messaging",
			"global.imageCredentials.registry":                           "docker-registry.io",
			"c11n-messaging.c11n-messaging-Runtime.replicaCount":         "1",
			"c11n-messaging.c11n-messaging-Runtime.c11n-messaging-image": "pega-docker.downloads.pega.com/constellation-messaging/docker-image:0.0.3-20220112124450309",
		},
			[]string{"charts/constellation-messaging/templates/c11n-messaging-deployment.yaml"}),
	)

	var c11nMessagingDeploymentObj appsv1.Deployment
	helmChartParser.getResourceYAML(SearchResourceOption{
		Name: "test-c11n-messaging",
		Kind: "Deployment",
	}, &c11nMessagingDeploymentObj)
	VerifyC11NMessagingDeployment(t, c11nMessagingDeploymentObj,
		c11nMessagingDeployment{
			"test-c11n-messaging",
			"c11n-messaging-service",
			int32(1),
			"pega-docker.downloads.pega.com/constellation-messaging/docker-image:0.0.3-20220112124450309",
			c11nMessagingpodResources{"1300m", "2Gi", "650m", "2Gi"},
		})
}

func TestC11NMessagingDeploymentVariables(t *testing.T) {

	helmChartParser := C11NHelmConfigParser(
		NewHelmTestFromTemplate(t, helmChartRelativePath, map[string]string{
			"c11n-messaging.enabled":                                          "true",
			"c11n-messaging.deploymentName":                                   "test-c11n-messaging",
			"global.imageCredentials.registry":                                "docker-registry.io",
			"c11n-messaging.c11n-messaging-Runtime.replicaCount":              "3",
			"c11n-messaging.c11n-messaging-Runtime.c11n-messaging-image":      "pega-docker.downloads.pega.com/constellation-messaging/docker-image:0.0.3-20220112124450309",
			"c11n-messaging.c11n-messaging-Runtime.resources.limits.cpu":      "2",
			"c11n-messaging.c11n-messaging-Runtime.resources.limits.memory":   "4Gi",
			"c11n-messaging.c11n-messaging-Runtime.resources.requests.cpu":    "1",
			"c11n-messaging.c11n-messaging-Runtime.resources.requests.memory": "2Gi",
		},
			[]string{"charts/constellation-messaging/templates/c11n-messaging-deployment.yaml"}),
	)

	var c11nMessagingDeploymentObj appsv1.Deployment
	helmChartParser.getResourceYAML(SearchResourceOption{
		Name: "test-c11n-messaging",
		Kind: "Deployment",
	}, &c11nMessagingDeploymentObj)
	VerifyC11NMessagingDeployment(t, c11nMessagingDeploymentObj,
		c11nMessagingDeployment{
			"test-c11n-messaging",
			"c11n-messaging-service",
			int32(3),
			"pega-docker.downloads.pega.com/constellation-messaging/docker-image:0.0.3-20220112124450309",
			c11nMessagingpodResources{"2", "4Gi", "1", "2Gi"},
		})
}

func VerifyC11NMessagingDeployment(t *testing.T, deploymentObj appsv1.Deployment, expectedDeployment c11nMessagingDeployment) {
	require.Equal(t, expectedDeployment.replicaCount, *deploymentObj.Spec.Replicas)
	require.Equal(t, expectedDeployment.appName, deploymentObj.Spec.Selector.MatchLabels["app.kubernetes.io/name"])
	// if expectedDeployment.internetEgress {
	// 	require.Equal(t, "true", deploymentObj.Spec.Selector.MatchLabels["networking/allow-internet-egress"])
	// }
	require.Equal(t, expectedDeployment.appName, deploymentObj.Spec.Template.Labels["app.kubernetes.io/name"])
	deploymentSpec := deploymentObj.Spec.Template.Spec
	VerifyC11NMessagingServiceDeployment(t, &deploymentSpec, expectedDeployment)
}

// VerifyDeployment - Performs common srs deployment assertions with the values as provided in default values.yaml
func VerifyC11NMessagingServiceDeployment(t *testing.T, pod *k8score.PodSpec, expectedSpec c11nMessagingDeployment) {
	require.Equal(t, pod.Containers[0].Name, "c11n-messaging-service")
	require.Equal(t, pod.Containers[0].Image, expectedSpec.imageURI)
	require.Equal(t, pod.Containers[0].Ports[0].Name, "c11n-messaging-port")
	require.Equal(t, pod.Containers[0].Ports[0].ContainerPort, int32(3000))

	var envIndex int32 = 0
	require.Equal(t, "--max-semi-space-size=1024", pod.Containers[0].Env[envIndex])
	envIndex++
	require.Equal(t, "root=/usr/src/app/environments", pod.Containers[0].Env[envIndex])
	envIndex++
	require.Equal(t, "port=3000", pod.Containers[0].Env[envIndex])

	require.Equal(t, expectedSpec.podLimits.cpuLimit, pod.Containers[0].Resources.Limits.Cpu().String())
	require.Equal(t, expectedSpec.podLimits.memoryLimit, pod.Containers[0].Resources.Limits.Memory().String())
	require.Equal(t, expectedSpec.podLimits.cpuRequest, pod.Containers[0].Resources.Requests.Cpu().String())
	require.Equal(t, expectedSpec.podLimits.memoryRequest, pod.Containers[0].Resources.Requests.Memory().String())

	require.Equal(t, pod.ImagePullSecrets[0].Name, expectedSpec.name+"-reg-secret")
}

type c11nMessagingDeployment struct {
	name         string
	appName      string
	replicaCount int32
	imageURI     string
	podLimits    c11nMessagingpodResources
}

type c11nMessagingpodResources struct {
	cpuLimit      string
	memoryLimit   string
	cpuRequest    string
	memoryRequest string
}
