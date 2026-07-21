# Autopilot Service Helm chart

The Pega `Autopilot Service` backing service provides GenAI capabilities for Pega Platform in client-managed cloud and on-premises deployments. This chart deploys the service as an external backing service that connects directly to supported LLM providers without requiring the Pega GenAI Hub gateway. The chart also supports optional OAuth authentication between Pega Platform and the Autopilot Service, a deployment-scoped model list through a bundled or external ConfigMap, and custom OpenAI-compatible providers.

## Configuring a backing service with your Pega environment

You can deploy the Autopilot Service in the `pega` namespace or in a separate namespace. After deployment, configure the service endpoint in Pega Platform so that the platform can route Autopilot requests to the backing service.

## Pega GenAI feature support matrix

| Pega Platform version | GenAI Connect | GenAI Coach | GenAI Agent |
|---|:---:|:---:|:---:|
| 24.2 | Yes | No | No |
| 24.2.4 | Yes | Yes | No |
| 25.1.1 (`HFIX-C4307` required) | Yes | Yes | Yes |
| 25.1.2 | Yes | Yes | Yes |

## Supported LLM providers

| Provider | Authentication Methods |
|---|---|
| Azure OpenAI | API Key, Existing Secret |
| AWS Bedrock | Access Key/Secret, Existing Secret |
| Google Vertex AI | Service Account JSON (base64-encoded), Existing Secret |
| Custom OpenAI-compatible | API Key (inline or existing Secret) |

## Autopilot Service configuration

The `values.yaml` file defines the deployment, networking, authentication, provider credentials, and model-list settings for the Autopilot Service.

### Configuration settings

| Configuration | Usage |
|---|---|
| `enabled` | Set this parameter to `true` to deploy the Autopilot Service as a backing service. |
| `deployment.name` | Specify the deployment name. The chart prefixes Autopilot resources with this value. Default: `autopilot` |
| `docker.registry.url` | Specify the image registry URL. |
| `docker.registry.username` | Specify the image registry username. |
| `docker.registry.password` | Specify the image registry password. |
| `docker.imagePullSecretNames` | Specify an array of existing image pull Secrets to use when pulling images. |
| `docker.autopilot.image` | Specify the Autopilot Service image and tag. |
| `docker.autopilot.imagePullPolicy` | Specify the image pull policy. Default: `IfNotPresent` |
| `replicas` | Specify the number of Autopilot pods. Default: `2` |
| `resources` | Define CPU and memory requests and limits for the Autopilot container. |
| `service.port` | Specify the Kubernetes Service port. Default: `80` |
| `service.targetPort` | Specify the container port exposed by the service. Default: `8080` |
| `service.serviceType` | Specify the Kubernetes Service type. Default: `ClusterIP` |
| `enableGenaiHub` | Leave this parameter as `false` for client-managed cloud or on-premises deployments that connect directly to providers. Set it to `true` only when requests must be routed through the Pega GenAI Hub gateway. |
| `authEnabled` | Set this parameter to `true` to require OAuth bearer tokens from Pega Platform. Default: `false` |
| `oauthPublicKeyURL` | Specify the identity provider public key endpoint used to validate incoming bearer tokens. This value is required when `authEnabled` is `true`. |
| `isInternalDeployment` | Leave this parameter as `false` for client-managed cloud or on-premises deployments. |
| `modelProviders` | Specify a comma-separated list of providers to expose, such as `"Azure,Vertex,Bedrock"`. If this value is empty, the service relies on the configured credentials and model list. For FedRAMP or PCFG environments, set `"Bedrock"` as the provider. |
| `awsRegion` | Specify the AWS region used by Bedrock. Default: `us-east-1` |
| `affinity` | Define pod affinity rules. |
| `tolerations` | Define pod tolerations. |

## Provider credentials

The Autopilot Service supports two credential patterns for Azure OpenAI, AWS Bedrock, and Google Vertex AI.

### Option 1: Inline credentials

Specify provider credentials directly in your values file. When you use inline credentials, the chart automatically creates a Kubernetes Secret named `<deployment.name>-provider-credentials`.

