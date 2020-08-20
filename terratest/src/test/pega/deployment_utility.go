package pega

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	k8score "k8s.io/api/core/v1"
	k8sv1beta1 "k8s.io/api/extensions/v1beta1"
	intstr "k8s.io/apimachinery/pkg/util/intstr"
	autoscaling "k8s.io/kubernetes/pkg/apis/autoscaling"
	api "k8s.io/kubernetes/pkg/apis/core"
)

var replicas int32 = 1
var replicasPtr = &replicas
var ProgressDeadlineSeconds int32 = 2147483647
var ProgressDeadlineSecondsPtr = &ProgressDeadlineSeconds
var rollingUpdate intstr.IntOrString = intstr.FromString("25%")
var rollingUpdatePtr = &rollingUpdate
var terminationGracePeriodSeconds int64 = 300
var terminationGracePeriodSecondsPtr = &terminationGracePeriodSeconds

type pegaDeployment struct {
	name               string
	initContainers     []string
	nodeType           string
	passivationTimeout string
}

// VerifyPegaStandardTierDeployment - Verifies Pega standard tier deployment for values as provided in default values.yaml.
// It ensures syntax of web deployment, batch deployment, stream statefulset, search service, hpa, rolling update, web services, ingresses and config maps
func VerifyPegaStandardTierDeployment(t *testing.T, helmChartPath string, options *helm.Options, initContainers []string) {

	// Verify Deployment objects
	SplitAndVerifyPegaDeployments(t, helmChartPath, options, initContainers)

	// Verify tier config
	VerifyTierConfg(t, helmChartPath, options)

	// Verify environment config
	VerifyEnvironmentConfig(t, helmChartPath, options)

	// Verify search service
	VerifySearchService(t, helmChartPath, options)

	if options.SetValues["constellation.enabled"] == "true" {
		// Verify constellation service
		VerifyConstellationService(t, helmChartPath, options)
	}

	// Verfiy Pega deployed services
	SplitAndVerifyPegaServices(t, helmChartPath, options)

	if options.SetValues["global.provider"] != "openshift" {
		// Verify pega deployed ingresses
		SplitAndVerifyPegaIngresses(t, helmChartPath, options)
	}
	// Verify Pega HPAObjects
	SplitAndVerifyPegaHPAs(t, helmChartPath, options)

	// Verify search transport service
	VerifySearchTransportService(t, helmChartPath, options)

}

// SplitAndVerifyPegaDeployments - Splits the deployments from the rendered template and asserts each deployment/statefulset objects
func SplitAndVerifyPegaDeployments(t *testing.T, helmChartPath string, options *helm.Options, initContainers []string) {
	deployment := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-deployment.yaml"})
	var deploymentObj appsv1.Deployment
	var statefulsetObj appsv1beta2.StatefulSet
	deploymentSlice := strings.Split(deployment, "---")
	for index, deploymentInfo := range deploymentSlice {
		if index >= 1 && index <= 3 {

			if index == 1 {
				helm.UnmarshalK8SYaml(t, deploymentInfo, &deploymentObj)
				VerifyPegaDeployment(t, &deploymentObj,
					pegaDeployment{"pega-web", initContainers, "WebUser", "900"},
					options)
			} else if index == 2 {
				helm.UnmarshalK8SYaml(t, deploymentInfo, &deploymentObj)
				VerifyPegaDeployment(t, &deploymentObj,
					pegaDeployment{"pega-batch", initContainers, "BackgroundProcessing,Search,Batch,RealTime,Custom1,Custom2,Custom3,Custom4,Custom5,BIX", ""},
					options)
			} else if index == 3 {
				helm.UnmarshalK8SYaml(t, deploymentInfo, &statefulsetObj)
				VerifyPegaStatefulSet(t, &statefulsetObj,
					pegaDeployment{"pega-stream", initContainers, "Stream", "900"},
					options)

			}
		}
	}
}

