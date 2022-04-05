{{- define "pegaCustomArtifactoryCertificatesConfigTemplate" }}
# Config map used for certificate of custom artifactory
{{ if (eq (include "customArtifactorySSLVerificationEnabled" .) "true") }}
{{- if .Values.global.customArtifactory.certificate }}
kind: ConfigMap
apiVersion: v1
metadata:
  name: {{ template "pegaCustomArtifactoryCertificateConfig" $ }}
  namespace: {{ .Release.Namespace }}
data:
  # cert File
{{- if .Values.global.customArtifactory.certificate }}
  # import certificate
  {{ .Values.global.customArtifactory.certificate | toYaml | nindent 2 -}}
{{- end }}
{{- end }}
{{- end }}
---
{{- end }}
