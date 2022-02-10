package addons

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	v12 "k8s.io/api/apps/v1"
)

func Test_shouldNotContainTraefikResourcesWhenDisabled(t *testing.T) {
	helmChartParser := NewHelmConfigParser(
		NewHelmTest(t, helmChartRelativePath, map[string]string{
			"traefik.enabled": "false",
		}),
	)

	for _, i := range traefikResources {
		require.False(t, helmChartParser.Contains(SearchResourceOption{
			Name: i.Name,
			Kind: i.Kind,
		}))
	}
}

func Test_shouldContainTraefikResourcesWhenEnabled(t *testing.T) {
	helmChartParser := NewHelmConfigParser(
		NewHelmTest(t, helmChartRelativePath, map[string]string{
			"traefik.enabled": "true",
		}),
	)

	for _, i := range traefikResources {
		require.True(t, helmChartParser.Contains(SearchResourceOption{
			Name: i.Name,
			Kind: i.Kind,
		}))
	}
}

func Test_shouldBeAbleToSetUpServiceTypeAsLoadBalancer(t *testing.T) {
	helmChartParser := NewHelmConfigParser(
		NewHelmTest(t, helmChartRelativePath, map[string]string{
			"traefik.enabled":      "true",
			"traefik.service.type": "LoadBalancer",
		}),
	)

	var d DeploymentMetadata
	var list string
	for _, slice := range helmChartParser.SlicedResource {
		helm.UnmarshalK8SYaml(helmChartParser.T, slice, &d)
		if d.Kind == "List" {
			list = slice
			break
		}
	}

	require.True(t, len(list) != 0)
	require.Contains(t, list, "LoadBalancer")
}

func Test_shouldBeAbleToSetUpServiceTypeAsNodePort(t *testing.T) {
	helmChartParser := NewHelmConfigParser(
		NewHelmTest(t, helmChartRelativePath, map[string]string{
			"traefik.enabled":      "true",
			"traefik.service.type": "NodePort",
		}),
	)

	var d DeploymentMetadata
	var list string
	for _, slice := range helmChartParser.SlicedResource {
		helm.UnmarshalK8SYaml(helmChartParser.T, slice, &d)
		if d.Kind == "List" {
			list = slice
			break
		}
	}

	require.True(t, len(list) != 0)
	require.Contains(t, list, "NodePort")
}

func Test_hasRoleWhenRbacEnabled(t *testing.T) {
	helmChartParser := NewHelmConfigParser(
		NewHelmTest(t, helmChartRelativePath, map[string]string{
			"traefik.enabled":      "true",
			"traefik.rbac.enabled": "true",
		}),
	)

	require.True(t, helmChartParser.Contains(SearchResourceOption{
		Name: "pega-traefik",
		Kind: "ClusterRole",
	}))

	require.True(t, helmChartParser.Contains(SearchResourceOption{
		Name: "pega-traefik",
		Kind: "ServiceAccount",
	}))

	require.True(t, helmChartParser.Contains(SearchResourceOption{
		Name: "pega-traefik",
		Kind: "ClusterRoleBinding",
	}))
}

func Test_noRoleWhenRbacDisabled(t *testing.T) {
	helmChartParser := NewHelmConfigParser(
		NewHelmTest(t, helmChartRelativePath, map[string]string{
			"traefik.enabled":      "true",
			"traefik.rbac.enabled": "false",
		}),
	)

	require.False(t, helmChartParser.Contains(SearchResourceOption{
		Name: "pega-traefik",
		Kind: "ClusterRole",
	}))
}

func Test_hasSecretWhenSSLEnabled(t *testing.T) {
	helmChartParser := NewHelmConfigParser(
		NewHelmTest(t, helmChartRelativePath, map[string]string{
			"traefik.enabled":                     "true",
			"traefik.ports.websecure.tls.enabled": "true",
		}),
	)

	var deployment v12.Deployment
	helmChartParser.Find(SearchResourceOption{
		Name: "pega-traefik",
		Kind: "Deployment",
	}, &deployment)

	require.Contains(t, deployment.Spec.Template.Spec.Containers[0].Args, "--entrypoints.websecure.http.tls=true")
}

