{{- define "pegaCertificatesConfigTemplate" }}
# Config map used for common configuration between Pega nodes
{{ if (eq (include "performDeployment" .) "true") }}
{{- if .Values.global.certificates }}
kind: ConfigMap
apiVersion: v1
metadata:
  name: {{ template "pegaImportCertificatesConfig" $ }}
  namespace: {{ .Release.Namespace }}
data:
  # cert Files
{{- if .Values.global.certificates }}
  # import certificates from values
{{- range $k, $v := .Values.global }}
  {{- if eq $k "certificates" }}
  {{ $v | toYaml | nindent 2 -}}
  {{- end }}
{{- end }}
{{- end }}
{{- end }}
{{ end }}
{{- end}}