// VerifyDeployment - Performs common pega deployment/statefulset assertions with the values as provided in default values.yaml
func VerifyDeployment(t *testing.T, pod *k8score.PodSpec, expectedSpec pegaDeployment, options *helm.Options) {
	require.Equal(t, pod.Volumes[0].Name, "pega-volume-config")
	require.Equal(t, expectedSpec.name, pod.Volumes[0].VolumeSource.ConfigMap.LocalObjectReference.Name)
	require.Equal(t, pod.Volumes[0].VolumeSource.ConfigMap.DefaultMode, volumeDefaultModePtr)
	require.Equal(t, pod.Volumes[1].Name, "pega-volume-credentials")
	require.Equal(t, pod.Volumes[1].VolumeSource.Secret.SecretName, "pega-credentials-secret")
	require.Equal(t, pod.Volumes[1].VolumeSource.Secret.DefaultMode, volumeDefaultModePtr)

	actualInitContainers := pod.InitContainers
	count := len(actualInitContainers)
	actualInitContainerNames := make([]string, count)
	for i := 0; i < count; i++ {
		actualInitContainerNames[i] = actualInitContainers[i].Name
	}

	require.Equal(t, expectedSpec.initContainers, actualInitContainerNames)
	VerifyInitContinerData(t, actualInitContainers, options)
	require.Equal(t, pod.Containers[0].Name, "pega-web-tomcat")
	require.Equal(t, pod.Containers[0].Image, "pegasystems/pega")
	require.Equal(t, pod.Containers[0].Ports[0].Name, "pega-web-port")
	require.Equal(t, pod.Containers[0].Ports[0].ContainerPort, int32(8080))
	var envIndex int32 = 0
	require.Equal(t, pod.Containers[0].Env[envIndex].Name, "NODE_TYPE")
	require.Equal(t, expectedSpec.nodeType, pod.Containers[0].Env[envIndex].Value)
	envIndex++
	require.Equal(t, pod.Containers[0].Env[envIndex].Name, "PEGA_APP_CONTEXT_PATH")
	require.Equal(t, pod.Containers[0].Env[envIndex].Value, "prweb")
	if expectedSpec.name == "pega-web" || expectedSpec.name == "pega-stream" {
		envIndex++
		require.Equal(t, pod.Containers[0].Env[envIndex].Name, "REQUESTOR_PASSIVATION_TIMEOUT")
		require.Equal(t, expectedSpec.passivationTimeout, pod.Containers[0].Env[envIndex].Value)
	}
	if options.SetValues["constellation.enabled"] == "true" && expectedSpec.name == "pega-web" {
		envIndex++
		require.Equal(t, pod.Containers[0].Env[envIndex].Name, "COSMOS_SETTINGS")
		require.Equal(t, "Pega-UIEngine/cosmosservicesURI=/c11n", pod.Containers[0].Env[envIndex].Value)
	}

	if options.SetValues["constellation.enabled"] == "true" && expectedSpec.name != "pega-web" {
		envIndex++
		require.Equal(t, pod.Containers[0].Env[envIndex].Name, "COSMOS_SETTINGS")
		require.Equal(t, "Pega-UIEngine/cosmosservicesURI=/c11n", pod.Containers[0].Env[envIndex].Value)
	}
	envIndex++
	require.Equal(t, pod.Containers[0].Env[envIndex].Name, "JAVA_OPTS")
	require.Equal(t, pod.Containers[0].Env[envIndex].Value, "")
	envIndex++
	require.Equal(t, pod.Containers[0].Env[envIndex].Name, "INITIAL_HEAP")
	require.Equal(t, pod.Containers[0].Env[envIndex].Value, "4096m")
	envIndex++
	require.Equal(t, pod.Containers[0].Env[envIndex].Name, "MAX_HEAP")
	require.Equal(t, pod.Containers[0].Env[envIndex].Value, "7168m")
	require.Equal(t, pod.Containers[0].EnvFrom[0].ConfigMapRef.LocalObjectReference.Name, "pega-environment-config")
	require.Equal(t, "4", pod.Containers[0].Resources.Limits.Cpu().String())
	require.Equal(t, "8Gi", pod.Containers[0].Resources.Limits.Memory().String())
	require.Equal(t, "2", pod.Containers[0].Resources.Requests.Cpu().String())
	require.Equal(t, "6Gi", pod.Containers[0].Resources.Requests.Memory().String())

	require.Equal(t, pod.Containers[0].VolumeMounts[0].Name, "pega-volume-config")
	require.Equal(t, pod.Containers[0].VolumeMounts[0].MountPath, "/opt/pega/config")

	require.Equal(t, pod.Containers[0].LivenessProbe.InitialDelaySeconds, int32(300))
	require.Equal(t, pod.Containers[0].LivenessProbe.TimeoutSeconds, int32(20))
	require.Equal(t, pod.Containers[0].LivenessProbe.PeriodSeconds, int32(30))
	require.Equal(t, pod.Containers[0].LivenessProbe.SuccessThreshold, int32(1))
	require.Equal(t, pod.Containers[0].LivenessProbe.FailureThreshold, int32(3))
	require.Equal(t, pod.Containers[0].LivenessProbe.HTTPGet.Path, "/prweb/PRRestService/monitor/pingService/ping")
	require.Equal(t, pod.Containers[0].LivenessProbe.HTTPGet.Port, intstr.FromInt(8080))
	require.Equal(t, pod.Containers[0].LivenessProbe.HTTPGet.Scheme, k8score.URIScheme("HTTP"))

	require.Equal(t, pod.Containers[0].ReadinessProbe.InitialDelaySeconds, int32(300))
	require.Equal(t, pod.Containers[0].ReadinessProbe.TimeoutSeconds, int32(20))
	require.Equal(t, pod.Containers[0].ReadinessProbe.PeriodSeconds, int32(30))
	require.Equal(t, pod.Containers[0].ReadinessProbe.SuccessThreshold, int32(1))
	require.Equal(t, pod.Containers[0].ReadinessProbe.FailureThreshold, int32(3))
	require.Equal(t, pod.Containers[0].ReadinessProbe.HTTPGet.Path, "/prweb/PRRestService/monitor/pingService/ping")
	require.Equal(t, pod.Containers[0].ReadinessProbe.HTTPGet.Port, intstr.FromInt(8080))
	require.Equal(t, pod.Containers[0].ReadinessProbe.HTTPGet.Scheme, k8score.URIScheme("HTTP"))

	require.Equal(t, pod.ImagePullSecrets[0].Name, "pega-registry-secret")
	require.Equal(t, pod.RestartPolicy, k8score.RestartPolicy("Always"))
	require.Equal(t, pod.TerminationGracePeriodSeconds, terminationGracePeriodSecondsPtr)
	require.Equal(t, pod.Containers[0].VolumeMounts[0].Name, "pega-volume-config")
	require.Equal(t, pod.Containers[0].VolumeMounts[0].MountPath, "/opt/pega/config")
	require.Equal(t, pod.Volumes[0].Name, "pega-volume-config")
	require.Equal(t, pod.Volumes[1].Name, "pega-volume-credentials")
	require.Equal(t, pod.Volumes[1].Secret.SecretName, "pega-credentials-secret")

}

