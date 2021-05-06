{{- define "pega.config" -}}
{{ template "pega.config.inner" dict "root" .root "node" .dep "name" .name "mode" (include "deployConfig" .root) "custom" .dep.custom }}
{{- end -}}


{{- define "pega.config.inner" -}}
{{- $arg := .mode -}}
# Node type specific configuration for {{ .name }}
kind: ConfigMap
apiVersion: v1
metadata:
  name: {{ .name }}
  namespace: {{ .root.Release.Namespace }}
data:
{{- if eq $arg "deploy-config" }}
{{- $custom_config := .custom }}
{{- $global_config := .root.Values.global.config }}

  # Pega deployment prconfig.xml file
  prconfig.xml: |-
{{- if $global_config.prconfig }}
{{ $global_config.prconfig | indent 6 }}
{{- else if $custom_config.prconfig }}
{{ $custom_config.prconfig | indent 6 }}
{{ else }}
{{ .root.Files.Get "config/deploy/prconfig.xml" | indent 6 }}
{{- end }}

  # Pega deployment prlog4j2.xml file
  prlog4j2.xml: |-
{{- if $global_config.prlog4j2 }}
{{ $global_config.prlog4j2 | indent 6 }}
{{- else if $custom_config.prlog4j2 }}
{{ $custom_config.prlog4j2 | indent 6 }}
{{ else }}
{{ .root.Files.Get "config/deploy/prlog4j2.xml" | indent 6 }}
{{- end }}

  # Pega deployment prweb.xml file
  prweb.xml: |-
{{- if $global_config.prwebXML }}
{{ $global_config.prwebXML | indent 6 }}
{{- else if $custom_config.prwebXML }}
{{ $custom_config.prwebXML | indent 6 }}
{{- end }}

  # Pega deployment contextXML template file
  context.xml.tmpl: |-
{{- if $global_config.contextXML }}
{{ $global_config.contextXML | indent 6 }}
{{- else if $custom_config.contextXML }}
{{ $custom_config.contextXML | indent 6 }}
{{ else }}
{{ .root.Files.Get "config/deploy/context.xml.tmpl" | indent 6 }}
{{- end }}

  # Pega deployment server.xml file
  server.xml: |-
{{- if $global_config.serverXML }}
{{ $global_config.serverXML | indent 6 }}
{{- else if $custom_config.serverXML }}
{{ $custom_config.serverXML | indent 6 }}
{{ else }}
{{ .root.Files.Get "config/deploy/server.xml" | indent 6 }}
{{- end }}

 # Pega deployment web.xml file
  web.xml: |-
{{- if $global_config.webXML }}
{{ $global_config.webXML | indent 6 }}
{{- else if $custom_config.webXML }}
{{ $custom_config.webXML | indent 6 }}
{{ else }}
{{- $web := .root.Files.Glob "config/deploy/web.xml" }}
{{ if $web }}
{{ .root.Files.Get "config/deploy/web.xml" | indent 6 }}
{{- end }}
{{- end }}
{{- end }}
---
{{- end }}
