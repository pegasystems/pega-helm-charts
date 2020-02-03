package addons

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	"path/filepath"
	"strings"
	"test"
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

	require.True(t, helmChart.contains(SearchResourceOption{
		name: "release-name-traefik-test",
		kind: "ConfigMap",
	}))
	require.True(t, helmChart.contains(SearchResourceOption{
		name: "release-name-traefik",
		kind: "ClusterRole",
	}))
	require.True(t, helmChart.contains(SearchResourceOption{
		name: "release-name-traefik",
		kind: "ConfigMap",
	}))
	require.True(t, helmChart.contains(SearchResourceOption{
		name: "release-name-traefik",
		kind: "ServiceAccount",
	}))
	require.True(t, helmChart.contains(SearchResourceOption{
		name: "release-name-traefik",
		kind: "ClusterRoleBinding",
	}))
	require.True(t, helmChart.contains(SearchResourceOption{
		name: "release-name-traefik",
		kind: "Service",
	}))
	require.True(t, helmChart.contains(SearchResourceOption{
		name: "release-name-traefik-test",
		kind: "Pod",
	}))
	require.True(t, helmChart.contains(SearchResourceOption{
		name: "release-name-traefik",
		kind: "Deployment",
	}))

	//var deployment appsv1.Deployment

	//
	//deploymentSlice := strings.Split(deployment, "---")
	//var traefikEnabled

	//func VerifyTraefikEnabled (t testing.T, helmChartPath string, options helm.Options){
	//	traefik :=
	//	}
}

type HelmChartParser struct {
	t              *testing.T
	slicedResource []string
}

func NewHelmConfigParser(t *testing.T, helmOptions *helm.Options, chartPath string) *HelmChartParser {
	parcedChart := helm.RenderTemplate(t, helmOptions, chartPath, []string{})
	slicedResource := strings.Split(parcedChart, "---")
	return &HelmChartParser{t: t, slicedResource: slicedResource}
}

func (p *HelmChartParser) find(searchOptions SearchResourceOption, resource *interface{}) {

}
func (p *HelmChartParser) contains(searchOptions SearchResourceOption) bool {
	var d test.DeploymentMetadata
	for _, slice := range p.slicedResource {
		helm.UnmarshalK8SYaml(p.t, slice, &d)
		if (searchOptions.kind != "" && searchOptions.kind == d.Kind) && (searchOptions.name != "" && searchOptions.name == d.Name) {
			return true
		}
	}
	return false
}

type SearchResourceOption struct {
	name string
	kind string
}

var traefikResources = []SearchResourceOption{
	{
		name: "",
		kind: "",
	},
}
