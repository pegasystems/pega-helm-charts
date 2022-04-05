{{- define  "pega.installer" -}}
{{- $arg := .action -}}
kind: Job
apiVersion: batch/v1
metadata:
  name: {{ .name }}
  namespace: {{ .root.Release.Namespace }}
  annotations:
{{- if  (eq .root.Values.waitForJobCompletion "true")   }}
    "helm.sh/hook-weight": "0"
    "helm.sh/hook-delete-policy": {{ if .root.Values.cleanAfterInstall -}} before-hook-creation,hook-succeeded {{- else -}} before-hook-creation {{- end }}
{{- if  (eq .root.Values.global.actions.execute "install") }}
    # Forces Helm to wait for the install to complete.
    "helm.sh/hook": post-install
{{- end }}
{{- if (eq .root.Values.global.actions.execute "upgrade") }}
    # Forces Helm to wait for the upgrade to complete.
    "helm.sh/hook": post-install, post-upgrade
{{- end }}
{{- end }}
{{- if .root.Values.global.pegaJob }}{{- if .root.Values.global.pegaJob.annotations }}
{{ toYaml .root.Values.global.pegaJob.annotations | indent 4 }}
{{- end }}{{- end }}
  labels:
    app: {{ .name }}
spec:
  backoffLimit: 0
  template:
    metadata:
      annotations:
{{- if .root.Values.podAnnotations}}
{{ toYaml .root.Values.podAnnotations | indent 8 }}
{{- end }}     
    spec:
      shareProcessNamespace: {{ .root.Values.shareProcessNamespace }}
{{- if .root.Values.serviceAccountName }}
      serviceAccountName: {{ .root.Values.serviceAccountName }}
{{- end }}   
      volumes:
{{- if and .root.Values.distributionKitVolumeClaimName (not .root.Values.distributionKitURL) }}
      - name: {{ template "pegaDistributionKitVolume" }}
        persistentVolumeClaim:
          claimName: {{ .root.Values.distributionKitVolumeClaimName }}
{{- end }}
{{- if .root.Values.custom }}{{- if .root.Values.custom.volumes }}
{{ toYaml .root.Values.custom.volumes | indent 6 }}          
{{- end }}{{- end }}  
{{- include "pegaCredentialVolumeTemplate" .root | indent 6 }}
      - name: {{ template "pegaVolumeInstall" }}
        configMap:
          # This name will be referred in the volume mounts kind.
     {{- if or (eq $arg "install") (eq $arg "install-deploy") }}
          name: {{ template "pegaInstallConfig"}}
     {{- else }}
          name: {{ template "pegaUpgradeConfig"}}
     {{- end }}
          # Used to specify permissions on files within the volume.
          defaultMode: 420
{{ if (eq (include "customArtifactorySSLVerificationEnabled" .root) "true") }}
{{- if .root.Values.global.customArtifactory.certificate }}
{{- include "pegaCustomArtifactoryCertificateTemplate" .root | indent 6 }}
{{- end }}
{{- end }}
      initContainers:
{{- range $i, $val := .initContainers }}
{{ include $val $.root | indent 6 }}
{{- end }}
      containers:
      - name: {{ template "pegaDBInstallerContainer" }}
        image: {{ .root.Values.image }}
        ports:
        - containerPort: 8080
        resources:
          # CPU and Memory that the containers for {{ .name }} request
          requests:
            cpu: "{{ .root.Values.resources.requests.cpu }}"
            memory: "{{ .root.Values.resources.requests.memory }}"
          limits:
            cpu: "{{ .root.Values.resources.limits.cpu }}"
            memory: "{{ .root.Values.resources.limits.memory }}"
        volumeMounts:
        # The given mountpath is mapped to volume with the specified name.  The config map files are mounted here.
        - name: {{ template "pegaVolumeInstall" }}
          mountPath: "/opt/pega/config"
        - name: {{ template "pegaVolumeCredentials" }}
          mountPath: "/opt/pega/secrets"
{{- if and .root.Values.distributionKitVolumeClaimName (not .root.Values.distributionKitURL) }}          
        - name: {{ template "pegaDistributionKitVolume" }}
          mountPath: "/opt/pega/mount/kit"                           
{{- end }}
{{ if (eq (include "customArtifactorySSLVerificationEnabled" .root) "true") }}
{{- if .root.Values.global.customArtifactory.certificate }}
        - name: {{ template "pegaVolumeCustomArtifactoryCertificate" }}
          mountPath: "/opt/pega/artifactory/cert"
{{- end }}
{{- end }}
{{- if or (eq $arg "pre-upgrade") (eq $arg "post-upgrade") (eq $arg "upgrade")  }}
        env:
        -  name: ACTION
           value: {{ .action }}
        envFrom:
        - configMapRef:
            name: {{ template "pegaUpgradeEnvironmentConfig" }}
{{- end }}
{{- if (eq $arg "install") }}
        envFrom:
        - configMapRef:
            name: {{ template "pegaInstallEnvironmentConfig" }}
{{- end }}
{{- if .root.Values.sidecarContainers }}
{{ toYaml .root.Values.sidecarContainers | indent 6 }}
{{- end }}                
      restartPolicy: Never
      imagePullSecrets:
      - name: {{ template "pegaRegistrySecret" .root }}
---
{{- end -}}
