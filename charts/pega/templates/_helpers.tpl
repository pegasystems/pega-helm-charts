{{- define "pegaEnvironmentConfig" }}pega-environment-config{{- end }}
{{- define "pegaVolumeConfig" }}pega-volume-config{{- end }}
{{- define "pegaVolumeCredentials" }}pega-volume-credentials{{- end }}
{{- define "pegaCredentialsSecret" }}pega-credentials-secret{{- end }}
{{- define "pegaRegistrySecret" }}pega-registry-secret{{- end }}
{{- define "deployConfig" -}}deploy-config{{- end -}}

{{- define "imagePullSecret" }}
{{- printf "{\"auths\": {\"%s\": {\"auth\": \"%s\"}}}" .Values.global.docker.registry.url (printf "%s:%s" .Values.global.docker.registry.username .Values.global.docker.registry.password | b64enc) | b64enc }}
{{- end }}

{{- define "performOnlyDeployment" }}
  {{- if (eq .Values.global.actions.execute "deploy") -}}
    true
  {{- else -}}
    false
  {{- end -}}
{{- end }}

{{- define "performDeployment" }}
  {{- if or (eq .Values.global.actions.execute "deploy") (eq .Values.global.actions.execute "install-deploy") (eq .Values.global.actions.execute "upgrade-deploy") -}}
    true
  {{- else -}}
    false
  {{- end -}}
{{- end }}

{{- define "performInstallAndDeployment" }}
  {{- if (eq .Values.global.actions.execute "install-deploy") -}}
    true
  {{- else -}}
    false
  {{- end -}}
{{- end }}

{{- define "performUpgradeAndDeployment" }}
  {{- if (eq .Values.global.actions.execute "upgrade-deploy") -}}
    true
  {{- else -}}
    false
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
  image: busybox:1.31.0
  # Init container for waiting for Elastic Search to initialize.  The URL should point at your Elastic Search instance.
  command: ['sh', '-c', 'until $(wget -q -S --spider --timeout=2 -O /dev/null {{ .Values.pegasearch.externalURL }}); do echo Waiting for search to become live...; sleep 10; done;']
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
# Initial JVM heap size, equivalent to -Xms
- name: INITIAL_HEAP
{{- if .node.initialHeap }}
  value: "{{ .node.initialHeap }}"
{{- else }}
  value: "4096m"
{{- end }}
# Maximum JVM heap size, equivalent to -Xmx
- name: MAX_HEAP
{{- if .node.maxHeap }}
  value: "{{ .node.maxHeap }}"
{{- else }}
  value: "7168m"
{{- end }}
{{- end -}}

{{- define "pega.health.probes" -}}
# LivenessProbe: indicates whether the container is live, i.e. running.
# If the LivenessProbe fails, the kubelet will kill the container and
# the container will be subjected to its RestartPolicy.
# The default state of Liveness before the initial delay is Success
livenessProbe:
  httpGet:
    # Path that is pinged to check for liveness.
    path: "/prweb/PRRestService/monitor/pingService/ping"
    port: 8080
    scheme: HTTP
  # Number of seconds after the container has started before liveness or readiness probes are initiated.
  initialDelaySeconds: 300
  # Number of seconds after which the probe times out. Defaults to 1 second. Minimum value is 1.
  timeoutSeconds: 20
  # How often (in seconds) to perform the probe. Default to 10 seconds. Minimum value is 1.
  periodSeconds: 10
  # Minimum consecutive successes for the probe to be considered successful after having failed. Defaults to 1.
  # Must be 1 for liveness. Minimum value is 1.
  successThreshold: 1
  # When a Pod starts and the probe fails, Kubernetes will try failureThreshold times before giving up.
  # Giving up in case of liveness probe means restarting the Pod. In case of readiness probe the
  # Pod will be marked Unready. Defaults to 3. Minimum value is 1.
  failureThreshold: 3
# ReadinessProbe: indicates whether the container is ready to service requests.
# If the ReadinessProbe fails, the endpoints controller will remove the
# pod's IP address from the endpoints of all services that match the pod.
# The default state of Readiness before the initial delay is Failure.
readinessProbe:
  httpGet:
    # Path that is pinged to check for readiness.
    path: "/prweb/PRRestService/monitor/pingService/ping"
    port: 8080
    scheme: HTTP
  # Number of seconds after the container has started before liveness or readiness probes are initiated.
  initialDelaySeconds: 300
  # Number of seconds after which the probe times out. Defaults to 1 second. Minimum value is 1.
  timeoutSeconds: 20
  # How often (in seconds) to perform the probe. Default to 10 seconds. Minimum value is 1.
  periodSeconds: 10
  # Minimum consecutive successes for the probe to be considered successful after having failed. Defaults to 1.
  # Must be 1 for liveness. Minimum value is 1.
  successThreshold: 1
  # When a Pod starts and the probe fails, Kubernetes will try failureThreshold times before giving up.
  # Giving up in case of liveness probe means restarting the Pod. In case of readiness probe the
  # Pod will be marked Unready. Defaults to 3. Minimum value is 1.
  failureThreshold: 3
{{- end }}

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
