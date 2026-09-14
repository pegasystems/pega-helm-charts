{{- define "distributionKitUrl" }}
{{- if and .Values.distributionKit .Values.distributionKit.url }}
{{- .Values.distributionKit.url }}
{{- else if .Values.distributionKitURL }}
{{- .Values.distributionKitURL }}
{{- end }}
{{- end }}

{{- define "usesKitICDownload" }}
{{- $usesKitICDownload := "false" -}}
{{- if .Values.distributionKit }}
{{- if .Values.distributionKit.downloadContainer }}
{{- if .Values.distributionKit.downloadContainer.image }}
{{- $usesKitICDownload = "true" }}
{{- end }}
{{- end }}
{{- end }}
{{- $usesKitICDownload -}}
{{- end }}

{{- define "pegaKitDownloadScriptConfig" }}
{{- printf "%s" (include "deploymentName" $) -}}-installer-kit-download-script-config
{{- end }}

{{- define "usesKitPVC" }}
{{- $usesKitPVC := "false" -}}
{{- if and .Values.distributionKitVolumeClaimName (eq (include "distributionKitUrl" .) "") }}
{{- $usesKitPVC = "true" }}
{{- end }}
{{- $usesKitPVC -}}
{{- end }}

{{- define "kitVolume" }}
{{- if eq (include "usesKitPVC" .) "true" }}
- name: {{ template "pegaDistributionKitVolume" }}
  persistentVolumeClaim:
    claimName: {{ .Values.distributionKitVolumeClaimName }}
{{- else if eq (include "usesKitICDownload" .) "true" }}
- name: {{ template "pegaDistributionKitVolume" }}
  emptyDir:
    sizeLimit: {{ default "5Gi" .Values.distributionKit.downloadContainer.sharedVolumeSize }}
{{- end }}
{{- end }}

{{- define "kitVolumeMount" }}
{{- if or (eq (include "usesKitPVC" .) "true") (eq (include "usesKitICDownload" .) "true") }}
- name: {{ template "pegaDistributionKitVolume" }}
  mountPath: "/opt/pega/mount/kit"
{{- end }}
{{- end }}

{{- define "kitDownloadScriptVolume" }}
{{- if eq (include "usesKitICDownload" .) "true" }}
- name: kit-download-script-volume
  configMap:
    name: {{ template "pegaKitDownloadScriptConfig" $ }}
    defaultMode: 0555
{{- end }}
{{- end }}

{{- define "kitDownloadScriptConfigMap" }}
{{- if eq (include "usesKitICDownload" .) "true" }}
---
# This template contains the script to download the Pega distribution kit from the specified location and copy it to
# the shared volume. This script will be executed by the init container before starting the Pega installer container.
kind: ConfigMap
apiVersion: v1
metadata:
  name: {{ template "pegaKitDownloadScriptConfig" . }}
  namespace: {{ .Release.Namespace }}
