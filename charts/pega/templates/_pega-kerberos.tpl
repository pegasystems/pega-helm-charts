{{- define "pegaKerberosTemplate" }}
{{ if (eq (include "performDeployment" .) "true") }}
{{- if .Values.global.kerberos }}
# Node type specific configuration for {{ .name }}
kind: ConfigMap
apiVersion: v1
metadata:
  name: {{ template "pegaImportKerberosConfigMap" $ }}
  namespace: {{ .Release.Namespace }}
data:
{{- $kerberos_value := .Values.global.kerberos }}
{{ $kerberos_value | toYaml | nindent 2 -}}
{{- end }}
{{ end }}
---
{{- end}}