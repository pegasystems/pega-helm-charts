{{- define "searchDeploymentName" }}{{ $deploymentNamePrefix := "pega" }}{{ if (.Values.global.deployment) }}{{ if (.Values.global.deployment.name) }}{{ $deploymentNamePrefix = .Values.global.deployment.name }}{{ end }}{{ end }}{{ $deploymentNamePrefix }}{{- end }}

{{- define "searchName" -}}{{ $depName := printf "%s" (include "searchDeploymentName" $) }}{{- printf "%s" $depName -}}-search{{- end -}}

{{- define "defaultSearchURL" }}http://{{ template "searchName" $}}{{- end }}

{{- define "searchURL" -}}
{{ $override := get $ "overrideURL" -}}
{{- $searchURLOverride := $.Values.externalURL -}}
{{- if not (empty $override) -}}
{{- $searchURLOverride = $override -}}
{{- end -}}
{{- $defaultSearchURL := printf "%s" (include "defaultSearchURL" $) -}}
{{- $deploymentName := printf "%s" (include "searchDeploymentName" $) -}}
{{- if not (empty $searchURLOverride) -}}
{{- if and (ne $deploymentName "pega") (eq $searchURLOverride "http://pega-search") -}}
{{- $defaultSearchURL -}}
{{- else -}}
{{- $searchURLOverride -}}
{{- end -}}
{{- else -}}
{{- $defaultSearchURL -}}
{{- end -}}
{{- end -}}

{{- /* From the perspective of the Pega Infinity nodes, Pega deployments always assume that the search instances are external to the deployment when you use these charts; 
       however this check determines if search is deployed externally as defined in the backing search service configuration or if it is defined in the Pega-provided 
       Docker search image.*/ -}}
{{- define "isExternalSearch" -}}
  {{- $defaultSearchURL := printf "%s" (include "defaultSearchURL" $) -}}
  {{- $searchURL := printf "%s" (include "searchURL" $) -}}
  {{- if or (ne $searchURL $defaultSearchURL) $.Values.externalSearchService -}}
    true
  {{- else -}}
    false
  {{- end -}}
{{- end -}}


{{- define "performDeployment" }}
  {{- if or (eq .Values.global.actions.execute "deploy") (eq .Values.global.actions.execute "install-deploy") (eq .Values.global.actions.execute "upgrade-deploy") -}}
    true
  {{- else -}}
    false
  {{- end -}}
{{- end }}