func Test_hasNoSecretWhenSSLEnabled(t *testing.T) {
	helmChartParser := NewHelmConfigParser(
		NewHelmTest(t, helmChartRelativePath, map[string]string{
			"traefik.enabled":                     "true",
			"traefik.ports.websecure.tls.enabled": "false",
		}),
	)

	var deployment v12.Deployment
	helmChartParser.Find(SearchResourceOption{
		Name: "pega-traefik",
		Kind: "Deployment",
	}, &deployment)

	require.NotContains(t, deployment.Spec.Template.Spec.Containers[0].Args, "--entrypoints.websecure.http.tls=true")
}

func Test_checkResourceRequests(t *testing.T) {
	helmChartParser := NewHelmConfigParser(
		NewHelmTest(t, helmChartRelativePath, map[string]string{
			"traefik.enabled":                   "true",
			"traefik.resources.requests.cpu":    "300m",
			"traefik.resources.requests.memory": "300Mi",
		}),
	)

	var deployment v12.Deployment
	helmChartParser.Find(SearchResourceOption{
		Name: "pega-traefik",
		Kind: "Deployment",
	}, &deployment)

	require.Equal(t, "300m", deployment.Spec.Template.Spec.Containers[0].Resources.Requests.Cpu().String())
	require.Equal(t, "300Mi", deployment.Spec.Template.Spec.Containers[0].Resources.Requests.Memory().String())
}

func Test_checkResourceLimits(t *testing.T) {
	helmChartParser := NewHelmConfigParser(
		NewHelmTest(t, helmChartRelativePath, map[string]string{
			"traefik.enabled":                 "true",
			"traefik.resources.limits.cpu":    "600m",
			"traefik.resources.limits.memory": "600Mi",
		}),
	)

	var deployment v12.Deployment
	helmChartParser.Find(SearchResourceOption{
		Name: "pega-traefik",
		Kind: "Deployment",
	}, &deployment)

	require.Equal(t, "600m", deployment.Spec.Template.Spec.Containers[0].Resources.Limits.Cpu().String())
	require.Equal(t, "600Mi", deployment.Spec.Template.Spec.Containers[0].Resources.Limits.Memory().String())
}

func Test_checkDefaultResourceRequests(t *testing.T) {
	helmChartParser := NewHelmConfigParser(
		NewHelmTest(t, helmChartRelativePath, map[string]string{
			"traefik.enabled": "true",
		}),
	)

	var deployment v12.Deployment
	helmChartParser.Find(SearchResourceOption{
		Name: "pega-traefik",
		Kind: "Deployment",
	}, &deployment)

	require.Equal(t, "200m", deployment.Spec.Template.Spec.Containers[0].Resources.Requests.Cpu().String())
	require.Equal(t, "200Mi", deployment.Spec.Template.Spec.Containers[0].Resources.Requests.Memory().String())
}

func Test_checkDefaultResourceLimits(t *testing.T) {
	helmChartParser := NewHelmConfigParser(
		NewHelmTest(t, helmChartRelativePath, map[string]string{
			"traefik.enabled": "true",
		}),
	)

	var deployment v12.Deployment
	helmChartParser.Find(SearchResourceOption{
		Name: "pega-traefik",
		Kind: "Deployment",
	}, &deployment)

	require.Equal(t, "500m", deployment.Spec.Template.Spec.Containers[0].Resources.Limits.Cpu().String())
	require.Equal(t, "500Mi", deployment.Spec.Template.Spec.Containers[0].Resources.Limits.Memory().String())
}

var traefikResources = []SearchResourceOption{
	{
		Name: "pega-traefik",
		Kind: "ServiceAccount",
	},
	{
		Name: "pega-traefik",
		Kind: "ClusterRole",
	},
	{
		Name: "pega-traefik",
		Kind: "ClusterRoleBinding",
	},
	{
		Name: "pega-traefik",
		Kind: "Deployment",
	},
	{
		Name: "pega-traefik",
		Kind: "List",
	},
	{
		Name: "pega-traefik-dashboard",
		Kind: "IngressRoute",
	},
}
