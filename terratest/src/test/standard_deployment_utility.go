package test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	k8score "k8s.io/api/core/v1"
	k8sv1beta1 "k8s.io/api/extensions/v1beta1"
	intstr "k8s.io/apimachinery/pkg/util/intstr"
	autoscaling "k8s.io/kubernetes/pkg/apis/autoscaling"
)

var replicas int32 = 1
var replicasPtr = &replicas
var ProgressDeadlineSeconds int32 = 2147483647
var ProgressDeadlineSecondsPtr = &ProgressDeadlineSeconds
var rollingUpdate intstr.IntOrString = intstr.FromString("25%")
var rollingUpdatePtr = &rollingUpdate
var volumeDefaultMode int32 = 420
var volumeDefaultModePtr = &volumeDefaultMode
var terminationGracePeriodSeconds int64 = 300
var terminationGracePeriodSecondsPtr = &terminationGracePeriodSeconds

type pegaDeployment struct {
	name           string
	initContainers []string
	nodeType       string
}

func VerifyInitContinerData(t *testing.T, containers []k8score.Container) {

	for i := 0; i < len(containers); i++ {
		container := containers[i]
		name := container.Name
		if name == "wait-for-pegainstall" {
			require.Equal(t, "dcasavant/k8s-wait-for", container.Image)
			require.Equal(t, []string{"job", "pega-db-install"}, container.Args)
		} else if name == "wait-for-pegasearch" {
			require.Equal(t, "busybox:1.31.0", container.Image)
			require.Equal(t, []string{"sh", "-c", "until $(wget -q -S --spider --timeout=2 -O /dev/null http://pega-search); do echo Waiting for search to become live...; sleep 10; done;"}, container.Command)
		} else if name == "wait-for-cassandra" {
			require.Equal(t, "cassandra:3.11.3", container.Image)
			require.Equal(t, []string{"sh", "-c", "until cqlsh -u \"dnode_ext\" -p \"dnode_ext\" -e \"describe cluster\" release-name-cassandra 9042 ; do echo Waiting for cassandra to become live...; sleep 10; done;"}, container.Command)
		} else {
			fmt.Println("in last else", name)
			t.Fail()
		}
	}
}

func VerifyPegaStandardTierDeployment(t *testing.T, helmChartPath string, options *helm.Options, initContainers []string) {

	// Deployment objects
	VerifyPegaDeployments(t, helmChartPath, options, initContainers)

	// Verify tier config
	VerifyTierConfg(t, helmChartPath, options)

	// Verify environment config
	VerifyEnvironmentConfig(t, helmChartPath, options)

	// Verify search service
	VerifySearchService(t, helmChartPath, options)

	// Verfiy pega services
	VerifyPegaServices(t, helmChartPath, options)

	// Verify pega ingresses
	VerifyPegaIngresses(t, helmChartPath, options)

	// Verify Pega HPAObjects
	VerifyPegaHPAs(t, helmChartPath, options)

	// Verify search transport service
	VerifySearchTransportService(t, helmChartPath, options)

}

func VerifyPegaDeployments(t *testing.T, helmChartPath string, options *helm.Options, initContainers []string) {
	deployment := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-deployment.yaml"})
	var deploymentObj appsv1.Deployment
	deploymentSlice := strings.Split(deployment, "---")
	for index, deploymentInfo := range deploymentSlice {
		if index >= 1 && index <= 3 {
			helm.UnmarshalK8SYaml(t, deploymentInfo, &deploymentObj)
			if index == 1 {
				VerifyPegaDeployment(t, &deploymentObj, pegaDeployment{"pega-web", initContainers, "WebUser"})
			} else if index == 2 {
				VerifyPegaDeployment(t, &deploymentObj, pegaDeployment{"pega-batch", initContainers, "BackgroundProcessing,Search,Batch,RealTime,Custom1,Custom2,Custom3,Custom4,Custom5,BIX"})
			}
		}
	}
}

