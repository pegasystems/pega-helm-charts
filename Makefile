dependencies:
	helm repo add traefik https://kubernetes-charts.storage.googleapis.com
	helm repo add cassandra https://kubernetes-charts-incubator.storage.googleapis.com/
	helm repo add elasticsearch https://kubernetes-charts.storage.googleapis.com/
	helm repo add fluentd-elasticsearch https://kubernetes-charts.storage.googleapis.com/
	helm repo add kibana https://kubernetes-charts.storage.googleapis.com/
	helm repo list
	helm dependency update ./charts/pega/

examples: dependencies
	mkdir -p ./build/kubernetes
	helm template ./charts/pega/ --output-dir ./build/kubernetes --values ./charts/pega/values.yaml --namespace example --set provider=k8s --set actions.execute=deploy
	tar -zcf ./pega-kubernetes-example.tar.gz ./build/kubernetes/pega/templates/*

	mkdir -p ./build/openshift
	helm template ./charts/pega/ --output-dir ./build/openshift --values ./charts/pega/values.yaml --namespace example --set provider=openshift --set actions.execute=deploy
	tar -zcf ./pega-openshift-example.tar.gz ./build/openshift/pega/templates/*

	mkdir -p ./build/aws-eks
	helm template ./charts/pega/ --output-dir ./build/aws-eks --values ./charts/pega/values.yaml --namespace example --set provider=eks --set actions.execute=deploy
	tar -zcf ./pega-eks-example.tar.gz ./build/aws-eks/pega/templates/*

clean:
	rm -rf build/
	rm -rf charts/pega/charts/*
