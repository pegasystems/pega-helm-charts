{{- define "pegaVolumeInstall" }}pega-volume-installer{{- end }}
{{- define "pegaInstallConfig" }}pega-installer-config{{- end }}
{{- define "pegaDBInstall" -}}pega-db-install{{- end -}}
{{- define "pegaDBUpgrade" -}}pega-db-upgrade{{- end -}}
{{- define "installerConfig" -}}installer-config{{- end -}}
{{- define "installerJobReaderRole" -}}jobs-reader{{- end -}}
{{- define "pegaPreDBUpgrade" -}}pega-pre-upgrade{{- end -}}
{{- define "pegaPostDBUpgrade" -}}pega-post-upgrade{{- end -}}
{{- define "pegaInstallEnvironmentConfig" -}}pega-install-environment-config{{- end -}}
{{- define "pegaUpgradeEnvironmentConfig" -}}pega-upgrade-environment-config{{- end -}}
{{- define "pegaDistributionKitVolume" -}}pega-distribution-kit-volume{{- end -}}

{{- define "performInstall" }}
  {{- if or (eq .Values.global.actions.execute "install") (eq .Values.global.actions.execute "install-deploy") -}}
    true
  {{- else -}}
    false
  {{- end -}}
{{- end }}

{{- define "performUpgrade" }}
  {{- if or (eq .Values.global.actions.execute "upgrade") (eq .Values.global.actions.execute "upgrade-deploy") -}}
    true
  {{- else -}}
    false
  {{- end -}}
{{- end }}

{{- define "performOnlyUpgrade" }}
  {{- if (eq .Values.global.actions.execute "upgrade") -}}
    true
  {{- else -}}
    false
  {{- end -}}
{{- end }}

{{- define "waitForPegaDBInstall" -}}
- name: wait-for-pegainstall
  image: dcasavant/k8s-wait-for
  args: [ 'job', '{{ template "pegaDBInstall" }}']
{{- end }}

{{- define "waitForPegaDBUpgrade" -}}
- name: wait-for-pegaupgrade
  image: dcasavant/k8s-wait-for
  args: [ 'job', '{{ template "pegaDBUpgrade" }}']
{{- include "initContainerEnvs" $ }}
{{- end }}

{{- define "waitForPreDBUpgrade" -}}
- name: wait-for-pre-dbupgrade
  image: dcasavant/k8s-wait-for
  args: [ 'job', '{{ template "pegaPreDBUpgrade" }}']
{{- end }}

{{- define "waitForRollingUpdates" -}}
{{- $rolloutCommand := "" }}
{{- $kindName := "" }}
{{- $lastIndex := sub (len .Values.global.tier) 1 }}
{{- $namespace := .Release.Namespace }}
{{- range $index, $dep := .Values.global.tier }}
{{- if ($dep.volumeClaimTemplate) }}
{{- $kindName = "statefulset" }}
{{- else -}}
{{- $kindName = "deployment" }}
{{- end }}
{{- $constructCommand := cat "kubectl rollout status" $kindName "/" "pega-" $dep.name "--namespace" $namespace }}
{{- if ne $index $lastIndex }}
{{- $rolloutCommand = cat $rolloutCommand $constructCommand "&&" }}
{{- else }}
{{- $rolloutCommand = cat $rolloutCommand $constructCommand }}
{{- end }}
{{- $rolloutCommand = regexReplaceAllLiteral " / " $rolloutCommand "/" }}
{{- $rolloutCommand = regexReplaceAllLiteral "pega- " $rolloutCommand "pega-" }}
{{- end -}}
- name: wait-for-rolling-updates
  image: dcasavant/k8s-wait-for
  command: ['sh', '-c',  '{{ $rolloutCommand }}' ]
{{- include "initContainerEnvs" $ }}
{{- end }}

{{- define "initContainerEnvs" -}}
{{- if or (eq .Values.global.provider "aks") (eq .Values.global.provider "pks") -}}
{{ $apiserver := index .Values.global.upgrade "kube-apiserver" }}
  env:
  - name: KUBERNETES_SERVICE_HOST
    value: {{ $apiserver.serviceHost | quote }}
  - name: KUBERNETES_SERVICE_PORT_HTTPS
    value: {{ $apiserver.httpsServicePort | quote }}
  - name: KUBERNETES_SERVICE_PORT
    value: {{ $apiserver.httpsServicePort | quote }}
{{- end }}
{{- end }}