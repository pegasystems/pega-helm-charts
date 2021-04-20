#!/bin/bash

# CHART_VERSION is computed from the TAG details of the commit. Every Github release creates tag with the release name.
# Release name (or) Tag name should be in vX.X.X format. Helm CHART_VERSION would be X.X.X
export TRAVIS_TAG=v12.5.0
export CHART_VERSION=$(expr ${TRAVIS_TAG:1})
export PEGA_FILE_NAME=pega-${CHART_VERSION}.tgz
export ADDONS_FILE_NAME=addons-${CHART_VERSION}.tgz
export BACKINGSERVICES_FILE_NAME=backingservices-${CHART_VERSION}.tgz
cat descriptor-template.json | jq '.files[0].includePattern=env.PEGA_FILE_NAME' | jq '.files[0].uploadPattern=env.PEGA_FILE_NAME' | jq '.files[1].includePattern=env.ADDONS_FILE_NAME' | jq '.files[1].uploadPattern=env.ADDONS_FILE_NAME' | jq '.files[2].includePattern=env.BACKINGSERVICES_FILE_NAME' | jq '.files[2].uploadPattern=env.BACKINGSERVICES_FILE_NAME' > descriptor.json
# Get the latest index.yaml from github.io
curl -o index.yaml https://pegasystems.github.io/pega-helm-charts/index.yaml
# Clone the versions from gh-pages to a temp directory - xyz
# The versions will be re-installed in gh-pages
git clone -b gh-pages https://github.com/${TRAVIS_REPO_SLUG} xyz
cp xyz/* .
rm -rf xyz
# Package up the changed version
helm package --version ${CHART_VERSION} ./charts/pega/
helm package --version ${CHART_VERSION} ./charts/addons/
helm package --version ${CHART_VERSION} ./charts/backingservices/
# and merge it
helm repo index --merge index.yaml --url https://pegasystems.github.io/pega-helm-charts/ .
