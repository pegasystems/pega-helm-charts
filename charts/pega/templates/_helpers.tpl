{{- define "pegaEnvironmentConfig" }}pega-environment-config{{- end }}
{{- define "pegaVolumeConfig" }}pega-volume-config{{- end }}
{{- define "pegaVolumeCredentials" }}pega-volume-credentials{{- end }}
{{- define "pegaStorageClassEBS" }}pega-ebs-storage-class{{- end }}
{{- define "pegaDatabaseSecret" }}pega-database-secret{{- end }}
{{- define "pegaRegistrySecret" }}pega-registry-secret{{- end }}
{{- define "pegaWebName" -}}pega-web{{- end -}}
{{- define "pegaBatchName" -}}pega-batch{{- end -}}
{{- define "pegaStreamName" -}}pega-stream{{- end -}}
{{- define "searchName" -}}pega-search{{- end -}}

{{- define "imagePullSecret" }}
{{- printf "{\"auths\": {\"%s\": {\"auth\": \"%s\"}}}" .Values.docker.registry.url (printf "%s:%s" .Values.docker.registry.username .Values.docker.registry.password | b64enc) | b64enc }}
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

{{- define "properPegaSearchURL" }}
  {{- if .Values.search.externalURL -}}
    {{ .Values.search.externalURL }}
  {{- else -}}
    http://{{ template "searchName" . }}
  {{- end -}}
{{- end }}

{{- define "waitForPegaSearch" -}}
- name: wait-for-pegasearch
  image: busybox:1.27.2
  # Init container for waiting for Elastic Search to initialize.  The URL should point at your Elastic Search instance.
  command: ['sh', '-c', 'until $(wget -q -S --spider --timeout=2 -O /dev/null {{ include "properPegaSearchURL" . }}); do echo Waiting for search to become live...; sleep 10; done;']
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
  value: "{{ .node.initialHeap }}"
# Maximum JVM heap size, equivalent to -Xmx
- name: MAX_HEAP
  value: "{{ .node.maxHeap }}"
{{- end -}}

{{- define "commonEnvironmentVariables" -}}
- name: CASSANDRA_CLUSTER
  valueFrom:
    configMapKeyRef:
      name: {{ template "pegaEnvironmentConfig" }}
      key: CASSANDRA_CLUSTER
- name: CASSANDRA_NODES
  valueFrom:
    configMapKeyRef:
      name: {{ template "pegaEnvironmentConfig" }}
      key: CASSANDRA_NODES
- name: CASSANDRA_PORT
  valueFrom:
    configMapKeyRef:
      name: {{ template "pegaEnvironmentConfig" }}
      key: CASSANDRA_PORT
- name: CASSANDRA_USERNAME
  valueFrom:
    configMapKeyRef:
      name: {{ template "pegaEnvironmentConfig" }}
      key: CASSANDRA_USERNAME
- name: CASSANDRA_PASSWORD
  valueFrom:
    configMapKeyRef:
      name: {{ template "pegaEnvironmentConfig" }}
      key: CASSANDRA_PASSWORD
- name: PEGA_SEARCH_URL
  valueFrom:
    configMapKeyRef:
      name: {{ template "pegaEnvironmentConfig" }}
      key: PEGA_SEARCH_URL
- name: JDBC_URL
  valueFrom:
    configMapKeyRef:
      name: {{ template "pegaEnvironmentConfig" }}
      key: JDBC_URL
- name: JDBC_CLASS
  valueFrom:
    configMapKeyRef:
      name: {{ template "pegaEnvironmentConfig" }}
      key: JDBC_CLASS
- name: JDBC_DRIVER_URI
  valueFrom:
    configMapKeyRef:
      name: {{ template "pegaEnvironmentConfig" }}
      key: JDBC_DRIVER_URI
- name: RULES_SCHEMA
  valueFrom:
    configMapKeyRef:
      name: {{ template "pegaEnvironmentConfig" }}
      key: RULES_SCHEMA
- name: DATA_SCHEMA
  valueFrom:
    configMapKeyRef:
      name: {{ template "pegaEnvironmentConfig" }}
      key: DATA_SCHEMA
- name: DL-NAME
  value: EMPTY
{{- end }}

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
