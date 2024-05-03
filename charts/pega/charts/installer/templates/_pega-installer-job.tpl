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
      labels:
        app: "installer"
        installer-job: {{ .name }}
        {{- if .root.Values.podLabels }}
{{ toYaml .root.Values.podLabels | indent 8 }}
        {{- end -}}
{{ include "generatedInstallerPodLabels" .root | indent 8 }}
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
{{- if .root.Values.installerMountVolumeClaimName }}
      - name: {{ template "pegaInstallerMountVolume" }}
        persistentVolumeClaim:
          claimName: {{ .root.Values.installerMountVolumeClaimName }}
{{- end }}
{{- if and .root.Values.distributionKitVolumeClaimName (not .root.Values.distributionKitURL) }}
      - name: {{ template "pegaDistributionKitVolume" }}
        persistentVolumeClaim:
          claimName: {{ .root.Values.distributionKitVolumeClaimName }}
{{- end }}
{{- if .root.Values.custom }}{{- if .root.Values.custom.volumes }}
{{ toYaml .root.Values.custom.volumes | indent 6 }}
{{- end }}{{- end }}
      - name: {{ template "pegaInstallerCredentialsVolume" }}
        projected:
          defaultMode: 420
          sources:
          {{- $d := dict "deploySecret" "deployDBSecret" "deployNonExtsecret" "deployNonExtDBSecret" "extSecretName" .root.Values.global.jdbc.external_secret_name "nonExtSecretName" "pega-db-secret-name" "context" .root  -}}
          {{ include "secretResolver" $d | indent 10}}

          {{- $artifactoryDict := dict "deploySecret" "deployArtifactorySecret" "deployNonExtsecret" "deployNonExtArtifactorySecret" "extSecretName" .root.Values.global.customArtifactory.authentication.external_secret_name "nonExtSecretName" "pega-custom-artifactory-secret-name" "context" .root -}}
          {{ include "secretResolver" $artifactoryDict | indent 10}}

          # Fix it, Below peace of code always uses secret created from hz username & password. It cannot resolve hz external secret due to helm sub chart limitations. Modify it once hazelcast deployment is isolated.
          {{- if ( eq .root.Values.upgrade.isHazelcastClientServer "true" ) }}
          - secret:
              name: {{ include  "pega-hz-secret-name" .root}}
          {{- end }}
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
{{- if .root.Values.nodeSelector }}
      nodeSelector:
{{ toYaml .root.Values.nodeSelector | indent 8 }}
{{- end }}
      containers:
      - name: {{ template "pegaDBInstallerContainer" }}
        image: {{ .root.Values.image }}
{{- if .root.Values.imagePullPolicy }}
        imagePullPolicy: {{ .root.Values.imagePullPolicy  }}
{{- end }}
        ports:
        - containerPort: 8080
{{- if .root.Values.securityContext }}
        securityContext:
{{ toYaml .root.Values.securityContext | indent 10 }}
{{- end }}
        resources:
          # CPU and Memory that the containers for {{ .name }} request
          requests:
            cpu: "{{ .root.Values.resources.requests.cpu }}"
            memory: "{{ .root.Values.resources.requests.memory }}"
          limits:
            cpu: "{{ .root.Values.resources.limits.cpu }}"
            memory: "{{ .root.Values.resources.limits.memory }}"
        volumeMounts:
{{- if .root.Values.installerMountVolumeClaimName }}
        - name: {{ template "pegaInstallerMountVolume" }}
          mountPath: "/opt/pega/mount/installer"
{{- end }}
        # The given mountpath is mapped to volume with the specified name.  The config map files are mounted here.
        - name: {{ template "pegaVolumeInstall" }}
          mountPath: "/opt/pega/config"
        - name: {{ template "pegaInstallerCredentialsVolume" }}
          mountPath: "/opt/pega/secrets"
{{- if and .root.Values.distributionKitVolumeClaimName (not .root.Values.distributionKitURL) }}
        - name: {{ template "pegaDistributionKitVolume" }}
          mountPath: "/opt/pega/mount/kit"
{{- end }}
{{- if .root.Values.custom }}
{{- if .root.Values.custom.volumeMounts }}
{{ toYaml .root.Values.custom.volumeMounts | indent 8 }}
{{- end }}
{{- end }}
{{ if (eq (include "customArtifactorySSLVerificationEnabled" .root) "true") }}
{{- if .root.Values.global.customArtifactory.certificate }}
        - name: {{ template "pegaVolumeCustomArtifactoryCertificate" }}
          mountPath: "/opt/pega/artifactory/cert"
{{- end }}
{{- end }}
        env:
        - name: ACTION
          value: {{ .action }}
{{- if .root.Values.custom }}
{{- if .root.Values.custom.env }}
        # Additional custom env vars
{{ toYaml .root.Values.custom.env | indent 8 }}
{{- end }}
{{- end }}
{{- if or (eq $arg "pre-upgrade") (eq $arg "post-upgrade") (eq $arg "upgrade")  }}
{{- if (eq .root.Values.upgrade.isHazelcastClientServer "true") }}
        -  name: HZ_VERSION
           valueFrom:
            configMapKeyRef:
              name: {{ template "pegaEnvironmentConfig" }}
              key: HZ_VERSION
        -  name: HZ_CLUSTER_NAME
           valueFrom:
            configMapKeyRef:
              name: {{ template "pegaEnvironmentConfig" }}
              key: HZ_CLUSTER_NAME
        -  name: HZ_SERVER_HOSTNAME
           valueFrom:
            configMapKeyRef:
              name: {{ template "pegaEnvironmentConfig" }}
              key: HZ_SERVER_HOSTNAME
{{- end }}
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
{{- include "imagePullSecrets" .root | indent 6 }}
---
{{- end -}}