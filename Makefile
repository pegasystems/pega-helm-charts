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