package test

import (
	"path/filepath"
	"testing"

	//"fmt"

	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	k8score "k8s.io/api/core/v1"
	k8srbac "k8s.io/api/rbac/v1"

	//k8sresource "k8s.io/apimachinery/pkg/api/resource"

	"github.com/gruntwork-io/terratest/modules/helm"
)

const pegaHelmChartPath = "../../../charts/pega"

var options = &helm.Options{
	SetValues: map[string]string{
		"global.actions.execute": "install-deploy",
	},
}

func TestInstallDeployActionSkippedTemplates(t *testing.T) {
	t.Parallel()

	// Path to the helm chart we will test
	helmChartPath, err := filepath.Abs(pegaHelmChartPath)
	require.NoError(t, err)

	// with action as 'install-deploy' below templates should not be rendered
	output := helm.RenderTemplate(t, options, helmChartPath, []string{
		"templates/pega-action-validate.yaml",
		"charts/installer/templates/pega-upgrade-environment-config.yaml",
	})

	var deployment appsv1.Deployment
	helm.UnmarshalK8SYaml(t, output, &deployment)
	// assert that above templates are not rendered
	require.Empty(t, deployment)
}

func TestInstallDeployActionInstallerRole(t *testing.T) {
	t.Parallel()
	// Path to the helm chart we will test
	helmChartPath, err := filepath.Abs(pegaHelmChartPath)
	require.NoError(t, err)

	deployRole := helm.RenderTemplate(t, options, helmChartPath, []string{"charts/installer/templates/pega-installer-role.yaml"})
	var deployRoleObj k8srbac.Role
	helm.UnmarshalK8SYaml(t, deployRole, &deployRoleObj)
	require.Equal(t, deployRoleObj.Rules[0].APIGroups, []string{"", "batch", "extensions", "apps"})
	require.Equal(t, deployRoleObj.Rules[0].Resources, []string{"jobs", "deployments", "statefulsets"})
	require.Equal(t, deployRoleObj.Rules[0].Verbs, []string{"get", "watch", "list"})

}

func TestInstallDeployActionInstallerRoleBinding(t *testing.T) {
	t.Parallel()

	// Path to the helm chart we will test
	helmChartPath, err := filepath.Abs(pegaHelmChartPath)
	require.NoError(t, err)

	installerRoleBinding := helm.RenderTemplate(t, options, helmChartPath, []string{"charts/installer/templates/pega-installer-status-rolebinding.yaml"})
	var installerRoleBindingObj k8srbac.RoleBinding
	helm.UnmarshalK8SYaml(t, installerRoleBinding, &installerRoleBindingObj)
	require.Equal(t, installerRoleBindingObj.RoleRef.APIGroup, "rbac.authorization.k8s.io")
	require.Equal(t, installerRoleBindingObj.RoleRef.Kind, "Role")
	require.Equal(t, installerRoleBindingObj.RoleRef.Name, "jobs-reader")

	require.Equal(t, installerRoleBindingObj.Subjects[0].Kind, "ServiceAccount")
	require.Equal(t, installerRoleBindingObj.Subjects[0].Name, "default")
	require.Equal(t, installerRoleBindingObj.Subjects[0].Namespace, "default")
}

func TestInstallDeployActionInstallerJob(t *testing.T) {
	t.Parallel()
	// Path to the helm chart we will test
	helmChartPath, err := filepath.Abs(pegaHelmChartPath)
	require.NoError(t, err)

	VerifyPegaJob(t, helmChartPath, options, pegaJob{"pega-db-install", []string{}, "pega-install-environment-config"})

}

func TestInstallDeployActionInstallerConfig(t *testing.T) {
	t.Parallel()
	t.Skip()
	// Path to the helm chart we will test
	helmChartPath, err := filepath.Abs(pegaHelmChartPath)
	require.NoError(t, err)

	installerConfig := helm.RenderTemplate(t, options, helmChartPath, []string{"charts/installer/templates/pega-installer-config.yaml"})
	var installConfigMap k8score.ConfigMap
	helm.UnmarshalK8SYaml(t, installerConfig, &installConfigMap)

	installConfigData := installConfigMap.Data
	compareConfigMapData(t, []byte(installConfigData["prconfig.xml.tmpl"]), "expectedPrconfig.xml")
	compareConfigMapData(t, []byte(installConfigData["setupDatabase.properties.tmpl"]), "expectedsetupDatabase.properties")
	compareConfigMapData(t, []byte(installConfigData["prbootstrap.properties.tmpl"]), "expectedPRbootstrap.properties")
	compareConfigMapData(t, []byte(installConfigData["prlog4j2.xml"]), "expectedPRlog4j2.xml")

}

func TestInstallDeployActionInstallerEnvConfig(t *testing.T) {
	t.Parallel()
	//.Skip()
	// Path to the helm chart we will test
	helmChartPath, err := filepath.Abs(pegaHelmChartPath)
	require.NoError(t, err)
	installEnvConfig := helm.RenderTemplate(t, options, helmChartPath, []string{"charts/installer/templates/pega-install-environment-config.yaml"})
	var installEnvConfigMap k8score.ConfigMap
	helm.UnmarshalK8SYaml(t, installEnvConfig, &installEnvConfigMap)

	installEnvConfigData := installEnvConfigMap.Data
	require.Equal(t, installEnvConfigData["DB_TYPE"], "YOUR_DATABASE_TYPE")
	require.Equal(t, installEnvConfigData["JDBC_URL"], "YOUR_JDBC_URL")
	require.Equal(t, installEnvConfigData["JDBC_CLASS"], "YOUR_JDBC_DRIVER_CLASS")
	require.Equal(t, installEnvConfigData["JDBC_DRIVER_URI"], "YOUR_JDBC_DRIVER_URI")
	require.Equal(t, installEnvConfigData["RULES_SCHEMA"], "YOUR_RULES_SCHEMA")
	require.Equal(t, installEnvConfigData["DATA_SCHEMA"], "YOUR_DATA_SCHEMA")
	require.Equal(t, installEnvConfigData["CUSTOMERDATA_SCHEMA"], "")
	require.Equal(t, installEnvConfigData["SYSTEM_NAME"], "pega")
	require.Equal(t, installEnvConfigData["PRODUCTION_LEVEL"], "2")
	require.Equal(t, installEnvConfigData["MULTITENANT_SYSTEM"], "false")
	require.Equal(t, "ADMIN_PASSWORD", installEnvConfigData["ADMIN_PASSWORD"])
	require.Equal(t, "", installEnvConfigData["STATIC_ASSEMBLER"])
	require.Equal(t, installEnvConfigData["BYPASS_UDF_GENERATION"], "false")
	require.Equal(t, installEnvConfigData["BYPASS_TRUNCATE_UPDATESCACHE"], "false")
	require.Equal(t, installEnvConfigData["JDBC_CUSTOM_CONNECTION"], "")
	require.Equal(t, installEnvConfigData["MAX_IDLE"], "5")
	require.Equal(t, installEnvConfigData["MAX_WAIT"], "-1")
	require.Equal(t, installEnvConfigData["MAX_ACTIVE"], "10")
	require.Equal(t, installEnvConfigData["ZOS_PROPERTIES"], "/opt/pega/config/DB2SiteDependent.properties")
	require.Equal(t, installEnvConfigData["DB2ZOS_UDF_WLM"], "")
	require.Equal(t, installEnvConfigData["ACTION"], "install-deploy")

}