// VerifyPegaDeployment - Performs specific Pega deployment assertions with the values as provided in default values.yaml
func VerifyPegaDeployment(t *testing.T, deploymentObj *appsv1.Deployment, expectedDeployment pegaDeployment, options *helm.Options) {
	require.Equal(t, deploymentObj.Spec.Replicas, replicasPtr)
	require.Equal(t, deploymentObj.Spec.ProgressDeadlineSeconds, ProgressDeadlineSecondsPtr)
	require.Equal(t, expectedDeployment.name, deploymentObj.Spec.Selector.MatchLabels["app"])
	require.Equal(t, deploymentObj.Spec.Strategy.RollingUpdate.MaxSurge, rollingUpdatePtr)
	require.Equal(t, deploymentObj.Spec.Strategy.RollingUpdate.MaxUnavailable, rollingUpdatePtr)
	require.Equal(t, deploymentObj.Spec.Strategy.Type, appsv1.DeploymentStrategyType("RollingUpdate"))
	require.Equal(t, expectedDeployment.name, deploymentObj.Spec.Template.Labels["app"])
	require.NotEmpty(t, deploymentObj.Spec.Template.Annotations["config-check"])
	deploymentSpec := deploymentObj.Spec.Template.Spec
	VerifyDeployment(t, &deploymentSpec, expectedDeployment, options)
}

