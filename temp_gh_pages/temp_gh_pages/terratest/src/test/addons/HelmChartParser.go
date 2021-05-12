package addons

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	"strings"
	"testing"
)

type HelmChartParser struct {
	T              *testing.T
	SlicedResource []string
}

func NewHelmConfigParser(helmTest *HelmTest) *HelmChartParser {
	parsedChart := helm.RenderTemplate(helmTest.T, helmTest.HelmOptions, helmTest.ChartPath, "pega",[]string{})
	slicedResource := strings.Split(parsedChart, "---")
	return &HelmChartParser{T: helmTest.T, SlicedResource: slicedResource}
}

func (p *HelmChartParser) Find(searchOptions SearchResourceOption, resource interface{}) {
	var d DeploymentMetadata
	for _, slice := range p.SlicedResource {
		helm.UnmarshalK8SYaml(p.T, slice, &d)
		if (searchOptions.Kind != "" && searchOptions.Kind == d.Kind) && (searchOptions.Name != "" && searchOptions.Name == d.Name) {
			helm.UnmarshalK8SYaml(p.T, slice, &resource)
			return
		}
	}
	p.T.Log("Resource not found.", searchOptions)
	p.T.FailNow()
}
func (p *HelmChartParser) Contains(searchOptions SearchResourceOption) bool {
	var d DeploymentMetadata
	for _, slice := range p.SlicedResource {
		helm.UnmarshalK8SYaml(p.T, slice, &d)
		if (searchOptions.Kind != "" && searchOptions.Kind == d.Kind) && (searchOptions.Name != "" && searchOptions.Name == d.Name) {
			return true
		}
	}
	return false
}
