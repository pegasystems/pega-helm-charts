package addons

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	"strings"
	"test"
	"testing"
)

type HelmChartParser struct {
	t              *testing.T
	slicedResource []string
}

func NewHelmConfigParser(t *testing.T, helmOptions *helm.Options, chartPath string) *HelmChartParser {
	parcedChart := helm.RenderTemplate(t, helmOptions, chartPath, []string{})
	slicedResource := strings.Split(parcedChart, "---")
	return &HelmChartParser{t: t, slicedResource: slicedResource}
}

func (p *HelmChartParser) find(searchOptions SearchResourceOption, resource interface{}) {
	var d test.DeploymentMetadata
	for _, slice := range p.slicedResource {
		helm.UnmarshalK8SYaml(p.t, slice, &d)
		if (searchOptions.kind != "" && searchOptions.kind == d.Kind) && (searchOptions.name != "" && searchOptions.name == d.Name) {
			helm.UnmarshalK8SYaml(p.t, slice, &resource)
			return
		}
	}
	p.t.Log("Resource not found.", searchOptions)
	p.t.FailNow()
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
