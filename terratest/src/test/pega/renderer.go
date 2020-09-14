package pega

import (
	
	"github.com/gruntwork-io/terratest/modules/helm"
	"testing"

)


func RenderTemplate(t *testing.T, options *helm.Options, helmChartPath string, templates []string, extraHelmArgs ...string) string{

	return helm.RenderTemplate(t, options, helmChartPath, PegaHelmRelease, templates, extraHelmArgs...)

}

func RenderTemplateE(t *testing.T, options *helm.Options, helmChartPath string, templates []string) (string, error) {

	return helm.RenderTemplateE(t, options, helmChartPath, PegaHelmRelease, templates)

}

func UnmarshalK8SYaml(t *testing.T, yamlData string, destinationObj interface{}){

	helm.UnmarshalK8SYaml(t,yamlData,destinationObj)

}