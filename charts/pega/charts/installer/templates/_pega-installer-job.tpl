{{- define  "pega.installer" -}}
{{- $arg := .action -}}
kind: Job
apiVersion: batch/v1
metadata:
  name: {{ .name }}
  namespace: {{ .root.Release.Namespace }}
{{- if and .root.Values.waitForJobCompletion (or (eq .root.Values.global.actions.execute "install") (eq .root.Values.global.actions.execute "upgrade")) }}
  annotations:
    # Forces Helm to wait for the install or upgrade to complete.
    "helm.sh/hook": post-install
    "helm.sh/hook-weight": "0"
    "helm.sh/hook-delete-policy": before-hook-creation
{{- end }}
spec:
  backoffLimit: 0
  template:
    spec:
      volumes:
{{- if and .root.Values.distributionKitVolumeClaimName (not .root.Values.distributionKitURL) }}
      - name: {{ template "pegaDistributionKitVolume" }}
        persistentVolumeClaim:
          claimName: {{ .root.Values.distributionKitVolumeClaimName }}
{{- end }}      
      - name: {{ template "pegaVolumeCredentials" }}
        secret:
          # This name will be referred in the volume mounts kind.
          secretName: {{ template "pegaCredentialsSecret" }}
          # Used to specify permissions on files within the volume.
          defaultMode: 420
      - name: {{ template "pegaVolumeInstall" }}
        configMap:
          # This name will be referred in the volume mounts kind.
          name: {{ template "pegaInstallConfig"}}
          # Used to specify permissions on files within the volume.
          defaultMode: 420          
      initContainers:
{{- range $i, $val := .initContainers }}
{{ include $val $.root | indent 6 }}
{{- end }}
      containers:
      - name: {{ .name }}
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
      restartPolicy: Never
      imagePullSecrets:
      - name: {{ template "pegaRegistrySecret" }}
---
{{- end -}}
