{{- define  "pega.deployment" -}}
kind: {{ .kind }}
apiVersion: {{ .apiVersion }}
metadata:
  name: {{ .name }}
  namespace: {{ .root.Release.Namespace }}
  labels:
    app: {{ template "pegaWebName" .root }} {{/* This is intentionally always the web name because that's what we call our "app" */}}
    component: Pega
spec:
  # Replicas specify the number of copies for {{ .name }}
  replicas: {{ .node.replicas }}
  selector:
    matchLabels:
      app: {{ .name }}
  template:
    metadata:
      labels:
        app: {{ .name }}
    spec:
      volumes:
      # Volume used to mount config files.
      - name: {{ template "pegaVolumeConfig" }}
        configMap:
          # This name will be referred in the volume mounts kind.
          name: {{ .name }}
          # Used to specify permissions on files within the volume.
          defaultMode: 420
      - name: {{ template "pegaVolumeCredentials" }}
        secret:
          # This name will be referred in the volume mounts kind.
          secretName: {{ template "pegaDatabaseSecret" }}
          # Used to specify permissions on files within the volume.
          defaultMode: 420
      initContainers:
{{- range $i, $val := .initContainers }}
{{ include $val $.root | indent 6 }}
{{- end }}
      containers:
      # Name of the container
      - name: pega-web-tomcat
        # The pega image, you may use the official pega distribution or you may extend
        # and host it yourself.  See the image documentation for more information.
        image: {{ .root.Values.docker.image }}
        # Pod (app instance) listens on this port
        ports:
        - containerPort: 8080
        # Specify any of the container environment variables here
        env:
        # Node type of the Pega nodes for {{ .name }}
        - name: NODE_TYPE
          value: {{ .nodeType }}
{{ include "pega.jvmconfig" (dict "node" .node) | indent 8 }}
{{ include "commonEnvironmentVariables" .root | indent 8 }}
        resources:
          # Maximum CPU and Memory that the containers for {{ .name }} can use
          limits:
            cpu: "{{ .node.cpuLimit }}"
            memory: "{{ .node.memLimit }}"
          # CPU and Memory that the containers for {{ .name }} request
          requests:
            cpu: "200m"
            memory: "2Gi"
        volumeMounts:
        # The given mountpath is mapped to volume with the specified name.  The config map files are mounted here.
        - name: {{ template "pegaVolumeConfig" }}
          mountPath: "/opt/pega/config"
{{- if .extraVolume }}
{{ include .extraVolume .root | indent 8 }}
{{- end }}
        - name: {{ template "pegaVolumeCredentials" }}
          mountPath: "/opt/pega/secrets"
{{ include "pega.health.probes" .root | indent 8 }}
      # Mentions the restart policy to be followed by the pod.  'Always' means that a new pod will always be created irrespective of type of the failure.
      restartPolicy: Always
      # Amount of time in which container has to gracefully shutdown.
      terminationGracePeriodSeconds: 300
      # Secret which is used to pull the image from the repository.  This secret contains docker login details for the particular user.
      # If the image is in a protected registry, you must specify a secret to access it.
      imagePullSecrets:
      - name: {{ template "pegaRegistrySecret" }}
{{- if .extraSpecData }}
{{ include .extraSpecData .root | indent 2 }}
{{- end }}
{{- end -}}