package addons

type SearchResourceOption struct {
	name string
	kind string
}

var traefikResources = []SearchResourceOption{
	{
		name: "release-name-traefik",
		kind: "ConfigMap",
	},
	{
		name: "release-name-traefik",
		kind: "ServiceAccount",
	},
	{
		name: "release-name-traefik",
		kind: "ClusterRole",
	},
	{
		name: "release-name-traefik",
		kind: "Deployment",
	},
	{
		name: "release-name-traefik",
		kind: "ClusterRoleBinding",
	},
	{
		name: "release-name-traefik",
		kind: "Service",
	},
	{
		name: "release-name-traefik-test",
		kind: "Pod",
	},
	{
		name: "release-name-traefik-test",
		kind: "ConfigMap",
	},
}
