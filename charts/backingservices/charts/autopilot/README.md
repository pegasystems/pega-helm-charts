# Autopilot Service Helm chart

The Pega `Autopilot Service` backing service provides GenAI-powered capabilities for Pega Infinity Platform by connecting directly to LLM providers (Azure OpenAI, AWS Bedrock, Google Vertex AI). This chart deploys the Autopilot Service for on-premise environments.

## Pega GenaAI Features Support Matrix

| Pega Version | GenAI Connect | GenAI Coach | GenAI Agent |
|---|:---:|:---:|:---:|
| 24.2 | Yes | | |
| 24.2.4 | Yes | Yes | |
| 25.1.1 *(requires HFIX-C4307)* | Yes | Yes | Yes |
| 25.1.2 | Yes | Yes | Yes |

## Configuring a backing service with your pega environment

You can provision the Autopilot Service into your `pega` environment namespace or any namesapce, with the autopilot service endpoint configured for your Pega Infinity environment.

## Supported LLM Providers

| Provider | Authentication Methods |
|---|---|
| Azure OpenAI | API Key, Pre-existing Secret |
| AWS Bedrock | Access Key/Secret, Pre-existing Secret |
| Google Vertex AI | Service Account JSON (base64-encoded), Pre-existing Secret |
| Custom OpenAI-compatible | API Key (inline or pre-existing Secret) |

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
| `service.port` | Defines the port used by the Service. Default is `80`. |
| `service.targetPort` | Defines the port used by the Pod and Container. Default is `8080`. |
| `service.serviceType` | The [type of service](https://kubernetes.io/docs/concepts/services-networking/service/#publishing-services-service-types) you wish to expose. Default is `ClusterIP`. |
| `enableGenaiHub` | Set to `false` for on-prem deployments with direct provider connectivity. Default is `false`. |
| `authEnabled` | Enable or disable authentication for the service. Default is `false`. |
| `isInternalDeployment` | Set to `false` for on-prem deployments. Default is `false`. |
| `modelProviders` | All models from the model list are displayed regardless of provider. For PCFG or Fedramp set the Bedrock as provider. |
| `awsRegion` | AWS region for Bedrock. Default is `us-east-1`. |
| `affinity` | Define pod affinity so that it is restricted to run on particular node(s), or to prefer to run on particular nodes. |
| `tolerations` | Define pod tolerations so that it is allowed to run on node(s) with particular taints. |

## Provider credentials

The Autopilot Service supports two methods for providing LLM provider credentials.

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
| `vertex.credentials` | Base64-encoded Google service account JSON key. The `project_id` is automatically extracted from the JSON. |
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
    location: "us-central1"
```

#### Updating credentials after deployment (Option 1)

When credentials change, run `helm upgrade` with the updated values file. The chart will re-encode and replace the Kubernetes Secret automatically:

```bash
helm upgrade <RELEASE-NAME> <CHART-PATH> \
  --namespace <NAMESPACE> \
  -f <VALUES-FILE> \
  --wait --timeout 5m
```

> **Do not manually patch the Secret when using inline credentials.** The chart uses `b64enc` to encode values into the Secret. Patching the Secret directly will cause double base64 encoding and result in invalid credentials.

### Option 2: Pre-existing Kubernetes Secret

Use a secret that you create and manage outside of this chart. Set `providerCredentialsSecret` to the name of your secret. The secret should contain keys matching the environment variable names used by the service.

| Configuration | Usage |
|---|---|
| `providerCredentialsSecret` | Name of an existing Kubernetes Secret containing provider credentials. When set, inline credentials are ignored. |

The secret can contain credentials for all providers in a single secret. Only the keys relevant to the providers you are using need to be present — all keys are mounted as `optional: true` so missing keys are silently ignored.

| Key | Provider | Description |
|---|---|---|
| `AZURE_ENDPOINT` | Azure OpenAI | Azure OpenAI endpoint URL (e.g. `https://my-openai.openai.azure.com/`) |
| `AZURE_OPENAI_KEY` | Azure OpenAI | Azure OpenAI API key |
| `AWS_ACCESS_KEY_ID` | AWS Bedrock | AWS access key ID |
| `AWS_SECRET_ACCESS_KEY` | AWS Bedrock | AWS secret access key |
| `AWS_SESSION_TOKEN` | AWS Bedrock | AWS session token (optional, for temporary credentials) |
| `VERTEX_AUTH` | Google Vertex AI | Base64-encoded Google service account JSON. The Autopilot service base64-decodes this value itself, so the secret must hold the base64 string — not the raw JSON. |

```yaml
autopilot:
  enabled: true
  providerCredentialsSecret: "my-provider-credentials"
```

Create the secret with credentials for all providers you intend to use. For `VERTEX_AUTH`, pass the base64-encoded JSON using `$(base64 -w0 ...)` so the pod receives the encoded string the service expects:

```bash
kubectl create secret generic my-provider-credentials \
  --namespace <NAMESPACE> \
  --from-literal=AZURE_ENDPOINT="https://my-openai.openai.azure.com/" \
  --from-literal=AZURE_OPENAI_KEY="your-azure-api-key" \
  --from-literal=AWS_ACCESS_KEY_ID="AKIAIOSFODNN7EXAMPLE" \
  --from-literal=AWS_SECRET_ACCESS_KEY="wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY" \
  --from-literal=AWS_SESSION_TOKEN="your-session-token" \
  --from-literal=VERTEX_AUTH="$(base64 -w0 /path/to/gcp-service-account.json)"
```

You can also create the secret from a YAML manifest to manage it in source control. For `VERTEX_AUTH`, put the base64-encoded JSON string directly in `stringData`:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: my-provider-credentials
  namespace: <NAMESPACE>
type: Opaque
stringData:
  AZURE_ENDPOINT: "https://my-openai.openai.azure.com/"
  AZURE_OPENAI_KEY: "your-azure-api-key"
  AWS_ACCESS_KEY_ID: "AKIAIOSFODNN7EXAMPLE"
  AWS_SECRET_ACCESS_KEY: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
  AWS_SESSION_TOKEN: "your-session-token"   # omit if not using temporary credentials
  VERTEX_AUTH: "<base64-encoded-gcp-service-account-json>"  # base64 encode the JSON file: base64 -w0 gcp-service-account.json
```

#### Updating credentials after deployment (Option 2)

When credentials change, patch the existing Secret directly and restart the pod to pick up the new values:

```bash
# Replace the secret with updated values
kubectl create secret generic my-provider-credentials \
  --namespace <NAMESPACE> \
  --from-literal=AZURE_ENDPOINT="https://my-openai.openai.azure.com/" \
  --from-literal=AZURE_OPENAI_KEY="your-azure-api-key" \
  --from-literal=AWS_ACCESS_KEY_ID="AKIAIOSFODNN7EXAMPLE" \
  --from-literal=AWS_SECRET_ACCESS_KEY="wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY" \
  --from-literal=AWS_SESSION_TOKEN="your-session-token" \
  --from-literal=VERTEX_AUTH="$(base64 -w0 /path/to/gcp-service-account.json)" \
  --dry-run=client -o yaml | kubectl apply -f -

# Restart the pod to pick up the updated secret
kubectl rollout restart deployment/<DEPLOYMENT-NAME> -n <NAMESPACE>
kubectl rollout status deployment/<DEPLOYMENT-NAME> -n <NAMESPACE> --timeout=120s
```

## Custom models configuration

The Autopilot Service uses a model list to determine which LLM models are available. The service routes requests to the appropriate provider endpoint based on the model metadata in each entry. The `model_id` field format differs by provider — see [Building model list](#building-model-list) below for the rules per provider.

## Building model list

### For Azure OpenAI and Vertex AI

- **`name`** — Must match the deployment name as it appears in the LLM provider console (e.g., the Azure OpenAI Studio deployment name, or the Vertex AI model ID shown in Model Garden). This value is used as the display name and routing key.
- **`model_path`** — Must be provided as an array of API endpoint paths relative to the provider base URL. For Azure OpenAI the path embeds the Azure portal deployment name (e.g., `["/openai/deployments/gpt-5/chat/completions"]`). For Vertex AI it embeds the model identifier (e.g., `["/google/deployments/gemini-2.5-pro/chat/completions"]`).

### For AWS Bedrock

- **`model_id`** — Must be the **exact model ID as shown in the AWS Bedrock console**, including the cross-region inference prefix and version suffix (e.g., `us.anthropic.claude-3-7-sonnet-20250219-v1:0`, `us.amazon.nova-pro-v1:0`, `amazon.titan-embed-text-v2:0`).


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

#### Updating the model list after deployment (Option 1)

After editing `files/default-models.json`, apply the changes to the running deployment without a full `helm upgrade`:

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

#### Updating the model list after deployment (Option 2)

When the contents of your external ConfigMap change, re-apply it and restart the pod:

```bash
# Re-apply the updated ConfigMap
kubectl create configmap my-models-configmap \
  --namespace <NAMESPACE> \
  --from-file=models.json=./my-models.json \
  --dry-run=client -o yaml | kubectl apply -f -

# Restart the pod to pick up the new ConfigMap
kubectl rollout restart deployment/<DEPLOYMENT-NAME> -n <NAMESPACE>
kubectl rollout status deployment/<DEPLOYMENT-NAME> -n <NAMESPACE> --timeout=120s
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

### Model JSON format

Each model entry in the models file requires the following fields:

| Field | Required | Description |
|---|---|---|
| `provider` | Yes | Cloud provider (`azure`, `bedrock`, `vertex`). |
| `creator` | Yes | Model creator (e.g., `openai`, `anthropic`, `google`, `amazon`). |
| `model_name` | Yes | Model identifer name. |
| `name` | Yes | **For Azure OpenAI and Vertex AI:** must match the deployment name as shown in the LLM provider console. Used as the display name and routing key. |
| `model_id` | Yes | The provider-native model identifier. **For Azure OpenAI and Vertex AI:** the deployment name exactly as shown in the provider console — same value as `model_mapping_id` (e.g., `gpt-5-2025-08-07`, `gemini-2.5-pro`). **For Bedrock:** the exact model ID from the AWS Bedrock console, including cross-region prefix and version suffix (e.g., `us.anthropic.claude-3-7-sonnet-20250219-v1:0`, `us.amazon.nova-pro-v1:0`, `amazon.titan-embed-text-v2:0`). |
| `model_mapping_id` | Yes | Provider-specific deployment name used for API routing and for constructing `model_path` entries. |
| `model_path` | Yes | Array of API endpoint paths for the model, embedding the `model_mapping_id` value. **Azure OpenAI:** `["/openai/deployments/<model_mapping_id>/chat/completions"]`. **Vertex AI:** `["/google/deployments/<model_mapping_id>/chat/completions"]`. **Bedrock:** `["/anthropic/deployments/<model_mapping_id>/converse", "/anthropic/deployments/<model_mapping_id>/converse-stream"]` (adjust creator prefix to match). |
| `type` | Yes | Model type: `chat_completion`, `embedding`, or `image`. |
| `version` | Yes | Model version string. |
| `input_tokens` | No | Maximum input token count. |
| `output_tokens` | No | Maximum output token count. |
| `supported_capabilities` | No | Object describing model capabilities (streaming, functions, multimodal, etc.). |
| `parameters` | No | Object describing tunable parameters (temperature, top_p, max_tokens, etc.). |

## Default fast and smart models

The `fast` and `smart` default models are now defined explicitly in a dedicated `default_models` section at the top level of the models JSON file.

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

### Liveness and readiness probes

The Autopilot Service uses liveness and readiness probes on the `/v1/health` endpoint. Configure probes as part of a `livenessProbe` or `readinessProbe` configuration.

Parameter             | Description    | Default `livenessProbe` | Default `readinessProbe`
---                   | ---            | ---                     | ---
`initialDelaySeconds` | Number of seconds after the container has started before probes are initiated. | `10` | `10`
`timeoutSeconds`      | Number of seconds after which the probe times out. | `10` | `10`
`periodSeconds`       | How often (in seconds) to perform the probe. | `10` | `10`
`successThreshold`    | Minimum consecutive successes for the probe to be considered successful after it determines a failure. | `1` | `1`
`failureThreshold`    | The number of consecutive failures for the pod to be terminated by Kubernetes. | `3` | `3`

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

## Custom OpenAI-compatible providers

The Autopilot Service supports any OpenAI-compatible LLM provider (such as OpenRouter, private LLM servers, or other hosted endpoints) through the `customOpenAI` configuration. Each provider is identified by a `creator` name that maps to environment variables and to the `creator` field in the model JSON.

### Configuration settings

| Configuration | Usage |
|---|---|
| `customOpenAI.providers` | List of custom OpenAI-compatible provider configurations. |
| `customOpenAI.providers[].creator` | Unique identifier for the provider (e.g., `openrouter`, `my-llm-server`). Used to derive env var names and must match the `creator` field in model entries. |
| `customOpenAI.providers[].baseUrl` | Base URL of the provider's OpenAI-compatible API (e.g., `https://openrouter.ai/api/v1`). |
| `customOpenAI.providers[].apiKey` | API key for the provider. Stored in a Kubernetes Secret. |
| `customOpenAI.existingSecret` | Name of a pre-existing Kubernetes Secret containing API keys. When set, `apiKey` fields are not required. |

### How it works

For each entry in `customOpenAI.providers`, the chart injects three environment variables into the pod:

- `CUSTOM_OPENAI_PROVIDERS` — comma-separated list of all configured creator names (e.g., `openrouter,my-llm-server`)
- `CUSTOM_OPENAI_<CREATOR>_BASE_URL` — the provider base URL (plain env var)
- `CUSTOM_OPENAI_<CREATOR>_API_KEY` — the API key (read from a Kubernetes Secret)

`<CREATOR>` is the `creator` value uppercased with hyphens and dots replaced by underscores. For example, `creator: openrouter` → `CUSTOM_OPENAI_OPENROUTER_BASE_URL` / `CUSTOM_OPENAI_OPENROUTER_API_KEY`.

### Option 1: Inline API key (auto-creates Kubernetes Secret)

Provide credentials directly in your values file. The chart automatically creates a Kubernetes Secret.

```yaml
customOpenAI:
  providers:
    - creator: openrouter
      baseUrl: https://openrouter.ai/api/v1
      apiKey: "sk-or-v1-..."
```

Multiple providers can be configured at once:

```yaml
customOpenAI:
  providers:
    - creator: openrouter
      baseUrl: https://openrouter.ai/api/v1
      apiKey: "sk-or-v1-..."
    - creator: my-llm-server
      baseUrl: https://my-llm.internal/v1
      apiKey: "my-internal-key"
```

### Option 2: Pre-existing Kubernetes Secret

Create a secret manually and reference it via `existingSecret`. The secret must contain keys in the format `CUSTOM_OPENAI_<CREATOR>_API_KEY`.

```yaml
customOpenAI:
  existingSecret: "my-custom-openai-secret"
  providers:
    - creator: openrouter
      baseUrl: https://openrouter.ai/api/v1
      # apiKey omitted — read from existingSecret
```

Create the secret:

```bash
kubectl create secret generic my-custom-openai-secret \
  --namespace <NAMESPACE> \
  --from-literal=CUSTOM_OPENAI_OPENROUTER_API_KEY="sk-or-v1-..."
```

### Adding custom OpenAI models to the model list

Models for custom OpenAI-compatible providers use `provider: custom-openai` in the model JSON. The `creator` field must match the `creator` value configured in `customOpenAI.providers`. The `name` field is the model identifier sent to the provider API.

> **Important:** `model_name` must not contain slashes. Use `name` for the actual API model identifier when it contains slashes (e.g., `openrouter/free`).

Example model entry for OpenRouter:

```json
{
  "provider": "custom-openai",
  "creator": "openrouter",
  "model_name": "Free",
  "model_mapping_id": "openrouter-free",
  "name": "openrouter/free",
  "input_tokens": 200000,
  "output_tokens": 8192,
  "type": "chat_completion",
  "model_id": "custom-openai/openrouter/Free/1",
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
    "fast": { "provider": "custom-openai", "creator": "openrouter", "model_id": "custom-openai/openrouter/Free/1", ... },
    "smart": { "provider": "custom-openai", "creator": "openrouter", "model_id": "custom-openai/openrouter/Free/1", ... }
  }
}
```

### Model JSON fields for custom-openai provider

| Field | Description |
|---|---|
| `provider` | Must be `custom-openai`. |
| `creator` | Must match the `creator` value in `customOpenAI.providers`. Used to look up `CUSTOM_OPENAI_<CREATOR>_BASE_URL` and `CUSTOM_OPENAI_<CREATOR>_API_KEY`. |
| `model_name` | Display name for the model. Must not contain slashes. |
| `name` | Actual model identifier sent to the provider API as the `model=` parameter (e.g., `openrouter/free`, `openai/gpt-4o`). |
| `model_path` | Array with the API path suffix appended to `baseUrl` (e.g., `["/chat/completions"]`). |
| `model_id` | Must follow the format `custom-openai/<creator>/<model_name>/<version>` (no slashes in `model_name`). |

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
    location: "us-central1"

  customOpenAI:
    providers:
      - creator: openrouter
        baseUrl: https://openrouter.ai/api/v1
        apiKey: "sk-or-v1-..."

  resources:
    requests:
      cpu: 500m
      memory: "1Gi"
    limits:
      cpu: "1"
      memory: "2Gi"
```
