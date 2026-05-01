{{- define "pega.serviceaccount" -}}
apiVersion: v1
kind: ServiceAccount
metadata:
  name: {{ .name | quote }}
  namespace: {{ .root.Release.Namespace }}
---
{{- end -}}