| Configuration | Usage |
|---|---|
| `azure.endpoint` | Azure OpenAI endpoint URL, for example `https://YOUR_AZURE_RESOURCE.openai.azure.com/`. |
| `azure.apiKey` | Azure OpenAI API key. |
| `azure.apiVersion` | Azure OpenAI API version. Default: `2024-10-21` |
| `aws.accessKeyId` | AWS access key ID for Bedrock. |
| `aws.secretAccessKey` | AWS secret access key for Bedrock. |
| `aws.sessionToken` | Optional AWS session token for temporary credentials. |
| `vertex.credentials` | Base64-encoded Google service account JSON. The `project_id` is automatically extracted from the JSON. |
| `vertex.location` | Vertex AI location. Default: `us-central1` |

```yaml
autopilot:
  enabled: true
  azure:
    endpoint: "https://YOUR_AZURE_RESOURCE.openai.azure.com/"
    apiKey: "YOUR_AZURE_OPENAI_API_KEY"
  aws:
    accessKeyId: "YOUR_AWS_ACCESS_KEY_ID"
    secretAccessKey: "YOUR_AWS_SECRET_ACCESS_KEY"
    sessionToken: ""
  vertex:
    credentials: "YOUR_BASE64_ENCODED_VERTEX_SERVICE_ACCOUNT_JSON"
    location: "us-central1"
```

#### Updating inline credentials after deployment

When credentials change, run `helm upgrade` with the updated values file. The chart automatically re-encodes and replaces the Kubernetes Secret:

```bash
helm upgrade <RELEASE-NAME> <CHART-PATH> \
  --namespace <NAMESPACE> \
  -f <VALUES-FILE> \
  --wait --timeout 5m
```

> **Do not manually patch the Secret when using inline credentials.** The chart uses `b64enc` to encode values into the Secret. Patching the Secret directly can cause double base64 encoding and result in invalid credentials.

### Option 2: Existing Kubernetes Secret

Use a Secret that you create and manage outside of this chart. Set `providerCredentialsSecret` to the name of your Secret. The Secret should contain keys that match the environment variable names used by the service.

| Configuration | Usage |
|---|---|
| `providerCredentialsSecret` | Name of an existing Kubernetes Secret that contains provider credentials. When set, inline credentials are ignored. |

The Secret can contain credentials for all providers in a single Secret. Only the keys that are relevant to the providers you are using need to be present. All keys are mounted as `optional: true`, so missing keys are silently ignored.

| Key | Provider | Description |
|---|---|---|
| `AZURE_ENDPOINT` | Azure OpenAI | Azure OpenAI endpoint URL |
| `AZURE_OPENAI_KEY` | Azure OpenAI | Azure OpenAI API key |
| `AWS_ACCESS_KEY_ID` | AWS Bedrock | AWS access key ID |
| `AWS_SECRET_ACCESS_KEY` | AWS Bedrock | AWS secret access key |
| `AWS_SESSION_TOKEN` | AWS Bedrock | Optional AWS session token |
| `VERTEX_AUTH` | Google Vertex AI | Base64-encoded Google service account JSON. The Autopilot Service base64-decodes this value itself, so the Secret must contain the base64 string rather than the raw JSON. |

```yaml
autopilot:
  enabled: true
  providerCredentialsSecret: "autopilot-provider-credentials"
```

Create the Secret with credentials for all providers that you intend to use. For `VERTEX_AUTH`, pass the base64-encoded JSON so the pod receives the encoded string the service expects:

```bash
kubectl create secret generic autopilot-provider-credentials \
  --namespace <NAMESPACE> \
  --from-literal=AZURE_ENDPOINT="https://YOUR_AZURE_RESOURCE.openai.azure.com/" \
  --from-literal=AZURE_OPENAI_KEY="YOUR_AZURE_OPENAI_API_KEY" \
  --from-literal=AWS_ACCESS_KEY_ID="YOUR_AWS_ACCESS_KEY_ID" \
  --from-literal=AWS_SECRET_ACCESS_KEY="YOUR_AWS_SECRET_ACCESS_KEY" \
  --from-literal=AWS_SESSION_TOKEN="YOUR_AWS_SESSION_TOKEN" \
  --from-literal=VERTEX_AUTH="$(base64 -w0 /path/to/gcp-service-account.json)"
```

