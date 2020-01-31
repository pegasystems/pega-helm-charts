# Minikube

Minikube runs a single-node Kubernetes cluster inside a Virtual Machine (VM) on your laptop for users looking to try out Kubernetes or develop with it day-to-day. For more information on minikube, see the [Minikube documentation](https://kubernetes.io/docs/setup/learning-environment/minikube/).


This document explains on how to deploy pega using minikube as a provider.

# Quick Start

1. For installing minikube - https://kubernetes.io/docs/tasks/tools/install-minikube/
2. Minikube Documentation - https://minikube.sigs.k8s.io/docs/overview/

# Basic Commands for Minikube

- Start a cluster by running:
 ```minikube start```

- Access the Kubernetes Dashboard running within the minikube cluster:
 ```minikube dashboard```

- Stop your local minikube cluster:
 ```minikube stop```

- Delete your local cluster:
```minikube delete```

- To start minikube with different version of kubernetes
```minikube start --kubernetes-version v1.15.0```


# FAQ's

1. How to increase the memory limit of a running minikube

	There is no direct way to increase the memory limit of a running minikube.

	``` minikube stop```
	
	```minikube delete ```
	
	```minikube start --cpus 4 --memory 12288 ```

2. How to start minikube with custom CPU/memory limits

	```minikube start --cpus 4 --memory 10240```

3. How to set default memory which is considered on each minikube start

	```minikube config set memory 5000``` followed by ```minikube start```

4. How to access Pega Designer Studio after deployment

	``` <minikube ip>:<Pega service nodePort>/prweb```

minikube ip can be fetched using command - ``` minikube ip``` and Pega service Nodeport can be fetched using below command
```kubectl get service -o go-template='{{range.spec.ports}}{{"Port to access: "}}{{.nodePort}}{{end}}' <service-name> --namespace <namespace name> ```

***Recommended Memory Limits***

Start minikube with at least 4 CPU’s and 10GB memory for complete pega deployment. As per the need increase the limits of minikube.

***Note***
1. Use “values-minimal.yaml” to deploy pega which is available in the [pega chart](../charts/pega) directory. 

	Example helm command to deploy
	
	```helm install . -n mypega --namespace myproject --values ./values-minimal.yaml```

2. As this runs on the personal laptop for a day-to-day project with minimal memory and CPU limits, minikube supports only "install", "deploy" and "install-deploy" actions. It is advisable to use this kind of cluster configuration for simple activities on Pega as it might spike with CPU and memory.
