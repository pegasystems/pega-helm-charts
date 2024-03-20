{{- define "pegaEnvironmentConfig" }}
{{- $depName := printf "%s" (include "deploymentName" $) -}}
{{- $depName -}}-environment-config
{{- end }}

{{- define "pegaImportCertificatesSecret" }}
{{- $depName := printf "%s" (include "deploymentName" $) -}}
{{- $depName -}}-import-certificates-secret
{{- end }}

{{- define "pegaImportKerberosConfigMap" }}
{{- $depName := printf "%s" (include "deploymentName" $) -}}
{{- $depName -}}-import-kerberos-configmap
{{- end }}

{{- define "pegaVolumeImportCertificates" }}pega-volume-import-certificates{{- end }}

{{- define "pegaImportCertificatesTemplate" }}
- name: {{ template "pegaVolumeImportCertificates" }}
  projected:
    defaultMode: 420
    sources:
  {{ if (.Values.global.certificatesSecrets) }}
  {{- range .Values.global.certificatesSecrets }}
    - secret:
        name: {{ . }}
  {{- end }}
  {{ else }}
    # This name will be referred in the volume mounts kind.
    - secret:
        name: {{ template "pegaImportCertificatesSecret" $ }}
  {{ end }}
{{- end}}

{{- define "pegaCustomArtifactoryCertificateConfig" }}
{{- $depName := printf "%s" (include "deploymentName" $) -}}
{{- $depName -}}-custom-artifactory-certificate-config
{{- end }}

{{- define "pegaVolumeCustomArtifactoryCertificate" }}pega-volume-custom-artifactory-certificate{{- end }}

{{- define "pegaCustomArtifactoryCertificateTemplate" }}
- name: {{ template "pegaVolumeCustomArtifactoryCertificate" }}
  configMap:
    # This name will be referred in the volume mounts kind.
    name: {{ template "pegaCustomArtifactoryCertificateConfig" $ }}
    # Used to specify permissions on files within the volume.
    defaultMode: 420
{{- end}}


{{- define "pegaTomcatKeystoreSecret" }}
{{- $depName := printf "%s" (include "deploymentName" .root) -}}
{{- $depName -}}-tomcat-keystore-secret
{{- end }}

{{- define "pegaVolumeTomcatKeystore" }}pega-volume-tomcat-keystore{{- end }}

{{- define "pegaVolumeTomcatKeystoreTemplate" }}
- name: {{ template "pegaVolumeTomcatKeystore" }}
  projected:
    defaultMode: 420
    sources:
  {{ if (((.node.service).tls).external_secret_names) }}
  {{- range (((.node.service).tls).external_secret_names) }}
    - secret:
        name: {{ . }}
  {{- end }}
  {{ else }}
    # This name will be referred in the volume mounts kind.
    - secret:
        name: {{ template "pegaTomcatKeystoreSecret" $ }}
  {{ end }}
{{- end}}

{{- define "pegaVolumeConfig" }}pega-volume-config{{- end }}

{{- define "pegaKerberosConfig" }}pega-import-kerberos{{- end }}



{{- define "deployConfig" -}}deploy-config{{- end -}}
{{- define "pegaBackendConfig" -}}pega-backend-config{{- end -}}


{{- define "performOnlyDeployment" }}
  {{- if (eq .Values.global.actions.execute "deploy") -}}
    true
  {{- else -}}
    false
  {{- end -}}
{{- end }}

{{- define "useBasicAuthForCustomArtifactory" }}
  {{- if (.Values.global.customArtifactory) }}
    {{- if (.Values.global.customArtifactory.authentication) }}
      {{- if (.Values.global.customArtifactory.authentication.basic) }}
        {{- if and (.Values.global.customArtifactory.authentication.basic.username) (.Values.global.customArtifactory.authentication.basic.password) -}}
          true
        {{- else -}}
          false
        {{- end -}}
      {{- end -}}
    {{- end }}
  {{- end }}
{{- end }}

{{- define "imagePullSecret" }}
{{- if .Values.global.docker.registry }}
{{- printf "{\"auths\": {\"%s\": {\"auth\": \"%s\"}}}" .Values.global.docker.registry.url (printf "%s:%s" .Values.global.docker.registry.username .Values.global.docker.registry.password | b64enc) | b64enc }}
{{- end }}
{{- end }}

