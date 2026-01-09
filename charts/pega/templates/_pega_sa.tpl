{{- define "pega.serviceaccount" -}}
{{- if .sa.create -}}
apiVersion: v1
kind: ServiceAccount
metadata:
  name: {{ .sa.name | quote }}
  namespace: {{ .root.Release.Namespace }}
---
{{- end -}}
{{- end -}}