{{- define  "pega.deployment" -}}
{{- $useStartupProbe := false }}
{{- $livenessProbe := .node.livenessProbe }}
{{- $readinessProbe := .node.readinessProbe }}
{{- $livenessProbeInitialDelaySeconds := $livenessProbe.initialDelaySeconds | default 200 }}
{{- $livenessProbeFailureThreshold := $livenessProbe.failureThreshold | default 3 }}
{{- $livenessProbePeriodSeconds := $livenessProbe.periodSeconds | default 30 }}
{{- $readinessProbeInitialDelaySeconds := $readinessProbe.initialDelaySeconds | default 30 }}
{{- if (semverCompare ">= 1.18.0-0" (trimPrefix "v" .root.Capabilities.KubeVersion.GitVersion)) }}
  {{- $useStartupProbe = true }}
  {{- $livenessProbeInitialDelaySeconds = $livenessProbe.initialDelaySeconds | default 0 }}
  {{- $readinessProbeInitialDelaySeconds = $readinessProbe.initialDelaySeconds | default 0 }}
{{- end }}

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
{{- if .root.Values.global.pegaTier }}{{- if .root.Values.global.pegaTier.labels }}
{{ toYaml .root.Values.global.pegaTier.labels | indent 4 }}
{{- end }}{{- end }}
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
{{- include "generatedPodLabels" .root | indent 8 }}
{{- end }}
      annotations:
{{- if .node.podAnnotations }}
{{ toYaml .node.podAnnotations | indent 8 }}
{{- end }}
        config-check: {{ include (print .root.Template.BasePath "/pega-environment-config.yaml") .root | sha256sum }}
        config-tier-check: {{ include "pega.config" (dict "root" .root "dep" .node) | sha256sum }}
        certificate-check: {{ include (print .root.Template.BasePath "/pega-certificates-secret.yaml") .root | sha256sum }}
{{- include "generatedPodAnnotations" .root | indent 8 }}

    spec:
{{- include "generatedDNSConfigAnnotations" .root | indent 6 }}
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
{{- include "pegaCredentialVolumeTemplate" .root | indent 6 }}
{{ if or (.root.Values.global.certificates) (.root.Values.global.certificatesSecrets) }}
{{- include "pegaImportCertificatesTemplate" .root | indent 6 }}
{{ end }}
{{ if (eq (include "customArtifactorySSLVerificationEnabled" .root) "true") }}
{{- if .root.Values.global.customArtifactory.certificate }}
{{- include "pegaCustomArtifactoryCertificateTemplate" .root | indent 6 }}
{{- end }}
{{- end }}
{{- if ((.node.service).tls).enabled }}
{{- $data := dict "root" .root "node" .node }}
{{- include "pegaVolumeTomcatKeystoreTemplate" $data | indent 6 }}
{{ end }}
{{- if .root.Values.global.kerberos }}
{{- include "pegaKerberosVolumeTemplate" .root | indent 6 }}
{{- end }}
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
      securityContext:
{{- if (ne .root.Values.global.provider "openshift") }}
        runAsUser: 9001
        fsGroup: 0
{{- end }}
{{- if .node.securityContext }}
{{ toYaml .node.securityContext | indent 8 }}
{{- end }}
{{- if .node.topologySpreadConstraints }}
      topologySpreadConstraints:
{{ toYaml .node.topologySpreadConstraints | indent 8 }}
{{- end }}
{{- if .node.tolerations }}
      tolerations:
{{ toYaml .node.tolerations | indent 8 }}
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
        - containerPort: 8443
          name: pega-tls-port
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
        - name: POD_NAME
          valueFrom:
            fieldRef:
                apiVersion: v1
                fieldPath: metadata.name
{{- if .node.requestor }}
        - name: REQUESTOR_PASSIVATION_TIMEOUT
          value: "{{ .node.requestor.passivationTimeSec }}"
{{- end }}
{{- if and .root.Values.constellation (eq .root.Values.constellation.enabled true) }}
        - name: COSMOS_SETTINGS
          value: "Pega-UIEngine/cosmosservicesURI=/c11n"
{{- end }}
{{- if ((.node.service).tls).enabled }}
        - name: EXTERNAL_KEYSTORE_NAME
          value: "{{ (((.node.service).tls).external_keystore_name) }}"
        - name: EXTERNAL_KEYSTORE_PASSWORD
          value: "{{ (((.node.service).tls).external_keystore_password) }}"
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
        - name: RETRY_TIMEOUT
          value: {{ include "tierClassloaderRetryTimeout" (dict "failureThreshold" $livenessProbeFailureThreshold "periodSeconds" $livenessProbePeriodSeconds ) | quote }}
        - name: MAX_RETRIES
          value: {{ include "tierClassloaderMaxRetries" (dict "failureThreshold" $livenessProbeFailureThreshold "periodSeconds" $livenessProbePeriodSeconds ) | quote }}
{{- if and (.root.Values.pegasearch.externalSearchService) ((.root.Values.pegasearch.srsAuth).enabled) }}
{{- if or (not .root.Values.pegasearch.srsAuth.authType) (eq .root.Values.pegasearch.srsAuth.authType "private_key_jwt") }}
        - name: SERV_AUTH_PRIVATE_KEY
          valueFrom:
            secretKeyRef:
{{- include "srsAuthEnvSecretFrom"  .root | indent 14 }}
{{- else if eq .root.Values.pegasearch.srsAuth.authType "client_secret_basic" }}
        - name: SERV_AUTH_CLIENT_SECRET
          valueFrom:
            secretKeyRef:
{{- include "srsAuthEnvSecretFrom"  .root | indent 14 }}
{{- else }}
  {{- fail "pegasearch.srsAuth.authType must be either private_key_jwt or client_secret_basic." }}
{{- end }}
{{- end }}
        envFrom:
        - configMapRef:
            name: {{ template "pegaEnvironmentConfig" .root }}
        resources:
{{- if .node.resources }}
{{ toYaml .node.resources | indent 10 }}
{{- else }}
          # Maximum CPU and Memory that the containers for {{ .name }} can use
          # Resources are configured through deprecated settings. Use .tier[].resources instead
          limits:
          {{- if .node.cpuLimit }}
            cpu: "{{ .node.cpuLimit }}"
          {{- else }}
            cpu: 4
          {{- end }}
          {{- if .node.memLimit }}
            memory: "{{ .node.memLimit }}"
          {{- else }}
            memory: "12Gi"
          {{- end }}
          {{- if .node.ephemeralStorageLimit }}
            ephemeral-storage: "{{ .node.ephemeralStorageLimit }}"
          {{- end }}
          # CPU and Memory that the containers for {{ .name }} request
          requests:
          {{- if .node.cpuRequest }}
            cpu: "{{ .node.cpuRequest }}"
          {{- else }}
            cpu: 3
          {{- end }}
          {{- if .node.memRequest }}
            memory: "{{ .node.memRequest }}"
          {{- else }}
            memory: "12Gi"
          {{- end }}
          {{- if .node.ephemeralStorageRequest }}
            ephemeral-storage: "{{ .node.ephemeralStorageRequest }}"
          {{- end }}
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
        #mount custom certificates
{{ if or (.root.Values.global.certificates) (.root.Values.global.certificatesSecrets) }}
        - name: {{ template "pegaVolumeImportCertificates" }}
          mountPath: "/opt/pega/certs"
{{ end }}
{{- if ((.node.service).tls).enabled }}
        - name: {{ template "pegaVolumeTomcatKeystore" }}
          mountPath: "/opt/pega/tomcatcertsmount"
{{ end }}
{{ if (eq (include "customArtifactorySSLVerificationEnabled" .root) "true") }}
{{- if .root.Values.global.customArtifactory.certificate }}
        - name: {{ template "pegaVolumeCustomArtifactoryCertificate" }}
          mountPath: "/opt/pega/artifactory/cert"
{{- end }}
{{- end }}
{{- if .root.Values.global.kerberos }}
        - name: {{ template "pegaKerberosConfig" }}-config
          mountPath: "/opt/pega/kerberos"
{{- end }}

        # LivenessProbe: indicates whether the container is live, i.e. running.
        livenessProbe:
          httpGet:
            path: "/{{ template "pega.applicationContextPath" . }}/PRRestService/monitor/pingService/ping"
            port: {{ $livenessProbe.port | default 8080 }}
            scheme: HTTP
          initialDelaySeconds: {{ $livenessProbeInitialDelaySeconds }}
          timeoutSeconds: {{ $livenessProbe.timeoutSeconds | default 20 }}
          periodSeconds: {{ $livenessProbePeriodSeconds }}
          successThreshold: {{ $livenessProbe.successThreshold | default 1 }}
          failureThreshold: {{ $livenessProbeFailureThreshold }}
        # ReadinessProbe: indicates whether the container is ready to service requests.
        readinessProbe:
          httpGet:
            path: "/{{ template "pega.applicationContextPath" . }}/PRRestService/monitor/pingService/ping"
            port: {{ $readinessProbe.port | default 8080 }}
            scheme: HTTP
          initialDelaySeconds: {{ $readinessProbeInitialDelaySeconds }}
          timeoutSeconds: {{ $readinessProbe.timeoutSeconds | default 10 }}
          periodSeconds: {{ $readinessProbe.periodSeconds | default 10 }}
          successThreshold: {{ $readinessProbe.successThreshold | default 1 }}
          failureThreshold: {{ $readinessProbe.failureThreshold | default 3 }}
        # StartupProbe: indicates whether the container has completed its startup process, and delays the LivenessProbe
{{- if ( $useStartupProbe ) }}
        {{- $startupProbe := .node.startupProbe }}
        startupProbe:
          httpGet:
            path: "/{{ template "pega.applicationContextPath" . }}/PRRestService/monitor/pingService/ping"
            port: {{ $startupProbe.port | default 8080 }}
            scheme: HTTP
          initialDelaySeconds: {{ $startupProbe.initialDelaySeconds | default 10 }}
          timeoutSeconds: {{ $startupProbe.timeoutSeconds | default 10 }}
          periodSeconds: {{ $startupProbe.periodSeconds | default 10 }}
          successThreshold: {{ $startupProbe.successThreshold | default 1 }}
          failureThreshold: {{ $startupProbe.failureThreshold | default 30 }}
{{- end }}

{{- if .custom }}
{{- if .custom.sidecarContainers }}
      # Additional custom sidecar containers
{{ toYaml .custom.sidecarContainers | indent 6 }}
{{- end }}
{{- end }}
      # Mentions the restart policy to be followed by the pod.  'Always' means that a new pod will always be created irrespective of type of the failure.
      restartPolicy: Always
      # Amount of time in which container has to gracefully shutdown.
      terminationGracePeriodSeconds: 300
      # Secret which is used to pull the image from the repository.  This secret contains docker login details for the particular user.
      # If the image is in a protected registry, you must specify a secret to access it.
      imagePullSecrets:
{{- include "imagePullSecrets" .root | indent 6 }}
{{- include "podAffinity" .node | indent 6 }}
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
{{- if ( .root.Values.global.storageClassName ) }}
      storageClassName: {{ .root.Values.global.storageClassName }}
{{ end }}
  serviceName: {{ .name }}
{{- end }}
---
{{- end -}}
