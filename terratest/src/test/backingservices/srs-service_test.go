package backingservices

import (
	"testing"

	"github.com/stretchr/testify/require"
	k8score "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func TestSRSService(t *testing.T) {

	helmChartParser := NewHelmConfigParser(
		NewHelmTestFromTemplate(t, helmChartRelativePath, map[string]string{
			"srs.enabled":                "true",
			"srs.deployment.name":        "test-srs",
			"srs.srsStorage.tls.enabled": "false",
		},
			[]string{"charts/srs/templates/srsservice_service.yaml"}),
	)

	var srsServiceObj k8score.Service
	helmChartParser.getResourceYAML(SearchResourceOption{
		Name: "test-srs",
		Kind: "Service",
	}, &srsServiceObj)
	VerifySRSService(t, &srsServiceObj)
}

func TestSRSServiceWithInternetEgress(t *testing.T) {

	helmChartParser := NewHelmConfigParser(
		NewHelmTestFromTemplate(t, helmChartRelativePath, map[string]string{
			"srs.enabled":                                "true",
			"srs.deployment.name":                        "test-srs",
			"srs.srsStorage.requireInternetAccess":       "true",
			"srs.srsStorage.provisionInternalESCluster":  "false",
			"srs.srsStorage.domain":                      "es.acme.io",
			"srs.srsStorage.port":                        "8008",
			"srs.srsStorage.tls.enabled":                 "false",
			"srs.srsStorage.protocol":                    "https",
			"srs.srsStorage.basicAuthentication.enabled": "false",
		},
			[]string{"charts/srs/templates/srsservice_service.yaml"}),
	)

	var srsServiceObj k8score.Service
	helmChartParser.getResourceYAML(SearchResourceOption{
		Name: "test-srs",
		Kind: "Service",
	}, &srsServiceObj)
	VerifySRSServiceWithEgress(t, &srsServiceObj)
}

func VerifySRSService(t *testing.T, serviceObj *k8score.Service) {
	require.Equal(t, "srs-service", serviceObj.Spec.Selector["app.kubernetes.io/name"])
	require.Equal(t, "", serviceObj.Spec.Selector["networking/allow-internet-egress"])
	require.Equal(t, "rest", serviceObj.Spec.Ports[0].Name)
	require.Equal(t, int32(8080), serviceObj.Spec.Ports[0].Port)
	require.Equal(t, intstr.FromInt(8080), serviceObj.Spec.Ports[0].TargetPort)
	require.Equal(t, "http80", serviceObj.Spec.Ports[1].Name)
	require.Equal(t, int32(80), serviceObj.Spec.Ports[1].Port)
	require.Equal(t, intstr.FromInt(8080), serviceObj.Spec.Ports[1].TargetPort)
}

func VerifySRSServiceWithEgress(t *testing.T, serviceObj *k8score.Service) {
	require.Equal(t, "srs-service", serviceObj.Spec.Selector["app.kubernetes.io/name"])
	require.Equal(t, "true", serviceObj.Spec.Selector["networking/allow-internet-egress"])
	require.Equal(t, "rest", serviceObj.Spec.Ports[0].Name)
	require.Equal(t, int32(8080), serviceObj.Spec.Ports[0].Port)
	require.Equal(t, intstr.FromInt(8080), serviceObj.Spec.Ports[0].TargetPort)
	require.Equal(t, "http80", serviceObj.Spec.Ports[1].Name)
	require.Equal(t, int32(80), serviceObj.Spec.Ports[1].Port)
	require.Equal(t, intstr.FromInt(8080), serviceObj.Spec.Ports[1].TargetPort)
}