func VerifyPegaDeployment(t *testing.T, deploymentObj *appsv1.Deployment, expectedDeployment pegaDeployment) {
	require.Equal(t, deploymentObj.Spec.Replicas, replicasPtr)
	require.Equal(t, deploymentObj.Spec.ProgressDeadlineSeconds, ProgressDeadlineSecondsPtr)
	require.Equal(t, expectedDeployment.name, deploymentObj.Spec.Selector.MatchLabels["app"])
	require.Equal(t, deploymentObj.Spec.Strategy.RollingUpdate.MaxSurge, rollingUpdatePtr)
	require.Equal(t, deploymentObj.Spec.Strategy.RollingUpdate.MaxUnavailable, rollingUpdatePtr)
	require.Equal(t, deploymentObj.Spec.Strategy.Type, appsv1.DeploymentStrategyType("RollingUpdate"))

	require.Equal(t, expectedDeployment.name, deploymentObj.Spec.Template.Labels["app"])
	require.NotEmpty(t, deploymentObj.Spec.Template.Annotations["config-check"])

	deploymentSpec := deploymentObj.Spec.Template.Spec

	require.Equal(t, deploymentSpec.Volumes[0].Name, "pega-volume-config")
	require.Equal(t, expectedDeployment.name, deploymentSpec.Volumes[0].VolumeSource.ConfigMap.LocalObjectReference.Name)
	require.Equal(t, deploymentSpec.Volumes[0].VolumeSource.ConfigMap.DefaultMode, volumeDefaultModePtr)
	require.Equal(t, deploymentSpec.Volumes[1].Name, "pega-volume-credentials")
	require.Equal(t, deploymentSpec.Volumes[1].VolumeSource.Secret.SecretName, "pega-credentials-secret")
	require.Equal(t, deploymentSpec.Volumes[1].VolumeSource.Secret.DefaultMode, volumeDefaultModePtr)

	actualInitContainers := deploymentSpec.InitContainers
	count := len(actualInitContainers)
	actualInitContainerNames := make([]string, count)
	for i := 0; i < count; i++ {
		actualInitContainerNames[i] = actualInitContainers[i].Name
	}

	require.Equal(t, expectedDeployment.initContainers, actualInitContainerNames)
	VerifyInitContinerData(t, actualInitContainers)
	require.Equal(t, deploymentSpec.Containers[0].Name, "pega-web-tomcat")
	require.Equal(t, deploymentSpec.Containers[0].Image, "YOUR_PEGA_DEPLOY_IMAGE:TAG")
	require.Equal(t, deploymentSpec.Containers[0].Ports[0].Name, "pega-web-port")
	require.Equal(t, deploymentSpec.Containers[0].Ports[0].ContainerPort, int32(8080))
	require.Equal(t, deploymentSpec.Containers[0].Env[0].Name, "NODE_TYPE")
	require.Equal(t, expectedDeployment.nodeType, deploymentSpec.Containers[0].Env[0].Value)
	require.Equal(t, deploymentSpec.Containers[0].Env[1].Name, "JAVA_OPTS")
	require.Equal(t, deploymentSpec.Containers[0].Env[1].Value, "")
	require.Equal(t, deploymentSpec.Containers[0].Env[2].Name, "INITIAL_HEAP")
	require.Equal(t, deploymentSpec.Containers[0].Env[2].Value, "4096m")
	require.Equal(t, deploymentSpec.Containers[0].Env[3].Name, "MAX_HEAP")
	require.Equal(t, deploymentSpec.Containers[0].Env[3].Value, "7168m")
	require.Equal(t, deploymentSpec.Containers[0].EnvFrom[0].ConfigMapRef.LocalObjectReference.Name, "pega-environment-config")

	require.Equal(t, "2", deploymentSpec.Containers[0].Resources.Limits.Cpu().String())
	require.Equal(t, "8Gi", deploymentSpec.Containers[0].Resources.Limits.Memory().String())
	require.Equal(t, "200m", deploymentSpec.Containers[0].Resources.Requests.Cpu().String())
	require.Equal(t, "6Gi", deploymentSpec.Containers[0].Resources.Requests.Memory().String())

	require.Equal(t, deploymentSpec.Containers[0].VolumeMounts[0].Name, "pega-volume-config")
	require.Equal(t, deploymentSpec.Containers[0].VolumeMounts[0].MountPath, "/opt/pega/config")
	require.Equal(t, deploymentSpec.Containers[0].VolumeMounts[1].Name, "pega-volume-credentials")
	require.Equal(t, deploymentSpec.Containers[0].VolumeMounts[1].MountPath, "/opt/pega/secrets")

	require.Equal(t, deploymentSpec.Containers[0].LivenessProbe.InitialDelaySeconds, int32(300))
	require.Equal(t, deploymentSpec.Containers[0].LivenessProbe.TimeoutSeconds, int32(20))
	require.Equal(t, deploymentSpec.Containers[0].LivenessProbe.PeriodSeconds, int32(10))
	require.Equal(t, deploymentSpec.Containers[0].LivenessProbe.SuccessThreshold, int32(1))
	require.Equal(t, deploymentSpec.Containers[0].LivenessProbe.FailureThreshold, int32(3))
	require.Equal(t, deploymentSpec.Containers[0].LivenessProbe.HTTPGet.Path, "/prweb/PRRestService/monitor/pingService/ping")
	require.Equal(t, deploymentSpec.Containers[0].LivenessProbe.HTTPGet.Port, intstr.FromInt(8080))
	require.Equal(t, deploymentSpec.Containers[0].LivenessProbe.HTTPGet.Scheme, k8score.URIScheme("HTTP"))

	require.Equal(t, deploymentSpec.Containers[0].ReadinessProbe.InitialDelaySeconds, int32(300))
	require.Equal(t, deploymentSpec.Containers[0].ReadinessProbe.TimeoutSeconds, int32(20))
	require.Equal(t, deploymentSpec.Containers[0].ReadinessProbe.PeriodSeconds, int32(10))
	require.Equal(t, deploymentSpec.Containers[0].ReadinessProbe.SuccessThreshold, int32(1))
	require.Equal(t, deploymentSpec.Containers[0].ReadinessProbe.FailureThreshold, int32(3))
	require.Equal(t, deploymentSpec.Containers[0].ReadinessProbe.HTTPGet.Path, "/prweb/PRRestService/monitor/pingService/ping")
	require.Equal(t, deploymentSpec.Containers[0].ReadinessProbe.HTTPGet.Port, intstr.FromInt(8080))
	require.Equal(t, deploymentSpec.Containers[0].ReadinessProbe.HTTPGet.Scheme, k8score.URIScheme("HTTP"))

	require.Equal(t, deploymentSpec.ImagePullSecrets[0].Name, "pega-registry-secret")
	require.Equal(t, deploymentSpec.RestartPolicy, k8score.RestartPolicy("Always"))
	require.Equal(t, deploymentSpec.TerminationGracePeriodSeconds, terminationGracePeriodSecondsPtr)

}

