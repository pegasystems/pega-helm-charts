{{- define "pegaLibDownloadScriptConfig" }}
{{- $depName := printf "%s" (include "deploymentName" $) -}}
{{- $depName -}}-{{- .chartType -}}-lib-download-script-config
{{- end }}

{{- define "jdbcLibVolume" }}
{{- if .Values.global.downloadContainer }}
{{- if .Values.global.downloadContainer.image }}
- name: jdbc-lib-volume
  emptyDir:
    sizeLimit: 5Mi
{{- end }}
{{- end }}
{{- end }}

{{- define "jdbcLibVolumeMount" }}
{{- if .Values.global.downloadContainer }}
{{- if .Values.global.downloadContainer.image }}
- name: jdbc-lib-volume
  mountPath: /opt/pega/lib
{{- end }}
{{- end }}
{{- end }}

{{- define "downloadScriptVolume" }}
{{- if .Values.global.downloadContainer }}
{{- if .Values.global.downloadContainer.image }}
- name: download-script-volume
  configMap:
    name: {{ template "pegaLibDownloadScriptConfig" $ }}
    defaultMode: 0550
{{- end }}
{{- end }}
{{- end }}

{{- define "jdbc-downloader-init-container" }}
{{- if .Values.global.downloadContainer }}
{{- if .Values.global.downloadContainer.image }}
- name: jdbc-lib-downloader
  image: {{ .Values.global.downloadContainer.image }}
  imagePullPolicy: {{ default "IfNotPresent" .Values.global.downloadContainer.imagePullPolicy }}
  command: ['sh', '-c', '/opt/pega/dlscripts/download-jdbc-lib.sh']
  env:
  - name: JDBC_DRIVER_URI
    value: {{ .Values.global.jdbc.driverUri | quote }}
  - name: ENABLE_CUSTOM_ARTIFACTORY_SSL_VERIFICATION
    value: "{{ .Values.global.customArtifactory.enableSSLVerification }}"
  volumeMounts:
  - name: jdbc-lib-volume
    mountPath: /opt/pega/lib
  - name: download-script-volume
    mountPath: /opt/pega/dlscripts
  - name: {{ .credVolumeName }}
    mountPath: "/opt/pega/secrets"
{{ if (eq (include "customArtifactorySSLVerificationEnabled" .root) "true") }}
{{- if .Values.global.customArtifactory }}
{{- if .Values.global.customArtifactory.certificate }}
  - name: {{ template "pegaVolumeCustomArtifactoryCertificate" }}
    mountPath: "/opt/pega/artifactory/cert"
{{- end }}
{{- end }}
{{- end }}
{{- end }}
{{- end }}
{{- end }}

{{- define "usesICDownload" }}
{{- $usesICDownload := "false" -}}
{{- if .Values.global.downloadContainer }}
{{- if .Values.global.downloadContainer.image }}
{{- $usesICDownload = "true" }}
{{- end }}
{{- end }}
{{- $usesICDownload -}}
{{- end }}

{{- define "downloadScriptConfigMap" }}
{{- if .Values.global.downloadContainer }}
{{- if .Values.global.downloadContainer.image }}
{{- $depName := printf "%s" (include "deploymentName" $) }}
---
# This template contains the script to download the JDBC driver from the specified location and copy it to the shared
# volume. This script will be executed by the init container before starting the Pega application container.
kind: ConfigMap
apiVersion: v1
metadata:
  name: {{ template "pegaLibDownloadScriptConfig" . }}
  namespace: {{ .Release.Namespace }}
data:
  download-jdbc-lib.sh: |
    #/bin/sh

    pega_root=/opt/pega
    lib_root=$pega_root/lib
    art_root=$pega_root/artifactory
    secret_root = $pega_root/secrets

    export CAU="$(base64 -d $secret_root/CUSTOM_ARTIFACTORY_USERNAME)"
    export CAP="$(base64 -d $secret_root/CUSTOM_ARTIFACTORY_PASSWORD)"
    export CAAH="$(base64 -d $secret_root/CUSTOM_ARTIFACTORY_APIKEY_HEADER)"
    export CAAK="$(base64 -d $secret_root/CUSTOM_ARTIFACTORY_APIKEY)"

    ca_auth=""
    if [ "$CAU" != "" ] || [ "$CAP" != "" ]; then
      if [ "$CAU" == "" ] || [ "$CAP" == "" ]; then
        echo "CUSTOM_ARTIFACTORY_USERNAME & CUSTOM_ARTIFACTORY_PASSWORD must be specified for artifactory basic auth."
        exit 1
      else
        echo "Using basic authentication for custom artifactory to download JDBC driver."
        ca_auth="-u "$CAU":"$CAP
      fi
    fi

    if [ "$ca_auth" == "" ]; then
      if [[ "$CAAH" != "" || "$CAAK" != "" ]]; then
        if [ "$CAAH" == "" ] || [ "$CAAK" == "" ]; then
          echo "CUSTOM_ARTIFACTORY_APIKEY_HEADER & CUSTOM_ARTIFACTORY_APIKEY must be specified for authentication using api key for custom artifactory."
          exit 1
        else
          echo "Using API key for artifactory authentication."
          ca_auth="-H $CAAH:$CAAK"
        fi
      fi
    fi

    ca_cert=""
    if [ "$(ls -A $art_root/cert/*)" ]; then
      if [ "$(ls $art_root/cert/* | wc -l)" ]; then
        echo "Certificate is provided for custom artifactory's domain ssl verification."
        certfilename="$(ls $art_root/cert)"
        ext="${certfilename##*.}"
        regex="^(cer|pem|crt|der|cert|jks|p7b|p7c|key)$"
        if [[ "$ext" =~ $regex ]]; then
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

    if [ "$JDBC_DRIVER_URI" != "" ]; then
      curl_cmd_options=""
      if [ "$ENABLE_CUSTOM_ARTIFACTORY_SSL_VERIFICATION" == true ]; then
        echo "Establishing a secure connection to download driver."
        curl_cmd_options="-sSL $custom_artifactory_auth $custom_artifactory_certificate"
      else
        echo "Establishing an insecure connection to download driver."
        curl_cmd_options="-ksSL $custom_artifactory_auth"
      fi
      urls=$(echo "$JDBC_DRIVER_URI" | tr "," "\n")
        for url in $urls
        do
         echo "Downloading database driver: ${url}";
         # jarabsurl="$(cut -d'?' -f1 <<<"$url")"
         jarabsurl="$(echo "$url" | cut -d'?' -f1)"
         echo "$jarabsurl"
         filename=$(basename "$jarabsurl")
         if curl $curl_cmd_options --output /dev/null --silent --fail -r 0-0 "$url"
         then
           curl $curl_cmd_options -o $lib_root/$filename "${url}"
         else
           echo "Could not download jar from ${url}"
           exit 1
         fi
        done
    fi

{{- end }}
{{- end }}
{{- end }}