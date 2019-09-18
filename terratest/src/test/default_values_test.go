package test

import (
	"path/filepath"
	"testing"
	"fmt"
	"encoding/json"
	"github.com/stretchr/testify/require"
	//appsv1 "k8s.io/api/apps/v1"
	"github.com/ghodss/yaml"
	"github.com/gruntwork-io/terratest/modules/helm"
)

// This file contains examples of how to use terratest to test helm chart template logic by rendering the templates
// using `helm template`, and then reading in the rendered templates.
// There are two tests:
// - TestHelmBasicExampleTemplateRenderedDeployment: An example of how to read in the rendered object and check the
//   computed values.
// - TestHelmBasicExampleTemplateRequiredTemplateArgs: An example of how to check that the required args are indeed
//   required for the template to render.

// An example of how to verify the rendered template object of a Helm Chart given various inputs.
func TestDefaultAction(t *testing.T) {
	t.Parallel()

	// Path to the helm chart we will test
	helmChartPath, err := filepath.Abs("../examples/pega-helm/src/main/helm/pega")
	require.NoError(t, err)

	options := &helm.Options{
		SetValues: map[string]string{
			"x": "y",
		},
	}

	values := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/values.yaml"})

	var body interface{}

	if err := yaml.Unmarshal([]byte(values), &body); err != nil {
        panic(err)
    }

    body = convert(body)

    if b, err := json.Marshal(body); err != nil {
        panic(err)
    } else {
		fmt.Printf("Output: %s\n", b)
	//	expectedAction := "deploy"
	    //require.Equal(t, b["provider"], expectedAction)
    }


	// Now we use kubernetes/client-go library to render the template output into the Deployment struct. This will
	// ensure the Deployment resource is rendered correctly.
	
//	helm.UnmarshalK8SYaml(t, values, &defaultValues)
	// Finally, we verify the deployment pod template spec is set to the expected container image value
	
			fmt.Println("askljASKLdhLKSdALXNDFKSKFSKDFSDKFSDNKFKDK")

}

func convert(i interface{}) interface{} {
    switch x := i.(type) {
    case map[interface{}]interface{}:
        m2 := map[string]interface{}{}
        for k, v := range x {
            m2[k.(string)] = convert(v)
        }
        return m2
    case []interface{}:
        for i, v := range x {
            x[i] = convert(v)
        }
    }
    return i
}