// VerifyPegaStatefulSet - Performs specific Pega statefulset assertions with the values as provided in default values.yaml
func VerifyPegaStatefulSet(t *testing.T, statefulsetObj *appsv1beta2.StatefulSet, expectedStatefulset pegaDeployment, options *helm.Options) {
	require.Equal(t, statefulsetObj.Spec.VolumeClaimTemplates[0].Name, "pega-stream")
	require.Equal(t, statefulsetObj.Spec.VolumeClaimTemplates[0].Spec.AccessModes[0], k8score.PersistentVolumeAccessMode("ReadWriteOnce"))
	require.Equal(t, statefulsetObj.Spec.ServiceName, "pega-stream")
	statefulsetSpec := statefulsetObj.Spec.Template.Spec
	require.Equal(t, statefulsetSpec.Containers[0].VolumeMounts[1].Name, "pega-stream")
	require.Equal(t, statefulsetSpec.Containers[0].VolumeMounts[1].MountPath, "/opt/pega/streamvol")
	require.Equal(t, statefulsetSpec.Containers[0].VolumeMounts[2].Name, "pega-volume-credentials")
	require.Equal(t, statefulsetSpec.Containers[0].VolumeMounts[2].MountPath, "/opt/pega/secrets")
	VerifyDeployment(t, &statefulsetSpec, expectedStatefulset, options)
}

type pegaServices struct {
	Name       string
	Port       int32
	TargetPort intstr.IntOrString
}

// SplitAndVerifyPegaServices - Splits the services from the rendered template and asserts each service objects
func SplitAndVerifyPegaServices(t *testing.T, helmChartPath string, options *helm.Options) {
	service := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-service.yaml"})
	var pegaServiceObj k8score.Service
	serviceSlice := strings.Split(service, "---")
	for index, serviceInfo := range serviceSlice {
		if index >= 1 && index <= 2 {
			helm.UnmarshalK8SYaml(t, serviceInfo, &pegaServiceObj)
			if index == 1 {
				VerifyPegaService(t, &pegaServiceObj, pegaServices{"pega-web", int32(80), intstr.IntOrString{IntVal: 8080}}, options)
			} else {
				VerifyPegaService(t, &pegaServiceObj, pegaServices{"pega-stream", int32(7003), intstr.IntOrString{IntVal: 7003}}, options)
			}
		}
	}
}

// VerifyPegaService - Performs Pega Service assertions with the values as provided in default values.yaml
func VerifyPegaService(t *testing.T, serviceObj *k8score.Service, expectedService pegaServices, options *helm.Options) {
	provider := options.SetValues["global.provider"]
	if provider == "k8s" {
		require.Equal(t, serviceObj.Annotations["traefik.ingress.kubernetes.io/affinity"], "true")
		require.Equal(t, serviceObj.Annotations["traefik.ingress.kubernetes.io/load-balancer-method"], "drr")
		require.Equal(t, serviceObj.Annotations["traefik.ingress.kubernetes.io/max-conn-amount"], "10")
		require.Equal(t, serviceObj.Annotations["traefik.ingress.kubernetes.io/session-cookie-name"], "UNIQUE-PEGA-COOKIE-NAME")
		require.Equal(t, serviceObj.Spec.Type, k8score.ServiceType("LoadBalancer"))
	} else if provider == "gke" {
		require.Equal(t, `{"ingress": true}`, serviceObj.Annotations["cloud.google.com/neg"])
		var expectedBackendConfig = fmt.Sprintf(`{"ports": {"%d": "pega-backend-config"}}`, expectedService.Port)
		require.Equal(t, expectedBackendConfig, serviceObj.Annotations["beta.cloud.google.com/backend-config"])
		require.Equal(t, serviceObj.Spec.Type, k8score.ServiceType("NodePort"))
	}
	require.Equal(t, serviceObj.Spec.Selector["app"], expectedService.Name)
	require.Equal(t, serviceObj.Spec.Ports[0].Port, expectedService.Port)
	require.Equal(t, serviceObj.Spec.Ports[0].TargetPort, expectedService.TargetPort)
}