{{- define "useApiKeyForCustomArtifactory" }}
  {{- if (.Values.global.customArtifactory) }}
    {{- if (.Values.global.customArtifactory.authentication) }}
      {{- if (.Values.global.customArtifactory.authentication.apiKey) }}
        {{- if and (.Values.global.customArtifactory.authentication.apiKey.headerName) (.Values.global.customArtifactory.authentication.apiKey.value) -}}
          true
        {{- else -}}
          false
        {{- end -}}
      {{- end }}
    {{- end }}
  {{- end -}}
{{- end }}

{{- define "tlssecretsnippet" }}
  tls:
   - hosts:
     - {{ template "domainName" dict "node" .node }}
     secretName: {{ .node.ingress.tls.secretName }}
{{- end }}

{{- define "hostPathType" }}
 {{- if .node.ingress.pathType -}}
   {{ .node.ingress.pathType }}
 {{- else -}}
   ImplementationSpecific
 {{- end }}
{{- end }}

{{- define "defaultIngressRule" }}
- pathType: {{ include "hostPathType" $ }}
  {{- if .node.ingress.path }}
  path: {{ .node.ingress.path }}
  {{- end }}
  backend:
    {{ include "ingressBackend" $ }}
{{- end }}

{{- define "customerDeploymentID" -}}
  {{- if .Values.global.customerDeploymentId -}}
    {{ .Values.global.customerDeploymentId}}
  {{- else -}}
    {{ .Release.Namespace }}
  {{- end -}}
{{- end }}

# list of either external or internal cassandra nodes
{{- define "cassandraNodes" }}
  {{- if .Values.dds.externalNodes -}}
    {{ .Values.dds.externalNodes }}
  {{- else -}}
    {{ template "getCassandraSubchartService" . }}
  {{- end -}}
{{- end }}

# whether or not cassandra is enabled at all (internally or externally)
{{- define "cassandraEnabled" }}
  {{- if .Values.dds.externalNodes -}}
    true
  {{- else -}}
    {{- if .Values.cassandra.enabled -}}
      true
    {{- else -}}
      false
    {{- end -}}
  {{- end -}}
{{- end }}

# whether we should create internal cassandra nodes
{{- define "internalCassandraEnabled" }}
  {{- if .Values.dds.externalNodes -}}
    false
  {{- else -}}
    {{- if .Values.cassandra.enabled -}}
      true
    {{- else -}}
      false
    {{- end -}}
  {{- end -}}
{{- end }}

{{- define "waitForPegaSearch" -}}
- name: wait-for-pegasearch
  image: {{ .Values.global.utilityImages.busybox.image }}
  imagePullPolicy: {{ .Values.global.utilityImages.busybox.imagePullPolicy }}
  # Init container for waiting for Elastic Search to initialize.  The URL should point at your Elastic Search instance.
  command: ['sh', '-c', 'until $(wget -q -S --spider --timeout=2 -O /dev/null {{ include "pegaSearchURL" $ }}); do echo Waiting for search to become live...; sleep 10; done;']
{{- include "initContainerResources" $ }}
{{- end }}

{{- define "waitForCassandra" -}}
  {{- if  eq (include "internalCassandraEnabled" .) "true" -}}
- name: wait-for-cassandra
  image: {{ .Values.cassandra.image.repo }}:{{ .Values.cassandra.image.tag}}
  # Init container for waiting for Cassndra to initialize.  For each node, a copy of the until loop should be made to check each node.
  # -u is username
  # -p is password
  # final 2 args for cqlsh are cassandra host and port respectively
  command: ['sh', '-c', '{{- template "waitForCassandraScript" dict "nodes" (include "getCassandraSubchartService" .) "node" .Values.dds -}}']
{{- include "initContainerResources" $ }}
 {{- end -}}
{{- end }}

{{- define "getCassandraSubchartService" -}}
  {{- if  eq (include "internalCassandraEnabled" .) "true" -}}
    {{- template "cassandra.fullname" dict "Values" .Values.cassandra "Release" .Release "Chart" (dict "Name" "cassandra") -}}
  {{- end -}}
{{- end -}}

