{{- define "pegaCertificatesConfigTemplate" }}
# Secret used for common configuration between Pega nodes
{{ if (eq (include "performDeployment" .) "true") }}
{{- if and (.Values.global.certificates) (not (.Values.global.certificatesSecret)) }}
kind: Secret
apiVersion: v1
metadata:
  name: {{ template "pegaImportCertificatesConfig" $ }}
  namespace: {{ .Release.Namespace }}
stringData:
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