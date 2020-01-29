package deployment

import (
	k8sapps "k8s.io/api/apps/v1"
	k8score "k8s.io/api/core/v1"
)

type K8sInformationExtractor interface {
	GetPod() *k8score.PodSpec
	GetDeploymentMetadata() DeploymentMetadata
}
type K8sDeploymentExtractor struct {
	K8sDeployment *k8sapps.Deployment
}

func (e K8sDeploymentExtractor) GetPod() *k8score.PodSpec {
	return &e.K8sDeployment.Spec.Template.Spec
}

func (e K8sDeploymentExtractor) GetDeploymentMetadata() DeploymentMetadata {
	return DeploymentMetadata{
		ObjectMeta: e.K8sDeployment.ObjectMeta,
		TypeMeta:   e.K8sDeployment.TypeMeta,
	}
}
