{{- define "pega.config" -}}
{{- $arg := .mode -}}
# Node type specific configuration for {{ .name }}
kind: ConfigMap
apiVersion: v1
metadata:
  name: {{ .name }}
  namespace: {{ .root.Release.Namespace }}
data:

# Start of Pega Deployment Configuration

{{ if eq $arg "deploy-config" }}

{{- $prconfigPath := "config/deploy/prconfig.xml" }}
{{- $contextXMLTemplate := "config/deploy/context.xml.tmpl" }}
{{- $prlog4j2Path := "config/deploy/prlog4j2.xml" }}

{{- if .custom }}
{{- if .custom.prconfig }}
 # CUSTOM prconfig file to be used by {{ .name }}
  prconfig.xml: |-
{{ .custom.prconfig | indent 6 }}
{{ else if $prconfig := .root.Files.Glob $prconfigPath }}
 # prconfig file to be used by {{ .name }}
  prconfig.xml: |-
{{ .root.Files.Get $prconfigPath | indent 6 }}
{{- end }}
{{ else if $prconfig := .root.Files.Glob $prconfigPath }}
 # prconfig file to be used by {{ .name }}
  prconfig.xml: |-
{{ .root.Files.Get $prconfigPath | indent 6 }}
{{- end }}

{{ if $contextXML := .root.Files.Glob $contextXMLTemplate }}
  # contextXMLTemplate to be used by {{ .name }}
  context.xml.tmpl: |-
{{ .root.Files.Get $contextXMLTemplate | indent 6 }}
{{- end }}

{{- if .custom }}
{{- if .custom.context }}
 # CUSTOM context file to be used by {{ .name }}
  context.xml: |-
{{ .custom.context | indent 6 }}
{{- end }}
{{- end }}

  # prlog4j2 file to be used by {{ .name }}
  prlog4j2.xml: |-
{{ .root.Files.Get $prlog4j2Path | indent 6 }}

{{- end }}
# End of Pega Deployment Configuration
---
{{- end }}

