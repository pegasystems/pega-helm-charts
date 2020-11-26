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

func VerifySRSServicePDB(t *testing.T, servicePDBObj *v1beta1.PodDisruptionBudget) {
	require.Equal(t, servicePDBObj.Spec.Selector.MatchLabels["app.kubernetes.io/name"], "srs-service")
	require.Equal(t, servicePDBObj.Spec.Selector.MatchLabels["networking/allow-internet-egress"], "true")
}