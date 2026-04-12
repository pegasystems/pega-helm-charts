# Autopilot Service Helm chart

The Pega `Autopilot Service` backing service provides AI-powered capabilities for Pega Infinity Platform by connecting directly to LLM providers (Azure OpenAI, AWS Bedrock, Google Vertex AI). This chart deploys the Autopilot Service for on-premise environments.

## Configuring a backing service with your pega environment

You can provision the Autopilot Service into your `pega` environment namespace, with the service endpoint configured for your Pega Infinity environment. When you include the Autopilot Service in your namespace, the service endpoint is included within your Pega Infinity environment network to ensure isolation of your application data.

## Supported LLM Providers

| Provider | Authentication Methods |
|---|---|
| Azure OpenAI | API Key, Pre-existing Secret |
| AWS Bedrock | Access Key/Secret, IAM Roles (IRSA) |
| Google Vertex AI | Service Account JSON Key, Workload Identity |

## Configuration settings

| Configuration | Usage |
|---|---|
| `enabled` | Enable the Autopilot Service deployment as a backing service. Set this parameter to `true` to deploy the service. |
| `deployment.name` | Specify the name of your Autopilot Service deployment. Your deployment creates resources prefixed with this string. |
| `docker.registry.url` | Specify the image registry URL. |
| `docker.registry.username` | Specify the username for the Docker registry. |
| `docker.registry.password` | Specify the password for the Docker registry. |
| `docker.imagePullSecretNames` | List pre-existing secrets to be used for pulling Docker images. |
| `docker.autopilot.image` | Specify the Autopilot Service Docker image and tag. |
| `docker.autopilot.imagePullPolicy` | Specify the image pull policy. Default is `Always`. |
| `replicas` | Number of pod replicas to provision. Default is `2`. |
| `service.port` | Defines the port used by the Service. Default is `8080`. |
| `service.targetPort` | Defines the port used by the Pod and Container. Default is `8080`. |
| `service.serviceType` | The [type of service](https://kubernetes.io/docs/concepts/services-networking/service/#publishing-services-service-types) you wish to expose. Default is `ClusterIP`. |
| `enableGenaiHub` | Set to `false` for on-prem deployments with direct provider connectivity. Default is `false`. |
| `authEnabled` | Enable or disable authentication for the service. Default is `false`. |
| `isInternalDeployment` | Set to `false` for on-prem deployments. Default is `false`. |
| `modelProviders` | Comma-separated list of providers to enable (e.g., `"Azure,Vertex,Bedrock"`). If empty, auto-detection based on available credentials is used. |
| `awsRegion` | AWS region for Bedrock. Default is `us-east-1`. |
| `serviceAccountName` | Kubernetes service account name for IAM role binding (IRSA or Workload Identity). |
| `affinity` | Define pod affinity so that it is restricted to run on particular node(s), or to prefer to run on particular nodes. |
| `tolerations` | Define pod tolerations so that it is allowed to run on node(s) with particular taints. |

## Provider credentials

The Autopilot Service supports three methods for providing LLM provider credentials.

### Option 1: Inline credentials (auto-creates Kubernetes Secret)

Specify provider credentials directly in your values file. The chart automatically creates a Kubernetes Secret containing these values.

| Configuration | Usage |
|---|---|
| `azure.endpoint` | Azure OpenAI endpoint URL (e.g., `https://my-openai.openai.azure.com/`). |
| `azure.apiKey` | Azure OpenAI API key. |
| `azure.apiVersion` | Azure OpenAI API version. Default is `2024-10-21`. |
| `aws.accessKeyId` | AWS access key ID for Bedrock. |
| `aws.secretAccessKey` | AWS secret access key for Bedrock. |
| `aws.sessionToken` | Optional AWS session token for temporary credentials. |
| `vertex.credentials` | Base64-encoded Google service account JSON key. |
| `vertex.applicationCredentialsFile` | Path to credentials file mounted in container. |
| `vertex.project` | Google Cloud project ID. |
| `vertex.location` | Vertex AI location. Default is `us-central1`. |

```yaml
autopilot:
  enabled: true
  azure:
    endpoint: "https://my-openai.openai.azure.com/"
    apiKey: "your-azure-api-key"
  aws:
    accessKeyId: "AKIAIOSFODNN7EXAMPLE"
    secretAccessKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
  vertex:
    credentials: "base64-encoded-service-account-json"
    project: "my-gcp-project"
    location: "us-central1"
```

### Option 2: Pre-existing Kubernetes Secret

Use a secret that you create and manage outside of this chart. Set `providerCredentialsSecret` to the name of your secret. The secret should contain keys matching the environment variable names used by the service.

| Configuration | Usage |
|---|---|
| `providerCredentialsSecret` | Name of an existing Kubernetes Secret containing provider credentials. When set, inline credentials are ignored. |

The secret should contain the applicable keys:

| Key | Description |
|---|---|
| `AZURE_ENDPOINT` | Azure OpenAI endpoint URL |
| `AZURE_OPENAI_KEY` | Azure OpenAI API key |
| `AWS_ACCESS_KEY_ID` | AWS access key ID |
| `AWS_SECRET_ACCESS_KEY` | AWS secret access key |
| `AWS_SESSION_TOKEN` | AWS session token (optional) |

```yaml
autopilot:
  enabled: true
  providerCredentialsSecret: "my-provider-credentials"
```

Create the secret manually:

```bash
kubectl create secret generic my-provider-credentials \
  --namespace <NAMESPACE> \
  --from-literal=AZURE_ENDPOINT="https://my-openai.openai.azure.com/" \
  --from-literal=AZURE_OPENAI_KEY="your-api-key" \
  --from-literal=AWS_ACCESS_KEY_ID="AKIAIOSFODNN7EXAMPLE" \
  --from-literal=AWS_SECRET_ACCESS_KEY="wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
```

### Option 3: IAM role-based authentication (recommended for production)

Use cloud-native IAM roles with no static credentials. For AWS, use IRSA (IAM Roles for Service Accounts). For GCP, use Workload Identity.

| Configuration | Usage |
|---|---|
| `serviceAccountName` | Name of the Kubernetes service account annotated with the IAM role. |

**AWS IRSA example:**

```yaml
autopilot:
  enabled: true
  serviceAccountName: "autopilot-sa"
```

Annotate the service account with the IAM role:

```bash
kubectl annotate serviceaccount autopilot-sa \
  --namespace <NAMESPACE> \
  eks.amazonaws.com/role-arn=arn:aws:iam::123456789012:role/autopilot-bedrock-role
```

**GCP Workload Identity example:**

```yaml
autopilot:
  enabled: true
  serviceAccountName: "autopilot-sa"
  vertex:
    project: "my-gcp-project"
    location: "us-central1"
```

Annotate the service account:

```bash
kubectl annotate serviceaccount autopilot-sa \
  --namespace <NAMESPACE> \
  iam.gke.io/gcp-service-account=autopilot@my-gcp-project.iam.gserviceaccount.com
```

## Custom models configuration

The Autopilot Service uses a model list to determine which LLM models are available. Each model entry follows the format used by the `model_id` field: `provider/creator/model_name/version` (e.g., `azure/openai/GPT-5/2025-08-07`). The service routes requests to the appropriate provider endpoint based on this model metadata.

### Option 1: Use the default models file bundled with the chart (recommended)

The chart includes a `files/default-models.json` file containing a curated list of models for all supported providers. Setting `deployModelsConfigMap: true` (the default) automatically creates a ConfigMap from this file during `helm install` or `helm upgrade`.

| Configuration | Usage |
|---|---|
| `deployModelsConfigMap` | Set to `true` to create a ConfigMap from the bundled `files/default-models.json`. Default is `true`. |

```yaml
autopilot:
  enabled: true
  deployModelsConfigMap: true
```

To customize the default model list before deployment, edit `files/default-models.json` in the chart directory.

### Option 2: Provide a pre-existing ConfigMap

Use a ConfigMap that you create and manage outside of this chart. The ConfigMap must contain a key named `models.json` with the model list as JSON content.

| Configuration | Usage |
|---|---|
| `customModels.existingConfigMap` | Name of an existing ConfigMap containing `models.json`. |

```yaml
autopilot:
  enabled: true
  deployModelsConfigMap: false
  customModels:
    existingConfigMap: "my-models-configmap"
```

Create the ConfigMap:

```bash
kubectl create configmap my-models-configmap \
  --namespace <NAMESPACE> \
  --from-file=models.json=./my-models.json
```

### Option 3: Inline model list in values

Provide the model JSON content directly in your values file. The chart creates a ConfigMap from this inline content.

| Configuration | Usage |
|---|---|
| `customModels.inline` | JSON string containing the model list. |

```yaml
autopilot:
  enabled: true
  deployModelsConfigMap: false
  customModels:
    inline: |
      [
        {
          "provider": "azure",
          "creator": "openai",
          "model_name": "GPT-5",
          "model_mapping_id": "gpt-5-2025-08-07",
          "name": "gpt-5-2025-08-07",
          "model_id": "azure/openai/GPT-5/2025-08-07",
          "input_tokens": 400000,
          "output_tokens": 128000,
          "type": "chat_completion",
          "version": "2025-08-07",
          "model_path": ["/openai/deployments/gpt-5/chat/completions"],
          "supported_capabilities": {
            "streaming": true,
            "functions": true,
            "json_mode": true
          }
        }
      ]
```

### Model JSON format

Each model entry in the models file requires the following fields:

| Field | Required | Description |
|---|---|---|
| `provider` | Yes | Cloud provider (`azure`, `bedrock`, `vertex`). |
| `creator` | Yes | Model creator (e.g., `openai`, `anthropic`, `google`, `amazon`). |
| `model_name` | Yes | Display name of the model (e.g., `GPT-5`, `Claude-37-Sonnet`). |
| `name` | Yes | Unique model identifier used internally by the service. |
| `model_id` | Yes | Model identifier in `provider/creator/model_name/version` format. |
| `model_mapping_id` | Yes | Provider-specific deployment name used for API routing. |
| `model_path` | Yes | Array of API endpoint paths for the model. |
| `type` | Yes | Model type: `chat_completion`, `embedding`, or `image`. |
| `version` | Yes | Model version string. |
| `input_tokens` | No | Maximum input token count. |
| `output_tokens` | No | Maximum output token count. |
| `supported_capabilities` | No | Object describing model capabilities (streaming, functions, multimodal, etc.). |
| `parameters` | No | Object describing tunable parameters (temperature, top_p, max_tokens, etc.). |

### Liveness and readiness probes

The Autopilot Service uses liveness and readiness probes on the `/v1/health` endpoint. Configure probes as part of a `livenessProbe` or `readinessProbe` configuration.

Parameter             | Description    | Default `livenessProbe` | Default `readinessProbe`
---                   | ---            | ---                     | ---
`initialDelaySeconds` | Number of seconds after the container has started before probes are initiated. | `10` | `5`
`timeoutSeconds`      | Number of seconds after which the probe times out. | `5` | `5`
`periodSeconds`       | How often (in seconds) to perform the probe. | `30` | `10`
`successThreshold`    | Minimum consecutive successes for the probe to be considered successful after it determines a failure. | `1` | `1`
`failureThreshold`    | The number of consecutive failures for the pod to be terminated by Kubernetes. | `3` | `3`

Example:

```yaml
livenessProbe:
  initialDelaySeconds: 10
  timeoutSeconds: 5
  periodSeconds: 30
  successThreshold: 1
  failureThreshold: 3
readinessProbe:
  initialDelaySeconds: 5
  timeoutSeconds: 5
  periodSeconds: 10
  successThreshold: 1
  failureThreshold: 3
```

## Ingress

To expose the Autopilot Service externally, enable the ingress configuration.

| Configuration | Usage |
|---|---|
| `ingress.enabled` | Set to `true` to deploy an ingress. Default is `false`. |
| `ingress.domain` | Specify your custom domain. |
| `ingress.ingressClassName` | Ingress class to be used. |
| `ingress.tls.enabled` | Specify the use of HTTPS for ingress connectivity. |
| `ingress.tls.secretName` | Specify the Kubernetes secret containing your SSL certificate. |
| `ingress.annotations` | Specify additional annotations to add to the ingress. |

```yaml
ingress:
  enabled: true
  domain: "autopilot.example.com"
  ingressClassName: "nginx"
  tls:
    enabled: true
    secretName: "autopilot-tls-secret"
  annotations:
    nginx.ingress.kubernetes.io/proxy-body-size: "10m"
```

## Connecting Pega Infinity with Autopilot Service

After deploying the Autopilot Service, configure your Pega Infinity environment to connect to it. There are two ways to do this.

### Option 1: Dynamic System Setting (DSS)

1. Log in to Pega Infinity as an administrator.
2. Navigate to **Records > SysAdmin > Dynamic System Settings** and create a new DSS with the following details:

| Field | Value |
|---|---|
| **Owning Ruleset** | `Pega-Engine` |
| **Setting Purpose** | `prconfig/services/genai/autopilot/servicebaseurl/default` |
| **Value** | `http://<servicename>.<namespace>.svc.cluster.local/` |

3. Replace `<servicename>` with the Autopilot Service name (default: `autopilot`) and `<namespace>` with the Kubernetes namespace where the service is deployed.
4. Save the DSS.

**Example:** If the Autopilot Service is deployed with the default name `autopilot` in the `autopilot` namespace:

```
http://autopilot.autopilot.svc.cluster.local/
```

5. **A restart of Pega Infinity nodes is required for the DSS change to take effect.**

### Option 2: prconfig.xml

Add the Autopilot service URL directly to `charts/pega/config/deploy/prconfig.xml` in the Pega Helm charts repository:

```xml
<env name="services/genai/autopilot/servicebaseurl" value="http://<servicename>.<namespace>.svc.cluster.local/"/>
```

**Example:**

```xml
<env name="services/genai/autopilot/servicebaseurl" value="http://autopilot.autopilot.svc.cluster.local/"/>
```

After editing `prconfig.xml`, apply the change with a `helm upgrade` followed by a rollout restart:

```bash
helm upgrade pega <pega-chart-path> -f my-values.yaml -n <pega-namespace>
kubectl rollout restart statefulset/<pega-deployment-name> -n <pega-namespace>
```

**A restart of Pega Infinity is required whenever the prconfig.xml entry is added or changed.**

### Precedence

If the Autopilot service URL is configured by both methods, **`prconfig.xml` takes precedence over the DSS**.

## OAuth authentication between Pega Infinity and the Autopilot service

You can enable OAuth authentication to secure requests between Pega Infinity and the Autopilot service using an Identity Provider (IdP). The Autopilot service only supports the `private_key_jwt` authentication type.

### How it works

- Pega Infinity obtains a Bearer token from your IdP using an OAuth 2.0 `client_credentials` grant.
- The token is attached as an `Authorization: Bearer <token>` header on every request to the Autopilot service.
- The Autopilot service validates incoming tokens against the IdP public key endpoint (`oauthPublicKeyURL`).

### Scopes

The Autopilot service does not require any OAuth scopes. Leave `autopilot.autopilotAuth.scopes` empty (or omit it entirely) when configuring the Pega chart.

### Shared credentials with SRS and token-minting precedence

Pega Infinity uses a single set of `SERV_AUTH_*` environment variables to mint Bearer tokens for backing services. When both SRS auth (`pegasearch.srsAuth`) and Autopilot auth (`autopilot.autopilotAuth`) are enabled at the same time, **the SRS credentials take precedence** and are used to mint tokens sent to both SRS and the Autopilot service.

This means:

- Both SRS and Autopilot can share the same IdP client application and credentials.
- You do not need separate client registrations for each backing service if both are pointed at the same IdP authorization server.
- If only `autopilot.autopilotAuth` is enabled (SRS auth is disabled), the Autopilot credentials are used to mint tokens.
- If both are enabled and credentials differ, the SRS credentials win — the Autopilot-specific credentials are not used for token minting.

In practice, configure both backing services to trust the same IdP public key endpoint and issue tokens from the same client application. The recommended setup when both services are deployed together:

```yaml
# pega chart values
pegasearch:
  srsAuth:
    enabled: true
    url: "https://your-idp-host/oauth2/v1/token"
    clientId: "your-shared-client-id"
    authType: "private_key_jwt"
    privateKey: "LS0tLS1CRUdJTiBSU0Eg...<base64-encoded-PKCS8-key>"

autopilot:
  autopilotAuth:
    enabled: true
    url: "https://your-idp-host/oauth2/v1/token"
    clientId: "your-shared-client-id"
    privateKey: "LS0tLS1CRUdJTiBSU0Eg...<base64-encoded-PKCS8-key>"
    scopes: ""   # Autopilot requires no scopes
```

Because SRS takes precedence, the token sent to Autopilot is minted using the SRS credentials above. Both services must therefore trust tokens issued for the same client.

### Autopilot service configuration (backingservices chart)

Enable auth on the Autopilot service and set the IdP public key URL so it can validate incoming tokens:

| Parameter | Description | Default |
|---|---|---|
| `authEnabled` | Enables token validation on incoming requests to the Autopilot service. | `false` |
| `oauthPublicKeyURL` | URL of the IdP public key endpoint used to verify Bearer tokens. Required when `authEnabled` is `true`. | `""` |

```yaml
autopilot:
  enabled: true
  authEnabled: true
  oauthPublicKeyURL: "https://your-idp-host/oauth2/v1/keys"
```

### Pega Infinity configuration (pega chart)

Configure the Pega chart to mint tokens and attach them to Autopilot requests. The Autopilot service URL itself is configured via a DSS in Pega Infinity (see "Connecting Pega Infinity with Autopilot Service" above) — no URL parameter is needed in the pega chart.

| Parameter | Description | Default |
|---|---|---|
| `autopilot.autopilotAuth.enabled` | Enables OAuth token minting on the Pega Infinity side. | `false` |
| `autopilot.autopilotAuth.url` | URL of the OAuth service endpoint to obtain a token. | `""` |
| `autopilot.autopilotAuth.clientId` | OAuth client ID. | `""` |
| `autopilot.autopilotAuth.authType` | Authentication type. Only `private_key_jwt` is supported. | `"private_key_jwt"` |
| `autopilot.autopilotAuth.privateKey` | Base64-encoded PKCS8 private key. | `""` |
| `autopilot.autopilotAuth.privateKeyAlgorithm` | Algorithm for the private key. Allowed values: `RS256`, `RS384`, `RS512`, `ES256`, `ES384`, `ES512`. Defaults to `RS256` if not set. | `""` |
| `autopilot.autopilotAuth.scopes` | OAuth scopes to request. The Autopilot service does not require any scopes — leave this empty. | `""` |
| `autopilot.autopilotAuth.external_secret_name` | Name of a pre-existing Kubernetes Secret containing the key (key: `AUTOPILOT_OAUTH_PRIVATE_KEY`). When set, `privateKey` is ignored and no secret is created by the chart. | `""` |

```yaml
autopilot:
  autopilotAuth:
    enabled: true
    url: "https://your-idp-host/oauth2/v1/token"
    clientId: "your-client-id"
    privateKey: "LS0tLS1CRUdJTiBSU0Eg...<base64-encoded-PKCS8-key>"
    privateKeyAlgorithm: "RS256"
    scopes: ""   # no scopes required for Autopilot
```

### Using a pre-existing secret

To avoid placing the private key directly in `values.yaml`, create a Kubernetes Secret beforehand and reference it:

```bash
kubectl create secret generic my-autopilot-auth-secret \
  --namespace <pega-namespace> \
  --from-literal=AUTOPILOT_OAUTH_PRIVATE_KEY="<base64-encoded-PKCS8-key>"
```

Then set in the pega chart values:

```yaml
autopilot:
  autopilotAuth:
    enabled: true
    external_secret_name: "my-autopilot-auth-secret"
```

## Example: Full deployment configuration

```yaml
autopilot:
  enabled: true
  deployment:
    name: "autopilot"

  docker:
    registry:
      url: YOUR_REGISTRY_URL
      username: YOUR_REGISTRY_USERNAME
      password: YOUR_REGISTRY_PASSWORD
    autopilot:
      image: YOUR_REGISTRY_URL/autopilot-service:latest
      imagePullPolicy: Always

  replicas: 2
  enableGenaiHub: false
  modelProviders: "Azure,Vertex,Bedrock"
  deployModelsConfigMap: true

  azure:
    endpoint: "https://my-openai.openai.azure.com/"
    apiKey: "your-azure-api-key"

  aws:
    accessKeyId: "your-access-key-id"
    secretAccessKey: "your-secret-access-key"

  vertex:
    credentials: "base64-encoded-service-account-json"
    project: "my-gcp-project"
    location: "us-central1"

  resources:
    requests:
      cpu: 500m
      memory: "2Gi"
    limits:
      cpu: 2000m
      memory: "4Gi"
```
