{{- define "pega-dds-secret-name" }}
{{- $depName := printf "%s" (include "deploymentName" $) -}}
{{- $depName -}}-dds-secret
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

{{- define "deployStreamSecret" }}
{{- if and (eq (include "performDeployment" .) "true") ((.Values.stream).enabled) -}}
true
{{- else -}}
false
{{- end -}}
{{- end }}

{{- define "deployNonExtStreamSecret" }}
{{- if and (eq (include "deployStreamSecret" .) "true") (not (.Values.stream).external_secret_name) -}}
true
{{- else -}}
false
{{- end -}}
{{- end -}}

{{- define "deployHzSecret" }}
{{- if (eq (include "hazelcastCSConfigRequired" .) "true") -}}
true
{{- else -}}
false
{{- end -}}
{{- end }}

{{- define "deployNonExtHzSecret" }}
{{- if and (eq (include "deployHzSecret" .) "true") (not (.Values.hazelcast).external_secret_name) -}}
true
{{- else -}}
false
{{- end -}}
{{- end -}}

{{- define "deployDDSSecret" }}
{{- if  and (eq (include "performDeployment" .) "true") (eq (include "cassandraEnabled" .) "true") (eq (include "internalCassandraEnabled" .) "false") -}}
true
{{- else -}}
false
{{- end -}}
{{- end }}

{{- define "deployNonExtDDSSecret" }}
{{- if and (eq (include "deployDDSSecret" .) "true") (not (.Values.dds).external_secret_name) -}}
true
{{- else -}}
false
{{- end -}}
{{- end -}}


{{- define "deployArtifactorySecret" }}
{{- if or (eq (include "useBasicAuthForCustomArtifactory" .) "true") (eq (include "useApiKeyForCustomArtifactory" .) "true") (.Values.global.customArtifactory.authentication.external_secret_name) -}}
true
{{- else -}}
false
{{- end -}}
{{- end }}

{{- define "deployNonExtArtifactorySecret" }}
{{- if and (eq (include "deployArtifactorySecret" .) "true") (not (.Values.global.customArtifactory.authentication).external_secret_name) -}}
true
{{- else -}}
false
{{- end -}}
{{- end -}}