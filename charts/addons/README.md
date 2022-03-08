# Addons Helm chart

The addons chart installs a collection of supporting services and tools for a Pega deployment.  The services you need to deploy will depend on your cloud environment - for example you may need a load balancer on Minikube, but not for EKS. These supporting services are deployed once per Kubernetes environment, regardless of how many Pega Infinity instances are deployed. This readme provides a detailed description of possible configurations and their default values as applicable.

Pega does **not** actively update the dependencies in `requirements.yaml`. Whether a dependency is enabled or disabled will depend on the service you choose for your environment. For any enabled dependencies listed in the `requirements.yaml` file, you should update its corresponding `version` value. Disabled dependencies do not require version updates.

## Load balancer

Pega Platform deployments by default assume that clients will use the load balancing tools featured in the Kubernetes environment of the deployment. The table below lists the default load balancer for each environment. Pega supports specifying the use of Traefik as a load balancer for deployments in GKE and AKS environments if you would prefer it; in these cases, use the Addon Helm chart to override the defaults. 

Environment                               | Suggested load balancer
---                                       | ---
Open-source Kubernetes                    | Traefik
Red Hat Openshift                         | HAProxy (Using the `roundrobin` load balancer strategy)
Amazon Elastic Kubernetes Service (EKS)   | Amazon Load Balancer (ALB)
Google Kubernetes Engine (GKE)            | Google Cloud Load Balancer (GCLB)
Pivotal Container Service (PKS)           | Traefik
Microsoft Azure Kubernetes Service (AKS)  | Application Gateway Ingress Controller (AGIC)

### Traefik

