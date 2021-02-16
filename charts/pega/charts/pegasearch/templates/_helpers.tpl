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

{{- /* From the perspective of the Pega Infinity nodes, the search instances are always external when deployed via these charts,
       but this check determines if search is being deployed externally from the perspective of the charts as a whole, such as 
       when the backing search service is being used.*/ -}}
{{- define "isExternalSearch" -}}
  {{- $defaultSearchURL := printf "%s" (include "defaultSearchURL" $) -}}
  {{- $searchURL := printf "%s" (include "searchURL" $) -}}
  {{- if or (ne $searchURL $defaultSearchURL) $.Values.externalSearchService -}}
    true
  {{- else -}}
    false
  {{- end -}}
{{- end -}}