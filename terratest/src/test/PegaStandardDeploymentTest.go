package test

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	k8sapps "k8s.io/api/apps/v1"
	"k8s.io/api/apps/v1beta2"
	"path/filepath"
	"strings"
	. "test/deployment"
	. "test/verifier"
	"testing"
)

type PegaStandardDeploymentTest struct {
	provider        string
	action          string
	verifiers       []VerifierImpl
	t               *testing.T
	_helmOptions    *helm.Options
	_initContainers []string
	_helmChartPath  string
}

func NewPegaStandardDeploymentTest(provider string, action string, initContainers []string, t *testing.T) *PegaStandardDeploymentTest {
	d := &PegaStandardDeploymentTest{
		provider:        provider,
		action:          action,
		_initContainers: initContainers,
		t:               t,
	}
	d.t.Parallel()

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
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

func (d *PegaStandardDeploymentTest) Run() {
	d.splitAndVerifyPegaDeployments()
}

func (d *PegaStandardDeploymentTest) splitAndVerifyPegaDeployments() {
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
