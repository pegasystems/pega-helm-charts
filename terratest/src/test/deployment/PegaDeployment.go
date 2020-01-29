package deployment

type PegaDeployment struct {
	deploymentName     string
	initContainers     []string
	nodeType           string
	passivationTimeout string
}
