package verifier

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	"k8s.io/api/apps/v1beta2"
	k8score "k8s.io/api/core/v1"
	"testing"
)

type PegaStreamDeployVerifier struct {
	StatefulSetVerifier
}

func NewPegaStreamDeployVerifier(t *testing.T, helmOptions *helm.Options, initContainers []string, deployment *v1beta2.StatefulSet) *PegaStreamDeployVerifier {
	v := &PegaStreamDeployVerifier{
		StatefulSetVerifier: StatefulSetVerifier{
			VerifierImpl:  *NewDeployVerifier(t, helmOptions, initContainers),
			k8sDeployment: deployment,
		},
	}
	v._nodeType = "Stream"
	v._passivationTimeout = "900"
	return v
}
func (p PegaStreamDeployVerifier) Verify() {
	require.Equal(p.t, p.k8sDeployment.Spec.VolumeClaimTemplates[0].Name, "pega-stream")
	require.Equal(p.t, p.k8sDeployment.Spec.VolumeClaimTemplates[0].Spec.AccessModes[0], k8score.PersistentVolumeAccessMode("ReadWriteOnce"))
	require.Equal(p.t, p.k8sDeployment.Spec.ServiceName, "pega-stream")
	statefulsetSpec := p.k8sDeployment.Spec.Template.Spec
	require.Equal(p.t, statefulsetSpec.Containers[0].VolumeMounts[1].Name, "pega-stream")
	require.Equal(p.t, statefulsetSpec.Containers[0].VolumeMounts[1].MountPath, "/opt/pega/streamvol")
	require.Equal(p.t, statefulsetSpec.Containers[0].VolumeMounts[2].Name, "pega-volume-credentials")
	require.Equal(p.t, statefulsetSpec.Containers[0].VolumeMounts[2].MountPath, "/opt/pega/secrets")
}