type pegaIngress struct {
	Name          string
	Port          intstr.IntOrString
	AlbStickiness int32
}

// VerifyPegaIngresses - Splits the ingresses from the rendered template and asserts each ingress object
func SplitAndVerifyPegaIngresses(t *testing.T, helmChartPath string, options *helm.Options) {
	ingress := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-ingress.yaml"})
	var pegaIngressObj k8sv1beta1.Ingress
	ingressSlice := strings.Split(ingress, "---")
	for index, ingressInfo := range ingressSlice {
		if index >= 1 && index <= 2 {
			helm.UnmarshalK8SYaml(t, ingressInfo, &pegaIngressObj)
			if index == 1 {
				VerifyPegaIngress(t, &pegaIngressObj,
					pegaIngress{"pega-web", intstr.IntOrString{IntVal: 80}, 1020},
					options)
			} else {
				VerifyPegaIngress(t, &pegaIngressObj,
					pegaIngress{"pega-stream", intstr.IntOrString{IntVal: 7003}, 1020},
					options)
			}

		}
	}
}

func VerifyPegaIngress(t *testing.T, ingressObj *k8sv1beta1.Ingress, expectedIngress pegaIngress, options *helm.Options) {
	provider := options.SetValues["global.provider"]
	if provider == "eks" {
		VerifyEKSIngress(t, ingressObj, expectedIngress)
	} else if provider == "gke" {
		VerifyGKEIngress(t, ingressObj, expectedIngress)
	} else if provider == "aks" {
		VerifyAKSIngress(t, ingressObj, expectedIngress)
	} else if provider != "eks" && options.SetValues["constellation.enabled"] == "false" {
		VerifyK8SIngress(t, ingressObj, expectedIngress)
	} else if provider != "aks" && options.SetValues["constellation.enabled"] == "false" {
		VerifyAKSIngress(t, ingressObj, expectedIngress)
	} else if provider != "gke" && options.SetValues["constellation.enabled"] == "false" {
		VerifyGKEIngress(t, ingressObj, expectedIngress)
	} else if options.SetValues["constellation.enabled"] == "true" {
		VerifyK8SIngressWithConstellationEnabled(t, ingressObj, expectedIngress)
	} else {
		VerifyK8SIngress(t, ingressObj, expectedIngress)
	}
}

func VerifyEKSIngress(t *testing.T, ingressObj *k8sv1beta1.Ingress, expectedIngress pegaIngress) {
	require.Equal(t, "alb", ingressObj.Annotations["kubernetes.io/ingress.class"])
	require.Equal(t, "[{\"HTTP\": 80}, {\"HTTPS\": 443}]", ingressObj.Annotations["alb.ingress.kubernetes.io/listen-ports"])
	require.Equal(t, "{\"Type\": \"redirect\", \"RedirectConfig\": { \"Protocol\": \"HTTPS\", \"Port\": \"443\", \"StatusCode\": \"HTTP_301\"}}", ingressObj.Annotations["alb.ingress.kubernetes.io/actions.ssl-redirect"])
	require.Equal(t, "internet-facing", ingressObj.Annotations["alb.ingress.kubernetes.io/scheme"])
	expectedStickiness := fmt.Sprint("stickiness.enabled=true,stickiness.lb_cookie.duration_seconds=", expectedIngress.AlbStickiness)
	require.Equal(t, expectedStickiness,
		ingressObj.Annotations["alb.ingress.kubernetes.io/target-group-attributes"])
	require.Equal(t, "ip", ingressObj.Annotations["alb.ingress.kubernetes.io/target-type"])
	require.Equal(t, "ssl-redirect", ingressObj.Spec.Rules[0].HTTP.Paths[0].Backend.ServiceName)
	require.Equal(t, "use-annotation", ingressObj.Spec.Rules[0].HTTP.Paths[0].Backend.ServicePort.StrVal)
	require.Equal(t, expectedIngress.Name, ingressObj.Spec.Rules[1].HTTP.Paths[0].Backend.ServiceName)
	require.Equal(t, expectedIngress.Port, ingressObj.Spec.Rules[1].HTTP.Paths[0].Backend.ServicePort)
}

