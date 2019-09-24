package test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	k8score "k8s.io/api/core/v1"
	k8sv1beta1 "k8s.io/api/extensions/v1beta1"
	intstr "k8s.io/apimachinery/pkg/util/intstr"
	autoscaling "k8s.io/kubernetes/pkg/apis/autoscaling"

	//k8sresource "k8s.io/apimachinery/pkg/api/resource"

	"github.com/gruntwork-io/terratest/modules/helm"
)

func VerifyPegaStandradTierDeployment(t *testing.T) {

}

func TestPegaDeployments(t *testing.T) {
	t.Parallel()

	// Path to the helm chart we will test
	helmChartPath, err := filepath.Abs("../../../charts/pega")
	require.NoError(t, err)

	// set action execute to deploy
	options := &helm.Options{
		SetValues: map[string]string{
			"global.actions.execute": "deploy",
		},
	}
	deployment := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-deployment.yaml"})
	var deploymentObj appsv1.Deployment
	deploymentSlice := strings.Split(deployment, "---")
	for index, deploymentInfo := range deploymentSlice {
		if index >= 1 && index <= 3 {
			helm.UnmarshalK8SYaml(t, deploymentInfo, &deploymentObj)
			if index == 1 {
				VerifyPegaDeployment(t, &deploymentObj, pegaDeployment{"pega-web", []string{"wait-for-pegasearch", "wait-for-cassandra"}, "WebUser"})
			} else if index == 2 {
				VerifyPegaDeployment(t, &deploymentObj, pegaDeployment{"pega-batch", []string{"wait-for-pegasearch", "wait-for-cassandra"}, "BackgroundProcessing,Search,Batch,RealTime,Custom1,Custom2,Custom3,Custom4,Custom5,BIX"})
			}
		}
	}
}

func TestTierConfigs(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath, err := filepath.Abs("../../../charts/pega")
	require.NoError(t, err)

	// set action execute to deploy
	options := &helm.Options{
		SetValues: map[string]string{
			"global.actions.execute": "deploy",
		},
	}
	config := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-config.yaml"})
	var pegaConfigMap k8score.ConfigMap
	configSlice := strings.Split(config, "---")
	for index, configData := range configSlice {
		if index >= 1 && index <= 3 {
			helm.UnmarshalK8SYaml(t, configData, &pegaConfigMap)
			VerifyTierConfig(t, pegaConfigMap)
		}
	}
}

func TestEnvironmentConfig(t *testing.T) {
	helmChartPath, err := filepath.Abs("../../../charts/pega")
	require.NoError(t, err)

	// set action execute to deploy
	options := &helm.Options{
		SetValues: map[string]string{
			"global.actions.execute": "deploy",
		},
	}
	// pega-environment-config.yaml
	envConfig := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/pega-environment-config.yaml"})
	var envConfigMap k8score.ConfigMap
	helm.UnmarshalK8SYaml(t, envConfig, &envConfigMap)
	VerifyEnvironmentConfig(t, envConfigMap)
}

func TestSearchService(t *testing.T) {

	helmChartPath, err := filepath.Abs("../../../charts/pega")
	require.NoError(t, err)

	// set action execute to deploy
	options := &helm.Options{
		SetValues: map[string]string{
			"global.actions.execute": "deploy",
		},
	}
	// pega-search-service.yaml
	searchService := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/pega-search-service.yaml"})
	var searchServiceObj k8score.Service
	helm.UnmarshalK8SYaml(t, searchService, &searchServiceObj)
	VerifySearchService(t, &searchServiceObj)

}

func TestPegaServices(t *testing.T) {
	helmChartPath, err := filepath.Abs("../../../charts/pega")
	require.NoError(t, err)
	// set action execute to deploy
	options := &helm.Options{
		SetValues: map[string]string{
			"global.actions.execute": "deploy",
		},
	}
	service := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-service.yaml"})
	var pegaServiceObj k8score.Service
	serviceSlice := strings.Split(service, "---")
	for index, serviceInfo := range serviceSlice {
		if index >= 1 && index <= 2 {
			helm.UnmarshalK8SYaml(t, serviceInfo, &pegaServiceObj)
			if index == 1 {
				VerifyPegaServices(t, &pegaServiceObj, pegaServices{"pega-web", int32(80), int32(8080)})
			} else {
				VerifyPegaServices(t, &pegaServiceObj, pegaServices{"pega-stream", int32(7003), int32(7003)})
			}
		}
	}
}

func TestPegaIngress(t *testing.T) {
	helmChartPath, err := filepath.Abs("../../../charts/pega")
	require.NoError(t, err)

	// set action execute to deploy
	options := &helm.Options{
		SetValues: map[string]string{
			"global.actions.execute": "deploy",
		},
	}
	ingress := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-ingress.yaml"})
	var pegaIngressObj k8sv1beta1.Ingress
	ingressSlice := strings.Split(ingress, "---")
	for index, ingressInfo := range ingressSlice {
		if index >= 1 && index <= 2 {
			helm.UnmarshalK8SYaml(t, ingressInfo, &pegaIngressObj)
			if index == 1 {
				VerifyPegaIngress(t, &pegaIngressObj, pegaIngress{"pega-web", int32(80)})
			} else {
				VerifyPegaIngress(t, &pegaIngressObj, pegaIngress{"pega-stream", int32(7003)})
			}

		}
	}

}

func TestPegaHpa(t *testing.T) {
	helmChartPath, err := filepath.Abs("../../../charts/pega")
	require.NoError(t, err)
	// set action execute to deploy
	options := &helm.Options{
		SetValues: map[string]string{
			"global.actions.execute": "deploy",
		},
	}
	pegaHpa := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-hpa.yaml"})
	var pegaHpaObj autoscaling.HorizontalPodAutoscaler
	hpaSlice := strings.Split(pegaHpa, "---")
	for index, hpaInfo := range hpaSlice {
		if index >= 1 && index <= 2 {
			helm.UnmarshalK8SYaml(t, hpaInfo, &pegaHpaObj)
			if index == 1 {
				VerifyPegaHpa(t, &pegaHpaObj, hpa{"pega-web-hpa", "pega-web", "Deployment", "extensions/v1beta1"})
			} else {
				VerifyPegaHpa(t, &pegaHpaObj, hpa{"pega-batch-hpa", "pega-batch", "Deployment", "extensions/v1beta1"})
			}
		}
	}
}

func TestSearchTransportService(t *testing.T) {

	helmChartPath, err := filepath.Abs("../../../charts/pega")
	require.NoError(t, err)
	// set action execute to deploy
	options := &helm.Options{
		SetValues: map[string]string{
			"global.actions.execute": "deploy",
		},
	}
	// pega-search-transport-service.yaml
	transportSearchService := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/pega-search-transport-service.yaml"})
	var transportSearchServiceObj k8score.Service
	helm.UnmarshalK8SYaml(t, transportSearchService, &transportSearchServiceObj)

	require.Equal(t, transportSearchServiceObj.Spec.Selector["component"], "Search")
	require.Equal(t, transportSearchServiceObj.Spec.Selector["app"], "pega-search")
	require.Equal(t, transportSearchServiceObj.Spec.ClusterIP, "None")
	require.Equal(t, transportSearchServiceObj.Spec.Ports[0].Name, "transport")
	require.Equal(t, transportSearchServiceObj.Spec.Ports[0].Port, int32(80))
	require.Equal(t, transportSearchServiceObj.Spec.Ports[0].TargetPort, intstr.FromInt(9300))
}

/*
func TestCredentialsSecret(){

}
func TestRegistrySecret(){

} */
