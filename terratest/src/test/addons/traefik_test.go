package addons

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	v12 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"path/filepath"
	"testing"
)

func TestShouldNotContainTraefikIfDisabled(t *testing.T) {
	t.Parallel()
	helmChartPath, err := filepath.Abs(helmChartRelativePath)

	require.NoError(t, err)
	options := &helm.Options{
		SetValues: map[string]string{
			"traefik.enabled": "false",
		},
	}

	helmChart := NewHelmConfigParser(t, options, helmChartPath)

	for _, i := range traefikResources {
		require.False(t, helmChart.contains(SearchResourceOption{
			name: i.name,
			kind: i.kind,
		}))
	}
}

func TestTraefikShouldContainAllResources(t *testing.T) {
	t.Parallel()
	helmChartPath, err := filepath.Abs(helmChartRelativePath)

	require.NoError(t, err)
	options := &helm.Options{
		SetValues: map[string]string{
			"traefik.enabled": "true",
		},
	}

	helmChart := NewHelmConfigParser(t, options, helmChartPath)

	for _, i := range traefikResources {
		require.True(t, helmChart.contains(SearchResourceOption{
			name: i.name,
			kind: i.kind,
		}))
	}
}

func Test_shouldBeLoadBalancer(t *testing.T) {
	helmChartPath, err := filepath.Abs(helmChartRelativePath)
	require.NoError(t, err)

	helmChartParser := NewHelmConfigParser(t, &helm.Options{
		SetValues: map[string]string{
			"traefik.enabled":     "true",
			"traefik.serviceType": "LoadBalancer",
		},
	}, helmChartPath)

	var service *v1.Service
	helmChartParser.find(SearchResourceOption{
		name: "release-name-traefik",
		kind: "Service",
	}, &service)

	serviceType := service.Spec.Type
	require.Equal(t, "LoadBalancer", string(serviceType))
}

func Test_shouldBeNodePort(t *testing.T) {
	helmChartPath, err := filepath.Abs(helmChartRelativePath)
	require.NoError(t, err)

	helmChartParser := NewHelmConfigParser(t, &helm.Options{
		SetValues: map[string]string{
			"traefik.enabled":     "true",
			"traefik.serviceType": "NodePort",
		},
	}, helmChartPath)

	var service *v1.Service
	helmChartParser.find(SearchResourceOption{
		name: "release-name-traefik",
		kind: "Service",
	}, &service)

	serviceType := service.Spec.Type
	require.Equal(t, "NodePort", string(serviceType))
	require.Equal(t, 30080, int(service.Spec.Ports[0].NodePort))
	require.Equal(t, 30443, int(service.Spec.Ports[1].NodePort))
}

func Test_hasRoleWhenRbacEnabled(t *testing.T) {
	helmChartPath, err := filepath.Abs(helmChartRelativePath)

	require.NoError(t, err)
	options := &helm.Options{
		SetValues: map[string]string{
			"traefik.enabled":      "true",
			"traefik.rbac.enabled": "true",
		},
	}

	helmChart := NewHelmConfigParser(t, options, helmChartPath)

	require.True(t, helmChart.contains(SearchResourceOption{
		name: "release-name-traefik",
		kind: "ClusterRole",
	}))

	require.True(t, helmChart.contains(SearchResourceOption{
		name: "release-name-traefik",
		kind: "ServiceAccount",
	}))

	require.True(t, helmChart.contains(SearchResourceOption{
		name: "release-name-traefik",
		kind: "ClusterRoleBinding",
	}))
}

func Test_noRoleWhenRbacDisabled(t *testing.T) {
	helmChartPath, err := filepath.Abs(helmChartRelativePath)

	require.NoError(t, err)
	options := &helm.Options{
		SetValues: map[string]string{
			"traefik.enabled":      "true",
			"traefik.rbac.enabled": "false",
		},
	}

	helmChart := NewHelmConfigParser(t, options, helmChartPath)

	require.False(t, helmChart.contains(SearchResourceOption{
		name: "release-name-traefik",
		kind: "ClusterRole",
	}))
}

func Test_hasSecretWhenSSLEnabled(t *testing.T) {
	helmChartPath, err := filepath.Abs(helmChartRelativePath)

	require.NoError(t, err)
	options := &helm.Options{
		SetValues: map[string]string{
			"traefik.enabled":     "true",
			"traefik.ssl.enabled": "true",
		},
	}

	helmChart := NewHelmConfigParser(t, options, helmChartPath)

	var deployment v12.Deployment
	helmChart.find(SearchResourceOption{
		name: "release-name-traefik",
		kind: "Deployment",
	}, &deployment)

	require.True(t, helmChart.contains(SearchResourceOption{
		name: "release-name-traefik-default-cert",
		kind: "Secret",
	}))

	require.Equal(t, "ssl", deployment.Spec.Template.Spec.Volumes[1].Name)
	require.Equal(t, "release-name-traefik-default-cert", deployment.Spec.Template.Spec.Volumes[1].Secret.SecretName)
}

