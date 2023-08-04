{{- define "pega-cassandra-secret-name" }}
{{- $depName := printf "%s" (include "deploymentName" $) -}}
{{- $depName -}}-cassandra-secret
{{- end -}}

{{- define "pega-stream-secret-name" }}
{{- $depName := printf "%s" (include "deploymentName" $) -}}
{{- $depName -}}-stream-secret
{{- end -}}

{{- define "pega-custom-artifactory-secret-name" }}
{{- $depName := printf "%s" (include "deploymentName" $) -}}
{{- $depName -}}-artifactory-secret
{{- end -}}

{{- define "pega-diagnostic-secret-name" }}
{{- $depName := printf "%s" (include "deploymentName" $) -}}
{{- $depName -}}-diagnostic-secret
{{- end -}}