{{/*
App name
*/}}
{{- define "srs.name" -}}
{{- default .Root.Chart.Name | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified deployment name. This is used to define the deployment resource name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "srs.fullname" -}}
{{- if .Values.deploymentName }}
{{- .Values.deploymentName | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.deploymentName }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "srs.chart" -}}
{{- $name := "srs" }}
{{- $version := "0.1.0" }}
{{- printf "%s-%s" $name $version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels for all k8s resources
*/}}
{{- define "srs.labels" }}
helm.sh/chart: {{ include "srs.chart" . }}
{{ include "srs.selectorLabels" . }}
{{- if .Root.Chart.AppVersion }}
app.kubernetes.io/version: {{ .Root.Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Root.Release.Service }}
{{- end }}

{{/*
Selector labels using app and release name
*/}}
{{- define "srs.selectorLabels" -}}
app.kubernetes.io/name: {{ include "srs.name" . }}
app.kubernetes.io/instance: {{ .Root.Release.Name }}
{{- end }}

{{/*
srs-service meta-data labels
*/}}
{{- define "srs.srs-service.labels" -}}
{{ $data := dict "Root" $ "Name" "srs-service" }}
{{- include "srs.labels" $data }}
{{- end -}}

{{/*
srs-service match labels
*/}}
{{- define "srs.srs-service.match-labels" }}
app.kubernetes.io/name: srs-service
app.kubernetes.io/instance: {{ .Release.Name }}
{{ if and (.Values.srsStorage.requireInternetAccess) (not .Values.srsStorage.provisionInternalESCluster) -}}
networking/allow-internet-egress: "true"
{{- end}}
{{- end -}}

{{/*
srs-ops labels
*/}}
{{- define "srs.srs-ops.labels" -}}
{{ $data := dict "Root" $ "Name" "srs-ops" }}
{{ include "srs.labels" $data }}
{{- end -}}

{{/*
srs-ops match labels
*/}}
{{- define "srs.srs-ops.match-labels" -}}
app.kubernetes.io/name: srs-ops
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}

{{- define "srsRegistrySecretName" -}}
{{ include "srs.fullname" . }}-reg-secret
{{- end }}

{{- define "imageRepositorySecret" }}
  {{- with .Values.global.imageCredentials }}
  {{- printf "{\"auths\":{\"%s\":{\"username\":\"%s\",\"password\":\"%s\",\"auth\":\"%s\"}}}" .registry .username .password (printf "%s:%s" .username .password | b64enc) | b64enc }}
  {{- end }}
{{- end }}

{{/*
  elasticsearch url details
*/}}
{{- define "elasticsearch.domain" -}}
{{- if .Values.srsStorage.provisionInternalESCluster }}
{{- $essvcname := "elasticsearch-master" }}
{{- printf "%s.%s.svc" $essvcname .Release.Namespace | trimSuffix "-" }}
{{- else }}
{{- required "A valid '.Values.srsStorage.domain' entry is required when connecting to external Elasticsearch!" .Values.srsStorage.domain }}
{{- end }}
{{- end }}

{{- define "elasticsearch.port" -}}
{{- if .Values.srsStorage.provisionInternalESCluster -}}
{{- quote 9200 }}
{{- else }}
{{- required "A valid '.Values.srsStorage.port' entry is required when connecting to external Elasticsearch!" .Values.srsStorage.port | quote}}
{{- end }}
{{- end }}

{{- define "elasticsearch.protocol" -}}
{{- if .Values.srsStorage.provisionInternalESCluster }}
{{- "http" | quote }}
{{- else }}
{{- required "A valid ''.Values.srsStorage.protocol' entry is required when connecting to external Elasticsearch!" .Values.srsStorage.protocol | quote  }}
{{- end }}
{{- end }}

{{- define "elasticsearch.authProvider" -}}
{{- if (.Values.srsStorage.provisionInternalESCluster) -}}
    {{- "basic-authentication" }}
{{- else }}
    {{- if and  (.Values.srsStorage.basicAuthentication) (not .Values.srsStorage.awsIAM ) -}}
    {{- "basic-authentication" }}
    {{- else if and  (.Values.srsStorage.awsIAM)  (not .Values.srsStorage.basicAuthentication ) -}}
    {{- "aws-iam"}}
    {{- else if and  (not .Values.srsStorage.basicAuthentication ) (not .Values.srsStorage.awsIAM )}}
    {{- "none" }}
    {{- else if and ( .Values.srsStorage.basicAuthentication ) ( .Values.srsStorage.awsIAM )}}
    {{- fail "Only one authentication can be enabled, please try to disable .Values.srsStorage.basicAuthentication/.Values.srsStorage.awsIAM when .Values.srsStorage.provisionInternalESCluster is false" | quote  }}
{{- end }}
{{- end }}
{{- end }}

{{- define "elasticsearch.region" -}}
{{- if  .Values.srsStorage.awsIAM -}}
{{- .Values.srsStorage.awsIAM.region }}
{{- end }}
{{- end }}

{{- define "elasticsearchBasicAuthNUsername" -}}
{{- if  .Values.srsStorage.provisionInternalESCluster  }}
{{- "elastic" |  b64enc }}
{{- else if and (.Values.srsStorage.basicAuthentication) (not .Values.srsStorage.awsIAM) }}
{{- .Values.srsStorage.basicAuthentication.username | b64enc }}
{{- end }}
{{- end}}

{{- define "elasticsearchBasicAuthNPassword" -}}
{{- if  .Values.srsStorage.provisionInternalESCluster }}
{{- randAlphaNum 20 | b64enc}}
{{- else if and (.Values.srsStorage.basicAuthentication) (not .Values.srsStorage.awsIAM) }}
{{- .Values.srsStorage.basicAuthentication.password | b64enc }}
{{- end }}
{{- end}}

{{/*
Network policy: kube-dns
*/}}
{{- define "srs.netpol.kube-dns" -}}
- namespaceSelector:
    matchLabels:
      name: kube-system
- podSelector:
    matchExpressions:
      - key: k8s-app
        operator: In
        values: ["kube-dns", "coredns"]
ports:
- protocol: TCP
  port: 53
- protocol: TCP
  port: 1053
- protocol: TCP
  port: 80
- protocol: TCP
  port: 8080
{{- end -}}
