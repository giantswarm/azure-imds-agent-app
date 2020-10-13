<!--

    TODO:

    - Add the project to the CircleCI:
      https://circleci.com/setup-project/gh/giantswarm/azure-imds-agent-app

    - Change the badge (with style=shield):
      https://circleci.com/gh/giantswarm/azure-imds-agent-app/edit#badges
      If this is a private repository token with scope `status` will be needed.

    - Update CODEOWNERS file according to the needs for this project

    - Run `devctl replace -i "REPOSITORY_NAME" "$(basename $(git rev-parse --show-toplevel))" *.md`
      and commit your changes.

    - If the repository is public consider adding godoc badge. This should be
      the first badge separated with a single space.
      [![GoDoc](https://godoc.org/github.com/giantswarm/azure-imds-agent-app?status.svg)](http://godoc.org/github.com/giantswarm/azure-imds-agent-app)

-->
[![CircleCI](https://circleci.com/gh/giantswarm/azure-imds-agent-app.svg?style=shield&circle-token=cbabd7d13186f190fca813db4f0c732b026f5f6c)](https://circleci.com/gh/giantswarm/azure-imds-agent-app)

# azure-imds-agent-app

App that talks to Azure Instance Metadata Service (IMDS) and exposes instance
metadata in Kubernetes cluster.
 