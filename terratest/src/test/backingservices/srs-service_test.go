package backingservices

import (
	"github.com/stretchr/testify/require"
	k8score "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"testing"
)

func TestSRSService(t *testing.T){

	helmChartParser := NewHelmConfigParser(
		NewHelmTestFromTemplate(t, helmChartRelativePath, map[string]string{
			"srs.enabled": "true",
			"srs.deploymentName": "test-srs",
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

func VerifySRSService(t *testing.T, serviceObj *k8score.Service) {
	require.Equal(t, serviceObj.Spec.Selector["app.kubernetes.io/name"], "srs-service")
	require.Equal(t, serviceObj.Spec.Selector["networking/allow-internet-egress"], "true")
	require.Equal(t, serviceObj.Spec.Ports[0].Name, "rest")
	require.Equal(t, serviceObj.Spec.Ports[0].Port, int32(8080))
	require.Equal(t, serviceObj.Spec.Ports[0].TargetPort, intstr.FromInt(8080))
	require.Equal(t, serviceObj.Spec.Ports[1].Name, "management")
	require.Equal(t, serviceObj.Spec.Ports[1].Port, int32(8081))
	require.Equal(t, serviceObj.Spec.Ports[1].TargetPort, intstr.FromInt(8081))
}