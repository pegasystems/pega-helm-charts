package addons

import (
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/apps/v1"
	"testing"
)

func TestShouldNotContainAlbIngressIfDisabled(t *testing.T) {
	helmChartParser := NewHelmConfigParser(
		NewHelmTest(t, helmChartRelativePath, map[string]string{
			"aws-load-balancer-controller.enabled": "false",
		}),
	)

	for _, i := range albIngressResources {
		require.False(t, helmChartParser.Contains(SearchResourceOption{
			Name: i.Name,
			Kind: i.Kind,
		}))
	}
}

func TestAlbIngressShouldContainAllResources(t *testing.T) {
	helmChartParser := NewHelmConfigParser(
		NewHelmTest(t, helmChartRelativePath, map[string]string{
			"aws-load-balancer-controller.enabled": "true",
		}),
	)

	for _, i := range albIngressResources {
		require.True(t, helmChartParser.Contains(SearchResourceOption{
			Name: i.Name,
			Kind: i.Kind,
		}))
	}
}

func Test_checkSetAwsRegion(t *testing.T) {
	helmChart := NewHelmConfigParser(
		NewHelmTest(t, helmChartRelativePath, map[string]string{
			"aws-load-balancer-controller.enabled":               "true",
			"aws-load-balancer-controller.autoDiscoverAwsRegion": "false",
			"aws-load-balancer-controller.region":             "YOUR_EKS_CLUSTER_REGION",
		}),
	)

	var deployment *v1.Deployment
	helmChart.Find(SearchResourceOption{
		Name: "pega-aws-load-balancer-controller",
		Kind: "Deployment",
	}, &deployment)

	require.Contains(t, deployment.Spec.Template.Spec.Containers[0].Args, "--aws-region=YOUR_EKS_CLUSTER_REGION")
}

func Test_checkAutoDiscoverAwsRegion(t *testing.T) {
	helmChart := NewHelmConfigParser(
		NewHelmTest(t, helmChartRelativePath, map[string]string{
			"aws-load-balancer-controller.enabled":               "true",
			"aws-load-balancer-controller.autoDiscoverAwsRegion": "true",
			"aws-load-balancer-controller.region":             "YOUR_EKS_CLUSTER_REGION",
		}),
	)

	var deployment *v1.Deployment
	helmChart.Find(SearchResourceOption{
		Name: "pega-aws-load-balancer-controller",
		Kind: "Deployment",
	}, &deployment)

	require.NotContains(t, deployment.Spec.Template.Spec.Containers[0].Args, "--aws-region=YOUR_EKS_CLUSTER_REGION")
}

func Test_checkSetAwsVpcID(t *testing.T) {
	helmChart := NewHelmConfigParser(
		NewHelmTest(t, helmChartRelativePath, map[string]string{
			"aws-load-balancer-controller.enabled":              "true",
			"aws-load-balancer-controller.autoDiscoverAwsVpcID": "false",
			"aws-load-balancer-controller.vpcId":             "YOUR_EKS_CLUSTER_VPC_ID",
		}),
	)

	var deployment *v1.Deployment
	helmChart.Find(SearchResourceOption{
		Name: "pega-aws-load-balancer-controller",
		Kind: "Deployment",
	}, &deployment)

	require.Contains(t, deployment.Spec.Template.Spec.Containers[0].Args, "--aws-vpc-id=YOUR_EKS_CLUSTER_VPC_ID")
}

func Test_checkAutoDiscoverAwsVpcID(t *testing.T) {
	helmChart := NewHelmConfigParser(
		NewHelmTest(t, helmChartRelativePath, map[string]string{
			"aws-load-balancer-controller.enabled":              "true",
			"aws-load-balancer-controller.autoDiscoverAwsVpcID": "true",
			"aws-load-balancer-controller.vpcId":             "YOUR_EKS_CLUSTER_VPC_ID",
		}),
	)

	var deployment *v1.Deployment
	helmChart.Find(SearchResourceOption{
		Name: "pega-aws-load-balancer-controller",
		Kind: "Deployment",
	}, &deployment)

	require.NotContains(t, deployment.Spec.Template.Spec.Containers[0].Args, "--aws-vpc-id=YOUR_EKS_CLUSTER_VPC_ID")
}

func Test_checkSetClusterName(t *testing.T) {
	helmChart := NewHelmConfigParser(
		NewHelmTest(t, helmChartRelativePath, map[string]string{
			"aws-load-balancer-controller.enabled":     "true",
			"aws-load-balancer-controller.clusterName": "YOUR_EKS_CLUSTER_NAME",
		}),
	)

	var deployment *v1.Deployment
	helmChart.Find(SearchResourceOption{
		Name: "pega-aws-load-balancer-controller",
		Kind: "Deployment",
	}, &deployment)

	require.Contains(t, deployment.Spec.Template.Spec.Containers[0].Args, "--cluster-name=YOUR_EKS_CLUSTER_NAME")
}

var albIngressResources = []SearchResourceOption{
	{
		Name: "pega-aws-load-balancer-controller",
		Kind: "ServiceAccount",
	},
	{
		Name: "pega-aws-load-balancer-controller",
		Kind: "ClusterRole",
	},
	{
		Name: "pega-aws-load-balancer-controller",
		Kind: "ClusterRoleBinding",
	},
	{
		Name: "pega-aws-load-balancer-controller",
		Kind: "Deployment",
	},
}