You can also create the Secret from a YAML manifest to manage it in source control. For `VERTEX_AUTH`, put the base64-encoded JSON string directly in `stringData`:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: autopilot-provider-credentials
  namespace: <NAMESPACE>
type: Opaque
stringData:
  AZURE_ENDPOINT: "https://YOUR_AZURE_RESOURCE.openai.azure.com/"
  AZURE_OPENAI_KEY: "YOUR_AZURE_OPENAI_API_KEY"
  AWS_ACCESS_KEY_ID: "YOUR_AWS_ACCESS_KEY_ID"
  AWS_SECRET_ACCESS_KEY: "YOUR_AWS_SECRET_ACCESS_KEY"
  AWS_SESSION_TOKEN: "YOUR_AWS_SESSION_TOKEN"   # omit if not using temporary credentials
  VERTEX_AUTH: "YOUR_BASE64_ENCODED_VERTEX_SERVICE_ACCOUNT_JSON"  # base64 encode the JSON file: base64 -w0 gcp-service-account.json
```

#### Updating a credentials Secret after deployment

When credentials change, patch the existing Secret directly and restart the pod to pick up the new values:

```bash
# Replace the Secret with updated values
kubectl create secret generic autopilot-provider-credentials \
  --namespace <NAMESPACE> \
  --from-literal=AZURE_ENDPOINT="https://YOUR_AZURE_RESOURCE.openai.azure.com/" \
  --from-literal=AZURE_OPENAI_KEY="YOUR_AZURE_OPENAI_API_KEY" \
  --from-literal=AWS_ACCESS_KEY_ID="YOUR_AWS_ACCESS_KEY_ID" \
  --from-literal=AWS_SECRET_ACCESS_KEY="YOUR_AWS_SECRET_ACCESS_KEY" \
  --from-literal=AWS_SESSION_TOKEN="YOUR_AWS_SESSION_TOKEN" \
  --from-literal=VERTEX_AUTH="$(base64 -w0 /path/to/gcp-service-account.json)" \
  --dry-run=client -o yaml | kubectl apply -f -

# Restart the pod to pick up the updated Secret
kubectl rollout restart deployment/<DEPLOYMENT-NAME> -n <NAMESPACE>
kubectl rollout status deployment/<DEPLOYMENT-NAME> -n <NAMESPACE> --timeout=120s
```

## Model list configuration

The Autopilot Service uses a model list to determine which models are available and how requests are routed to each provider. When a model list is mounted by the chart, the container receives it through the `LOCAL_MODELS_FILE` environment variable at `/config/models.json`. The `model_id` field format differs by provider.

### Building model list

#### For Azure OpenAI and Vertex AI

- **`name`** must match the deployment name as it appears in the LLM provider console (for example, the Azure OpenAI Studio deployment name or the Vertex AI model ID shown in Model Garden). This value is used as the display name and routing key.
- **`model_path`** must be provided as an array of API endpoint paths relative to the provider base URL. For Azure OpenAI, the path embeds the Azure portal deployment name (for example, `["/openai/deployments/gpt-5/chat/completions"]`). For Vertex AI, it embeds the model identifier (for example, `["/google/deployments/gemini-2.5-pro/chat/completions"]`).

#### For AWS Bedrock

- **`model_id`** must be the exact model ID as shown in the AWS Bedrock console, including the cross-region inference prefix and version suffix (for example, `us.anthropic.claude-3-7-sonnet-20250219-v1:0`, `us.amazon.nova-pro-v1:0`, `amazon.titan-embed-text-v2:0`).

### Option 1: Bundled models file (recommended)

By default, the chart creates a ConfigMap from `files/default-models.json` during `helm install` or `helm upgrade`.

| Configuration | Usage |
|---|---|
| `deployModelsConfigMap` | Set to `true` to create a ConfigMap from the bundled `files/default-models.json` file. Default: `true` |

```yaml
autopilot:
  enabled: true
  deployModelsConfigMap: true
```

To customize the default model list before deployment, edit `files/default-models.json` in the chart directory.

#### Updating the bundled model list after deployment

After you edit `files/default-models.json`, apply the changes to the running deployment without a full `helm upgrade`:

```bash
# Update the ConfigMap from the bundled file
kubectl create configmap autopilot-models \
  --namespace <NAMESPACE> \
  --from-file=models.json=./files/default-models.json \
  --dry-run=client -o yaml | kubectl apply -f -

