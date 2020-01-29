package verifier

import (
	"k8s.io/api/apps/v1beta2"
)

type StatefulSetVerifier struct {
	VerifierImpl
	k8sDeployment *v1beta2.StatefulSet
}
