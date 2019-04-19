dependencies:
	helm repo add traefik https://kubernetes-charts.storage.googleapis.com
	helm repo add cassandra https://kubernetes-charts-incubator.storage.googleapis.com/
	helm repo add elasticsearch https://kubernetes-charts.storage.googleapis.com/
	helm repo add fluentd-elasticsearch https://kubernetes-charts.storage.googleapis.com/
	helm repo add kibana https://kubernetes-charts.storage.googleapis.com/
	helm repo list
	helm dependency update ./charts/pega/

examples: dependencies
	mkdir -p ./examples/kubernetes
	helm template ./charts/pega/ --output-dir ./examples/kubernetes --values ./charts/pega/values.yaml --namespace example --set provider=k8s --set actions.execute=deploy
	mv ./examples/kubernetes/pega/templates/* ./examples/kubernetes
	tar -zcf ./examples-kubernetes.tar.gz ./examples/kubernetes
	rm -rf ./examples/kubernetes/pega

	mkdir -p ./examples/openshift
	helm template ./charts/pega/ --output-dir ./examples/openshift --values ./charts/pega/values.yaml --namespace example --set provider=openshift --set actions.execute=deploy
	mv ./examples/openshift/pega/templates/* ./examples/openshift
	rm -rf ./examples/openshift/pega

	mkdir -p ./examples/aws-eks
	helm template ./charts/pega/ --output-dir ./examples/aws-eks --values ./charts/pega/values.yaml --namespace example --set provider=eks --set actions.execute=deploy
	mv ./examples/aws-eks/pega/templates/* ./examples/aws-eks
	rm -rf ./examples/aws-eks/pega

clean:
	rm -rf examples/
	rm -rf charts/pega/charts/*
