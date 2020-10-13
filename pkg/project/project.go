package project

var (
	description = "App that talks to Azure Instance Metadata Service (IMDS) and exposes instance metadata in Kubernetes cluster "
	gitSHA      = "n/a"
	name        = "azure-imds-agent-app"
	source      = "https://github.com/giantswarm/azure-imds-agent-app"
	version     = "0.1.0-dev"
)

func Description() string {
	return description
}

func GitSHA() string {
	return gitSHA
}

func Name() string {
	return name
}

func Source() string {
	return source
}

func Version() string {
	return version
}