func VerifyGKEIngress(t *testing.T, ingressObj *k8sv1beta1.Ingress, expectedIngress pegaIngress) {
	require.Equal(t, "false", ingressObj.Annotations["kubernetes.io/ingress.allow-http"])
	require.Equal(t, expectedIngress.Name, ingressObj.Spec.Backend.ServiceName)
	require.Equal(t, expectedIngress.Port, ingressObj.Spec.Backend.ServicePort)
	require.Equal(t, 1, len(ingressObj.Spec.Rules))
	require.Equal(t, 1, len(ingressObj.Spec.Rules[0].HTTP.Paths))
	require.Equal(t, expectedIngress.Name, ingressObj.Spec.Rules[0].HTTP.Paths[0].Backend.ServiceName)
	require.Equal(t, expectedIngress.Port, ingressObj.Spec.Rules[0].HTTP.Paths[0].Backend.ServicePort)
}

func VerifyAKSIngress(t *testing.T, ingressObj *k8sv1beta1.Ingress, expectedIngress pegaIngress) {
	require.Equal(t, "azure/application-gateway", ingressObj.Annotations["kubernetes.io/ingress.class"])
	require.Equal(t, "true", ingressObj.Annotations["appgw.ingress.kubernetes.io/cookie-based-affinity"])
	require.Equal(t, expectedIngress.Name, ingressObj.Spec.Rules[0].HTTP.Paths[0].Backend.ServiceName)
	require.Equal(t, expectedIngress.Port, ingressObj.Spec.Rules[0].HTTP.Paths[0].Backend.ServicePort)
}

// VerifyPegaIngress - Performs Pega Ingress assertions with the values as provided in default values.yaml
func VerifyK8SIngress(t *testing.T, ingressObj *k8sv1beta1.Ingress, expectedIngress pegaIngress) {
	require.Equal(t, "traefik", ingressObj.Annotations["kubernetes.io/ingress.class"])
	require.Equal(t, expectedIngress.Name, ingressObj.Spec.Rules[0].HTTP.Paths[0].Backend.ServiceName)
	require.Equal(t, expectedIngress.Port, ingressObj.Spec.Rules[0].HTTP.Paths[0].Backend.ServicePort)
}

// VerifyPegaIngressWithconstellationEnabled - Performs Pega Ingress assertions with the values as provided in default values.yaml when the constellation service is enabled.
func VerifyK8SIngressWithConstellationEnabled(t *testing.T, ingressObj *k8sv1beta1.Ingress, expectedIngress pegaIngress) {
	if ingressObj.Spec.Rules[0].HTTP.Paths[0].Backend.ServiceName == "constellation" {
		require.Equal(t, "traefik", ingressObj.Annotations["kubernetes.io/ingress.class"])
		require.Equal(t, "/c11n", ingressObj.Spec.Rules[0].HTTP.Paths[0].Path)
		require.Equal(t, "constellation", ingressObj.Spec.Rules[0].HTTP.Paths[0].Backend.ServiceName)
		require.Equal(t, intstr.FromInt(3000), ingressObj.Spec.Rules[0].HTTP.Paths[0].Backend.ServicePort)
		require.Equal(t, expectedIngress.Name, ingressObj.Spec.Rules[0].HTTP.Paths[1].Backend.ServiceName)
		require.Equal(t, expectedIngress.Port, ingressObj.Spec.Rules[0].HTTP.Paths[1].Backend.ServicePort)
	} else {
		require.Equal(t, expectedIngress.Name, ingressObj.Spec.Rules[0].HTTP.Paths[0].Backend.ServiceName)
		require.Equal(t, expectedIngress.Port, ingressObj.Spec.Rules[0].HTTP.Paths[0].Backend.ServicePort)
	}
}

