{{- define  "pega.deployment" -}}
kind: {{ .kind }}
apiVersion: {{ .apiVersion }}
metadata:
  annotations: 
{{- if .root.Values.global.pegaTier }}{{- if .root.Values.global.pegaTier.annotations }}
{{ toYaml .root.Values.global.pegaTier.annotations | indent 4 }}
{{- end }}{{- end }}
  name: {{ .name }}
  namespace: {{ .root.Release.Namespace }}
  labels:
    app: {{ .name }} {{/* This is intentionally always the web name because that's what we call our "app" */}}
    component: Pega
spec:
  # Replicas specify the number of copies for {{ .name }}
  replicas: {{ .node.replicas }}
{{- if (eq .kind "Deployment") }}
  progressDeadlineSeconds: 2147483647
{{- end }}
  selector:
    matchLabels:
      app: {{ .name }}
{{- if .node.deploymentStrategy }}
{{- if (ne .kind "Deployment") }}
  {{- $error := printf "tier[%s] may not specify a deploymentStrategy because it uses a volumeClaimTemplate which requires it be a StatefulSet" .name -}}
  {{ required $error nil }}
{{- end }}
  strategy:
{{ toYaml .node.deploymentStrategy | indent 4 }}
{{- end }}
  template:
    metadata:
      labels:
        app: {{ .name }}
{{- if .node.podLabels }}
{{ toYaml .node.podLabels | indent 8 }}
{{- end }}
      annotations:
{{- if .node.podAnnotations }}
{{ toYaml .node.podAnnotations | indent 8 }}
{{- end }}
        config-check: {{ include (print .root.Template.BasePath "/pega-environment-config.yaml") .root | sha256sum }}
        revision: "{{ .root.Release.Revision }}"
{{- include "generatedPodAnnotations" .root | indent 8 }}

    spec:
{{- if .custom }}
{{- if .custom.serviceAccountName }}
      serviceAccountName: {{ .custom.serviceAccountName }}
{{- end }}
{{- end }}
      volumes:
      # Volume used to mount config files.
      - name: {{ template "pegaVolumeConfig" }}
        configMap:
          # This name will be referred in the volume mounts kind.
          name: {{ .name }}
          # Used to specify permissions on files within the volume.
          defaultMode: 420
{{- include "pegaCredentialVolumeTemplate" . | indent 6 }}
{{- if .custom }}
{{- if .custom.volumes }}
      # Additional custom volumes
{{ toYaml .custom.volumes | indent 6 }}
{{- end }}
{{- end }}
      initContainers:
{{- range $i, $val := .initContainers }}
{{ include $val $.root | indent 6 }}
{{- end }}
{{- if .custom }}
{{- if .custom.initContainers }}
        # Additional custom init containers
{{ toYaml .custom.initContainers | indent 6 }}
{{- end }}
{{- end }}
{{- if .node.nodeSelector }}
      nodeSelector:
{{ toYaml .node.nodeSelector | indent 8 }}
{{- end }}
{{- if (ne .root.Values.global.provider "openshift") }}
      securityContext:
        fsGroup: 0
{{- if .node.securityContext }}
        runAsUser: {{ .node.securityContext.runAsUser }}
{{- else }}
        runAsUser: 9001
{{- end }}
{{- end }}
      containers:
      # Name of the container
      - name: pega-web-tomcat
        # The pega image, you may use the official pega distribution or you may extend
        # and host it yourself.  See the image documentation for more information.
        image: {{ .root.Values.global.docker.pega.image }}
{{- if (.root.Values.global.docker.pega.imagePullPolicy) }}
        imagePullPolicy: {{ .root.Values.global.docker.pega.imagePullPolicy }}
{{- end }}
        # Pod (app instance) listens on this port
        ports:
        - containerPort: 8080
          name: pega-web-port
{{- if .custom }}
{{- if .custom.ports }}
        # Additional custom ports
{{ toYaml .custom.ports | indent 8 }}
{{- end }}
{{- end }}
        # Specify any of the container environment variables here
        env:	
        # Node type of the Pega nodes for {{ .name }}
{{- if .root.Values.stream }}
{{- if .root.Values.stream.url }}
{{- if contains "Stream" .nodeType }}
{{ fail "Cannot have 'Stream' nodeType when Stream url is provided" }}
{{- end }}
{{- end }}
{{- end }}
        - name: NODE_TYPE
          value: {{ .nodeType }}
        - name: PEGA_APP_CONTEXT_PATH
          value: {{ template "pega.applicationContextPath" . }}
{{- if .node.requestor }}
        - name: REQUESTOR_PASSIVATION_TIMEOUT
          value: "{{ .node.requestor.passivationTimeSec }}"
{{- end }}
{{- if and .root.Values.constellation (eq .root.Values.constellation.enabled true) }}
        - name: COSMOS_SETTINGS
          value: "Pega-UIEngine/cosmosservicesURI=/c11n"
{{- end }}
{{- if .custom }}
{{- if .custom.env }}
        # Additional custom env vars
{{ toYaml .custom.env | indent 8 }}
{{- end }}
{{- end }}
{{ include "pega.jvmconfig" (dict "node" .node) | indent 8 }}
        # Tier of the Pega node
        - name: NODE_TIER
          value: {{ .tierName }}
        envFrom:
        - configMapRef:
            name: {{ template "pegaEnvironmentConfig" }}
        resources:
          # Maximum CPU and Memory that the containers for {{ .name }} can use
          limits:
          {{- if .node.cpuLimit }}
            cpu: "{{ .node.cpuLimit }}"
          {{- else }}
            cpu: 4
          {{- end }}
          {{- if .node.memLimit }}
            memory: "{{ .node.memLimit }}"
          {{- else }}
            memory: "8Gi"
          {{- end }}
          # CPU and Memory that the containers for {{ .name }} request
          requests:
          {{- if .node.cpuRequest }}
            cpu: "{{ .node.cpuRequest }}"
          {{- else }}
            cpu: 2
          {{- end }}
          {{- if .node.memRequest }}
            memory: "{{ .node.memRequest }}"
          {{- else }}
            memory: "6Gi"
          {{- end }}
        volumeMounts:
        # The given mountpath is mapped to volume with the specified name.  The config map files are mounted here.
        - name: {{ template "pegaVolumeConfig" }}
          mountPath: "/opt/pega/config"
{{- if (.node.volumeClaimTemplate) }}
        - name: {{ .name }}
          mountPath: "/opt/pega/kafkadata"
{{- end }}
{{- if .custom }}
{{- if .custom.volumeMounts }}
        # Additional custom mounts
{{ toYaml .custom.volumeMounts | indent 8 }}
{{- end }}
{{- end }}
        - name: {{ template "pegaVolumeCredentials" }}
          mountPath: "/opt/pega/secrets"
        # LivenessProbe: indicates whether the container is live, i.e. running.
        {{- $livenessProbe := .node.livenessProbe }}
        livenessProbe:
          httpGet:
            path: "/{{ template "pega.applicationContextPath" . }}/PRRestService/monitor/pingService/ping"
            port: 8080
            scheme: HTTP
          initialDelaySeconds: {{ $livenessProbe.initialDelaySeconds | default 300 }}
          timeoutSeconds: {{ $livenessProbe.timeoutSeconds | default 20 }}
          periodSeconds: {{ $livenessProbe.periodSeconds | default 30 }}
          successThreshold: {{ $livenessProbe.successThreshold | default 1 }}
          failureThreshold: {{ $livenessProbe.failureThreshold | default 3 }}
        # ReadinessProbe: indicates whether the container is ready to service requests.
        {{- $readinessProbe := .node.readinessProbe }}
        readinessProbe:
          httpGet:
            path: "/{{ template "pega.applicationContextPath" . }}/PRRestService/monitor/pingService/ping"
            port: 8080
            scheme: HTTP
          initialDelaySeconds: {{ $readinessProbe.initialDelaySeconds | default 300 }}
          timeoutSeconds: {{ $readinessProbe.timeoutSeconds | default 20 }}
          periodSeconds: {{ $readinessProbe.periodSeconds | default 30 }}
          successThreshold: {{ $readinessProbe.successThreshold | default 1 }}
          failureThreshold: {{ $readinessProbe.failureThreshold | default 3 }}
      # Mentions the restart policy to be followed by the pod.  'Always' means that a new pod will always be created irrespective of type of the failure.
      restartPolicy: Always
      # Amount of time in which container has to gracefully shutdown.
      terminationGracePeriodSeconds: 300
      # Secret which is used to pull the image from the repository.  This secret contains docker login details for the particular user.
      # If the image is in a protected registry, you must specify a secret to access it.
      imagePullSecrets:
      - name: {{ template "pegaRegistrySecret" }}
{{- if (.node.volumeClaimTemplate) }}
  volumeClaimTemplates:
  - metadata:
      name: {{ .name }}
      creationTimestamp:
    spec:
      accessModes:
      - ReadWriteOnce
      resources:
        requests:
          storage: {{ .node.volumeClaimTemplate.resources.requests.storage }}
  serviceName: {{ .name }}
{{- end }}
---
{{- end -}}
