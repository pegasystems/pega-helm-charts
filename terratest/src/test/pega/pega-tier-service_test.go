package pega

import (
	"fmt"
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	k8score "k8s.io/api/core/v1"
	intstr "k8s.io/apimachinery/pkg/util/intstr"
	"path/filepath"
	"strings"
	"testing"
)

func TestPegaTierService(t *testing.T) {
	var supportedVendors = []string{"k8s", "openshift", "eks", "gke", "aks", "pks"}
	var supportedOperations = []string{"deploy", "install-deploy", "upgrade-deploy"}
	var deploymentNames = []string{"pega", "myapp-dev"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for _, vendor := range supportedVendors {

		for _, operation := range supportedOperations {

			for _, depName := range deploymentNames {

				fmt.Println(vendor + "-" + operation)

				var options = &helm.Options{
					SetValues: map[string]string{
						"global.deployment.name":        depName,
						"global.provider":               vendor,
						"global.actions.execute":        operation,
						"installer.upgrade.upgradeType": "zero-downtime",
					},
				}

				yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-service.yaml"})
				VerifyPegaServices(t, yamlContent, options)
			}
		}
	}

}

// SplitAndVerifyPegaServices - Splits the services from the rendered template and asserts each service objects
func VerifyPegaServices(t *testing.T, yamlContent string, options *helm.Options) {
	var pegaServiceObj k8score.Service
	serviceSlice := strings.Split(yamlContent, "---")
	for index, serviceInfo := range serviceSlice {
		if index >= 1 && index <= 2 {
			UnmarshalK8SYaml(t, serviceInfo, &pegaServiceObj)

			if index == 1 {
				require.Equal(t, getObjName(options, "-web"), pegaServiceObj.ObjectMeta.Name)
				VerifyPegaService(t, &pegaServiceObj, pegaServices{getObjName(options, "-web"), int32(80), intstr.IntOrString{IntVal: 8080}}, options)
			} else {
				require.Equal(t, getObjName(options, "-stream"), pegaServiceObj.ObjectMeta.Name)
				VerifyPegaService(t, &pegaServiceObj, pegaServices{getObjName(options, "-stream"), int32(7003), intstr.IntOrString{IntVal: 7003}}, options)
			}
		}
	}
}

// VerifyPegaService - Performs Pega Service assertions with the values as provided in default values.yaml
func VerifyPegaService(t *testing.T, serviceObj *k8score.Service, expectedService pegaServices, options *helm.Options) {
	provider := options.SetValues["global.provider"]
	if provider == "k8s" {
		require.Equal(t, serviceObj.Annotations["traefik.ingress.kubernetes.io/service.sticky.cookie"], "true")
		require.Equal(t, serviceObj.Annotations["traefik.ingress.kubernetes.io/load-balancer-method"], "drr")
		require.Equal(t, serviceObj.Annotations["traefik.ingress.kubernetes.io/max-conn-amount"], "10")
		require.Equal(t, serviceObj.Annotations["traefik.ingress.kubernetes.io/service.sticky.cookie.name"], "UNIQUE-PEGA-COOKIE-NAME")
		require.Equal(t, serviceObj.Spec.Type, k8score.ServiceType("LoadBalancer"))
	} else if provider == "gke" {
		require.Equal(t, `{"ingress": true}`, serviceObj.Annotations["cloud.google.com/neg"])
		var expectedBackendConfig = fmt.Sprintf(`{"ports": {"%d": "%s"}}`, expectedService.Port, expectedService.Name)
		require.Equal(t, expectedBackendConfig, serviceObj.Annotations["cloud.google.com/backend-config"])
		require.Equal(t, serviceObj.Spec.Type, k8score.ServiceType("NodePort"))
	}
	require.Equal(t, serviceObj.Spec.Selector["app"], expectedService.Name)
	require.Equal(t, serviceObj.Spec.Ports[0].Port, expectedService.Port)
	require.Equal(t, serviceObj.Spec.Ports[0].TargetPort, expectedService.TargetPort)
}

type pegaServices struct {
	Name       string
	Port       int32
	TargetPort intstr.IntOrString
}

func TestPegaTierServiceWithCustomPorts(t *testing.T) {
	var supportedVendors = []string{"k8s", "openshift", "eks", "gke", "aks", "pks"}
	var supportedOperations = []string{"deploy", "install-deploy", "upgrade-deploy"}
	var deploymentNames = []string{"pega", "myapp-dev"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)
	var serviceObj k8score.Service
	for _, vendor := range supportedVendors {
		for _, operation := range supportedOperations {
			for _, depName := range deploymentNames {

				var options = &helm.Options{
					ValuesFiles: []string{"data/values_service_custom_ports.yaml"},
					SetValues: map[string]string{
						"global.deployment.name":        depName,
						"global.provider":               vendor,
						"global.actions.execute":        operation,
						"installer.upgrade.upgradeType": "zero-downtime",
					},
				}
				yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-service.yaml"})
				yamlSplit := strings.Split(yamlContent, "---")
				UnmarshalK8SYaml(t, yamlSplit[1], &serviceObj)
				ports := serviceObj.Spec.Ports
				require.Equal(t, 2, len(ports))
				require.Equal(t, "http", ports[0].Name)
				require.Equal(t, int32(80), ports[0].Port)
				require.Equal(t, int32(8080), ports[0].TargetPort.IntVal)
				require.Equal(t, "port1", ports[1].Name)
				require.Equal(t, int32(5005), ports[1].Port)
				require.Equal(t, int32(5005), ports[1].TargetPort.IntVal)
			}
		}
	}
}