data:
  download-kit.sh: |
    #!/bin/sh

    pega_root="/opt/pega"
    kit_root="$pega_root/mount/kit"
    art_root="$pega_root/artifactory"
    secret_root="$pega_root/secrets"

    cau="$(cat $secret_root/CUSTOM_ARTIFACTORY_USERNAME)"
    cap="$(cat $secret_root/CUSTOM_ARTIFACTORY_PASSWORD)"
    caah="$(cat $secret_root/CUSTOM_ARTIFACTORY_APIKEY_HEADER)"
    caak="$(cat $secret_root/CUSTOM_ARTIFACTORY_APIKEY)"

    ca_auth=""
    if [ "$cau" != "" ] || [ "$cap" != "" ]; then
      if [ "$cau" = "" ] || [ "$cap" = "" ]; then
        echo "CUSTOM_ARTIFACTORY_USERNAME & CUSTOM_ARTIFACTORY_PASSWORD must be specified for artifactory basic auth."
        exit 1
      else
        echo "Using basic authentication for custom artifactory to download the distribution kit."
        ca_auth="-u $cau:$cap"
      fi
    fi

    if [ "$ca_auth" = "" ]; then
      if [ "$caah" != "" ] || [ "$caak" != "" ]; then
        if [ "$caah" = "" ] || [ "$caak" = "" ]; then
          echo "CUSTOM_ARTIFACTORY_APIKEY_HEADER & CUSTOM_ARTIFACTORY_APIKEY must be specified for authentication using api key for custom artifactory."
          exit 1
        else
          echo "Using API key for artifactory authentication."
          ca_auth="-H $caah:$caak"
        fi
      fi
    fi

    cert_type_list="cer
    pem
    crt
    der
    cert
    jks
    p7b
    p7c
    key"
    ca_cert=""
    if [ "$(ls -A $art_root/cert/*)" ]; then
      if [ "$(ls $art_root/cert/* | wc -l)" = "1" ]; then
        echo "Certificate is provided for custom artifactory's domain ssl verification."
        certfilename="$(ls $art_root/cert)"
        ext="${certfilename##*.}"

        echo "false" > ${kit_root}/cert_type_is_valid.txt
        echo "$cert_type_list" | while IFS= read -r line; do
          if [ "$ext" = "$line" ]; then
           echo "true" > ${kit_root}/cert_type_is_valid.txt
          fi
        done
        isValid="$(cat ${kit_root}/cert_type_is_valid.txt)"
        rm ${kit_root}/cert_type_is_valid.txt
        if [ "$isValid" = "true" ]; then
          echo "Using $certfilename"
          ca_cert="--cacert $art_root/cert/$certfilename"
        else
          echo "curl needs valid format certificate file for ssl verification."
          exit 1
        fi
      else
        echo "Provide one certificate file. The file may contain multiple CA certificates."
        exit 1
      fi
    fi

    if [ "$DISTRIBUTION_KIT_URL" != "" ]; then
      curl_cmd_options=""
      if [ "$ENABLE_CUSTOM_ARTIFACTORY_SSL_VERIFICATION" = "true" ]; then
        echo "Establishing a secure connection to download the distribution kit."
        curl_cmd_options="-sSL $ca_auth $ca_cert"
      else
        echo "Establishing an insecure connection to download the distribution kit."
        curl_cmd_options="-ksSL $ca_auth"
      fi
      echo "Downloading distribution kit: ${DISTRIBUTION_KIT_URL}";
      kitabsurl="$(echo "$DISTRIBUTION_KIT_URL" | cut -d'?' -f1)"
      filename=$(basename "$kitabsurl")
      if curl $curl_cmd_options --output /dev/null --silent --fail -r 0-0 "$DISTRIBUTION_KIT_URL"
      then
        curl $curl_cmd_options -o $kit_root/$filename "${DISTRIBUTION_KIT_URL}"
      else
        echo "Could not download distribution kit from ${DISTRIBUTION_KIT_URL}"
        exit 1
      fi
    fi
{{- end }}
{{- end }}

{{- define "kit-downloader-init-container" }}
{{- if eq (include "usesKitICDownload" .) "true" }}
- name: kit-downloader
  image: {{ .Values.distributionKit.downloadContainer.image }}
  imagePullPolicy: {{ default "IfNotPresent" .Values.distributionKit.downloadContainer.imagePullPolicy }}
  command: ['sh', '-c', '/opt/pega/dlscripts/download-kit.sh']
  env:
  - name: DISTRIBUTION_KIT_URL
    value: {{ include "distributionKitUrl" . | quote }}
  - name: ENABLE_CUSTOM_ARTIFACTORY_SSL_VERIFICATION
    value: "{{ .Values.global.customArtifactory.enableSSLVerification }}"
  volumeMounts:
  - name: {{ template "pegaDistributionKitVolume" }}
    mountPath: /opt/pega/mount/kit
  - name: kit-download-script-volume
    mountPath: /opt/pega/dlscripts
  - name: {{ .credVolumeName }}
    mountPath: "/opt/pega/secrets"
{{ if (eq (include "customArtifactorySSLVerificationEnabled" .) "true") }}
{{- if .Values.global.customArtifactory }}
{{- if .Values.global.customArtifactory.certificate }}
  - name: {{ template "pegaVolumeCustomArtifactoryCertificate" }}
    mountPath: "/opt/pega/artifactory/cert"
{{- end }}
{{- end }}
{{- end }}
{{- end }}
{{- end }}
