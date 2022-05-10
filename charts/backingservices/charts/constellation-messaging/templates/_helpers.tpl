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
{{- $version := "1.0.0" }}
{{- printf "%s-%s" $name $version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

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
