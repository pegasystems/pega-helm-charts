package addons

import (
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
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

func Test_checkSetAwsVpcID(t *testing.T) {
	helmChart := NewHelmConfigParser(
		NewHelmTest(t, helmChartRelativePath, map[string]string{
			"aws-load-balancer-controller.enabled":              "true",
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

func Test_checkSetServiceAnnotation(t *testing.T) {
	helmChart := NewHelmConfigParser(
		NewHelmTest(t, helmChartRelativePath, map[string]string{
			"aws-load-balancer-controller.enabled":              "true",
		}),
	)
	var serviceAccount *corev1.ServiceAccount
	helmChart.Find(SearchResourceOption{
		Name: "pega-aws-load-balancer-controller",
		Kind: "ServiceAccount",
	}, &serviceAccount)

	require.Contains(t, serviceAccount.ObjectMeta.Annotations["eks.amazonaws.com/role-arn"], "YOUR_IAM_ROLE_ARN")
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
		Name: "pega-aws-load-balancer-controller-role",
		Kind: "ClusterRole",
	},
	{
		Name: "pega-aws-load-balancer-controller-rolebinding",
		Kind: "ClusterRoleBinding",
	},
	{
		Name: "pega-aws-load-balancer-controller",
		Kind: "Deployment",
	},
}