Deploying Pega Platform with more than one Pod typically requires a load balancer to ensure that traffic is routed equally.  Some IaaS and PaaS providers supply a load balancer and some do not. If a native load balancer is not provided and configured, or the load balancer does not support cookie based session affinity, Traefik may be used instead.  If you do not wish to deploy Traefik, set `traefik.enabled` to `false` in the addons values.yaml configuration. For more configuration options available for Traefik, see the [Traefik Helm chart](https://github.com/traefik/traefik-helm-chart/blob/master/traefik/values.yaml).

Example:

```yaml
traefik:
  enabled: true
  ssl:
    enabled: false
  rbac:
    enabled: true
  service:
    type: NodePort
    nodePorts:
      http: 30080
      https: 30443
  resources:
    requests:
      cpu: 200m
      memory: 200Mi
    limits:
      cpu: 500m
      memory: 500Mi
```

### Amazon ALB

When deploying on AWS EKS, set the parameters, install `aws-load-balancer-controller.enabled` and `metrics-server` to true, `traefik` to false and then fill in the remaining parameters with your EKS environment details.

Configuration   | Usage
---             | ---
`clusterName`   | The name of your EKS cluster.  Resources created by the ALB Ingress controller will be prefixed with this string.
`region`     | AWS region of the EKS cluster. Required if if ec2metadata is unavailable from the controller Pod.
`vpcId`      | VPC ID of EKS cluster, required if ec2metadata is unavailable from controller pod.
`serviceAccount.annotations`  | Annotate the service account with `eks.amazonaws.com/role-arn` IAM Role that provides access to AWS resources.

Example:

```yaml
aws-load-balancer-controller:
  enabled: true
  clusterName: "YOUR_EKS_CLUSTER_NAME"
  region: "YOUR_EKS_CLUSTER_REGION"
  vpcId: "YOUR_EKS_CLUSTER_VPC_ID"
  serviceAccount:
    annotations:
      eks.amazonaws.com/role-arn: "YOUR_IAM_ROLE_ARN"
```

### Azure AGIC

When deploying on Azure AKS, you can use an Application Gateway Ingress Controller (AGIC) for the deployment load balancer. The AGIC is a pod within your AKS cluster that monitors the Kubernetes Ingress resources, which creates and applies the Application Gateway configuration based on the status of the Kubernetes cluster. For details, see [Azure Resource Manager Authentication](https://docs.microsoft.com/en-us/azure/application-gateway/ingress-controller-install-existing#azure-resource-manager-authentication).

After you create the deployment ingress controller, in the Addons Helm chart, disable Traefik (set `traefik.enabled` to `false`), enable AGIC (set `ingress-azure.enabled` to `true`) and add the AGIC gateway configuration details from your AKS deployment.

To authenticate with the AGIC in your AKS cluster, generate a kubernetes secret from an Active Directory Service Principal that is based on your AKS subscription ID. You must encode the Service Principal with base64 and add the result to the `armAuth.secretJSON` field. For details, see the comments in the addons [values.yaml](/values.yaml) or the [AKS runbook](../../docs/Deploying-Pega-on-AKS.md).

As an authentication alternative, you can configure an AAD Pod Identity to manage authentication access with the AGIC in your cluster via the Azure Resource Manager. For details, see [Set up AAD Pod Identity](https://docs.microsoft.com/en-us/azure/application-gateway/ingress-controller-install-existing#set-up-aad-pod-identity).

It is a recommended best practice to enable RBAC on your AKS cluster and match the setting in the Addons Helm chart.

Example:

```yaml
ingress-azure:
  enabled: true
    appgw:
      subscriptionId: <YOUR.SUBSCRIPTION_ID>
      resourceGroup: <RESOURCE_GROUP_NAME>
      name: <APPLICATION_GATEWAY_NAME>
      usePrivateIP: true
    armAuth:
    type: servicePrincipal
    secretJSON: <SECRET_JSON_CREATED_USING_ABOVE_COMMAND>
  rbac:
    enabled: true
```

### Google (GCLB)

If deploying on GKE, you can use Google Cloud Load Balancer to route your traffic.  In the Addons Helm chart, disable Traefik (set `traefik.enabled` to `false`).  All other GCLB configurations are automatic.

## Aggregated logging

Environment                               | Suggested logging tools
---                                       | ---
Open-source Kubernetes                    | EFK
Red Hat Openshift                         | Built-in EFK
Amazon Elastic Kubernetes Service (EKS)   | Built-in EFK
Google Kubernetes Engine (GKE)            | Google Cloud Operations
Pivotal Container Service (PKS)           | EFK
Microsoft Azure Kubernetes Service (AKS)  | Azure Monitor

## Logging with Elasticsearch-Fluentd-Kibana (EFK)

EFK is a standard logging stack that is provided as an example for ease of getting started in environments that do not have aggregated logging configured such as open-source Kubernetes. Other IaaS and PaaS providers typically include a logging system out of the box. You may enable the three components of EFK ([Elasticsearch](https://github.com/elastic/helm-charts/blob/v7.16.3/elasticsearch/values.yaml),[Fluentd](https://github.com/helm/charts/tree/master/stable/fluentd-elasticsearch/values.yaml), and [Kibana](https://github.com/elastic/helm-charts/blob/v7.16.3/kibana/values.yaml)) in the addons values.yaml file to deploy EFK automatically. For more configuration options available for each of the components, see their Helm Charts.

Example:

```yaml

deploy_efk: &deploy_efk true

elasticsearch:
  enabled: *deploy_efk
  fullnameOverride: "elastic-search"

kibana:
  enabled: *deploy_efk
  files:
    kibana.yml:
      elasticsearch.url: http://elastic-search-client:9200
  service:
    externalPort: 80
  ingress:
 
    enabled: true
    # Enter the domain name to access kibana via a load balancer.
    hosts:
      - "YOUR_WEB.KIBANA.EXAMPLE.COM"

fluentd-elasticsearch:
  enabled: *deploy_efk
  elasticsearch:
    host: elastic-search-client
    buffer_chunk_limit: 250M
    buffer_queue_limit: 30

```

## Metrics

Environment             | Suggested metrics server
---                     | ---
Open-source Kubernetes  | Metrics server
All others              | Built-in metrics server

Autoscaling in Kubernetes requires the use of a metrics server, a cluster-wide aggregator of resource usage data.  Most PaaS and IaaS providers supply a metrics server, but if you wish to deploy into open source kubernetes, you will need to supply your own.

See the [metrics-server Helm chart](https://github.com/kubernetes-sigs/metrics-server/tree/master/charts/metrics-server) for additional parameters.

Example:

```yaml
metrics-server:
  enabled: true
  args:
    - --logtostderr
```
