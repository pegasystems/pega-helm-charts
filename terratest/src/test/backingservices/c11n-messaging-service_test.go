package backingservices

import (
	"testing"

	"github.com/stretchr/testify/require"
	k8score "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func TestC11NMessagingService(t *testing.T) {

	helmChartParser := C11NHelmConfigParser(
		NewHelmTestFromTemplate(t, helmChartRelativePath, map[string]string{
			"c11n-messaging.enabled":        "true",
			"c11n-messaging.deploymentName": "c11n-messaging",
		},
			[]string{"charts/constellation-messaging/templates/c11n-messaging-service.yaml"}),
	)

	var c11nMessagingServiceObj k8score.Service
	helmChartParser.getResourceYAML(SearchResourceOption{
		Name: "c11n-messaging",
		Kind: "Service",
	}, &c11nMessagingServiceObj)
	VerifyC11NMessagingService(t, &c11nMessagingServiceObj)
}

func VerifyC11NMessagingService(t *testing.T, serviceObj *k8score.Service) {
	require.Equal(t, "http", serviceObj.Spec.Ports[0].Name)
	require.Equal(t, int32(443), serviceObj.Spec.Ports[0].Port)
	require.Equal(t, intstr.FromInt(3000), serviceObj.Spec.Ports[0].TargetPort)
}
