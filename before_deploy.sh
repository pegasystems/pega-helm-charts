#!/bin/bash

# CHART_VERSION is computed from the TAG details of the commit. Every Github release creates tag with the release name.
# Release name (or) Tag name should be in vX.X.X format. Helm CHART_VERSION would be X.X.X
tagVersion=""
if [ ${GITHUB_REF_TYPE} == "tag" ]
then
    tagVersion=${GITHUB_REF_NAME}
fi
export CHART_VERSION=$(expr ${tagVersion:1})
export PEGA_FILE_NAME=pega-${CHART_VERSION}.tgz
export ADDONS_FILE_NAME=addons-${CHART_VERSION}.tgz
export BACKINGSERVICES_FILE_NAME=backingservices-${CHART_VERSION}.tgz
export DEPLOY_CONFIGURATIONS_FILE_NAME=deploy-config-${CHART_VERSION}.tgz
export INSTALLER_CONFIGURATIONS_FILE_NAME=installer-config-${CHART_VERSION}.tgz
# Get the latest index.yaml from github.io
curl -o index.yaml https://pegasystems.github.io/pega-helm-charts/index.yaml
# Clone the versions from gh-pages to a temp directory - xyz
# The versions will be re-installed in temporary directory - temp_gh_pages
git clone --single-branch --branch gh-pages https://github.com/${GITHUB_REPOSITORY} temp_gh_pages
cp temp_gh_pages/*.tgz .
rm -rf temp_gh_pages
ls
# Package up the changed version
helm package --version ${CHART_VERSION} ./charts/pega/
helm package --version ${CHART_VERSION} ./charts/addons/
helm package --version ${CHART_VERSION} ./charts/backingservices/
tar -czvf ${DEPLOY_CONFIGURATIONS_FILE_NAME} --directory=./charts/pega/config deploy/context.xml.tmpl deploy/server.xml deploy/prconfig.xml deploy/prlog4j2.xml
mkdir -p ./charts/pega/charts/installer/config/installer && cp ./charts/pega/charts/installer/config/migrateSystem.properties.tmpl ./charts/pega/charts/installer/config/installer && cp ./charts/pega/charts/installer/config/prbootstrap.properties.tmpl ./charts/pega/charts/installer/config/installer && cp ./charts/pega/charts/installer/config/prconfig.xml.tmpl ./charts/pega/charts/installer/config/installer && cp ./charts/pega/charts/installer/config/prlog4j2.xml ./charts/pega/charts/installer/config/installer && cp ./charts/pega/charts/installer/config/prpcUtils.properties.tmpl ./charts/pega/charts/installer/config/installer && cp ./charts/pega/charts/installer/config/setupDatabase.properties.tmpl ./charts/pega/charts/installer/config/installer && tar -czvf ${INSTALLER_CONFIGURATIONS_FILE_NAME} --directory=./charts/pega/charts/installer/config installer/migrateSystem.properties.tmpl installer/prbootstrap.properties.tmpl installer/prconfig.xml.tmpl installer/prlog4j2.xml installer/prpcUtils.properties.tmpl installer/setupDatabase.properties.tmpl
# and merge it
helm repo index --merge index.yaml --url https://pegasystems.github.io/pega-helm-charts/ .