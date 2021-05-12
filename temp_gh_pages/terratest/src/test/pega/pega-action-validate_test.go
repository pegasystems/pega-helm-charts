package pega

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	"path/filepath"
	"testing"
)



func TestPegaActonValidate_WithValidAction(t *testing.T){
	
	var supportedOperations =  []string{"install","upgrade","deploy","install-deploy","upgrade-deploy"}
	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)

	for _,operation := range supportedOperations{

		var options = &helm.Options{			
			SetValues: map[string]string{
				"global.provider":        "k8s",
				"global.actions.execute": operation,
		 	},
	    }

		yamlContent, err := RenderTemplateE(t, options, helmChartPath, []string{"templates/pega-action-validate.yaml"})
		require.Contains(t,yamlContent,"could not find template templates/pega-action-validate.yaml")
		require.Contains(t,err.Error(),"could not find template templates/pega-action-validate.yaml")
		
	}

}

func TestPegaActonValidate_WithInValidAction(t *testing.T){

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)


	var options = &helm.Options{			
		SetValues: map[string]string{
			"global.provider":        "k8s",
			"global.actions.execute": "invalidAction",
	 	},
	}

	yamlContent, err := RenderTemplateE(t, options, helmChartPath, []string{"templates/pega-action-validate.yaml"})
	require.Contains(t,yamlContent,"Action value is not correct")
	require.Contains(t,err.Error(),"Action value is not correct")
		
}

