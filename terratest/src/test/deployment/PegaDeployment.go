package deployment

type PegaDeploymentTestOptions struct {
	deploymentName     string
	initContainers     []string
	nodeType           string
	passivationTimeout string
}
