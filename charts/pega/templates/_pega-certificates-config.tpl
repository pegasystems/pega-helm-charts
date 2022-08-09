{{- define "pegaCertificatesConfigTemplate" }}
# Config map used for common configuration between Pega nodes
{{ if (eq (include "performDeployment" .) "true") }}
{{- if .Values.global.certificates }}
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
  {{- if and (eq $k "certificates") (not (hasKey $v "external_secret_name")) }}
  {{ $v | toYaml | nindent 2 -}}
  {{- end }}
{{- end }}
{{- end }}
{{- end }}
{{ end }}
{{- end}}