func TestInstallDeployActionStandardDeployment(t *testing.T) {
	t.Parallel()
	//t.Skip()
	// Path to the helm chart we will test
	helmChartPath, err := filepath.Abs(pegaHelmChartPath)
	require.NoError(t, err)

	VerifyPegaStandardTierDeployment(t, helmChartPath, options, []string{"wait-for-pegainstall", "wait-for-pegasearch", "wait-for-cassandra"})

}

/*func TestInstallDeployActionShouldNotRenderDeployments(t *testing.T) {
	t.Parallel()
	t.Skip()

	// Path to the helm chart we will test
	helmChartPath, err := filepath.Abs(pegaHelmChartPath)
	require.NoError(t, err)

	// with action as 'install-deploy' below templates should not be rendered
	output := helm.RenderTemplate(t, options, helmChartPath, []string{
		"templates/pega-action-validate.yaml",
		"templates/pega-upgrade-environment-config.yaml",
	})

	var deployment appsv1.Deployment
	helm.UnmarshalK8SYaml(t, output, &deployment)
	// assert that above templates are not rendered
	require.Empty(t, deployment)

	// pega-batch-config.yaml
	batchConfig := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/pega-batch-config.yaml"})
	var batchConfigMap k8score.ConfigMap
	helm.UnmarshalK8SYaml(t, batchConfig, &batchConfigMap)

	batchConfigData := batchConfigMap.Data
	compareConfigMapData(t, []byte(batchConfigData["prconfig.xml"]), "expectedInstallDeployPrconfig.xml")
	compareConfigMapData(t, []byte(batchConfigData["context.xml.tmpl"]), "expectedInstallDeployContext.xml")
	compareConfigMapData(t, []byte(batchConfigData["prlog4j2.xml"]), "expectedInstallDeployPRlog4j2.xml")

	// pega-deploy-role.yaml

	// pega-environment-config.yaml
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
	require.Equal(t, envConfigData["JDBC_CONNECTION_PROPERTIES"], "socketTimeout=90")
	require.Equal(t, envConfigData["PEGA_SEARCH_URL"], "http://pega-search")
	require.Equal(t, envConfigData["CASSANDRA_CLUSTER"], "true")
	require.Equal(t, envConfigData["CASSANDRA_NODES"], "release-name-cassandra")
	require.Equal(t, envConfigData["CASSANDRA_PORT"], "9042")

	//pega-installer-status-rolebinding.yaml

	// pega-search-service.yaml
	searchService := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/pega-search-service.yaml"})
	var searchServiceObj k8score.Service
	helm.UnmarshalK8SYaml(t, searchService, &searchServiceObj)
	var servicePort int32 = 80

	require.Equal(t, searchServiceObj.Spec.Selector["component"], "Search")
	require.Equal(t, searchServiceObj.Spec.Selector["app"], "pega-search")
	require.Equal(t, searchServiceObj.Spec.Ports[0].Name, "http")
	require.Equal(t, searchServiceObj.Spec.Ports[0].Port, servicePort)
	require.Equal(t, searchServiceObj.Spec.Ports[0].TargetPort, intstr.FromInt(9200))

	// pega-search-transport-service.yaml
	transportSearchService := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/pega-search-transport-service.yaml"})
	var transportSearchServiceObj k8score.Service
	helm.UnmarshalK8SYaml(t, transportSearchService, &transportSearchServiceObj)

	require.Equal(t, transportSearchServiceObj.Spec.Selector["component"], "Search")
	require.Equal(t, transportSearchServiceObj.Spec.Selector["app"], "pega-search")
	require.Equal(t, transportSearchServiceObj.Spec.ClusterIP, "None")
	require.Equal(t, transportSearchServiceObj.Spec.Ports[0].Name, "transport")
	require.Equal(t, transportSearchServiceObj.Spec.Ports[0].Port, servicePort)
	require.Equal(t, transportSearchServiceObj.Spec.Ports[0].TargetPort, intstr.FromInt(9300))

	// pega-stream-service.yaml
	streamService := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/pega-stream-service.yaml"})
	var streamServiceObj k8score.Service
	helm.UnmarshalK8SYaml(t, streamService, &streamServiceObj)

	require.Equal(t, streamServiceObj.Spec.Selector["app"], "pega-stream")
	require.Equal(t, streamServiceObj.Spec.Ports[0].Name, "http")
	require.Equal(t, streamServiceObj.Spec.Ports[0].Port, int32(7003))
	require.Equal(t, streamServiceObj.Spec.Ports[0].TargetPort, intstr.FromInt(7003))
	require.Equal(t, streamServiceObj.Annotations["traefik.ingress.kubernetes.io/affinity"], "true")
	require.Equal(t, streamServiceObj.Annotations["traefik.ingress.kubernetes.io/load-balancer-method"], "drr")
	require.Equal(t, streamServiceObj.Annotations["traefik.ingress.kubernetes.io/max-conn-amount"], "10")
	require.Equal(t, streamServiceObj.Annotations["traefik.ingress.kubernetes.io/session-cookie-name"], "UNIQUE-PEGA-COOKIE-NAME")

	// pega-web-service.yaml
	webService := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/pega-web-service.yaml"})
	var webServiceObj k8score.Service
	helm.UnmarshalK8SYaml(t, webService, &webServiceObj)

	require.Equal(t, webServiceObj.Spec.Selector["app"], "pega-web")
	require.Equal(t, webServiceObj.Spec.Ports[0].Name, "http")
	require.Equal(t, webServiceObj.Spec.Ports[0].Port, int32(80))
	require.Equal(t, webServiceObj.Spec.Ports[0].TargetPort, intstr.FromInt(8080))
	require.Equal(t, webServiceObj.Annotations["traefik.ingress.kubernetes.io/affinity"], "true")
	require.Equal(t, webServiceObj.Annotations["traefik.ingress.kubernetes.io/load-balancer-method"], "drr")
	require.Equal(t, webServiceObj.Annotations["traefik.ingress.kubernetes.io/max-conn-amount"], "10")
	require.Equal(t, webServiceObj.Annotations["traefik.ingress.kubernetes.io/session-cookie-name"], "UNIQUE-PEGA-COOKIE-NAME")

	// pega-stream-config.yaml
	streamConfig := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/pega-stream-config.yaml"})
	var streamConfigMap k8score.ConfigMap
	helm.UnmarshalK8SYaml(t, streamConfig, &streamConfigMap)

	streamConfigData := streamConfigMap.Data
	compareConfigMapData(t, []byte(streamConfigData["prconfig.xml"]), "expectedInstallDeployPrconfig.xml")
	compareConfigMapData(t, []byte(streamConfigData["context.xml.tmpl"]), "expectedInstallDeployContext.xml")
	compareConfigMapData(t, []byte(streamConfigData["prlog4j2.xml"]), "expectedInstallDeployPRlog4j2.xml")

	// pega-web-config.yaml
	webConfig := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/pega-web-config.yaml"})
	var webConfigMap k8score.ConfigMap
	helm.UnmarshalK8SYaml(t, webConfig, &webConfigMap)

	webConfigData := webConfigMap.Data
	compareConfigMapData(t, []byte(webConfigData["prconfig.xml"]), "expectedInstallDeployPrconfig.xml")
	compareConfigMapData(t, []byte(webConfigData["context.xml.tmpl"]), "expectedInstallDeployContext.xml")
	compareConfigMapData(t, []byte(webConfigData["prlog4j2.xml"]), "expectedInstallDeployPRlog4j2.xml")

	// pega-stream-ingress.yaml
	streamIngress := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/pega-stream-ingress.yaml"})
	var streamIngressObj k8sv1beta1.Ingress
	helm.UnmarshalK8SYaml(t, streamIngress, &streamIngressObj)
	require.Equal(t, streamIngressObj.Annotations["kubernetes.io/ingress.class"], "traefik")
	require.Equal(t, streamIngressObj.Spec.Rules[0].Host, "YOUR_STREAM_NODE_DOMAIN")
	require.Equal(t, streamIngressObj.Spec.Rules[0].HTTP.Paths[0].Backend.ServiceName, "pega-stream")
	require.Equal(t, streamIngressObj.Spec.Rules[0].HTTP.Paths[0].Backend.ServicePort, intstr.FromInt(7003))

	// pega-web-ingress.yaml
	webIngress := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/pega-web-ingress.yaml"})
	var webIngressObj k8sv1beta1.Ingress
	helm.UnmarshalK8SYaml(t, webIngress, &webIngressObj)
	require.Equal(t, webIngressObj.Annotations["kubernetes.io/ingress.class"], "traefik")
	require.Equal(t, webIngressObj.Spec.Rules[0].Host, "YOUR_WEB_NODE_DOMAIN")
	require.Equal(t, webIngressObj.Spec.Rules[0].HTTP.Paths[0].Backend.ServiceName, "pega-web")
	require.Equal(t, webIngressObj.Spec.Rules[0].HTTP.Paths[0].Backend.ServicePort, intstr.FromInt(80))

	// pega-batch-deployment.yaml
	batchDeploymemt := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/pega-batch-deployment.yaml"})
	var batchDeploymemtObj appsv1.Deployment
	helm.UnmarshalK8SYaml(t, batchDeploymemt, &batchDeploymemtObj)
	var replicas int32 = 1
	var replicasPtr = &replicas
	var ProgressDeadlineSeconds int32 = 2147483647
	var ProgressDeadlineSecondsPtr = &ProgressDeadlineSeconds
	require.Equal(t, batchDeploymemtObj.Spec.Replicas, replicasPtr)
	require.Equal(t, batchDeploymemtObj.Spec.ProgressDeadlineSeconds, ProgressDeadlineSecondsPtr)
	require.Equal(t, batchDeploymemtObj.Spec.Selector.MatchLabels["app"], "pega-batch")

	var rollingUpdate intstr.IntOrString = intstr.FromString("25%")
	var rollingUpdatePtr = &rollingUpdate

	require.Equal(t, batchDeploymemtObj.Spec.Strategy.RollingUpdate.MaxSurge, rollingUpdatePtr)
	require.Equal(t, batchDeploymemtObj.Spec.Strategy.RollingUpdate.MaxUnavailable, rollingUpdatePtr)
	require.Equal(t, batchDeploymemtObj.Spec.Strategy.Type, appsv1.DeploymentStrategyType("RollingUpdate"))

	batchDeploymemtSpec := batchDeploymemtObj.Spec.Template.Spec
	var volumeDefaultMode int32 = 420
	var volumeDefaultModePtr = &volumeDefaultMode

	require.Equal(t, batchDeploymemtSpec.Volumes[0].Name, "pega-volume-config")
	require.Equal(t, batchDeploymemtSpec.Volumes[0].VolumeSource.ConfigMap.LocalObjectReference.Name, "pega-batch")
	require.Equal(t, batchDeploymemtSpec.Volumes[0].VolumeSource.ConfigMap.DefaultMode, volumeDefaultModePtr)
	require.Equal(t, batchDeploymemtSpec.Volumes[1].Name, "pega-volume-credentials")
	require.Equal(t, batchDeploymemtSpec.Volumes[1].VolumeSource.Secret.SecretName, "pega-credentials-secret")
	require.Equal(t, batchDeploymemtSpec.Volumes[1].VolumeSource.Secret.DefaultMode, volumeDefaultModePtr)

	require.Equal(t, batchDeploymemtSpec.InitContainers[0].Name, "wait-for-pegainstall")
	require.Equal(t, batchDeploymemtSpec.InitContainers[0].Image, "dcasavant/k8s-wait-for")
	require.Equal(t, batchDeploymemtSpec.InitContainers[0].Args, []string{"job", "pega-db-install"})
	require.Equal(t, batchDeploymemtSpec.InitContainers[1].Name, "wait-for-pegasearch")
	require.Equal(t, batchDeploymemtSpec.InitContainers[1].Image, "busybox:1.27.2")
	require.Equal(t, batchDeploymemtSpec.InitContainers[1].Command, []string{"sh", "-c", "until $(wget -q -S --spider --timeout=2 -O /dev/null http://pega-search); do echo Waiting for search to become live...; sleep 10; done;"})
	require.Equal(t, batchDeploymemtSpec.InitContainers[2].Name, "wait-for-cassandra")
	require.Equal(t, batchDeploymemtSpec.InitContainers[2].Image, "cassandra:3.11.3")
	require.Equal(t, batchDeploymemtSpec.InitContainers[2].Command, []string{"sh", "-c", "until cqlsh -u \"dnode_ext\" -p \"dnode_ext\" -e \"describe cluster\" release-name-cassandra 9042 ; do echo Waiting for cassandra to become live...; sleep 10; done;"})

	require.Equal(t, batchDeploymemtSpec.Containers[0].Name, "pega-web-tomcat")
	require.Equal(t, batchDeploymemtSpec.Containers[0].Image, "YOUR_PEGA_DEPLOY_IMAGE:TAG")
	require.Equal(t, batchDeploymemtSpec.Containers[0].Ports[0].Name, "pega-web-port")
	require.Equal(t, batchDeploymemtSpec.Containers[0].Ports[0].ContainerPort, int32(8080))
	require.Equal(t, batchDeploymemtSpec.Containers[0].Env[0].Name, "NODE_TYPE")
	require.Equal(t, batchDeploymemtSpec.Containers[0].Env[0].Value, "BackgroundProcessing,Search,Batch,RealTime,Custom1,Custom2,Custom3,Custom4,Custom5,BIX,ADM,RTDG")
	require.Equal(t, batchDeploymemtSpec.Containers[0].Env[1].Name, "JAVA_OPTS")
	require.Equal(t, batchDeploymemtSpec.Containers[0].Env[1].Value, "")
	require.Equal(t, batchDeploymemtSpec.Containers[0].Env[2].Name, "INITIAL_HEAP")
	require.Equal(t, batchDeploymemtSpec.Containers[0].Env[2].Value, "4096m")
	require.Equal(t, batchDeploymemtSpec.Containers[0].Env[3].Name, "MAX_HEAP")
	require.Equal(t, batchDeploymemtSpec.Containers[0].Env[3].Value, "7168m")
	require.Equal(t, batchDeploymemtSpec.Containers[0].EnvFrom[0].ConfigMapRef.LocalObjectReference.Name, "pega-environment-config")

	//degah
	//require.Equal(t, batchDeploymemtSpec.Containers[0].Resources.Limits[k8score.ResourceName("cpu")].Fromat, k8sresource.Quantity)
	//require.Equal(t, batchDeploymemtSpec.Containers[0].Resources.Limits[k8score.ResourceName("memory")], "8Gi")
	//require.Equal(t, batchDeploymemtSpec.Containers[0].Resources.Requests[k8score.ResourceName("cpu")].Fromat, k8sresource.Quantity)
	//require.Equal(t, batchDeploymemtSpec.Containers[0].Resources.Requests[k8score.ResourceName("memory")], "8Gi")

	require.Equal(t, batchDeploymemtSpec.Containers[0].VolumeMounts[0].Name, "pega-volume-config")
	require.Equal(t, batchDeploymemtSpec.Containers[0].VolumeMounts[0].MountPath, "/opt/pega/config")
	require.Equal(t, batchDeploymemtSpec.Containers[0].VolumeMounts[1].Name, "pega-volume-credentials")
	require.Equal(t, batchDeploymemtSpec.Containers[0].VolumeMounts[1].MountPath, "/opt/pega/secrets")

	require.Equal(t, batchDeploymemtSpec.Containers[0].LivenessProbe.InitialDelaySeconds, int32(300))
	require.Equal(t, batchDeploymemtSpec.Containers[0].LivenessProbe.TimeoutSeconds, int32(20))
	require.Equal(t, batchDeploymemtSpec.Containers[0].LivenessProbe.PeriodSeconds, int32(10))
	require.Equal(t, batchDeploymemtSpec.Containers[0].LivenessProbe.SuccessThreshold, int32(1))
	require.Equal(t, batchDeploymemtSpec.Containers[0].LivenessProbe.FailureThreshold, int32(3))
	require.Equal(t, batchDeploymemtSpec.Containers[0].LivenessProbe.HTTPGet.Path, "/prweb/PRRestService/monitor/pingService/ping")
	require.Equal(t, batchDeploymemtSpec.Containers[0].LivenessProbe.HTTPGet.Port, intstr.FromInt(8080))
	require.Equal(t, batchDeploymemtSpec.Containers[0].LivenessProbe.HTTPGet.Scheme, k8score.URIScheme("HTTP"))

	require.Equal(t, batchDeploymemtSpec.Containers[0].ReadinessProbe.InitialDelaySeconds, int32(300))
	require.Equal(t, batchDeploymemtSpec.Containers[0].ReadinessProbe.TimeoutSeconds, int32(20))
	require.Equal(t, batchDeploymemtSpec.Containers[0].ReadinessProbe.PeriodSeconds, int32(10))
	require.Equal(t, batchDeploymemtSpec.Containers[0].ReadinessProbe.SuccessThreshold, int32(1))
	require.Equal(t, batchDeploymemtSpec.Containers[0].ReadinessProbe.FailureThreshold, int32(3))
	require.Equal(t, batchDeploymemtSpec.Containers[0].ReadinessProbe.HTTPGet.Path, "/prweb/PRRestService/monitor/pingService/ping")
	require.Equal(t, batchDeploymemtSpec.Containers[0].ReadinessProbe.HTTPGet.Port, intstr.FromInt(8080))
	require.Equal(t, batchDeploymemtSpec.Containers[0].ReadinessProbe.HTTPGet.Scheme, k8score.URIScheme("HTTP"))

	require.Equal(t, batchDeploymemtSpec.ImagePullSecrets[0].Name, "pega-registry-secret")
	require.Equal(t, batchDeploymemtSpec.RestartPolicy, k8score.RestartPolicy("Always"))
	var terminationGracePeriodSeconds int64 = 300
	var terminationGracePeriodSecondsPtr = &terminationGracePeriodSeconds
	require.Equal(t, batchDeploymemtSpec.TerminationGracePeriodSeconds, terminationGracePeriodSecondsPtr)

	// pega-search-deployment.yaml
	searchDeployment := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/pega-search-deployment.yaml"})
	var searchDeploymentObj appsv1.StatefulSet
	helm.UnmarshalK8SYaml(t, searchDeployment, &searchDeploymentObj)

	require.Equal(t, searchDeploymentObj.Spec.Replicas, replicasPtr)
	require.Equal(t, searchDeploymentObj.Spec.Selector.MatchLabels["app"], "pega-search")
	searchDeploymentSpec := searchDeploymentObj.Spec.Template.Spec

	var secContext int64 = 1000
	var secContextPtr = &secContext
	var boolTrue bool = true
	var boolTruePtr = &boolTrue
	require.Equal(t, searchDeploymentSpec.SecurityContext.FSGroup, secContextPtr)
	require.Equal(t, searchDeploymentSpec.InitContainers[0].Name, "set-max-map-count")
	require.Equal(t, searchDeploymentSpec.InitContainers[0].Image, "busybox:1.27.2")
	require.Equal(t, searchDeploymentSpec.InitContainers[0].Command, []string{"sysctl", "-w", "vm.max_map_count=262144"})
	require.Equal(t, searchDeploymentSpec.InitContainers[0].SecurityContext.Privileged, boolTruePtr)

	require.Equal(t, searchDeploymentSpec.Containers[0].Name, "search")
	require.Equal(t, searchDeploymentSpec.Containers[0].Image, "YOUR_ELASTICSEARCH_IMAGE:TAG")
	require.Equal(t, searchDeploymentSpec.Containers[0].SecurityContext.RunAsUser, secContextPtr)
	require.Equal(t, searchDeploymentSpec.Containers[0].Env[0].Name, "HOST_LIST")
	require.Equal(t, searchDeploymentSpec.Containers[0].Env[0].Value, "pega-search-transport")
	require.Equal(t, searchDeploymentSpec.Containers[0].Env[1].Name, "ES_JAVA_OPTS")
	require.Equal(t, searchDeploymentSpec.Containers[0].Env[1].Value, "-Xmx2g -Xms2g")
	require.Equal(t, searchDeploymentSpec.Containers[0].Env[2].Name, "UNICAST_HOSTS")
	require.Equal(t, searchDeploymentSpec.Containers[0].Env[2].Value, "pega-search-transport")
	require.Equal(t, searchDeploymentSpec.Containers[0].Env[3].Name, "NUMBER_OF_MASTERS")
	require.Equal(t, searchDeploymentSpec.Containers[0].Env[3].Value, "1")

	//degah
	//require.Equal(t, searchDeploymentSpec.Containers[0].Resources.Limits[k8score.ResourceName("cpu")].Fromat, k8sresource.Quantity)
	//require.Equal(t, searchDeploymentSpec.Containers[0].Resources.Requests[k8score.ResourceName("cpu")].Fromat, k8sresource.Quantity)

	require.Equal(t, searchDeploymentSpec.Containers[0].Ports[0].Name, "http")
	require.Equal(t, searchDeploymentSpec.Containers[0].Ports[0].ContainerPort, int32(9200))
	require.Equal(t, searchDeploymentSpec.Containers[0].Ports[1].Name, "transport")
	require.Equal(t, searchDeploymentSpec.Containers[0].Ports[1].ContainerPort, int32(9300))

	require.Equal(t, searchDeploymentSpec.Containers[0].VolumeMounts[0].Name, "esstorage")
	require.Equal(t, searchDeploymentSpec.Containers[0].VolumeMounts[0].MountPath, "/usr/share/elasticsearch/data")

	require.Equal(t, searchDeploymentSpec.Containers[0].LivenessProbe.InitialDelaySeconds, int32(5))
	require.Equal(t, searchDeploymentSpec.Containers[0].LivenessProbe.PeriodSeconds, int32(10))
	require.Equal(t, searchDeploymentSpec.Containers[0].LivenessProbe.TCPSocket.Port, intstr.FromString("transport"))

	require.Equal(t, searchDeploymentSpec.Containers[0].ReadinessProbe.InitialDelaySeconds, int32(20))
	require.Equal(t, searchDeploymentSpec.Containers[0].ReadinessProbe.TimeoutSeconds, int32(5))
	require.Equal(t, searchDeploymentSpec.Containers[0].ReadinessProbe.HTTPGet.Path, "/_cat")
	require.Equal(t, searchDeploymentSpec.Containers[0].ReadinessProbe.HTTPGet.Port, intstr.FromString("http"))

	require.Equal(t, searchDeploymentSpec.ImagePullSecrets[0].Name, "pega-registry-secret")

	require.Equal(t, searchDeploymentObj.Spec.VolumeClaimTemplates[0].Name, "esstorage")
	require.Equal(t, searchDeploymentObj.Spec.VolumeClaimTemplates[0].Spec.AccessModes[0], k8score.PersistentVolumeAccessMode("ReadWriteOnce"))
	//require.Equal(t, searchDeploymentObj.Spec.VolumeClaimTemplates[0].Spec.Resources.Requests[k8score.ResourceName("storage")], "5Gi")

	// pega-stream-deployment.yaml
	streamDeployment := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/pega-stream-deployment.yaml"})
	var streamDeploymentObj appsv1.StatefulSet
	helm.UnmarshalK8SYaml(t, streamDeployment, &streamDeploymentObj)
	var streamReplicas int32 = 2
	var streamReplicasPtr = &streamReplicas
	require.Equal(t, streamDeploymentObj.Spec.Replicas, streamReplicasPtr)
	require.Equal(t, streamDeploymentObj.Spec.Selector.MatchLabels["app"], "pega-stream")
	streamDeploymentSpec := streamDeploymentObj.Spec.Template.Spec

	require.Equal(t, streamDeploymentSpec.Volumes[0].Name, "pega-volume-config")
	require.Equal(t, streamDeploymentSpec.Volumes[0].VolumeSource.ConfigMap.LocalObjectReference.Name, "pega-stream")
	require.Equal(t, streamDeploymentSpec.Volumes[0].VolumeSource.ConfigMap.DefaultMode, volumeDefaultModePtr)
	require.Equal(t, streamDeploymentSpec.Volumes[1].Name, "pega-volume-credentials")
	require.Equal(t, streamDeploymentSpec.Volumes[1].VolumeSource.Secret.SecretName, "pega-credentials-secret")
	require.Equal(t, streamDeploymentSpec.Volumes[1].VolumeSource.Secret.DefaultMode, volumeDefaultModePtr)

	require.Equal(t, streamDeploymentSpec.InitContainers[0].Name, "wait-for-pegainstall")
	require.Equal(t, streamDeploymentSpec.InitContainers[0].Image, "dcasavant/k8s-wait-for")
	require.Equal(t, streamDeploymentSpec.InitContainers[0].Args, []string{"job", "pega-db-install"})
	require.Equal(t, streamDeploymentSpec.InitContainers[1].Name, "wait-for-pegasearch")
	require.Equal(t, streamDeploymentSpec.InitContainers[1].Image, "busybox:1.27.2")
	require.Equal(t, streamDeploymentSpec.InitContainers[1].Command, []string{"sh", "-c", "until $(wget -q -S --spider --timeout=2 -O /dev/null http://pega-search); do echo Waiting for search to become live...; sleep 10; done;"})

	require.Equal(t, streamDeploymentSpec.Containers[0].Name, "pega-web-tomcat")
	require.Equal(t, streamDeploymentSpec.Containers[0].Image, "YOUR_PEGA_DEPLOY_IMAGE:TAG")
	require.Equal(t, streamDeploymentSpec.Containers[0].Ports[0].Name, "pega-web-port")
	require.Equal(t, streamDeploymentSpec.Containers[0].Ports[0].ContainerPort, int32(8080))
	require.Equal(t, streamDeploymentSpec.Containers[0].Env[0].Name, "NODE_TYPE")
	require.Equal(t, streamDeploymentSpec.Containers[0].Env[0].Value, "Stream")
	require.Equal(t, streamDeploymentSpec.Containers[0].Env[1].Name, "JAVA_OPTS")
	require.Equal(t, streamDeploymentSpec.Containers[0].Env[1].Value, "")
	require.Equal(t, streamDeploymentSpec.Containers[0].Env[2].Name, "INITIAL_HEAP")
	require.Equal(t, streamDeploymentSpec.Containers[0].Env[2].Value, "4096m")
	require.Equal(t, streamDeploymentSpec.Containers[0].Env[3].Name, "MAX_HEAP")
	require.Equal(t, streamDeploymentSpec.Containers[0].Env[3].Value, "7168m")
	require.Equal(t, streamDeploymentSpec.Containers[0].EnvFrom[0].ConfigMapRef.LocalObjectReference.Name, "pega-environment-config")

	//degah
	//require.Equal(t, streamDeploymentSpec.Containers[0].Resources.Limits[k8score.ResourceName("cpu")].Fromat, k8sresource.Quantity)
	//require.Equal(t, streamDeploymentSpec.Containers[0].Resources.Limits[k8score.ResourceName("memory")], "8Gi")
	//require.Equal(t, streamDeploymentSpec.Containers[0].Resources.Requests[k8score.ResourceName("cpu")].Fromat, k8sresource.Quantity)
	//require.Equal(t, streamDeploymentSpec.Containers[0].Resources.Requests[k8score.ResourceName("memory")], "8Gi")

	require.Equal(t, streamDeploymentSpec.Containers[0].VolumeMounts[0].Name, "pega-volume-config")
	require.Equal(t, streamDeploymentSpec.Containers[0].VolumeMounts[0].MountPath, "/opt/pega/config")
	require.Equal(t, streamDeploymentSpec.Containers[0].VolumeMounts[1].Name, "pega-stream")
	require.Equal(t, streamDeploymentSpec.Containers[0].VolumeMounts[1].MountPath, "/opt/pega/streamvol")
	require.Equal(t, streamDeploymentSpec.Containers[0].VolumeMounts[2].Name, "pega-volume-credentials")
	require.Equal(t, streamDeploymentSpec.Containers[0].VolumeMounts[2].MountPath, "/opt/pega/secrets")

	require.Equal(t, streamDeploymentSpec.Containers[0].LivenessProbe.InitialDelaySeconds, int32(300))
	require.Equal(t, streamDeploymentSpec.Containers[0].LivenessProbe.TimeoutSeconds, int32(20))
	require.Equal(t, streamDeploymentSpec.Containers[0].LivenessProbe.PeriodSeconds, int32(10))
	require.Equal(t, streamDeploymentSpec.Containers[0].LivenessProbe.SuccessThreshold, int32(1))
	require.Equal(t, streamDeploymentSpec.Containers[0].LivenessProbe.FailureThreshold, int32(3))
	require.Equal(t, streamDeploymentSpec.Containers[0].LivenessProbe.HTTPGet.Path, "/prweb/PRRestService/monitor/pingService/ping")
	require.Equal(t, streamDeploymentSpec.Containers[0].LivenessProbe.HTTPGet.Port, intstr.FromInt(8080))
	require.Equal(t, streamDeploymentSpec.Containers[0].LivenessProbe.HTTPGet.Scheme, k8score.URIScheme("HTTP"))

	require.Equal(t, streamDeploymentSpec.Containers[0].ReadinessProbe.InitialDelaySeconds, int32(300))
	require.Equal(t, streamDeploymentSpec.Containers[0].ReadinessProbe.TimeoutSeconds, int32(20))
	require.Equal(t, streamDeploymentSpec.Containers[0].ReadinessProbe.PeriodSeconds, int32(10))
	require.Equal(t, streamDeploymentSpec.Containers[0].ReadinessProbe.SuccessThreshold, int32(1))
	require.Equal(t, streamDeploymentSpec.Containers[0].ReadinessProbe.FailureThreshold, int32(3))
	require.Equal(t, streamDeploymentSpec.Containers[0].ReadinessProbe.HTTPGet.Path, "/prweb/PRRestService/monitor/pingService/ping")
	require.Equal(t, streamDeploymentSpec.Containers[0].ReadinessProbe.HTTPGet.Port, intstr.FromInt(8080))
	require.Equal(t, streamDeploymentSpec.Containers[0].ReadinessProbe.HTTPGet.Scheme, k8score.URIScheme("HTTP"))

	require.Equal(t, streamDeploymentSpec.ImagePullSecrets[0].Name, "pega-registry-secret")
	require.Equal(t, streamDeploymentSpec.RestartPolicy, k8score.RestartPolicy("Always"))
	require.Equal(t, streamDeploymentSpec.TerminationGracePeriodSeconds, terminationGracePeriodSecondsPtr)

	require.Equal(t, streamDeploymentObj.Spec.VolumeClaimTemplates[0].Name, "pega-stream")
	require.Equal(t, streamDeploymentObj.Spec.VolumeClaimTemplates[0].Spec.AccessModes[0], k8score.PersistentVolumeAccessMode("ReadWriteOnce"))
	//require.Equal(t, streamDeploymentObj.Spec.VolumeClaimTemplates[0].Spec.Resources.Requests[k8score.ResourceName("storage")], "5Gi")

	require.Equal(t, streamDeploymentObj.Spec.ServiceName, "pega-stream")

	// pega-web-deployment.yaml
	webDeployment := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/pega-web-deployment.yaml"})
	var webDeploymentObj appsv1.Deployment
	helm.UnmarshalK8SYaml(t, webDeployment, &webDeploymentObj)
	require.Equal(t, webDeploymentObj.Spec.Replicas, replicasPtr)
	require.Equal(t, webDeploymentObj.Spec.ProgressDeadlineSeconds, ProgressDeadlineSecondsPtr)
	require.Equal(t, webDeploymentObj.Spec.Selector.MatchLabels["app"], "pega-web")
	require.Equal(t, webDeploymentObj.Spec.Strategy.RollingUpdate.MaxSurge, rollingUpdatePtr)
	require.Equal(t, webDeploymentObj.Spec.Strategy.RollingUpdate.MaxUnavailable, rollingUpdatePtr)
	require.Equal(t, webDeploymentObj.Spec.Strategy.Type, appsv1.DeploymentStrategyType("RollingUpdate"))

	require.Equal(t, webDeploymentObj.Spec.Template.Labels["app"], "pega-web")
	require.NotEmpty(t, webDeploymentObj.Spec.Template.Annotations["config-check"])

	webDeploymentSpec := webDeploymentObj.Spec.Template.Spec

	require.Equal(t, webDeploymentSpec.Volumes[0].Name, "pega-volume-config")
	require.Equal(t, webDeploymentSpec.Volumes[0].VolumeSource.ConfigMap.LocalObjectReference.Name, "pega-web")
	require.Equal(t, webDeploymentSpec.Volumes[0].VolumeSource.ConfigMap.DefaultMode, volumeDefaultModePtr)
	require.Equal(t, webDeploymentSpec.Volumes[1].Name, "pega-volume-credentials")
	require.Equal(t, webDeploymentSpec.Volumes[1].VolumeSource.Secret.SecretName, "pega-credentials-secret")
	require.Equal(t, webDeploymentSpec.Volumes[1].VolumeSource.Secret.DefaultMode, volumeDefaultModePtr)

	require.Equal(t, webDeploymentSpec.InitContainers[0].Name, "wait-for-pegainstall")
	require.Equal(t, webDeploymentSpec.InitContainers[0].Image, "dcasavant/k8s-wait-for")
	require.Equal(t, webDeploymentSpec.InitContainers[0].Args, []string{"job", "pega-db-install"})
	require.Equal(t, webDeploymentSpec.InitContainers[1].Name, "wait-for-pegasearch")
	require.Equal(t, webDeploymentSpec.InitContainers[1].Image, "busybox:1.27.2")
	require.Equal(t, webDeploymentSpec.InitContainers[1].Command, []string{"sh", "-c", "until $(wget -q -S --spider --timeout=2 -O /dev/null http://pega-search); do echo Waiting for search to become live...; sleep 10; done;"})

	require.Equal(t, webDeploymentSpec.Containers[0].Name, "pega-web-tomcat")
	require.Equal(t, webDeploymentSpec.Containers[0].Image, "YOUR_PEGA_DEPLOY_IMAGE:TAG")
	require.Equal(t, webDeploymentSpec.Containers[0].Ports[0].Name, "pega-web-port")
	require.Equal(t, webDeploymentSpec.Containers[0].Ports[0].ContainerPort, int32(8080))
	require.Equal(t, webDeploymentSpec.Containers[0].Env[0].Name, "NODE_TYPE")
	require.Equal(t, webDeploymentSpec.Containers[0].Env[0].Value, "Foreground")
	require.Equal(t, webDeploymentSpec.Containers[0].Env[1].Name, "JAVA_OPTS")
	require.Equal(t, webDeploymentSpec.Containers[0].Env[1].Value, "")
	require.Equal(t, webDeploymentSpec.Containers[0].Env[2].Name, "INITIAL_HEAP")
	require.Equal(t, webDeploymentSpec.Containers[0].Env[2].Value, "4096m")
	require.Equal(t, webDeploymentSpec.Containers[0].Env[3].Name, "MAX_HEAP")
	require.Equal(t, webDeploymentSpec.Containers[0].Env[3].Value, "7168m")
	require.Equal(t, webDeploymentSpec.Containers[0].EnvFrom[0].ConfigMapRef.LocalObjectReference.Name, "pega-environment-config")

	//degah
	//require.Equal(t, webDeploymentSpec.Containers[0].Resources.Limits[k8score.ResourceName("cpu")].Fromat, k8sresource.Quantity)
	//require.Equal(t, webDeploymentSpec.Containers[0].Resources.Limits[k8score.ResourceName("memory")], "8Gi")
	//require.Equal(t, webDeploymentSpec.Containers[0].Resources.Requests[k8score.ResourceName("cpu")].Fromat, k8sresource.Quantity)
	//require.Equal(t, webDeploymentSpec.Containers[0].Resources.Requests[k8score.ResourceName("memory")], "6Gi")

	require.Equal(t, webDeploymentSpec.Containers[0].VolumeMounts[0].Name, "pega-volume-config")
	require.Equal(t, webDeploymentSpec.Containers[0].VolumeMounts[0].MountPath, "/opt/pega/config")
	require.Equal(t, webDeploymentSpec.Containers[0].VolumeMounts[1].Name, "pega-volume-credentials")
	require.Equal(t, webDeploymentSpec.Containers[0].VolumeMounts[1].MountPath, "/opt/pega/secrets")

	require.Equal(t, webDeploymentSpec.Containers[0].LivenessProbe.InitialDelaySeconds, int32(300))
	require.Equal(t, webDeploymentSpec.Containers[0].LivenessProbe.TimeoutSeconds, int32(20))
	require.Equal(t, webDeploymentSpec.Containers[0].LivenessProbe.PeriodSeconds, int32(10))
	require.Equal(t, webDeploymentSpec.Containers[0].LivenessProbe.SuccessThreshold, int32(1))
	require.Equal(t, webDeploymentSpec.Containers[0].LivenessProbe.FailureThreshold, int32(3))
	require.Equal(t, webDeploymentSpec.Containers[0].LivenessProbe.HTTPGet.Path, "/prweb/PRRestService/monitor/pingService/ping")
	require.Equal(t, webDeploymentSpec.Containers[0].LivenessProbe.HTTPGet.Port, intstr.FromInt(8080))
	require.Equal(t, webDeploymentSpec.Containers[0].LivenessProbe.HTTPGet.Scheme, k8score.URIScheme("HTTP"))

	require.Equal(t, webDeploymentSpec.Containers[0].ReadinessProbe.InitialDelaySeconds, int32(300))
	require.Equal(t, webDeploymentSpec.Containers[0].ReadinessProbe.TimeoutSeconds, int32(20))
	require.Equal(t, webDeploymentSpec.Containers[0].ReadinessProbe.PeriodSeconds, int32(10))
	require.Equal(t, webDeploymentSpec.Containers[0].ReadinessProbe.SuccessThreshold, int32(1))
	require.Equal(t, webDeploymentSpec.Containers[0].ReadinessProbe.FailureThreshold, int32(3))
	require.Equal(t, webDeploymentSpec.Containers[0].ReadinessProbe.HTTPGet.Path, "/prweb/PRRestService/monitor/pingService/ping")
	require.Equal(t, webDeploymentSpec.Containers[0].ReadinessProbe.HTTPGet.Port, intstr.FromInt(8080))
	require.Equal(t, webDeploymentSpec.Containers[0].ReadinessProbe.HTTPGet.Scheme, k8score.URIScheme("HTTP"))

	require.Equal(t, webDeploymentSpec.ImagePullSecrets[0].Name, "pega-registry-secret")
	require.Equal(t, webDeploymentSpec.RestartPolicy, k8score.RestartPolicy("Always"))
	require.Equal(t, webDeploymentSpec.TerminationGracePeriodSeconds, terminationGracePeriodSecondsPtr)

	// pega-batch-hpa.yaml
	batchHPA := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/pega-batch-hpa.yaml"})
	var batchHPAobj k8sAutoscale.HorizontalPodAutoscaler
	helm.UnmarshalK8SYaml(t, batchHPA, &batchHPAobj)

	require.Equal(t, batchHPAobj.Name, "pega-batch-hpa")
	require.Equal(t, batchHPAobj.Spec.ScaleTargetRef.Name, "pega-batch")
	require.Equal(t, batchHPAobj.Spec.ScaleTargetRef.Kind, "Deployment")
	require.Equal(t, batchHPAobj.Spec.ScaleTargetRef.APIVersion, "extensions/v1beta1")

	require.Equal(t, batchHPAobj.Spec.MinReplicas, replicasPtr)
	require.Equal(t, batchHPAobj.Spec.MaxReplicas, int32(3))

	// pega-batch-hpa.yaml
	webHPA := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/pega-web-hpa.yaml"})
	var webHPAobj k8sAutoscale.HorizontalPodAutoscaler
	helm.UnmarshalK8SYaml(t, webHPA, &webHPAobj)

	require.Equal(t, webHPAobj.Name, "pega-web-hpa")
	require.Equal(t, webHPAobj.Spec.ScaleTargetRef.Name, "pega-web")
	require.Equal(t, webHPAobj.Spec.ScaleTargetRef.Kind, "Deployment")
	require.Equal(t, webHPAobj.Spec.ScaleTargetRef.APIVersion, "extensions/v1beta1")

	require.Equal(t, webHPAobj.Spec.MinReplicas, replicasPtr)
	require.Equal(t, webHPAobj.Spec.MaxReplicas, int32(5))

	// pega-credentials-secret.yaml
	secretOutput := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/pega-credentials-secret.yaml"})

	var secretobj k8score.Secret
	helm.UnmarshalK8SYaml(t, secretOutput, &secretobj)
	secretData := secretobj.Data
	require.Equal(t, string(secretData["DB_USERNAME"]), "YOUR_JDBC_USERNAME")
	require.Equal(t, string(secretData["DB_PASSWORD"]), "YOUR_JDBC_PASSWORD")
	require.Equal(t, string(secretData["CASSANDRA_USERNAME"]), "dnode_ext")
	require.Equal(t, string(secretData["CASSANDRA_PASSWORD"]), "dnode_ext")
	require.Equal(t, string(secretData["PEGA_DIAGNOSTIC_USER"]), "")
	require.Equal(t, string(secretData["PEGA_DIAGNOSTIC_PASSWORD"]), "")

	// pega-install-environment-config.yaml

	// pega-registry-secret.yaml
	registrySecret := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/pega-registry-secret.yaml"})

	var registrySecretObj k8score.Secret
	helm.UnmarshalK8SYaml(t, registrySecret, &registrySecretObj)
	reqgistrySecretData := registrySecretObj.Data

	require.Contains(t, string(reqgistrySecretData[".dockerconfigjson"]), "YOUR_DOCKER_REGISTRY")

	// pega-installer-config.yaml

	// pega-installer-job.yaml

}*/