# Restart the pod to pick up the new ConfigMap
kubectl rollout restart deployment/<DEPLOYMENT-NAME> -n <NAMESPACE>
kubectl rollout status deployment/<DEPLOYMENT-NAME> -n <NAMESPACE> --timeout=120s
```

### Option 2: Existing ConfigMap

Use a ConfigMap that you create and manage outside of this chart. The ConfigMap must include a `models.json` key.

| Configuration | Usage |
|---|---|
| `customModels.existingConfigMap` | Name of an existing ConfigMap that contains `models.json`. |

```yaml
autopilot:
  enabled: true
  deployModelsConfigMap: false
  customModels:
    existingConfigMap: "autopilot-models"
```

Create the ConfigMap:

```bash
kubectl create configmap autopilot-models \
  --namespace <NAMESPACE> \
  --from-file=models.json=./my-models.json
```

#### Updating a model list ConfigMap after deployment

When the contents of your external ConfigMap change, re-apply it and restart the pod:

```bash
# Re-apply the updated ConfigMap
kubectl create configmap autopilot-models \
  --namespace <NAMESPACE> \
  --from-file=models.json=./my-models.json \
  --dry-run=client -o yaml | kubectl apply -f -

# Restart the pod to pick up the updated ConfigMap
kubectl rollout restart deployment/<DEPLOYMENT-NAME> -n <NAMESPACE>
kubectl rollout status deployment/<DEPLOYMENT-NAME> -n <NAMESPACE> --timeout=120s
```

### Option 3: Inline model list

Specify the JSON content directly in `values.yaml`. The chart creates a ConfigMap from the inline content.

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
          "model_id": "gpt-5-2025-08-07",
          "input_tokens": 400000,
          "output_tokens": 128000,
          "type": "chat_completion",
          "version": "2025-08-07",
          "model_path": ["/openai/deployments/gpt-5-2025-08-07/chat/completions"],
          "supported_capabilities": {
            "streaming": true,
            "functions": true,
            "json_mode": true
          }
        }
      ]
```

If you set `deployModelsConfigMap: false` and do not provide `customModels.existingConfigMap` or `customModels.inline`, the service uses its built-in model list.

### Model JSON format

| Field | Required | Description |
|---|---|---|
| `provider` | Yes | Provider name: `azure`, `bedrock`, `vertex`, or `custom`. |
| `creator` | Yes | Model creator, for example `openai`, `anthropic`, `google`, or `amazon`. |
| `model_name` | Yes | Display name for the model. |
| `name` | Yes | Routing name used by the Autopilot Service. For Azure OpenAI and Vertex AI, this value must match the deployment name shown in the provider console. |
| `model_id` | Yes | Provider-native model identifier. For Azure OpenAI and Vertex AI, use the deployment name as shown in the provider console. For AWS Bedrock, use the exact model ID from the AWS Bedrock console, including any cross-region prefix or version suffix. |
| `model_mapping_id` | Yes | Deployment or provider-specific identifier used to build `model_path` entries. |
| `model_path` | Yes | Array of API endpoint paths used by the service to reach the provider API. For Azure OpenAI, use `["/openai/deployments/<model_mapping_id>/chat/completions"]`. For Vertex AI, use `["/google/deployments/<model_mapping_id>/chat/completions"]`. For Bedrock, use `["/anthropic/deployments/<model_mapping_id>/converse", "/anthropic/deployments/<model_mapping_id>/converse-stream"]` (adjust the creator prefix to match the model creator). |
| `type` | Yes | Model type, for example `chat_completion`, `embedding`, or `image`. |
| `version` | Yes | Model version string. |
| `input_tokens` | No | Maximum supported input tokens. |
| `output_tokens` | No | Maximum supported output tokens. |
| `supported_capabilities` | No | Object that describes supported capabilities such as streaming, function calling, or multimodal support. |
| `parameters` | No | Object that describes tunable model parameters, such as `temperature`, `top_p`, or `max_tokens`. |

## Default fast and smart models

The `fast` and `smart` default models are defined explicitly in a dedicated `default_models` section at the top level of the models JSON file.

