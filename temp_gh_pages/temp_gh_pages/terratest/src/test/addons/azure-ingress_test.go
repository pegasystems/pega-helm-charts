package addons

import (
	b64 "encoding/base64"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	"testing"
)

func TestShouldNotContainAzureIngressIfDisabled(t *testing.T) {
	helmChartParser := NewHelmConfigParser(
		NewHelmTest(t, helmChartRelativePath, map[string]string{
			"ingress-azure.enabled": "false",
		}),
	)

	for _, i := range azureIngressResources {
		require.False(t, helmChartParser.Contains(SearchResourceOption{
			Name: i.Name,
			Kind: i.Kind,
		}))
	}
}

func TestAzureIngressShouldContainAllResources(t *testing.T) {
	helmChartParser := NewHelmConfigParser(
		NewHelmTest(t, helmChartRelativePath, map[string]string{
			"ingress-azure.enabled": "true",
		}),
	)

	for _, i := range azureIngressResources {
		require.True(t, helmChartParser.Contains(SearchResourceOption{
			Name: i.Name,
			Kind: i.Kind,
		}))
	}
}

func TestSetValuesForAppGW(t *testing.T) {
	helmChart := NewHelmConfigParser(
		NewHelmTest(t, helmChartRelativePath, map[string]string{
			"ingress-azure.enabled":              "true",
			"ingress-azure.appgw.subscriptionId": "<YOUR.SUBSCRIPTION_ID>",
			"ingress-azure.appgw.resourceGroup":  "<RESOURCE_GROUP_NAME>",
			"ingress-azure.appgw.name":           "<APPLICATION_GATEWAY_NAME>",
		}),
	)

	var configMap *v1.ConfigMap
	helmChart.Find(SearchResourceOption{
		Name: "pega-cm-ingress-azure",
		Kind: "ConfigMap",
	}, &configMap)

	require.Equal(t, "<YOUR.SUBSCRIPTION_ID>", configMap.Data["APPGW_SUBSCRIPTION_ID"])
	require.Equal(t, "<RESOURCE_GROUP_NAME>", configMap.Data["APPGW_RESOURCE_GROUP"])
	require.Equal(t, "<APPLICATION_GATEWAY_NAME>", configMap.Data["APPGW_NAME"])
	require.Equal(t, "<APPLICATION_GATEWAY_NAME>-subnet", configMap.Data["APPGW_SUBNET_NAME"])
}

func TestSetValuesForArmAuth(t *testing.T) {
	helmChart := NewHelmConfigParser(
		NewHelmTest(t, helmChartRelativePath, map[string]string{
			"ingress-azure.enabled":            "true",
			"ingress-azure.armAuth.type":       "servicePrincipal",
			"ingress-azure.armAuth.secretJSON": b64.StdEncoding.EncodeToString([]byte("<SECRET_JSON_CREATED_USING_ABOVE_COMMAND>")),
		}),
	)

	var secret *v1.Secret
	helmChart.Find(SearchResourceOption{
		Name: "networking-appgw-k8s-azure-service-principal",
		Kind: "Secret",
	}, &secret)

	require.Equal(t, "<SECRET_JSON_CREATED_USING_ABOVE_COMMAND>", string(secret.Data["armAuth.json"]))
}

var azureIngressResources = []SearchResourceOption{
	{
		Name: "networking-appgw-k8s-azure-service-principal",
		Kind: "Secret",
	},
	{
		Name: "pega-cm-ingress-azure",
		Kind: "ConfigMap",
	},
	{
		Name: "pega-sa-ingress-azure",
		Kind: "ServiceAccount",
	},
	{
		Name: "pega-ingress-azure",
		Kind: "ClusterRole",
	},
	{
		Name: "pega-ingress-azure",
		Kind: "ClusterRoleBinding",
	},
	{
		Name: "pega-ingress-azure",
		Kind: "Deployment",
	},
}
