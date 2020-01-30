package deployment

import (
	"k8s.io/api/apps/v1beta2"
)

type StatefulSetVerifier struct {
	Verifier
	k8sDeployment *v1beta2.StatefulSet
}
