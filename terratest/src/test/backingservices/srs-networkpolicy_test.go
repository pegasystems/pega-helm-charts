package backingservices

import (
	"github.com/stretchr/testify/require"
	networkingv1 "k8s.io/api/networking/v1"
	"testing"
)

func TestSRSServiceNetworkPolicy(t *testing.T){

	helmChartParser := NewHelmConfigParser(
		NewHelmTestFromTemplate(t, helmChartRelativePath, map[string]string{
			"srs.enabled": "true",
			"srs.deploymentName": "test-srs",
			"srs.srsStorage.tls.enabled": "false",
		},
			[]string{"charts/srs/templates/srsservice_networkpolicy.yaml"}),
	)

	var networkPolicyObj networkingv1.NetworkPolicy
	helmChartParser.getResourceYAML(SearchResourceOption{
		Name: "test-srs-networkpolicy",
		Kind: "NetworkPolicy",
	}, &networkPolicyObj)
	VerifySRSServiceNetworkPolicy(t, &networkPolicyObj, false)
}

func TestSRSServiceNetworkPolicyWithProvisionInternalESCluster(t *testing.T){

	helmChartParser := NewHelmConfigParser(
		NewHelmTestFromTemplate(t, helmChartRelativePath, map[string]string{
			"srs.enabled": "true",
			"srs.deploymentName": "test-srs",
			"srs.srsStorage.tls.enabled": "false",
			"srs.srsStorage.requireInternetAccess": "true",
			"srs.srsStorage.provisionInternalESCluster": "true",
		},
			[]string{"charts/srs/templates/srsservice_networkpolicy.yaml"}),
	)

	var networkPolicyObj networkingv1.NetworkPolicy
	helmChartParser.getResourceYAML(SearchResourceOption{
		Name: "test-srs-networkpolicy",
		Kind: "NetworkPolicy",
	}, &networkPolicyObj)
	VerifySRSServiceNetworkPolicy(t, &networkPolicyObj, false)
}

func TestSRSServiceNetworkPolicyWithProvisionInternalESClusterFalse(t *testing.T){

	helmChartParser := NewHelmConfigParser(
		NewHelmTestFromTemplate(t, helmChartRelativePath, map[string]string{
			"srs.enabled": "true",
			"srs.deploymentName": "test-srs",
			"srs.srsStorage.requireInternetAccess": "true",
			"srs.srsStorage.provisionInternalESCluster": "false",
			"srs.srsStorage.tls.enabled": "false",
			"srs.srsStorage.domain": "es.acme.io",
			"srs.srsStorage.port": "8008",
			"srs.srsStorage.protocol": "https",
			"srs.srsStorage.basicAuthentication.enabled": "false",
		},
			[]string{"charts/srs/templates/srsservice_networkpolicy.yaml"}),
	)

	var networkPolicyObj networkingv1.NetworkPolicy
	helmChartParser.getResourceYAML(SearchResourceOption{
		Name: "test-srs-networkpolicy",
		Kind: "NetworkPolicy",
	}, &networkPolicyObj)
	VerifySRSServiceNetworkPolicy(t, &networkPolicyObj, true)
}

func VerifySRSServiceNetworkPolicy(t *testing.T, networkPolicySvcObj *networkingv1.NetworkPolicy, internetEgress bool) {
	require.Equal(t, "srs-service", networkPolicySvcObj.Spec.PodSelector.MatchLabels["app.kubernetes.io/name"])


	require.ElementsMatch(t, []networkingv1.PolicyType{ "Ingress","Egress"}, networkPolicySvcObj.Spec.PolicyTypes)
	// ingress
	require.Equal(t, "0.0.0.0/0", networkPolicySvcObj.Spec.Ingress[0].From[0].IPBlock.CIDR)
	require.Equal(t, int32(8080), networkPolicySvcObj.Spec.Ingress[0].Ports[0].Port.IntVal)
	require.Equal(t, "srs-ops", networkPolicySvcObj.Spec.Ingress[1].From[0].PodSelector.MatchLabels["app.kubernetes.io/name"])
	require.Equal(t, int32(8080), networkPolicySvcObj.Spec.Ingress[1].Ports[0].Port.IntVal)
	//egress
	require.Equal(t, "elasticsearch-master", networkPolicySvcObj.Spec.Egress[0].To[0].PodSelector.MatchLabels["app"])
	require.Equal(t, int32(9200), networkPolicySvcObj.Spec.Egress[0].Ports[0].Port.IntVal)

	//verify internet egress
	if internetEgress {
		require.Equal(t, "true", networkPolicySvcObj.Spec.PodSelector.MatchLabels["networking/allow-internet-egress"])
		require.Equal(t, "true", networkPolicySvcObj.Spec.Egress[2].To[0].PodSelector.MatchLabels["networking/allow-internet-egress"])
	} else {
		require.Equal(t, "", networkPolicySvcObj.Spec.PodSelector.MatchLabels["networking/allow-internet-egress"])
	}

}