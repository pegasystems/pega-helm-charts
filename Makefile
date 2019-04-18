
examples:
	helm dependency update ./charts/pega/

	mkdir -p ./examples/kubernetes
	helm template ./charts/pega/ --output-dir ./examples/kubernetes --values ./charts/pega/values.yaml --namespace example --set provider=k8s --set actions.execute=deploy
	mv ./examples/kubernetes/pega/templates/* ./examples/kubernetes
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
	rm -rfv examples/
