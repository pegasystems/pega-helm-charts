{{- define "pega.config" -}}
# Node type specific configuration for {{ .name }}
kind: ConfigMap
apiVersion: v1
metadata:
  name: {{ .name }}
  namespace: {{ .root.Release.Namespace }}
data:
    # prconfig file to be used by {{ .name }}
  prconfig.xml: |-
{{ .root.Files.Get .node.prconfigPath | indent 6 }}
    # prlog4j2 file to be used by {{ .name }}
  prlog4j2.xml: |-
{{ .root.Files.Get .node.prlog4j2Path | indent 6 }}
{{- end }}
