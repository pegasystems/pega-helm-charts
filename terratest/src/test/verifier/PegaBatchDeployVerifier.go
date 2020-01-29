package verifier

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	k8sapps "k8s.io/api/apps/v1"
	"testing"
)

type PegaBatchDeployVerifier struct {
	DeploymentVerifier
}

func NewPegaBatchDeployVerifier(t *testing.T, helmOptions *helm.Options, initContainers []string, deployment *k8sapps.Deployment) *PegaBatchDeployVerifier {
	v := &PegaBatchDeployVerifier{
		DeploymentVerifier: *NewDeploymentVerifier(t, helmOptions, initContainers, deployment),
	}
	v._nodeType = "BackgroundProcessing,Search,Batch,RealTime,Custom1,Custom2,Custom3,Custom4,Custom5,BIX"
	return v
}
