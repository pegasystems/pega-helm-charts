{{/*
App name
*/}}
{{- define "c11n-messaging.name" -}}
{{- default .Root.Chart.Name | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified deployment name. This is used to define the deployment resource name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "c11n-messaging.fullname" -}}
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
{{- define "c11n-messaging.chart" -}}
{{- $name := "c11n-messaging" }}
{{- $version := "0.1.0" }}
{{- printf "%s-%s" $name $version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels for all k8s resources
*/}}
{{- define "c11n-messaging.labels" }}
helm.sh/chart: {{ include "c11n-messaging.chart" . }}
{{ include "c11n-messaging.selectorLabels" . }}
{{- if .Root.Chart.AppVersion }}
app.kubernetes.io/version: {{ .Root.Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Root.Release.Service }}
{{- end }}

{{/*
Selector labels using app and release name
*/}}
{{- define "c11n-messaging.selectorLabels" -}}
app.kubernetes.io/name: {{ include "c11n-messaging.name" . }}
app.kubernetes.io/instance: {{ .Root.Release.Name }}
{{- end }}

{{/*
c11n-messaging-service meta-data labels
*/}}
{{- define "c11n-messaging.c11n-messaging-service.labels" -}}
{{ $data := dict "Root" $ "Name" "c11n-messaging-service" }}
{{- include "c11n-messaging.labels" $data }}
{{- end -}}

{{/*
c11n-messaging-service match labels
*/}}
{{- define "c11n-messaging.c11n-messaging-service.match-labels" }}
app.kubernetes.io/name: c11n-messaging-service
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}

{{/*
c11n-messaging-service registry secret
*/}}
{{- define "c11n-messagingRegistrySecretName" -}}
{{ include "c11n-messaging.fullname" . }}-reg-secret
{{- end }}

{{/*
c11n-messaging-service image repository secret
*/}}
{{- define "imageRepositorySecret" }}
  {{- with .Values.global.imageCredentials }}
  {{- printf "{\"auths\":{\"%s\":{\"username\":\"%s\",\"password\":\"%s\",\"auth\":\"%s\"}}}" .registry .username .password (printf "%s:%s" .username .password | b64enc) | b64enc }}
  {{- end }}
{{- end }}