type pegaServices struct {
	Name       string
	Port       int32
	TargetPort intstr.IntOrString
}

func VerifyPegaServices(t *testing.T, helmChartPath string, options *helm.Options) {
	service := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-service.yaml"})
	var pegaServiceObj k8score.Service
	serviceSlice := strings.Split(service, "---")
	for index, serviceInfo := range serviceSlice {
		if index >= 1 && index <= 2 {
			helm.UnmarshalK8SYaml(t, serviceInfo, &pegaServiceObj)
			if index == 1 {
				VerifyPegaService(t, &pegaServiceObj, pegaServices{"pega-web", int32(80), intstr.IntOrString{IntVal: 8080}})
			} else {
				VerifyPegaService(t, &pegaServiceObj, pegaServices{"pega-stream", int32(7003), intstr.IntOrString{IntVal: 7003}})
			}
		}
	}
}

func VerifyPegaService(t *testing.T, serviceObj *k8score.Service, expectedService pegaServices) {
	require.Equal(t, serviceObj.Spec.Selector["app"], expectedService.Name)
	require.Equal(t, serviceObj.Spec.Ports[0].Port, expectedService.Port)
	require.Equal(t, serviceObj.Spec.Ports[0].TargetPort, expectedService.TargetPort)
	require.Equal(t, serviceObj.Annotations["traefik.ingress.kubernetes.io/affinity"], "true")
	require.Equal(t, serviceObj.Annotations["traefik.ingress.kubernetes.io/load-balancer-method"], "drr")
	require.Equal(t, serviceObj.Annotations["traefik.ingress.kubernetes.io/max-conn-amount"], "10")
	require.Equal(t, serviceObj.Annotations["traefik.ingress.kubernetes.io/session-cookie-name"], "UNIQUE-PEGA-COOKIE-NAME")
}

type pegaIngress struct {
	Name string
	Port intstr.IntOrString
}

func VerifyPegaIngress(t *testing.T, ingressObj *k8sv1beta1.Ingress, expectedIngress pegaIngress) {
	require.Equal(t, ingressObj.Annotations["kubernetes.io/ingress.class"], "traefik")
	require.Equal(t, ingressObj.Spec.Rules[0].HTTP.Paths[0].Backend.ServiceName, expectedIngress.Name)
	require.Equal(t, ingressObj.Spec.Rules[0].HTTP.Paths[0].Backend.ServicePort, expectedIngress.Port)
}

func VerifyPegaIngresses(t *testing.T, helmChartPath string, options *helm.Options) {
	ingress := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-ingress.yaml"})
	var pegaIngressObj k8sv1beta1.Ingress
	ingressSlice := strings.Split(ingress, "---")
	for index, ingressInfo := range ingressSlice {
		if index >= 1 && index <= 2 {
			helm.UnmarshalK8SYaml(t, ingressInfo, &pegaIngressObj)
			if index == 1 {
				VerifyPegaIngress(t, &pegaIngressObj, pegaIngress{"pega-web", intstr.IntOrString{IntVal: 80}})
			} else {
				VerifyPegaIngress(t, &pegaIngressObj, pegaIngress{"pega-stream", intstr.IntOrString{IntVal: 7003}})
			}

		}
	}
}

