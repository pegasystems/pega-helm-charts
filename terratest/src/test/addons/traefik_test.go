package addons

import (
	"github.com/stretchr/testify/require"
	v12 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"testing"
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
			"traefik.enabled":     "true",
			"traefik.serviceType": "LoadBalancer",
		}),
	)

	var service *v1.Service
	helmChartParser.Find(SearchResourceOption{
		Name: "pega-traefik",
		Kind: "Service",
	}, &service)

	serviceType := service.Spec.Type
	require.Equal(t, "LoadBalancer", string(serviceType))
}

func Test_shouldBeAbleToSetUpServiceTypeAsNodePort(t *testing.T) {
	helmChartParser := NewHelmConfigParser(
		NewHelmTest(t, helmChartRelativePath, map[string]string{
			"traefik.enabled":     "true",
			"traefik.serviceType": "NodePort",
		}),
	)

	var service *v1.Service
	helmChartParser.Find(SearchResourceOption{
		Name: "pega-traefik",
		Kind: "Service",
	}, &service)

	serviceType := service.Spec.Type
	require.Equal(t, "NodePort", string(serviceType))
	require.Equal(t, 30080, int(service.Spec.Ports[0].NodePort))
	require.Equal(t, 30443, int(service.Spec.Ports[1].NodePort))
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
			"traefik.enabled":     "true",
			"traefik.ssl.enabled": "true",
		}),
	)

	var deployment v12.Deployment
	helmChartParser.Find(SearchResourceOption{
		Name: "pega-traefik",
		Kind: "Deployment",
	}, &deployment)

	require.True(t, helmChartParser.Contains(SearchResourceOption{
		Name: "pega-traefik-default-cert",
		Kind: "Secret",
	}))

	require.Equal(t, "ssl", deployment.Spec.Template.Spec.Volumes[1].Name)
	require.Equal(t, "pega-traefik-default-cert", deployment.Spec.Template.Spec.Volumes[1].Secret.SecretName)
}

func Test_hasNoSecretWhenSSLEnabled(t *testing.T) {
	helmChartParser := NewHelmConfigParser(
		NewHelmTest(t, helmChartRelativePath, map[string]string{
			"traefik.enabled":     "true",
			"traefik.ssl.enabled": "false",
		}),
	)

	require.False(t, helmChartParser.Contains(SearchResourceOption{
		Name: "pega-traefik-default-cert",
		Kind: "Secret",
	}))
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
		Kind: "ConfigMap",
	},
	{
		Name: "pega-traefik",
		Kind: "Deployment",
	},
	{
		Name: "pega-traefik",
		Kind: "Service",
	},
	{
		Name: "pega-traefik-test",
		Kind: "Pod",
	},
	{
		Name: "pega-traefik-test",
		Kind: "ConfigMap",
	},
}
