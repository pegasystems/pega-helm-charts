package backingservices

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	k8score "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func TestSRSDeployment(t *testing.T) {

	helmChartParser := NewHelmConfigParser(
		NewHelmTestFromTemplate(t, helmChartRelativePath, map[string]string{
			"srs.enabled":                                "true",
			"srs.deploymentName":                         "test-srs",
			"global.imageCredentials.registry":           "docker-registry.io",
			"srs.srsRuntime.replicaCount":                "1",
			"srs.srsRuntime.srsImage":                    "platform-services/search-n-reporting-service:latest",
			"srs.srsRuntime.env.AuthEnabled":             "false",
			"srs.srsRuntime.env.OAuthPublicKeyURL":       "",
			"srs.srsStorage.tls.enabled":                 "true",
			"srs.srsStorage.basicAuthentication.enabled": "false",
		},
			[]string{"charts/srs/templates/srsservice_deployment.yaml"}),
	)

	var srsDeploymentObj appsv1.Deployment
	helmChartParser.getResourceYAML(SearchResourceOption{
		Name: "test-srs",
		Kind: "Deployment",
	}, &srsDeploymentObj)
	VerifySRSDeployment(t, srsDeploymentObj,
		srsDeployment{
			"test-srs",
			"srs-service",
			int32(1),
			"platform-services/search-n-reporting-service:latest",
			"false",
			"",
			false,
			podResources{"1300m", "2Gi", "650m", "2Gi"},
			esDomain{
				domain:   "elasticsearch-master.default.svc",
				port:     "9200",
				protocol: "https",
			},
			false,
		})
}

func TestSRSDeploymentVariables(t *testing.T) {

	helmChartParser := NewHelmConfigParser(
		NewHelmTestFromTemplate(t, helmChartRelativePath, map[string]string{
			"srs.enabled":                                "true",
			"srs.deploymentName":                         "test-srs-dev",
			"global.imageCredentials.registry":           "docker-registry.io",
			"srs.srsRuntime.replicaCount":                "3",
			"srs.srsRuntime.srsImage":                    "platform-services/search-n-reporting-service:1.0.0",
			"srs.srsRuntime.imagePullSecretNames":        "{secret1, secret2}",
			"srs.srsRuntime.env.AuthEnabled":             "true",
			"srs.srsRuntime.env.OAuthPublicKeyURL":       "https://acme.authenticator.com/OAuthPublicKeyURL",
			"srs.srsRuntime.resources.limits.cpu":        "2",
			"srs.srsRuntime.resources.limits.memory":     "4Gi",
			"srs.srsRuntime.resources.requests.cpu":      "1",
			"srs.srsRuntime.resources.requests.memory":   "2Gi",
			"srs.srsStorage.provisionInternalESCluster":  "false",
			"srs.srsStorage.tls.enabled":                 "false",
			"srs.srsStorage.domain":                      "es-id.managed.cloudiest.io",
			"srs.srsStorage.port":                        "443",
			"srs.srsStorage.protocol":                    "https",
			"srs.srsStorage.awsIAM.region":               "us-east-1",
			"srs.srsStorage.requireInternetAccess":       "true",
			"srs.srsStorage.basicAuthentication.enabled": "false",
		},
			[]string{"charts/srs/templates/srsservice_deployment.yaml"}),
	)

	var srsDeploymentObj appsv1.Deployment
	helmChartParser.getResourceYAML(SearchResourceOption{
		Name: "test-srs-dev",
		Kind: "Deployment",
	}, &srsDeploymentObj)
	VerifySRSDeployment(t, srsDeploymentObj,
		srsDeployment{
			"test-srs-dev",
			"srs-service",
			int32(3),
			"platform-services/search-n-reporting-service:1.0.0",
			"true",
			"https://acme.authenticator.com/OAuthPublicKeyURL",
			true,
			podResources{"2", "4Gi", "1", "2Gi"},
			esDomain{
				domain:   "es-id.managed.cloudiest.io",
				port:     "443",
				protocol: "https",
				region:   "us-east-1",
			},
			true,
		})
}

