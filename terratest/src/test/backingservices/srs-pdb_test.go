package backingservices

import (
	"github.com/stretchr/testify/require"
	"k8s.io/api/policy/v1beta1"
	"testing"
)

func TestSRSServicePDB(t *testing.T){

	helmChartParser := NewHelmConfigParser(
		NewHelmTestFromTemplate(t, helmChartRelativePath, map[string]string{
			"srs.enabled": "true",
			"srs.deploymentName": "test-srs",
			"srs.srsStorage.tls.enabled": "false",
		},
			[]string{"charts/srs/templates/srsservice_poddisruptionbudget.yaml"}),
	)

	var pdbObj v1beta1.PodDisruptionBudget
	helmChartParser.getResourceYAML(SearchResourceOption{
		Name: "test-srs",
		Kind: "PodDisruptionBudget",
	}, &pdbObj)
	VerifySRSServicePDB(t, &pdbObj)
}

func TestSRSServicePDBWithESInternetAccess(t *testing.T){

	helmChartParser := NewHelmConfigParser(
		NewHelmTestFromTemplate(t, helmChartRelativePath, map[string]string{
			"srs.enabled": "true",
			"srs.deploymentName": "test-srs",
			"srs.srsStorage.tls.enabled": "false",
			"srs.srsStorage.requireInternetAccess": "true",
		},
			[]string{"charts/srs/templates/srsservice_poddisruptionbudget.yaml"}),
	)

	var pdbObj v1beta1.PodDisruptionBudget
	helmChartParser.getResourceYAML(SearchResourceOption{
		Name: "test-srs",
		Kind: "PodDisruptionBudget",
	}, &pdbObj)
	VerifySRSServicePDB(t, &pdbObj)
}

func TestSRSServicePDBWithESInternetAccessWithExternalES(t *testing.T){

	helmChartParser := NewHelmConfigParser(
		NewHelmTestFromTemplate(t, helmChartRelativePath, map[string]string{
			"srs.enabled": "true",
			"srs.deploymentName": "test-srs",
			"srs.srsStorage.requireInternetAccess": "true",
			"srs.srsStorage.tls.enabled": "false",
			"srs.srsStorage.provisionInternalESCluster": "false",
			"srs.srsStorage.domain": "es.acme.io",
			"srs.srsStorage.port": "8008",
			"srs.srsStorage.protocol": "https",
			"srs.srsStorage.basicAuthentication.enabled": "false",
		},
			[]string{"charts/srs/templates/srsservice_poddisruptionbudget.yaml"}),
	)

	var pdbObj v1beta1.PodDisruptionBudget
	helmChartParser.getResourceYAML(SearchResourceOption{
		Name: "test-srs",
		Kind: "PodDisruptionBudget",
	}, &pdbObj)
	VerifySRSServicePDBWithEgressLabel(t, &pdbObj)
}

func VerifySRSServicePDB(t *testing.T, servicePDBObj *v1beta1.PodDisruptionBudget) {
	require.Equal(t, "srs-service", servicePDBObj.Spec.Selector.MatchLabels["app.kubernetes.io/name"] )
	require.Equal(t, "", servicePDBObj.Spec.Selector.MatchLabels["networking/allow-internet-egress"])
}

func VerifySRSServicePDBWithEgressLabel(t *testing.T, servicePDBObj *v1beta1.PodDisruptionBudget) {
	require.Equal(t, "srs-service", servicePDBObj.Spec.Selector.MatchLabels["app.kubernetes.io/name"] )
	require.Equal(t, "true", servicePDBObj.Spec.Selector.MatchLabels["networking/allow-internet-egress"])
}