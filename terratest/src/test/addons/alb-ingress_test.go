package addons

import (
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/apps/v1"
	"testing"
)

func TestShouldNotContainAlbIngressIfDisabled(t *testing.T) {
	helmChartParser := NewHelmConfigParser(
		NewHelmTest(t, helmChartRelativePath, map[string]string{
			"aws-alb-ingress-controller.enabled": "false",
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
			"aws-alb-ingress-controller.enabled": "true",
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
			"aws-alb-ingress-controller.enabled":               "true",
			"aws-alb-ingress-controller.autoDiscoverAwsRegion": "false",
			"aws-alb-ingress-controller.awsRegion":             "YOUR_EKS_CLUSTER_REGION",
		}),
	)

	var deployment *v1.Deployment
	helmChart.Find(SearchResourceOption{
		Name: "pega-aws-alb-ingress-controller",
		Kind: "Deployment",
	}, &deployment)

	require.Contains(t, deployment.Spec.Template.Spec.Containers[0].Args, "--aws-region=YOUR_EKS_CLUSTER_REGION")
}

func Test_checkAutoDiscoverAwsRegion(t *testing.T) {
	helmChart := NewHelmConfigParser(
		NewHelmTest(t, helmChartRelativePath, map[string]string{
			"aws-alb-ingress-controller.enabled":               "true",
			"aws-alb-ingress-controller.autoDiscoverAwsRegion": "true",
			"aws-alb-ingress-controller.awsRegion":             "YOUR_EKS_CLUSTER_REGION",
		}),
	)

	var deployment *v1.Deployment
	helmChart.Find(SearchResourceOption{
		Name: "pega-aws-alb-ingress-controller",
		Kind: "Deployment",
	}, &deployment)

	require.NotContains(t, deployment.Spec.Template.Spec.Containers[0].Args, "--aws-region=YOUR_EKS_CLUSTER_REGION")
}

func Test_checkSetAwsVpcID(t *testing.T) {
	helmChart := NewHelmConfigParser(
		NewHelmTest(t, helmChartRelativePath, map[string]string{
			"aws-alb-ingress-controller.enabled":              "true",
			"aws-alb-ingress-controller.autoDiscoverAwsVpcID": "false",
			"aws-alb-ingress-controller.awsVpcID":             "YOUR_EKS_CLUSTER_VPC_ID",
		}),
	)

	var deployment *v1.Deployment
	helmChart.Find(SearchResourceOption{
		Name: "pega-aws-alb-ingress-controller",
		Kind: "Deployment",
	}, &deployment)

	require.Contains(t, deployment.Spec.Template.Spec.Containers[0].Args, "--aws-vpc-id=YOUR_EKS_CLUSTER_VPC_ID")
}

func Test_checkAutoDiscoverAwsVpcID(t *testing.T) {
	helmChart := NewHelmConfigParser(
		NewHelmTest(t, helmChartRelativePath, map[string]string{
			"aws-alb-ingress-controller.enabled":              "true",
			"aws-alb-ingress-controller.autoDiscoverAwsVpcID": "true",
			"aws-alb-ingress-controller.awsVpcID":             "YOUR_EKS_CLUSTER_VPC_ID",
		}),
	)

	var deployment *v1.Deployment
	helmChart.Find(SearchResourceOption{
		Name: "pega-aws-alb-ingress-controller",
		Kind: "Deployment",
	}, &deployment)

	require.NotContains(t, deployment.Spec.Template.Spec.Containers[0].Args, "--aws-vpc-id=YOUR_EKS_CLUSTER_VPC_ID")
}

func Test_checkSetClusterName(t *testing.T) {
	helmChart := NewHelmConfigParser(
		NewHelmTest(t, helmChartRelativePath, map[string]string{
			"aws-alb-ingress-controller.enabled":     "true",
			"aws-alb-ingress-controller.clusterName": "YOUR_EKS_CLUSTER_NAME",
		}),
	)

	var deployment *v1.Deployment
	helmChart.Find(SearchResourceOption{
		Name: "pega-aws-alb-ingress-controller",
		Kind: "Deployment",
	}, &deployment)

	require.Contains(t, deployment.Spec.Template.Spec.Containers[0].Args, "--cluster-name=YOUR_EKS_CLUSTER_NAME")
}

func Test_checkSetAccessKey(t *testing.T) {
	helmChart := NewHelmConfigParser(
		NewHelmTest(t, helmChartRelativePath, map[string]string{
			"aws-alb-ingress-controller.enabled":                        "true",
			"aws-alb-ingress-controller.extraEnv.AWS_ACCESS_KEY_ID":     "YOUR_AWS_ACCESS_KEY_ID",
			"aws-alb-ingress-controller.extraEnv.AWS_SECRET_ACCESS_KEY": "YOUR_AWS_SECRET_ACCESS_KEY",
		}),
	)

	var deployment *v1.Deployment
	helmChart.Find(SearchResourceOption{
		Name: "pega-aws-alb-ingress-controller",
		Kind: "Deployment",
	}, &deployment)

	require.Equal(t, "AWS_ACCESS_KEY_ID", deployment.Spec.Template.Spec.Containers[0].Env[0].Name)
	require.Equal(t, "YOUR_AWS_ACCESS_KEY_ID", deployment.Spec.Template.Spec.Containers[0].Env[0].Value)
	require.Equal(t, "AWS_SECRET_ACCESS_KEY", deployment.Spec.Template.Spec.Containers[0].Env[1].Name)
	require.Equal(t, "YOUR_AWS_SECRET_ACCESS_KEY", deployment.Spec.Template.Spec.Containers[0].Env[1].Value)
}

var albIngressResources = []SearchResourceOption{
	{
		Name: "pega-aws-alb-ingress-controller",
		Kind: "ServiceAccount",
	},
	{
		Name: "pega-aws-alb-ingress-controller",
		Kind: "ClusterRole",
	},
	{
		Name: "pega-aws-alb-ingress-controller",
		Kind: "ClusterRoleBinding",
	},
	{
		Name: "pega-aws-alb-ingress-controller",
		Kind: "Deployment",
	},
}