func TestSRSDeploymentVariablesDefaultInternetEgress(t *testing.T) {

	helmChartParser := NewHelmConfigParser(
		NewHelmTestFromTemplate(t, helmChartRelativePath, map[string]string{
			"srs.enabled":                                "true",
			"srs.deploymentName":                         "test-srs-dev",
			"global.imageCredentials.registry":           "docker-registry.io",
			"srs.srsRuntime.replicaCount":                "3",
			"srs.srsRuntime.srsImage":                    "platform-services/search-n-reporting-service:1.0.0",
			"srs.srsRuntime.imagePullSecretNames":        "{secret1, secret2}",
			"srs.srsRuntime.env.AuthEnabled":             "true",
			"srs.srsRuntime.env.OAuthPublicKeyURL":       "https://acme.authenticator.com/OAuthPublicKeyURL",
			"srs.srsRuntime.resources.limits.cpu":        "2",
			"srs.srsRuntime.resources.limits.memory":     "4Gi",
			"srs.srsRuntime.resources.requests.cpu":      "1",
			"srs.srsRuntime.resources.requests.memory":   "2Gi",
			"srs.srsStorage.provisionInternalESCluster":  "false",
			"srs.srsStorage.domain":                      "es-id.managed.cloudiest.io",
			"srs.srsStorage.port":                        "443",
			"srs.srsStorage.protocol":                    "https",
			"srs.srsStorage.tls.enabled":                 "false",
			"srs.srsStorage.basicAuthentication.enabled": "false",
		},
			[]string{"charts/srs/templates/srsservice_deployment.yaml"}),
	)

	var srsDeploymentObj appsv1.Deployment
	helmChartParser.getResourceYAML(SearchResourceOption{
		Name: "test-srs-dev",
		Kind: "Deployment",
	}, &srsDeploymentObj)
	VerifySRSDeployment(t, srsDeploymentObj,
		srsDeployment{
			"test-srs-dev",
			"srs-service",
			int32(3),
			"platform-services/search-n-reporting-service:1.0.0",
			"true",
			"https://acme.authenticator.com/OAuthPublicKeyURL",
			false,
			podResources{"2", "4Gi", "1", "2Gi"},
			esDomain{
				domain:   "es-id.managed.cloudiest.io",
				port:     "443",
				protocol: "https",
			},
			true,
		})
}

func TestSRSDeploymentWithAffinity(t *testing.T) {

	var affintiyBasePath = "srs.srsRuntime.affinity.nodeAffinity.requiredDuringSchedulingIgnoredDuringExecution.nodeSelectorTerms[0].matchExpressions[0]."

	helmChartParser := NewHelmConfigParser(
		NewHelmTestFromTemplate(t, helmChartRelativePath, map[string]string{
			"srs.enabled":                                "true",
			"srs.deploymentName":                         "test-srs",
			"global.imageCredentials.registry":           "docker-registry.io",
			"srs.srsRuntime.replicaCount":                "1",
			"srs.srsRuntime.srsImage":                    "platform-services/search-n-reporting-service:latest",
			"srs.srsRuntime.env.AuthEnabled":             "false",
			"srs.srsRuntime.env.OAuthPublicKeyURL":       "",
			"srs.srsStorage.tls.enabled":                 "true",
			"srs.srsStorage.basicAuthentication.enabled": "false",
			affintiyBasePath + "key":                     "kubernetes.io/os",
			affintiyBasePath + "operator":                "In",
			affintiyBasePath + "values[0]":               "linux",
		},
			[]string{"charts/srs/templates/srsservice_deployment.yaml"}),
	)

	var srsDeploymentObj appsv1.Deployment
	helmChartParser.getResourceYAML(SearchResourceOption{
		Name: "test-srs",
		Kind: "Deployment",
	}, &srsDeploymentObj)

	deploymentSpec := srsDeploymentObj.Spec.Template.Spec
	require.Equal(t, deploymentSpec.Containers[0].Name, "srs-service")
	require.Equal(t, deploymentSpec.Containers[0].Image, "platform-services/search-n-reporting-service:latest")
	require.Equal(t, deploymentSpec.Containers[0].Ports[0].Name, "srs-port")
	deploymentAffinity := deploymentSpec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution
	require.Equal(t, "kubernetes.io/os", deploymentAffinity.NodeSelectorTerms[0].MatchExpressions[0].Key)
	require.Equal(t, "In", string(deploymentAffinity.NodeSelectorTerms[0].MatchExpressions[0].Operator))
	require.Equal(t, "linux", deploymentAffinity.NodeSelectorTerms[0].MatchExpressions[0].Values[0])
}