{{- define "waitForCassandraScript" -}}
  {{- $cassandraPort := .node.port -}}
  {{- $cassandraUser := .node.username -}}
  {{- $cassandraPassword := .node.password -}}
  {{- range $i, $val := splitList "," .nodes -}}
until cqlsh -u {{ $cassandraUser | quote }} -p {{ $cassandraPassword | quote }} -e "describe cluster" {{ $val | trim }} {{ $cassandraPort }} ; do echo Waiting for cassandra to become live...; sleep 10; done;
  {{- end -}}
{{- end -}}

{{- define "pega.jvmconfig" -}}
# Additional JVM arguments
- name: JAVA_OPTS
  value: "{{ .node.javaOpts }}"
# Additional CATALINA arguments
- name: CATALINA_OPTS
{{- if .node.catalinaOpts }}
  value: "{{ .node.catalinaOpts }}"
{{- else }}
  value: ""
{{- end }}
# Initial JVM heap size, equivalent to -Xms
- name: INITIAL_HEAP
{{- if .node.initialHeap }}
  value: "{{ .node.initialHeap }}"
{{- else }}
  value: "8192m"
{{- end }}
# Maximum JVM heap size, equivalent to -Xmx
- name: MAX_HEAP
{{- if .node.maxHeap }}
  value: "{{ .node.maxHeap }}"
{{- else }}
  value: "8192m"
{{- end }}
{{- end -}}

# Evaluate background node types based on cassandra enabled or not(internally or externally)
{{- define "evaluateBackgroundNodeTypes" }}
  {{- if  eq (include "cassandraEnabled" .) "true" -}}
    BackgroundProcessing,Search,Batch,RealTime,Custom1,Custom2,Custom3,Custom4,Custom5,BIX,ADM,RTDG  
  {{- else -}}
    Background
  {{- end -}}
{{- end }}

# Load balancer session cookie stickiness time in seconds,
# calculated as sum of .requestor.passivationTimeSec and passivation delay.
{{- define "lbSessionCookieStickiness" }}
  {{- $passivationTime := 3600 -}}
  {{- $passivationDelay := 120 -}}

  {{- if .node.requestor -}}
    {{- if .node.requestor.passivationTimeSec -}}
      {{- $passivationTime = .node.requestor.passivationTimeSec -}}
    {{- end -}}
  {{- else if .node.service.alb_stickiness_lb_cookie_duration_seconds -}}
    {{- $passivationTime = .node.service.alb_stickiness_lb_cookie_duration_seconds -}}
  {{- end -}}

  {{- add $passivationTime $passivationDelay -}}
{{- end -}}

# Determine application root context to use in pega tomcat nodes
{{- define "pega.applicationContextPath" -}}
   {{- if .node.ingress -}}
      {{- if .node.ingress.appContextPath -}}
         {{ trimAll "/" .node.ingress.appContextPath }}
      {{- else -}}
         prweb
      {{- end -}}	 
   {{- else -}}
      prweb
   {{- end -}}
{{- end }}

{{- define "gkemanagedcertificate" }}
{{- if (semverCompare ">= 1.19.0-0" (trimPrefix "v" .root.Capabilities.KubeVersion.GitVersion)) }}
apiVersion: networking.gke.io/v1
{{- else }}
apiVersion: networking.gke.io/v1beta1
{{- end }}
kind: ManagedCertificate
metadata:
  name: {{ .name }}
spec:
  domains:
    - {{ .domain }}
---
{{- end -}}

{{- define "domainName" }}
  {{- if .node.ingress -}}
  {{- if .node.ingress.domain -}}
    {{ .node.ingress.domain }}
  {{- end -}}
  {{- else if .node.service.domain -}}
    {{ .node.service.domain }}
  {{- end -}}
{{- end }}

{{- define "eksHttpsAnnotationSnippet" }}
    # specifies the ports that ALB used to listen on
    alb.ingress.kubernetes.io/listen-ports: '[{"HTTP": 80}, {"HTTPS": 443}]'
    # set the redirect action to redirect http traffic into https
    alb.ingress.kubernetes.io/actions.ssl-redirect: '{"Type": "redirect", "RedirectConfig": { "Protocol": "HTTPS", "Port": "443", "StatusCode": "HTTP_301"}}'
{{- end }}

