package test

import (
	"testing"
)

// TestPegaStandardTierDeployment - Test case to verify the standard pega tier deployment.
// Standard tier deployment includes web deployment, batch deployment, stream statefulset, search service, hpa, rolling update, web services, ingresses and config maps
func TestPegaStandardTierDeployment2(t *testing.T) {
	NewPegaStandardDeploymentTest(
		"K8s",
		"deploy",
		[]string{"wait-for-pegasearch", "wait-for-cassandra"},
		t,
	).Run()
}
