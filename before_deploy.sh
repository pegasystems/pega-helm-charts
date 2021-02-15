#!/bin/bash

# CHART_VERSION is computed from the TAG details of the commit. Every Github release creates tag with the release name.
# Release name (or) Tag name should be in vX.X.X format. Helm CHART_VERSION would be X.X.X
export CHART_VERSION=$(expr ${TRAVIS_TAG:1})
export PEGA_FILE_NAME=pega-${CHART_VERSION}.tgz
export ADDONS_FILE_NAME=addons-${CHART_VERSION}.tgz
export BACKINGSERVICES_FILE_NAME=backingservices-${CHART_VERSION}.tgz
cat descriptor-template.json | jq '.files[0].includePattern=env.PEGA_FILE_NAME' | jq '.files[0].uploadPattern=env.PEGA_FILE_NAME' | jq '.files[1].includePattern=env.ADDONS_FILE_NAME' | jq '.files[1].uploadPattern=env.ADDONS_FILE_NAME' | jq '.files[2].includePattern=env.BACKINGSERVICES_FILE_NAME' | jq '.files[2].uploadPattern=env.BACKINGSERVICES_FILE_NAME' > descriptor.json
curl -o index.yaml https://dl.bintray.com/pegasystems/pega-helm-charts/index.yaml
helm package --version ${CHART_VERSION} ./charts/pega/
helm package --version ${CHART_VERSION} ./charts/addons/
helm package --version ${CHART_VERSION} ./charts/backingservices/
helm repo index --merge index.yaml --url https://dl.bintray.com/pegasystems/pega-helm-charts/ .