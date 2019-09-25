package test

import (
	"testing"

	//"fmt"

	k8sbatch "k8s.io/api/batch/v1"

	//k8sresource "k8s.io/apimachinery/pkg/api/resource"

	"github.com/gruntwork-io/terratest/modules/helm"
)

const pegaHelmChartPath = "../../../charts/pega"

var options = &helm.Options{
	SetValues: map[string]string{
		"global.actions.execute": "upgrade-deploy",
	},
}

func TestValidateUpgradeJobs(t *testing.T) {
	var installerJobObj k8sbatch.Job
	var installerSlice = returnJobSlices(t, pegaHelmChartPath, options)
	println("***************************")
	println(len(installerSlice))
	var expectedJob pegaJob
	for index, installerInfo := range installerSlice {
		if index >= 1 && index <= 3 {
			if index == 1 {
				expectedJob = pegaJob{"pega-pre-upgrade", []string{}, "pega-upgrade-environment-config"}
			} else if index == 2 {
				expectedJob = pegaJob{"pega-db-upgrade", []string{"wait-for-pre-dbupgrade"}, "pega-upgrade-environment-config"}
			} else if index == 3 {
				expectedJob = pegaJob{"pega-post-upgrade", []string{"wait-for-pegaupgrade", "wait-for-rolling-updates"}, "pega-upgrade-environment-config"}
			}

			helm.UnmarshalK8SYaml(t, installerInfo, &installerJobObj)
			VerifyJob(t, options, &installerJobObj, expectedJob)
		}

	}
}
