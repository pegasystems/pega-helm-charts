package pega

import (
	"fmt"
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/networking/v1"
	"path/filepath"
	"strings"
	"testing"
)

func TestPegaTierIngress(t *testing.T) {
	var supportedOperations = []string{"deploy", "install-deploy", "upgrade-deploy"}
	var supportedVendors = []string{"k8s", "eks", "gke", "aks", "pks"}
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

				yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-ingress.yaml"})
				VerifyPegaIngresses(t, yamlContent, options)
			}
		}
	}

}

// VerifyPegaIngresses - Splits the ingresses from the rendered template and asserts each ingress object
func VerifyPegaIngresses(t *testing.T, yamlContent string, options *helm.Options) {
	var pegaIngressObj v1.Ingress
	ingressSlice := strings.Split(yamlContent, "---")
	for index, ingressInfo := range ingressSlice {
		if index >= 1 && index <= 2 {
			UnmarshalK8SYaml(t, ingressInfo, &pegaIngressObj)
			if index == 1 {
				VerifyPegaIngress(t, &pegaIngressObj,
					pegaIngress{getObjName(options, "-web"), 80, 1020},
					options)
			} else {
				VerifyPegaIngress(t, &pegaIngressObj,
					pegaIngress{getObjName(options, "-stream"), 7003, 1020},
					options)
			}

		}
	}
}

func VerifyPegaIngress(t *testing.T, ingressObj *v1.Ingress, expectedIngress pegaIngress, options *helm.Options) {
	require.Equal(t, ingressObj.ObjectMeta.Name, expectedIngress.Name)
	provider := options.SetValues["global.provider"]
	if provider == "eks" {
		VerifyEKSIngress(t, ingressObj, expectedIngress)
	} else if provider == "gke" {
		VerifyGKEIngress(t, ingressObj, expectedIngress)
	} else if provider == "aks" {
		VerifyAKSIngress(t, ingressObj, expectedIngress)
	} else {
		VerifyK8SIngress(t, ingressObj, expectedIngress)
	}
}

func VerifyEKSIngress(t *testing.T, ingressObj *v1.Ingress, expectedIngress pegaIngress) {
	require.Equal(t, "alb", ingressObj.Annotations["kubernetes.io/ingress.class"])
	require.Equal(t, "[{\"HTTP\": 80}, {\"HTTPS\": 443}]", ingressObj.Annotations["alb.ingress.kubernetes.io/listen-ports"])
	require.Equal(t, "{\"Type\": \"redirect\", \"RedirectConfig\": { \"Protocol\": \"HTTPS\", \"Port\": \"443\", \"StatusCode\": \"HTTP_301\"}}", ingressObj.Annotations["alb.ingress.kubernetes.io/actions.ssl-redirect"])
	require.Equal(t, "internet-facing", ingressObj.Annotations["alb.ingress.kubernetes.io/scheme"])
	expectedStickinessAndALBAlgo := fmt.Sprint("load_balancing.algorithm.type=least_outstanding_requests,stickiness.enabled=true,stickiness.lb_cookie.duration_seconds=", expectedIngress.AlbStickiness)
	require.Equal(t, expectedStickinessAndALBAlgo,
		ingressObj.Annotations["alb.ingress.kubernetes.io/target-group-attributes"])
	require.Equal(t, "ip", ingressObj.Annotations["alb.ingress.kubernetes.io/target-type"])
	require.Equal(t, "ssl-redirect", ingressObj.Spec.Rules[0].HTTP.Paths[0].Backend.Service.Name)
	require.Equal(t, "use-annotation", ingressObj.Spec.Rules[0].HTTP.Paths[0].Backend.Service.Port.Name)
	require.Equal(t, expectedIngress.Name, ingressObj.Spec.Rules[1].HTTP.Paths[0].Backend.Service.Name)
	require.Equal(t, expectedIngress.Port, ingressObj.Spec.Rules[1].HTTP.Paths[0].Backend.Service.Port.Number)
}

func VerifyGKEIngress(t *testing.T, ingressObj *v1.Ingress, expectedIngress pegaIngress) {
	require.Equal(t, "false", ingressObj.Annotations["kubernetes.io/ingress.allow-http"])
	require.Equal(t, expectedIngress.Name, ingressObj.Spec.DefaultBackend.Service.Name)
	require.Equal(t, expectedIngress.Port, ingressObj.Spec.DefaultBackend.Service.Port.Number)
	require.Equal(t, 1, len(ingressObj.Spec.Rules))
	require.Equal(t, 1, len(ingressObj.Spec.Rules[0].HTTP.Paths))
	require.Equal(t, expectedIngress.Name, ingressObj.Spec.Rules[0].HTTP.Paths[0].Backend.Service.Name)
	require.Equal(t, expectedIngress.Port, ingressObj.Spec.Rules[0].HTTP.Paths[0].Backend.Service.Port.Number)
}

func VerifyAKSIngress(t *testing.T, ingressObj *v1.Ingress, expectedIngress pegaIngress) {
	require.Equal(t, "azure/application-gateway", ingressObj.Annotations["kubernetes.io/ingress.class"])
	require.Equal(t, "true", ingressObj.Annotations["appgw.ingress.kubernetes.io/cookie-based-affinity"])
	require.Equal(t, expectedIngress.Name, ingressObj.Spec.Rules[0].HTTP.Paths[0].Backend.Service.Name)
	require.Equal(t, expectedIngress.Port, ingressObj.Spec.Rules[0].HTTP.Paths[0].Backend.Service.Port.Number)
}

// VerifyPegaIngress - Performs Pega Ingress assertions with the values as provided in default values.yaml
func VerifyK8SIngress(t *testing.T, ingressObj *v1.Ingress, expectedIngress pegaIngress) {
	require.Equal(t, "traefik", ingressObj.Annotations["kubernetes.io/ingress.class"])
	require.Equal(t, expectedIngress.Name, ingressObj.Spec.Rules[0].HTTP.Paths[0].Backend.Service.Name)
	require.Equal(t, expectedIngress.Port, ingressObj.Spec.Rules[0].HTTP.Paths[0].Backend.Service.Port.Number)
}

type pegaIngress struct {
	Name          string
	Port          int32
	AlbStickiness int32
}