// VerifySearchService - Verifies search service deployment used by search pod with the values as provided in default values.yaml
func VerifySearchService(t *testing.T, helmChartPath string, options *helm.Options) {

	searchService := helm.RenderTemplate(t, options, helmChartPath, []string{"charts/pegasearch/templates/pega-search-service.yaml"})
	var searchServiceObj k8score.Service
	helm.UnmarshalK8SYaml(t, searchService, &searchServiceObj)
	require.Equal(t, searchServiceObj.Spec.Selector["component"], "Search")
	require.Equal(t, searchServiceObj.Spec.Selector["app"], "pega-search")
	require.Equal(t, searchServiceObj.Spec.Ports[0].Name, "http")
	require.Equal(t, searchServiceObj.Spec.Ports[0].Port, int32(80))
	require.Equal(t, searchServiceObj.Spec.Ports[0].TargetPort, intstr.FromInt(9200))
}

// VerifyconstellationService - Verifies constellation service deployment used by constellation pod with the values as provided in default values.yaml
func VerifyConstellationService(t *testing.T, helmChartPath string, options *helm.Options) {

	constellationService := helm.RenderTemplate(t, options, helmChartPath, []string{"charts/constellation/templates/clln-service.yaml"})
	var constellationServiceObj k8score.Service
	helm.UnmarshalK8SYaml(t, constellationService, &constellationServiceObj)
	require.Equal(t, constellationServiceObj.Spec.Selector["app"], "constellation")
	// require.Equal(t, constellationServiceObj.Spec.Ports[0].Protocol, "TCP")
	require.Equal(t, constellationServiceObj.Spec.Ports[0].Port, int32(3000))
	require.Equal(t, constellationServiceObj.Spec.Ports[0].TargetPort, intstr.FromInt(3000))
}

// VerifyEnvironmentConfig - Verifies the environment configuration used by the pods with the values as provided in default values.yaml
func VerifyEnvironmentConfig(t *testing.T, helmChartPath string, options *helm.Options) {

	envConfig := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/pega-environment-config.yaml"})
	var envConfigMap k8score.ConfigMap
	helm.UnmarshalK8SYaml(t, envConfig, &envConfigMap)
	envConfigData := envConfigMap.Data
	require.Equal(t, envConfigData["DB_TYPE"], "YOUR_DATABASE_TYPE")
	require.Equal(t, envConfigData["JDBC_URL"], "YOUR_JDBC_URL")
	require.Equal(t, envConfigData["JDBC_CLASS"], "YOUR_JDBC_DRIVER_CLASS")
	require.Equal(t, envConfigData["JDBC_DRIVER_URI"], "YOUR_JDBC_DRIVER_URI")
	if options.SetValues["global.actions.execute"] == "upgrade-deploy" {
		require.Equal(t, envConfigData["RULES_SCHEMA"], "")
	} else {
		require.Equal(t, envConfigData["RULES_SCHEMA"], "YOUR_RULES_SCHEMA")
	}
	require.Equal(t, envConfigData["DATA_SCHEMA"], "YOUR_DATA_SCHEMA")
	require.Equal(t, envConfigData["CUSTOMERDATA_SCHEMA"], "")
	require.Equal(t, envConfigData["JDBC_CONNECTION_PROPERTIES"], "")
	require.Equal(t, envConfigData["PEGA_SEARCH_URL"], "http://pega-search")
	require.Equal(t, envConfigData["CASSANDRA_CLUSTER"], "true")
	require.Equal(t, envConfigData["CASSANDRA_NODES"], "release-name-cassandra")
	require.Equal(t, envConfigData["CASSANDRA_PORT"], "9042")
}