The `default_models` section contains two keys — `fast` and `smart` — each holding the full model object of the intended default:

```json
{
  "models": [ ... ],
  "default_models": {
    "fast": {
      "provider": "azure",
      "model_id": "gpt-5-chat-2025-08-07",
      ...
    },
    "smart": {
      "provider": "azure",
      "model_id": "gpt-5-mini-2025-08-07",
      ...
    }
  }
}
```

To change which models serve as defaults, update the `default_models.fast` and `default_models.smart` entries in your models JSON file to the desired model objects. The model referenced must also be present in the `models` array.

## Liveness and readiness probes

The Autopilot Service uses liveness and readiness probes on the `/v1/health` endpoint. Configure probes as part of a `livenessProbe` or `readinessProbe` configuration.

| Parameter | Description | Default `livenessProbe` | Default `readinessProbe` |
|---|---|---|---|
| `initialDelaySeconds` | Number of seconds to wait before starting probes. | `10` | `10` |
| `timeoutSeconds` | Probe timeout in seconds. | `10` | `10` |
| `periodSeconds` | Probe interval in seconds. | `10` | `10` |
| `successThreshold` | Minimum consecutive successes required after a failure. | `1` | `1` |
| `failureThreshold` | Number of failures before Kubernetes marks the probe as failed. | `3` | `3` |

Example:

```yaml
livenessProbe:
  initialDelaySeconds: 10
  timeoutSeconds: 10
  periodSeconds: 10
  successThreshold: 1
  failureThreshold: 3
readinessProbe:
  initialDelaySeconds: 10
  timeoutSeconds: 10
  periodSeconds: 10
  successThreshold: 1
  failureThreshold: 3
```

## Ingress

To expose the Autopilot Service outside the cluster, enable ingress.

| Configuration | Usage |
|---|---|
| `ingress.enabled` | Leave as `false` unless you need to expose the Autopilot Service outside the cluster. |
| `ingress.domain` | Specify the ingress host name. |
| `ingress.ingressClassName` | Specify the ingress class name. |
| `ingress.tls.enabled` | Set to `true` to enable TLS termination at the ingress. |
| `ingress.tls.secretName` | Specify the Kubernetes Secret that stores the TLS certificate. |
| `ingress.annotations` | Specify additional annotations for the ingress. |

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

## Connecting Pega Platform to the Autopilot Service

After you deploy the Autopilot Service, configure Pega Platform to use the service URL using Dynamic System Setting (DSS) or `prconfig.xml`.

> If both methods are configured, `prconfig.xml` takes precedence over the DSS.

### Option 1: Dynamic System Setting (DSS)

1. Log in to Pega Platform as an administrator.
2. Navigate to **Records > SysAdmin > Dynamic System Settings** and create a new DSS with the following values:

| Field | Value |
|---|---|
| Owning Ruleset | `Pega-Engine` |
| Setting Purpose | `prconfig/services/genai/autopilot/servicebaseurl/default` |
| Value | `http://<service-name>.<namespace>.svc.cluster.local:<service-port>/` |

3. Replace `<service-name>` with the Autopilot Service name (default: `autopilot`) and `<namespace>` with the Kubernetes namespace where the service is deployed.
4. Save the DSS.

**Example:** If the Autopilot Service is deployed with the default name `autopilot` in the `autopilot` namespace:

```text
http://autopilot.autopilot.svc.cluster.local/
```

5. Restart Pega Platform nodes for the DSS change to take effect.

### Option 2: `prconfig.xml`

Add the Autopilot Service URL directly to `charts/pega/config/deploy/prconfig.xml` in the Pega Helm charts repository:

```xml
<env name="services/genai/autopilot/servicebaseurl" value="http://<service-name>.<namespace>.svc.cluster.local:<service-port>/"/>
```

Example:

```xml
<env name="services/genai/autopilot/servicebaseurl" value="http://autopilot.autopilot.svc.cluster.local/"/>
```

After you edit `prconfig.xml`, apply the change with a `helm upgrade` followed by a rollout restart:

```bash
helm upgrade pega <PEGA-CHART-PATH> -f my-values.yaml -n <PEGA-NAMESPACE>
kubectl rollout restart statefulset/<PEGA-DEPLOYMENT-NAME> -n <PEGA-NAMESPACE>
```

