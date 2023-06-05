package pega

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	"path/filepath"
	"testing"
)



func TestPegaProviderValidate_WithValidProvider(t *testing.T){
	
	var supportedVendors = []string{"k8s","openshift","eks","gke","aks","pks"}
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for _,provider := range supportedVendors{

		var options = &helm.Options{			
			SetValues: map[string]string{
				"global.provider":        provider,
				"global.actions.execute": "deploy",
		 	},
	    }

		yamlContent, err := RenderTemplateE(t, options, helmChartPath, []string{"templates/pega-provider-validate.yaml"})
		require.Contains(t,yamlContent,"could not find template templates/pega-provider-validate.yaml")
		require.Contains(t,err.Error(),"could not find template templates/pega-provider-validate.yaml")
		
	}

}

func TestPegaProviderValidate_WithInvalidProvider(t *testing.T){

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)


	var options = &helm.Options{			
		SetValues: map[string]string{
			"global.provider":        "invalidProvider",
	 	},
	}

	yamlContent, err := RenderTemplateE(t, options, helmChartPath, []string{"templates/pega-provider-validate.yaml"})
	require.Contains(t,yamlContent,"global.provider must be one of")
	require.Contains(t,err.Error(),"global.provider must be one of")
		
}