{{- define "ingressTlsEnabled" }}
{{- if (.node.ingress) }}
{{- if (.node.ingress.tls) }}
{{- if (eq .node.ingress.tls.enabled true) }}
true
{{- end }}
{{- end }}
{{- end }}
{{- end }}

#Override this template to generate additional pod annotations that are dynamically composed during helm deployment (do not indent annotations)
{{- define "generatedPodAnnotations" }}
{{- end }}

#Override this template to generate additional pod labels that are dynamically composed during helm deployment (do not indent labels)
{{- define "generatedPodLabels" }}
{{- end }}

#Kerberos config map
{{- define "pegaKerberosVolumeTemplate" }}
# Volume used to mount config files.
- name: {{ template "pegaKerberosConfig" }}-config
  configMap:
    # This name will be referred in the volume mounts kind.
    name: {{ template "pegaImportKerberosConfigMap" $ }}
    # Used to specify permissions on files within the volume.
    defaultMode: 420
{{- end}}

{{- define "generatedDNSConfigAnnotations" }}
{{ if (.Values.global.privateHostedZoneDomainName) }}
dnsConfig:
  searches:
  - {{ .Values.global.privateHostedZoneDomainName }}
{{- end }}
{{- end }}

{{- define "pegaSearchURL" -}}
{{- $d1 := dict "overrideURL" $.Values.pegasearch.externalURL }}
{{- $d2 := merge $ $d1 }}
{{- template "searchURL" $d2 }}
{{- end -}}

{{- define "srsAuthPrivateKey" -}}
{{- if and (.Values.pegasearch.externalSearchService) ((.Values.pegasearch.srsAuth).enabled) }}
    {{- if (.Values.pegasearch.srsAuth).privateKey }}
        {{- .Values.pegasearch.srsAuth.privateKey | b64enc }}
    {{- else }}
        {{- fail "A valid entry is required for pegasearch.srsAuth.privateKey or pegasearch.srsAuth.external_secret_name, when request authentication mechanism(IDP) is enabled between SRS and Pega Infinity i.e. pegasearch.srsAuth.enabled is true." | quote }}
    {{- end }}
{{- end }}
{{- end }}

{{- define "srsAuthEnvSecretFrom" }}
{{- if .Values.pegasearch.srsAuth.external_secret_name }}
name: {{ .Values.pegasearch.srsAuth.external_secret_name }}
key: SRS_OAUTH_PRIVATE_KEY
{{- else }}
name: pega-srs-auth-secret
key: privateKey
{{- end }}
{{- end }}

{{- define "ingressApiVersion" }}
{{- if (semverCompare ">= 1.19.0-0" (trimPrefix "v" .root.Capabilities.KubeVersion.GitVersion)) }}
apiVersion: networking.k8s.io/v1
{{- else }}
apiVersion: extensions/v1beta1
{{- end }}
{{- end }}

{{- define "ingressService" }}
{{- if (semverCompare ">= 1.19.0-0" (trimPrefix "v" .root.Capabilities.KubeVersion.GitVersion)) }}
service:
  name: {{ .name }}
  port: 
    number: {{ .node.service.port }}
{{- else }}
serviceName: {{ .name }}
servicePort: {{ .node.service.port }}
{{- end }}
{{- end }}

{{- define "ingressServiceHttps" }}
{{- if (semverCompare ">= 1.19.0-0" (trimPrefix "v" .root.Capabilities.KubeVersion.GitVersion)) }}
service:
  name: {{ .name }}
  port:
    number: {{ .node.service.tls.port }}
{{- else }}
serviceName: {{ .name }}
servicePort: {{ .node.service.tls.port }}
{{- end }}
{{- end }}

{{- define "ingressBackend" }}
{{- if ((.node.service).tls).enabled }}
    {{ include "ingressServiceHttps" . | indent 10 }}
{{- else }}
    {{ include "ingressService" . | indent 10 }}
{{- end }}
{{- end }}

{{- define "ingressServiceC11n" }}
{{- if (semverCompare ">= 1.19.0-0" (trimPrefix "v" .root.Capabilities.KubeVersion.GitVersion)) }}
service:
  name: constellation
  port: 
    number: 3000
{{- else }}
serviceName: constellation
servicePort: 3000
{{- end }}
{{- end }}