A restart of Pega Platform is required whenever the `prconfig.xml` entry is added or changed.

## OAuth authentication between Pega Platform and the Autopilot Service

You can enable OAuth authentication to secure requests between Pega Platform and the Autopilot Service using an Identity Provider (IdP). The Autopilot Service only supports the `private_key_jwt` authentication type.

### How it works

- Pega Platform obtains a bearer token from your IdP using an OAuth 2.0 `client_credentials` grant.
- The token is attached as an `Authorization: Bearer <token>` header on every request to the Autopilot Service.
- The Autopilot Service validates incoming tokens against the IdP public key endpoint (`oauthPublicKeyURL`).

### Scopes

The Autopilot Service does not require any OAuth scopes. Leave `autopilot.autopilotAuth.scopes` empty (or omit it entirely) when configuring the Pega chart.

### Autopilot Service configuration in the backingservices chart

Enable auth on the Autopilot Service and set the IdP public key URL so it can validate incoming tokens:

| Parameter | Description | Default |
|---|---|---|
| `authEnabled` | Enables token validation on incoming requests to the Autopilot Service. | `false` |
| `oauthPublicKeyURL` | URL of the IdP public key endpoint used to verify bearer tokens. Required when `authEnabled` is `true`. | `""` |

```yaml
autopilot:
  enabled: true
  authEnabled: true
  oauthPublicKeyURL: "https://YOUR_IDP_HOST/oauth2/v1/keys"
```

### Pega Platform configuration in the `pega` chart

