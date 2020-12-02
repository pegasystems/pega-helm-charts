{{- define "hazelcastName" -}} pega-hazelcast {{- end -}}
{{- define "hazelcastEnvironmentConfig" -}} pega-hz-env-config {{- end -}}


{{- define "isHazelcastEnabled" }}
 {{- if .Values.enabled -}}
  true
 {{- else -}}
  false
 {{- end -}}
{{- end }}

# Override this template to generate additional pod annotations that are dynamically composed during helm deployment (do not indent annotations)
{{- define "generatedHazelcastServicePodAnnotations" }}
{{- end }}

{{- define "waitForHazelcast" -}}
- name: wait-for-hazelcast
  image: busybox:1.31.0
  # Init container for waiting for hazelcast service to initialize.
  args:
  - sh
  - -c
  - >
    set -e
    counter=0;
    while [ -z $(wget -S "http://{{ template "hazelcastName" . }}-service.{{ .Release.Namespace }}:5701/hazelcast/health/cluster-state" 2>&1 | grep "HTTP/" | awk '{print $2}') ] || [ $(wget -q -O - "http://{{ template "hazelcastName" . }}-service.{{ .Release.Namespace }}:5701/hazelcast/health/cluster-size" /dev/null) -ne {{ .Values.hazelcast.replicas }} ] || [ $(wget -S "http://{{ template "hazelcastName" . }}-service.{{ .Release.Namespace }}:5701/hazelcast/health/cluster-state" 2>&1 | grep "HTTP/" | awk '{print $2}') -ne 200 ]; do
    echo "waiting for hazelcast pods to start and join the cluster..." ;
    counter=$(($counter+5));
    sleep 5;
    if [ $counter -gt 150 ]; then
    echo "Timeout Reached. Hazelcast pods failed to join the cluster";
    exit 1;
    fi
    done;
    echo "Hazelcast cluster is up now";
    exit 0;
{{- end }}

