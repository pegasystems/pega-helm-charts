{{- define "pegaCredentialsSecretTemplate" }}
kind: Secret
apiVersion: v1
metadata:
  name: {{ template "pegaCredentialsSecret" $ }}
  namespace: {{ .Release.Namespace }}
  annotations:
    "helm.sh/hook": pre-install, pre-upgrade
    "helm.sh/hook-weight": "0"
    "helm.sh/hook-delete-policy": before-hook-creation
data:
  # Base64 encoded username for connecting to the Pega DB
  DB_USERNAME: {{ .Values.global.jdbc.username | b64enc }}
  # Base64 encoded password for connecting to the Pega DB
  DB_PASSWORD: {{ .Values.global.jdbc.password | b64enc }}

 {{ if (eq (include "performDeployment" .) "true") }}
  # Base64 encoded username for connecting to cassandra
  CASSANDRA_USERNAME: {{ .Values.dds.username | b64enc }}
  # Base64 encoded password for connecting to cassandra
  CASSANDRA_PASSWORD: {{ .Values.dds.password | b64enc }}
  {{ if .Values.dds.trustStorePassword -}}
  # Base64 encoded password for the cassandra trust store
  CASSANDRA_TRUSTSTORE_PASSWORD: {{ .Values.dds.trustStorePassword | b64enc }}
  {{- end }}
  {{ if .Values.dds.keyStorePassword -}}
  # Base64 encoded password for the cassandra key store
  CASSANDRA_KEYSTORE_PASSWORD: {{ .Values.dds.keyStorePassword | b64enc }}
  {{- end }}
  {{ if $.Values.hazelcast.enabled }}
  # Base64 encoded username used for authentication in hazelcast client server mode
  HZ_CS_AUTH_USERNAME: {{ .Values.hazelcast.username | b64enc }}
  # Base64 encoded password used for authentication in hazelcast client server mode
  HZ_CS_AUTH_PASSWORD: {{ .Values.hazelcast.password | b64enc }}
  {{ end }}
  {{ range $index, $dep := .Values.global.tier}}
  {{ if and ($dep.pegaDiagnosticUser) (eq $dep.name "web") }}
  # Base64 encoded username for a Tomcat user that will be created with the PegaDiagnosticUser role
  PEGA_DIAGNOSTIC_USER: {{ $dep.pegaDiagnosticUser | b64enc }}
  # Base64 encoded password for a Tomcat user that will be created with the PegaDiagnosticUser role
  PEGA_DIAGNOSTIC_PASSWORD: {{ $dep.pegaDiagnosticPassword | b64enc }}
  {{ end }}
  {{ end }}
{{ end }}
type: Opaque
{{- end }}