Configure the Pega chart to mint tokens and attach them to Autopilot requests. The Autopilot Service URL itself is configured via a DSS or `prconfig.xml` (see [Connecting Pega Platform to the Autopilot Service](#connecting-pega-platform-to-the-autopilot-service)) — no URL parameter is needed in the `pega` chart.

| Parameter | Description | Default |
|---|---|---|
| `autopilot.autopilotAuth.enabled` | Enables OAuth token minting on the Pega Platform side. | `false` |
| `autopilot.autopilotAuth.url` | URL of the OAuth service endpoint to obtain a token. | `""` |
| `autopilot.autopilotAuth.clientId` | OAuth client ID. | `""` |
| `autopilot.autopilotAuth.authType` | Authentication type. Only `private_key_jwt` is supported. | `"private_key_jwt"` |
| `autopilot.autopilotAuth.privateKey` | Base64-encoded PKCS8 private key. | `""` |
| `autopilot.autopilotAuth.privateKeyAlgorithm` | Algorithm for the private key. Allowed values: `RS256`, `RS384`, `RS512`, `ES256`, `ES384`, `ES512`. Default: `RS256` | `""` |
| `autopilot.autopilotAuth.scopes` | OAuth scopes to request. The Autopilot Service does not require any scopes — leave this empty. | `""` |
| `autopilot.autopilotAuth.external_secret_name` | Name of an existing Kubernetes Secret containing the key (key: `AUTOPILOT_OAUTH_PRIVATE_KEY`). When set, `privateKey` is ignored and no Secret is created by the chart. | `""` |

```yaml
autopilot:
  autopilotAuth:
    enabled: true
    url: "https://YOUR_IDP_HOST/oauth2/v1/token"
    clientId: "YOUR_CLIENT_ID"
    authType: "private_key_jwt"
    privateKey: "YOUR_BASE64_ENCODED_PKCS8_PRIVATE_KEY"
    privateKeyAlgorithm: "RS256"
    scopes: ""   # no scopes required for Autopilot
```

### Using an existing Secret for the private key

To avoid placing the private key directly in `values.yaml`, create a Kubernetes Secret beforehand and reference it:

```bash
kubectl create secret generic autopilot-auth-secret \
  --namespace <PEGA-NAMESPACE> \
  --from-literal=AUTOPILOT_OAUTH_PRIVATE_KEY="YOUR_BASE64_ENCODED_PKCS8_PRIVATE_KEY"
```

Then reference the Secret in the `pega` chart:

```yaml
autopilot:
  autopilotAuth:
    enabled: true
    external_secret_name: "autopilot-auth-secret"
```

### Shared credentials with SRS and token-minting precedence

Pega Platform uses a single set of `SERV_AUTH_*` environment variables to mint bearer tokens for backing services. When both SRS auth (`pegasearch.srsAuth`) and Autopilot auth (`autopilot.autopilotAuth`) are enabled at the same time, **the SRS credentials take precedence** and are used to mint tokens sent to both SRS and the Autopilot Service.

This means:

- Both SRS and Autopilot can share the same IdP client application and credentials.
- You do not need separate client registrations for each backing service if both are pointed at the same IdP authorization server.
- If only `autopilot.autopilotAuth` is enabled (SRS auth is disabled), the Autopilot credentials are used to mint tokens.
- If both are enabled and the credentials differ, the SRS credentials win — the Autopilot-specific credentials are not used for token minting.

In practice, configure both backing services to trust the same IdP public key endpoint and issue tokens from the same client application. The recommended setup when both services are deployed together:

```yaml
pegasearch:
  srsAuth:
    enabled: true
    url: "https://YOUR_IDP_HOST/oauth2/v1/token"
    clientId: "YOUR_SHARED_CLIENT_ID"
    authType: "private_key_jwt"
    privateKey: "YOUR_BASE64_ENCODED_PKCS8_PRIVATE_KEY"
autopilot:
  autopilotAuth:
    enabled: true
    url: "https://YOUR_IDP_HOST/oauth2/v1/token"
    clientId: "YOUR_SHARED_CLIENT_ID"
    privateKey: "YOUR_BASE64_ENCODED_PKCS8_PRIVATE_KEY"
    scopes: ""   # Autopilot requires no scopes
```

Because SRS takes precedence, the token sent to Autopilot is minted using the SRS credentials above. Both services must therefore trust tokens issued for the same client.

## Custom OpenAI-compatible providers

The Autopilot Service supports OpenAI-compatible providers (such as OpenRouter, private LLM servers, or other hosted endpoints) through the `customOpenAI` section. Each provider is identified by a `creator` value that maps to environment variables and to the `creator` field in the model JSON.

### Configuration settings

| Configuration | Usage |
|---|---|
| `customOpenAI.providers` | List of custom OpenAI-compatible provider configurations. Each entry produces a set of environment variables in the pod. |
| `customOpenAI.providers[].creator` | Provider identifier (for example, `openrouter`, `my-llm-server`). Must match the `creator` field in model entries. |
| `customOpenAI.providers[].baseUrl` | Base URL for the provider API (for example, `https://openrouter.ai/api/v1`). |
| `customOpenAI.providers[].apiKey` | Provider API key. The chart stores inline values in a Secret. |
| `customOpenAI.existingSecret` | Name of an existing Secret that contains provider API keys. When set, `apiKey` fields are not required. |

For each entry in `customOpenAI.providers`, the chart injects these environment variables into the pod:

- `CUSTOM_OPENAI_PROVIDERS` as a comma-separated list of all `creator` values
- `CUSTOM_OPENAI_<CREATOR>_BASE_URL` for the provider base URL
- `CUSTOM_OPENAI_<CREATOR>_API_KEY` for the provider API key

`<CREATOR>` is the `creator` value uppercased with hyphens and dots replaced by underscores. For example, `creator: openrouter` produces `CUSTOM_OPENAI_OPENROUTER_BASE_URL` and `CUSTOM_OPENAI_OPENROUTER_API_KEY`.

### Option 1: Inline API key

```yaml
autopilot:
  customOpenAI:
    providers:
      - creator: openrouter
        baseUrl: "https://openrouter.ai/api/v1"
        apiKey: "YOUR_OPENROUTER_API_KEY"
```

Multiple providers can be configured at once:

```yaml
autopilot:
  customOpenAI:
    providers:
      - creator: openrouter
        baseUrl: "https://openrouter.ai/api/v1"
        apiKey: "YOUR_OPENROUTER_API_KEY"
      - creator: my-llm-server
        baseUrl: "https://my-llm.internal/v1"
        apiKey: "YOUR_INTERNAL_API_KEY"
```

### Option 2: Existing Kubernetes Secret

Create a Secret manually and reference it via `existingSecret`. When set, inline `apiKey` values are not required.

```yaml
autopilot:
  customOpenAI:
    existingSecret: "custom-openai-credentials"
    providers:
      - creator: openrouter
        baseUrl: "https://openrouter.ai/api/v1"
        # apiKey omitted — read from existingSecret
```

The Secret must use keys in the format `CUSTOM_OPENAI_<CREATOR>_API_KEY`. For `creator: openrouter`, the required key is `CUSTOM_OPENAI_OPENROUTER_API_KEY`.

Create the Secret:

```bash
kubectl create secret generic custom-openai-credentials \
  --namespace <NAMESPACE> \
  --from-literal=CUSTOM_OPENAI_OPENROUTER_API_KEY="YOUR_OPENROUTER_API_KEY"
```

### Adding custom OpenAI models to the model list

Models for custom OpenAI-compatible providers use `provider: "custom"` in the model JSON. The `creator` field must match the `creator` value configured in `customOpenAI.providers`. The `name` field is the model identifier sent to the provider API.

> **Important:** `model_name` must not contain slashes. Use `name` for the actual API model identifier when it contains slashes (for example, `openrouter/free`).

Example model entry:

```json
{
  "provider": "custom",
  "creator": "openrouter",
  "model_name": "Free",
  "model_mapping_id": "openrouter-free",
  "name": "openrouter/free",
  "input_tokens": 200000,
  "output_tokens": 8192,
  "type": "chat_completion",
  "model_id": "custom/openrouter/Free/1",
  "version": "1",
  "model_path": ["/chat/completions"],
  "supported_capabilities": {
    "streaming": true,
    "functions": true,
    "json_mode": true,
    "is_multimodal": false
  }
}
```

To use this model as the default fast and smart model, set it in the `default_models` section of your models JSON:

```json
{
  "models": [ ... ],
  "default_models": {
    "fast": { "provider": "custom", "creator": "openrouter", "model_id": "custom/openrouter/Free/1", ... },
    "smart": { "provider": "custom", "creator": "openrouter", "model_id": "custom/openrouter/Free/1", ... }
  }
}
```

### Model JSON fields for custom providers

| Field | Description |
|---|---|
| `provider` | Must be `custom`. |
| `creator` | Must match the `creator` value in `customOpenAI.providers`. Used to look up `CUSTOM_OPENAI_<CREATOR>_BASE_URL` and `CUSTOM_OPENAI_<CREATOR>_API_KEY`. |
| `model_name` | Display name for the model. Must not contain slashes. |
| `name` | Actual model identifier sent to the provider API as the `model` parameter (for example, `openrouter/free`). |
| `model_path` | Array with the API path suffix appended to `baseUrl` (for example, `["/chat/completions"]`). |
| `model_id` | Must follow the format `custom/<creator>/<model_name>/<version>` (no slashes in `model_name`). |

## Example deployment configuration

```yaml
autopilot:
  enabled: true
  deployment:
    name: "autopilot"

  docker:
    registry:
      url: "YOUR_REGISTRY_URL"
      username: "YOUR_REGISTRY_USERNAME"
      password: "YOUR_REGISTRY_PASSWORD"
    autopilot:
      image: "YOUR_REGISTRY_URL/autopilot-service:YOUR_TAG"
      imagePullPolicy: Always

  replicas: 2
  enableGenaiHub: false
  authEnabled: false
  isInternalDeployment: false
  modelProviders: "Azure,Vertex,Bedrock"
  deployModelsConfigMap: true

  azure:
    endpoint: "https://YOUR_AZURE_RESOURCE.openai.azure.com/"
    apiKey: "YOUR_AZURE_OPENAI_API_KEY"

  aws:
    accessKeyId: "YOUR_AWS_ACCESS_KEY_ID"
    secretAccessKey: "YOUR_AWS_SECRET_ACCESS_KEY"

  vertex:
    credentials: "YOUR_BASE64_ENCODED_VERTEX_SERVICE_ACCOUNT_JSON"
    location: "us-central1"

  customOpenAI:
    providers:
      - creator: openrouter
        baseUrl: "https://openrouter.ai/api/v1"
        apiKey: "YOUR_OPENROUTER_API_KEY"

  resources:
    requests:
      cpu: 500m
      memory: "1Gi"
    limits:
      cpu: "1"
      memory: "2Gi"
```