func VerifySRSDeployment(t *testing.T, deploymentObj appsv1.Deployment, expectedDeployment srsDeployment) {
	require.Equal(t, expectedDeployment.replicaCount, *deploymentObj.Spec.Replicas)
	require.Equal(t, expectedDeployment.appName, deploymentObj.Spec.Selector.MatchLabels["app.kubernetes.io/name"])
	if expectedDeployment.internetEgress {
		require.Equal(t, "true", deploymentObj.Spec.Selector.MatchLabels["networking/allow-internet-egress"])
	}
	require.Equal(t, expectedDeployment.appName, deploymentObj.Spec.Template.Labels["app.kubernetes.io/name"])
	deploymentSpec := deploymentObj.Spec.Template.Spec
	VerifyDeployment(t, &deploymentSpec, expectedDeployment)
}

// VerifyDeployment - Performs common srs deployment assertions with the values as provided in default values.yaml
func VerifyDeployment(t *testing.T, pod *k8score.PodSpec, expectedSpec srsDeployment) {
	require.Equal(t, pod.Containers[0].Name, "srs-service")
	require.Equal(t, pod.Containers[0].Image, expectedSpec.imageURI)
	require.Equal(t, pod.Containers[0].Ports[0].Name, "srs-port")
	require.Equal(t, pod.Containers[0].Ports[0].ContainerPort, int32(8080))
	var envIndex int32 = 0
	require.Equal(t, "ELASTICSEARCH_HOST", pod.Containers[0].Env[envIndex].Name)
	require.Equal(t, expectedSpec.elasticsearchEndPoint.domain, pod.Containers[0].Env[envIndex].Value)
	envIndex++
	require.Equal(t, "ELASTICSEARCH_PORT", pod.Containers[0].Env[envIndex].Name)
	require.Equal(t, expectedSpec.elasticsearchEndPoint.port, pod.Containers[0].Env[envIndex].Value)
	envIndex++
	require.Equal(t, "ELASTICSEARCH_PROTO", pod.Containers[0].Env[envIndex].Name)
	require.Equal(t, expectedSpec.elasticsearchEndPoint.protocol, pod.Containers[0].Env[envIndex].Value)
	envIndex++
	require.Equal(t, "ELASTICSEARCH_AUTH_PROVIDER", pod.Containers[0].Env[envIndex].Name)
	require.NotEmpty(t, pod.Containers[0].Env[envIndex].Value)
	var authProvider string = pod.Containers[0].Env[envIndex].Value
	envIndex++
	if strings.EqualFold("aws-iam", authProvider) {
		require.Equal(t, "ELASTICSEARCH_REGION", pod.Containers[0].Env[envIndex].Name)
		require.Equal(t, expectedSpec.elasticsearchEndPoint.region, pod.Containers[0].Env[envIndex].Value)
		envIndex++
	}
	if strings.EqualFold("basic-authentication", authProvider) {
		require.Equal(t, "ELASTICSEARCH_USERNAME", pod.Containers[0].Env[envIndex].Name)
		require.Equal(t, "srs-elastic-credentials", pod.Containers[0].Env[envIndex].ValueFrom.SecretKeyRef.Name)
		require.Equal(t, "username", pod.Containers[0].Env[envIndex].ValueFrom.SecretKeyRef.Key)
		envIndex++
		require.Equal(t, "ELASTICSEARCH_PASSWORD", pod.Containers[0].Env[envIndex].Name)
		require.Equal(t, "srs-elastic-credentials", pod.Containers[0].Env[envIndex].ValueFrom.SecretKeyRef.Name)
		require.Equal(t, "password", pod.Containers[0].Env[envIndex].ValueFrom.SecretKeyRef.Key)
		envIndex++
	}
	if strings.EqualFold("tls", authProvider) {
		require.Equal(t, "ELASTICSEARCH_USERNAME", pod.Containers[0].Env[envIndex].Name)
		require.Equal(t, "srs-elastic-credentials", pod.Containers[0].Env[envIndex].ValueFrom.SecretKeyRef.Name)
		require.Equal(t, "username", pod.Containers[0].Env[envIndex].ValueFrom.SecretKeyRef.Key)
		envIndex++
		require.Equal(t, "ELASTICSEARCH_PASSWORD", pod.Containers[0].Env[envIndex].Name)
		require.Equal(t, "srs-elastic-credentials", pod.Containers[0].Env[envIndex].ValueFrom.SecretKeyRef.Name)
		require.Equal(t, "password", pod.Containers[0].Env[envIndex].ValueFrom.SecretKeyRef.Key)
		envIndex++
		require.Equal(t, "PATH_TO_TRUSTSTORE", pod.Containers[0].Env[envIndex].Name)
		require.Equal(t, "/usr/share/elastic-certificates.p12", pod.Containers[0].Env[envIndex].Value)
		envIndex++
		require.Equal(t, "PATH_TO_KEYSTORE", pod.Containers[0].Env[envIndex].Name)
		require.Equal(t, "", pod.Containers[0].Env[envIndex].Value)
		envIndex++
	}
	require.Equal(t, "APPLICATION_HOST", pod.Containers[0].Env[envIndex].Name)
	require.Equal(t, "0.0.0.0", pod.Containers[0].Env[envIndex].Value)
	envIndex++
	require.Equal(t, "APPLICATION_PORT", pod.Containers[0].Env[envIndex].Name)
	require.Equal(t, "8080", pod.Containers[0].Env[envIndex].Value)
	envIndex++
	require.Equal(t, "AUTH_ENABLED", pod.Containers[0].Env[envIndex].Name)
	require.Equal(t, expectedSpec.authEnabled, pod.Containers[0].Env[envIndex].Value)
	envIndex++
	require.Equal(t, "OAUTH_PUBLIC_KEY_URL", pod.Containers[0].Env[envIndex].Name)
	require.Equal(t, expectedSpec.oauthPublicKeyURL, pod.Containers[0].Env[envIndex].Value)
	envIndex++

	require.Equal(t, expectedSpec.podLimits.cpuLimit, pod.Containers[0].Resources.Limits.Cpu().String())
	require.Equal(t, expectedSpec.podLimits.memoryLimit, pod.Containers[0].Resources.Limits.Memory().String())
	require.Equal(t, expectedSpec.podLimits.cpuRequest, pod.Containers[0].Resources.Requests.Cpu().String())
	require.Equal(t, expectedSpec.podLimits.memoryRequest, pod.Containers[0].Resources.Requests.Memory().String())

	require.Equal(t, pod.Containers[0].ReadinessProbe.InitialDelaySeconds, int32(2))
	require.Equal(t, pod.Containers[0].ReadinessProbe.TimeoutSeconds, int32(30))
	require.Equal(t, pod.Containers[0].ReadinessProbe.PeriodSeconds, int32(5))
	require.Equal(t, pod.Containers[0].ReadinessProbe.HTTPGet.Path, "/health")
	require.Equal(t, pod.Containers[0].ReadinessProbe.HTTPGet.Port, intstr.FromString("srs-port"))
	require.Equal(t, pod.Containers[0].ReadinessProbe.HTTPGet.Scheme, k8score.URIScheme("HTTP"))

	require.Equal(t, pod.ImagePullSecrets[0].Name, expectedSpec.name+"-reg-secret")
	if expectedSpec.imagePullSecretNames {
		require.Equal(t, pod.ImagePullSecrets[1].Name, "secret1")
		require.Equal(t, pod.ImagePullSecrets[2].Name, "secret2")
	}

	podAffinity := pod.Affinity
	require.Empty(t, podAffinity)
}

type srsDeployment struct {
	name                  string
	appName               string
	replicaCount          int32
	imageURI              string
	authEnabled           string
	oauthPublicKeyURL     string
	internetEgress        bool
	podLimits             podResources
	elasticsearchEndPoint esDomain
	imagePullSecretNames  bool
}

type podResources struct {
	cpuLimit      string
	memoryLimit   string
	cpuRequest    string
	memoryRequest string
}

type esDomain struct {
	domain   string
	port     string
	protocol string
	region   string
}
