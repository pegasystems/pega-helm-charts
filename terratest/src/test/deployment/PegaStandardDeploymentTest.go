package deployment

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	k8sapps "k8s.io/api/apps/v1"
	"k8s.io/api/apps/v1beta2"
	"path/filepath"
	"strings"
	"test"
	"testing"
)

type PegaDeploymentTest struct {
	provider        string
	action          string
	t               *testing.T
	_helmOptions    *helm.Options
	_initContainers []string
	_helmChartPath  string
}

func NewPegaDeploymentTest(provider string, action string, initContainers []string, t *testing.T) *PegaDeploymentTest {
	d := &PegaDeploymentTest{
		provider:        provider,
		action:          action,
		_initContainers: initContainers,
		t:               t,
	}
	d.t.Parallel()

	helmChartPath, err := filepath.Abs(test.PegaHelmChartPath)
	d._helmChartPath = helmChartPath
	require.NoError(d.t, err)

	d._helmOptions = &helm.Options{
		SetValues: map[string]string{
			"global.provider":        d.provider,
			"global.actions.execute": d.action,
		},
	}
	return d
}

func (d *PegaDeploymentTest) Run() {
	d.splitAndVerifyPegaDeployments()
}

func (d *PegaDeploymentTest) splitAndVerifyPegaDeployments() {
	deployment := helm.RenderTemplate(d.t, d._helmOptions, d._helmChartPath, []string{"templates/pega-tier-deployment.yaml"})
	deploymentSlice := strings.Split(deployment, "---")

	var processedDeployments int
	for _, deploymentInfo := range deploymentSlice {
		if len(deploymentInfo) == 0 {
			continue
		}
		var deployMeta DeploymentMetadata
		helm.UnmarshalK8SYaml(d.t, deploymentInfo, &deployMeta)
		switch deployMeta.Name {
		case "pega-web":
			var deployment k8sapps.Deployment
			helm.UnmarshalK8SYaml(d.t, deploymentInfo, &deployment)
			NewPegaWebDeployVerifier(d.t, d._helmOptions, d._initContainers, &deployment).Verify()
			processedDeployments++
			break
		case "pega-batch":
			var deployment k8sapps.Deployment
			helm.UnmarshalK8SYaml(d.t, deploymentInfo, &deployment)
			NewPegaBatchDeployVerifier(d.t, d._helmOptions, d._initContainers, &deployment).Verify()
			processedDeployments++
			break
		case "pega-stream":
			var deployment v1beta2.StatefulSet
			helm.UnmarshalK8SYaml(d.t, deploymentInfo, &deployment)
			NewPegaStreamDeployVerifier(d.t, d._helmOptions, d._initContainers, &deployment).Verify()
			processedDeployments++
			break
		}
	}
	require.Equal(d.t, 3, processedDeployments, "Should process all deployments")
}

func (d *PegaDeploymentTest) setPodValidator(verifier PodValidator) {

}
