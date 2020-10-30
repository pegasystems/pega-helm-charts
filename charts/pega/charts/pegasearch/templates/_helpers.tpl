{{- define "searchName" -}}
pega-search
{{- end -}}
{{- /* From the perspective of the Pega Infinity nodes, the search instances are always external when deployed via these charts,
       but this check determines if search is being deployed externally from the perspective of the charts as a whole, such as 
       when the backing search service is being used.*/ -}}
{{- define "isExternalSearch" }}
  {{- if or (and (.Values.externalURL) (ne .Values.externalURL "http://pega-search")) .Values.externalSearchService -}}
    true
  {{- else -}}
    false
  {{- end -}}
{{- end }} 