{{- define "ingressServiceSSLRedirect" }}
{{- if (semverCompare ">= 1.19.0-0" (trimPrefix "v" .root.Capabilities.KubeVersion.GitVersion)) }}
service:
  name: ssl-redirect
  port: 
    name: use-annotation
{{- else }}
serviceName: ssl-redirect
servicePort: use-annotation
{{- end }}
{{- end }}

{{- define "tierClassloaderRetryTimeout" }}
{{- if gt (add .periodSeconds 0) 180 -}}
180
{{- else -}}
{{- add .periodSeconds 0}}
{{- end -}}
{{- end -}}

{{- define "tierClassloaderMaxRetries" }}
{{- if gt (add .periodSeconds 0) 180 -}}
{{- add (round (div (mul .periodSeconds .failureThreshold) 180) 0) 1 -}}
{{- else -}}
{{- add .failureThreshold 1 -}}
{{- end -}}
{{- end -}}

{{- define "hzServiceName" -}}
  {{- if and (not .Values.hazelcast.enabled)  .Values.hazelcast.clusteringServiceEnabled -}}
    {{ template "clusteringServiceName" }}
  {{- else -}}
    {{ template "hazelcastName" }}
  {{- end -}}
{{- end -}}

{{- define "hzClusterName" -}}
  {{- if and (not .Values.hazelcast.enabled)  .Values.hazelcast.clusteringServiceEnabled -}}
    {{ .Values.hazelcast.server.clustering_service_group_name }}
  {{- else -}}
    {{ .Values.hazelcast.client.clusterName }}
  {{- end -}}
{{- end -}}

{{- define "hazelcastCSConfigRequired" }}
  {{- if and (or (.Values.hazelcast.enabled) (.Values.hazelcast.clusteringServiceEnabled)) (not (.Values.hazelcast.migration.embeddedToCSMigration)) -}}
    true
  {{- else -}}
    false
  {{- end -}}
{{- end -}}

{{- define "hazelcastVersion" }}
  {{- if and (not .Values.hazelcast.enabled)  .Values.hazelcast.clusteringServiceEnabled -}}
    v5
  {{- else -}}
    v4
  {{- end -}}
{{- end -}}

{{- define "pegaCredentialVolumeTemplate" }}
- name: {{ template "pegaVolumeCredentials" }}
  projected:
    defaultMode: 420
    sources:
    {{- $dbDict := dict "deploySecret" "deployDBSecret" "deployNonExtsecret" "deployNonExtDBSecret" "extSecretName" .Values.global.jdbc.external_secret_name "nonExtSecretName" "pega-db-secret-name" "context" $  -}}
    {{ include "secretResolver" $dbDict | indent 4}}

    {{- $hzDict := dict "deploySecret" "deployHzSecret" "deployNonExtsecret" "deployNonExtHzSecret" "extSecretName" .Values.hazelcast.external_secret_name "nonExtSecretName" "pega-hz-secret-name" "context" $ -}}
    {{ include "secretResolver" $hzDict | indent 4}}

    {{- $streamDict := dict "deploySecret" "deployStreamSecret" "deployNonExtsecret" "deployNonExtStreamSecret" "extSecretName" .Values.stream.external_secret_name "nonExtSecretName" "pega-stream-secret-name" "context" $ -}}
    {{ include "secretResolver" $streamDict | indent 4}}

    {{- $ddsDict := dict "deploySecret" "deployDDSSecret" "deployNonExtsecret" "deployNonExtDDSSecret" "extSecretName" .Values.dds.external_secret_name "nonExtSecretName" "pega-dds-secret-name" "context" $ -}}
    {{ include "secretResolver" $ddsDict | indent 4}}

    {{- $artifactoryDict := dict "deploySecret" "deployArtifactorySecret" "deployNonExtsecret" "deployNonExtArtifactorySecret" "extSecretName" .Values.global.customArtifactory.authentication.external_secret_name "nonExtSecretName" "pega-custom-artifactory-secret-name" "context" $ -}}
    {{ include "secretResolver" $artifactoryDict | indent 4}}

    - secret:
        name: {{ include "pega-diagnostic-secret-name" $}}

{{- end}}