// VerifyTierConfg - Performs the tier specific configuration assetions with the values as provided in default values.yaml
func VerifyTierConfg(t *testing.T, helmChartPath string, options *helm.Options) {
	config := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-config.yaml"})
	var pegaConfigMap k8score.ConfigMap
	configSlice := strings.Split(config, "---")
	for index, configData := range configSlice {
		if index >= 1 && index <= 3 {
			helm.UnmarshalK8SYaml(t, configData, &pegaConfigMap)
			pegaConfigMapData := pegaConfigMap.Data
			compareConfigMapData(t, pegaConfigMapData["prconfig.xml"], "data/expectedInstallDeployPrconfig.xml")
			compareConfigMapData(t, pegaConfigMapData["context.xml.tmpl"], "data/expectedInstallDeployContext.xml")
			compareConfigMapData(t, pegaConfigMapData["prlog4j2.xml"], "data/expectedInstallDeployPRlog4j2.xml")
		}
	}
}

type hpa struct {
	name          string
	targetRefName string
	kind          string
	apiversion    string
}

// VerifyPegaHPAs - Splits the HPA object from the rendered template and asserts each HPA object
func SplitAndVerifyPegaHPAs(t *testing.T, helmChartPath string, options *helm.Options) {
	pegaHpa := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-hpa.yaml"})
	var pegaHpaObj autoscaling.HorizontalPodAutoscaler
	hpaSlice := strings.SplitAfter(pegaHpa, "85")
	for index, hpaInfo := range hpaSlice {
		if index >= 0 && index <= 1 {
			helm.UnmarshalK8SYaml(t, hpaInfo, &pegaHpaObj)
			if index == 0 {
				VerifyPegaHpa(t, &pegaHpaObj, hpa{"pega-web-hpa", "pega-web", "Deployment", "apps/v1"})
			} else {
				VerifyPegaHpa(t, &pegaHpaObj, hpa{"pega-batch-hpa", "pega-batch", "Deployment", "apps/v1"})
			}
		}
	}
}

// VerifyPegaHpa - Performs Pega HPA assertions with the values as provided in default values.yaml
func VerifyPegaHpa(t *testing.T, hpaObj *autoscaling.HorizontalPodAutoscaler, expectedHpa hpa) {
	require.Equal(t, hpaObj.Spec.ScaleTargetRef.Name, expectedHpa.targetRefName)
	require.Equal(t, hpaObj.Spec.ScaleTargetRef.Kind, expectedHpa.kind)
	require.Equal(t, hpaObj.Spec.ScaleTargetRef.APIVersion, expectedHpa.apiversion)
	require.Equal(t, hpaObj.Spec.Metrics[0].Resource.Name, api.ResourceName("cpu"))
	require.Equal(t, hpaObj.Spec.Metrics[1].Resource.Name, api.ResourceName("memory"))
	require.Equal(t, hpaObj.Spec.MaxReplicas, int32(5))
}

// VerifySearchTransportService - Performs the search transport service assertions deployed with the values as provided in default values.yaml
func VerifySearchTransportService(t *testing.T, helmChartPath string, options *helm.Options) {
	transportSearchService := helm.RenderTemplate(t, options, helmChartPath, []string{"charts/pegasearch/templates/pega-search-transport-service.yaml"})
	var transportSearchServiceObj k8score.Service
	helm.UnmarshalK8SYaml(t, transportSearchService, &transportSearchServiceObj)

	require.Equal(t, transportSearchServiceObj.Spec.Selector["component"], "Search")
	require.Equal(t, transportSearchServiceObj.Spec.Selector["app"], "pega-search")
	require.Equal(t, transportSearchServiceObj.Spec.ClusterIP, "None")
	require.Equal(t, transportSearchServiceObj.Spec.Ports[0].Name, "transport")
	require.Equal(t, transportSearchServiceObj.Spec.Ports[0].Port, int32(80))
	require.Equal(t, transportSearchServiceObj.Spec.Ports[0].TargetPort, intstr.FromInt(9300))
}
