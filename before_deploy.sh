export CHART_VERSION=$(expr ${TRAVIS_TAG:1})
export PEGA_FILE_NAME=pega-${CHART_VERSION}.tgz
export ADDONS_FILE_NAME=addons-${CHART_VERSION}.tgz
cat descriptor-template.json | jq '.files[0].includePattern=env.PEGA_FILE_NAME' | jq '.files[0].uploadPattern=env.PEGA_FILE_NAME' | jq '.files[1].includePattern=env.ADDONS_FILE_NAME' | jq '.files[1].uploadPattern=env.ADDONS_FILE_NAME' > descriptor.json
curl -o index.yaml https://kishor.bintray.com/pega-helm-charts/index.yaml
helm package --version ${CHART_VERSION} ./charts/pega/
helm package --version ${CHART_VERSION} ./charts/addons/
helm repo index --merge index.yaml --url https://kishor.bintray.com/pega-helm-charts/ .