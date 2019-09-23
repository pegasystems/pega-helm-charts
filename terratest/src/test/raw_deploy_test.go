package test

import (
	"path/filepath"
	"testing"

	"fmt"

	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"

	//k8sresource "k8s.io/apimachinery/pkg/api/resource"

	"github.com/gruntwork-io/terratest/modules/helm"
)

const pegaHelmChartPath = "../../../charts/pega"

func TestWebDeployment(t *testing.T) {
	t.Parallel()
	
	fmt.Println("hello-1")

	// Path to the helm chart we will test
	helmChartPath, err := filepath.Abs(pegaHelmChartPath)
	require.NoError(t, err)

	// set action execute to install
	options := &helm.Options{
		SetValues: map[string]string{
			"global.actions.execute": "deploy",
		},
	}

	// with action as 'install-deploy' below templates should not be rendered
	output := helm.RenderTemplate(t, options, helmChartPath, []string{
		"templates/pega-tier-deployment.yaml",
	})
	//fmt.Println("hello")
	//fmt.Println(output)
	var deployment appsv1.Deployment
	helm.UnmarshalK8SYaml(t, output, &deployment)

	var initContainers = []string{"wait-for-pegasearch", "wait-for-cassandra"}
	excepectedDeployment := pegaDeployment{"pega-web", initContainers, "WebUser"}

	VerifyPegaDeployment(t, &deployment, &excepectedDeployment)
	// any specific assertions goes here
}

func TestBatchDeployment(t *testing.T) {
	t.Parallel()
	t.Skip()
	// Path to the helm chart we will test
	helmChartPath, err := filepath.Abs(pegaHelmChartPath)
	require.NoError(t, err)

	// set action execute to install
	options := &helm.Options{
		SetValues: map[string]string{
			"global.actions.execute": "deploy",
		},
	}

	// with action as 'install-deploy' below templates should not be rendered
	output := helm.RenderTemplate(t, options, helmChartPath, []string{
		"templates/pega-tier-deployment.yaml",
	})
	//fmt.Println("hello")
	//fmt.Println(output)
	var deployment appsv1.Deployment
	helm.UnmarshalK8SYaml(t, output, &deployment)

	var initContainers = []string{"wait-for-pegasearch", "wait-for-cassandra"}
	excepectedDeployment := pegaDeployment{"pega-batch", initContainers, "BackgroundProcessing,Search,Batch,RealTime,Custom1,Custom2,Custom3,Custom4,Custom5,BIX"}

	VerifyPegaDeployment(t, &deployment, &excepectedDeployment)
}
