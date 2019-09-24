dependencies:
	helm repo add incubator https://kubernetes-charts-incubator.storage.googleapis.com/
	helm repo add stable https://kubernetes-charts.storage.googleapis.com
	helm repo list	
	helm dependency update ./charts/pega/
	helm dependency update ./charts/addons/ 

examples: 
	mkdir -p ./build/kubernetes
	helm template ./charts/pega/ \
		--output-dir ./build/kubernetes \
		--values ./charts/pega/values.yaml \
		--namespace example \
		--set global.provider=k8s \
		--set global.actions.execute=deploy
	tar -C ./build/kubernetes/pega/templates -cvzf ./pega-kubernetes-example.tar.gz .

	mkdir -p ./build/openshift
	helm template ./charts/pega/ \
		--output-dir ./build/openshift \
		--values ./charts/pega/values.yaml \
		--namespace example \
		--set global.provider=openshift \
		--set global.actions.execute=deploy
	tar -C ./build/openshift/pega/templates -cvzf ./pega-openshift-example.tar.gz .

	mkdir -p ./build/aws-eks
	helm template ./charts/pega/ \
		--output-dir ./build/aws-eks \
		--values ./charts/pega/values.yaml \
		--namespace example \
		--set global.provider=eks \
		--set global.actions.execute=deploy
	tar -C ./build/aws-eks/pega/templates -cvzf ./pega-aws-eks-example.tar.gz .

	mkdir -p ./build/azure-aks
	helm template ./charts/pega/ \
		--output-dir ./build/azure-aks \
		--values ./charts/pega/values.yaml \
		--namespace example \
		--set global.provider=aks \
		--set global.actions.execute=deploy
	tar -C ./build/azure-aks/pega/templates -cvzf ./pega-azure-aks-example.tar.gz .

	mkdir -p ./build/google-gke
	helm template ./charts/pega/ \
		--output-dir ./build/google-gke \
		--values ./charts/pega/values.yaml \
		--namespace example \
		--set global.provider=gke \
		--set global.actions.execute=deploy
	tar -C ./build/google-gke/pega/templates -cvzf ./pega-google-gke-example.tar.gz .

	mkdir -p ./build/pivotal-pks
	helm template ./charts/pega/ \
		--output-dir ./build/pivotal-pks \
		--values ./charts/pega/values.yaml \
		--namespace example \
		--set global.provider=pks \
		--set global.actions.execute=deploy
	tar -C ./build/pivotal-pks/pega/templates -cvzf ./pega-pivotal-pks-example.tar.gz .

clean:
	rm -rf ./build
	rm -rf ./charts/pega/charts/*
	rm -rf ./*.tar.gz

packages:
	export CHART_VERSION=$(expr ${TRAVIS_TAG:1})
	export PEGA_FILE_NAME=pega-${CHART_VERSION}.tgz
	export ADDONS_FILE_NAME=addons-${CHART_VERSION}.tgz
	cat descriptor-template.json | jq '.files[0].includePattern=env.PEGA_FILE_NAME' | jq '.files[0].uploadPattern=env.PEGA_FILE_NAME' | jq '.files[1].includePattern=env.ADDONS_FILE_NAME' | jq '.files[1].uploadPattern=env.ADDONS_FILE_NAME' > descriptor.json
	curl -o index.yaml https://kishor.bintray.com/pega-helm-charts/index.yaml
	helm package --version ${CHART_VERSION} ./charts/pega/
	helm package --version ${CHART_VERSION} ./charts/addons/
	helm repo index --merge index.yaml --url https://kishor.bintray.com/pega-helm-charts/ .