func VerifyPegaStatefulset() {

}

// Just verify what is exposed in the values yaml & k8s objects
func VerifyCassandraService() {

}

// Just verify what is exposed in the values yaml & k8s objects
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

// VerifyEnvironmentConfig
func VerifyEnvironmentConfig(t *testing.T, helmChartPath string, options *helm.Options) {

	envConfig := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/pega-environment-config.yaml"})
	var envConfigMap k8score.ConfigMap
	helm.UnmarshalK8SYaml(t, envConfig, &envConfigMap)
	envConfigData := envConfigMap.Data
	require.Equal(t, envConfigData["DB_TYPE"], "YOUR_DATABASE_TYPE")
	require.Equal(t, envConfigData["JDBC_URL"], "YOUR_JDBC_URL")
	require.Equal(t, envConfigData["JDBC_CLASS"], "YOUR_JDBC_DRIVER_CLASS")
	require.Equal(t, envConfigData["JDBC_DRIVER_URI"], "YOUR_JDBC_DRIVER_URI")
	require.Equal(t, envConfigData["RULES_SCHEMA"], "YOUR_RULES_SCHEMA")
	require.Equal(t, envConfigData["DATA_SCHEMA"], "YOUR_DATA_SCHEMA")
	require.Equal(t, envConfigData["CUSTOMERDATA_SCHEMA"], "")
	require.Equal(t, envConfigData["JDBC_CONNECTION_PROPERTIES"], "")
	require.Equal(t, envConfigData["PEGA_SEARCH_URL"], "http://pega-search")
	require.Equal(t, envConfigData["CASSANDRA_CLUSTER"], "true")
	require.Equal(t, envConfigData["CASSANDRA_NODES"], "release-name-cassandra")
	require.Equal(t, envConfigData["CASSANDRA_PORT"], "9042")
}

func VerifyTierConfg(t *testing.T, helmChartPath string, options *helm.Options) {
	config := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-config.yaml"})
	var pegaConfigMap k8score.ConfigMap
	configSlice := strings.Split(config, "---")
	for index, configData := range configSlice {
		if index >= 1 && index <= 3 {
			helm.UnmarshalK8SYaml(t, configData, &pegaConfigMap)
			//pegaConfigMapData := pegaConfigMap.Data
			//compareConfigMapData(t, []byte(pegaConfigMapData["data/prconfig.xml"]), "expectedInstallDeployPrconfig.xml")
			//compareConfigMapData(t, []byte(pegaConfigMapData["data/context.xml.tmpl"]), "expectedInstallDeployContext.xml")
			//compareConfigMapData(t, []byte(pegaConfigMapData["data/prlog4j2.xml"]), "expectedInstallDeployPRlog4j2.xml")
		}
	}
}

type hpa struct {
	name          string
	targetRefName string
	kind          string
	apiversion    string
}

func VerifyPegaHpa(t *testing.T, hpaObj *autoscaling.HorizontalPodAutoscaler, expectedHpa hpa) {
	require.Equal(t, hpaObj.Name, expectedHpa.name)
	require.Equal(t, hpaObj.Spec.ScaleTargetRef.Name, expectedHpa.targetRefName)
	require.Equal(t, hpaObj.Spec.ScaleTargetRef.Kind, expectedHpa.kind)
	require.Equal(t, hpaObj.Spec.ScaleTargetRef.APIVersion, expectedHpa.apiversion)

	require.Equal(t, hpaObj.Spec.MinReplicas, replicasPtr)
	require.Equal(t, hpaObj.Spec.MaxReplicas, int32(5))
}

func VerifyPegaHPAs(t *testing.T, helmChartPath string, options *helm.Options) {
	pegaHpa := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-hpa.yaml"})
	var pegaHpaObj autoscaling.HorizontalPodAutoscaler
	hpaSlice := strings.Split(pegaHpa, "#")

	for index, hpaInfo := range hpaSlice {
		if index >= 2 && index <= 3 {
			helm.UnmarshalK8SYaml(t, hpaInfo, &pegaHpaObj)
			if index == 2 {
				VerifyPegaHpa(t, &pegaHpaObj, hpa{"pega-web-hpa", "pega-web", "Deployment", "extensions/v1beta1"})
			} else {
				VerifyPegaHpa(t, &pegaHpaObj, hpa{"pega-batch-hpa", "pega-batch", "Deployment", "extensions/v1beta1"})
			}
		}
	}
}

func VerifySearchTransportService(t *testing.T, helmChartPath string, options *helm.Options) {
	// pega-search-transport-service.yaml
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