func Test_hasNoSecretWhenSSLEnabled(t *testing.T) {
	helmChartPath, err := filepath.Abs(helmChartRelativePath)

	require.NoError(t, err)
	options := &helm.Options{
		SetValues: map[string]string{
			"traefik.enabled":     "true",
			"traefik.ssl.enabled": "false",
		},
	}

	helmChart := NewHelmConfigParser(t, options, helmChartPath)

	require.False(t, helmChart.contains(SearchResourceOption{
		name: "release-name-traefik-default-cert",
		kind: "Secret",
	}))
}

func Test_checkResourceRequests(t *testing.T) {
	helmChartPath, err := filepath.Abs(helmChartRelativePath)

	require.NoError(t, err)
	options := &helm.Options{
		SetValues: map[string]string{
			"traefik.enabled":                   "true",
			"traefik.resources.requests.cpu":    "300m",
			"traefik.resources.requests.memory": "300Mi",
		},
	}

	helmChart := NewHelmConfigParser(t, options, helmChartPath)

	var deployment v12.Deployment
	helmChart.find(SearchResourceOption{
		name: "release-name-traefik",
		kind: "Deployment",
	}, &deployment)

	require.Equal(t, "300m", deployment.Spec.Template.Spec.Containers[0].Resources.Requests.Cpu().String())
	require.Equal(t, "300Mi", deployment.Spec.Template.Spec.Containers[0].Resources.Requests.Memory().String())
}

func Test_checkResourceLimits(t *testing.T) {
	helmChartPath, err := filepath.Abs(helmChartRelativePath)

	require.NoError(t, err)
	options := &helm.Options{
		SetValues: map[string]string{
			"traefik.enabled":                 "true",
			"traefik.resources.limits.cpu":    "600m",
			"traefik.resources.limits.memory": "600Mi",
		},
	}

	helmChart := NewHelmConfigParser(t, options, helmChartPath)

	var deployment v12.Deployment
	helmChart.find(SearchResourceOption{
		name: "release-name-traefik",
		kind: "Deployment",
	}, &deployment)

	require.Equal(t, "600m", deployment.Spec.Template.Spec.Containers[0].Resources.Limits.Cpu().String())
	require.Equal(t, "600Mi", deployment.Spec.Template.Spec.Containers[0].Resources.Limits.Memory().String())
}

func Test_checkDefaultResourceRequests(t *testing.T) {
	helmChartPath, err := filepath.Abs(helmChartRelativePath)

	require.NoError(t, err)
	options := &helm.Options{
		SetValues: map[string]string{
			"traefik.enabled": "true",
		},
	}

	helmChart := NewHelmConfigParser(t, options, helmChartPath)

	var deployment v12.Deployment
	helmChart.find(SearchResourceOption{
		name: "release-name-traefik",
		kind: "Deployment",
	}, &deployment)

	require.Equal(t, "200m", deployment.Spec.Template.Spec.Containers[0].Resources.Requests.Cpu().String())
	require.Equal(t, "200Mi", deployment.Spec.Template.Spec.Containers[0].Resources.Requests.Memory().String())
}

func Test_checkDefaultResourceLimits(t *testing.T) {
	helmChartPath, err := filepath.Abs(helmChartRelativePath)

	require.NoError(t, err)
	options := &helm.Options{
		SetValues: map[string]string{
			"traefik.enabled": "true",
		},
	}

	helmChart := NewHelmConfigParser(t, options, helmChartPath)

	var deployment v12.Deployment
	helmChart.find(SearchResourceOption{
		name: "release-name-traefik",
		kind: "Deployment",
	}, &deployment)

	require.Equal(t, "500m", deployment.Spec.Template.Spec.Containers[0].Resources.Limits.Cpu().String())
	require.Equal(t, "500Mi", deployment.Spec.Template.Spec.Containers[0].Resources.Limits.Memory().String())
}

var traefikResources = []SearchResourceOption{
	{
		name: "release-name-traefik",
		kind: "ConfigMap",
	},
	{
		name: "release-name-traefik",
		kind: "Deployment",
	},
	{
		name: "release-name-traefik",
		kind: "Service",
	},
	{
		name: "release-name-traefik-test",
		kind: "Pod",
	},
	{
		name: "release-name-traefik-test",
		kind: "ConfigMap",
	},
}
