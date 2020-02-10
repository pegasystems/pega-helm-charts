# Addons Helm chart

The addons chart installs a collection of supporting services and tools required for a Pega deployment.  The services you will need to deploy will depend on your cloud environment - for example you may need a load balancer on Minikube, but not for EKS. These supporting services are deployed once per Kubernetes environment, regardless of how many Pega Infinity instances are deployed. This readme provides a detailed description of possible configurations and their default values as applicable.

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

Deploying Pega Platform with more than one Pod typically requires a load balancer to ensure that traffic is routed equally.  Some IaaS and PaaS providers supply a load balancer and some do not. If a native load balancer is not provided and configured, or the load balancer does not support cookie based session affinity, Traefik may be used instead.  If you do not wish to deploy Traefik, set `traefik.enabled` to `false` in the addons values.yaml configuration. For more configuration options available for Traefik, see the [Traefik Helm chart](https://github.com/helm/charts/blob/master/stable/traefik/values.yaml).

Example:

```yaml
traefik:
  enabled: true
  serviceType: NodePort
  ssl:
    enabled: false
  rbac:
    enabled: true
  service:
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

If deploying on EKS, you can use the native Amazon Load Balancer (ALB). In the Addons Helm chart, disable Traefik (set `traefik.enabled` to `false`) and enable ALB (set `aws-alb-ingress-controller.enabled` to `true`)

Configuration   | Usage
---             | ---
`clusterName`   | The name of your EKS cluster.  Resources created by the ALB Ingress controller will be prefixed with this string.
`autoDiscoverAwsRegion` | Auto discover awsRegion from ec2metadata, set this to true and omit awsRegion when ec2metadata is available.
`awsRegion`     | AWS region of the EKS cluster. Required if if ec2metadata is unavailable from the controller Pod or if `autoDiscoverAwsRegion` is not `true`.
`autoDiscoverAwsVpcID` | Auto discover awsVpcID from ec2metadata, set this to true and omit awsVpcID when ec2metadata is available.
`awsVpcID`      | VPC ID of EKS cluster, required if ec2metadata is unavailable from controller pod. Required if if ec2metadata is unavailable from the controller Pod or if `autoDiscoverAwsVpcID` is not `true`.
`extraEnv.AWS_ACCESS_KEY_ID` and `extraEnv.AWS_SECRET_ACCESS_KEY` | The access key and secret access key with access to configure AWS resources.

Example:

```yaml
aws-alb-ingress-controller:
  enabled: false
  clusterName: "YOUR_EKS_CLUSTER_NAME"
  autoDiscoverAwsRegion: true
  awsRegion: "YOUR_EKS_CLUSTER_REGION"
  autoDiscoverAwsVpcID: true
  awsVpcID: "YOUR_EKS_CLUSTER_VPC_ID"
  extraEnv:
    AWS_ACCESS_KEY_ID: "YOUR_AWS_ACCESS_KEY_ID"
    AWS_SECRET_ACCESS_KEY: "YOUR_AWS_SECRET_ACCESS_KEY"
```

### Azure AGIC

If deploying on AKS, you can use an Azure Application Gateway for the deployment load balancer. After you create the Application Gateway and an ingress controller, in the Addons Helm chart, disable Traefik (set `traefik.enabled` to `false`) and enable AGIC (set `ingress-azure.enabled` to `true`).

Example:

```yaml
ingress-azure:
  enabled: true
    appgw:
      subscriptionId: <YOUR.SUBSCRIPTION_ID>
      resourceGroup: <RESOURCE_GROUP_NAME>
      name: <APPLICATION_GATEWAY_NAME>
    armAuth:
    type: servicePrincipal
    secretJSON: <SECRET_JSON_CREATED_USING_ABOVE_COMMAND>
  rbac:
    enabled: true


## Aggregated logging

Environment                               | Suggested logging tools
---                                       | ---
Open-source Kubernetes                    | EFK
Red Hat Openshift                         | Built-in EFK
Amazon Elastic Kubernetes Service (EKS)   | Built-in EFK
Google Kubernetes Engine (GKE)            | Stackdriver
Pivotal Container Service (PKS)           | EFK
Microsoft Azure Kubernetes Service (AKS)  | Azure Monitor

## Logging with Elasticsearch-Fluentd-Kibana (EFK)

EFK is a standard logging stack that is provided as an example for ease of getting started in environments that do not have aggregated logging configured such as open-source Kubernetes. Other IaaS and PaaS providers typically include a logging system out of the box. You may enable the three components of EFK ([Elasticsearch](https://github.com/helm/charts/tree/master/stable/elasticsearch/values.yaml),[Fluentd](https://github.com/helm/charts/tree/master/stable/fluentd-elasticsearch/values.yaml), and [Kibana](https://github.com/helm/charts/tree/master/stable/kibana/values.yaml)) in the addons values.yaml file to deploy EFK automatically. For more configuration options available for each of the components, see their Helm Charts.

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

See the [metrics-server Helm chart](https://github.com/helm/charts/blob/master/stable/metrics-server/values.yaml) for additional parameters.

Example:

```yaml
metrics-server:
  enabled: true
  args:
    - --logtostderr
```
