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

{{- define "jdbc-downloader-init-container" }}
{{- if .Values.global.downloadContainer }}
{{- if .Values.global.downloadContainer.image }}
{{- $root := . }}

{{- $urls := regexFindAll "[^,]+" .Values.global.jdbc.driverUri -1 }}
{{- range $i, $url := $urls }}
- name: jdbc-lib-downloader{{ $i }}
  image: {{ $root.Values.global.downloadContainer.image }}
  imagePullPolicy: {{ default "IfNotPresent" $root.Values.global.downloadContainer.imagePullPolicy }}
  command: ['sh', '-c', '{{- template "jdbcDownloadCmd" (dict "url" $url "root" $root) }}']
  env:
  - name: JDBC_DRIVER_URI
    value: {{ $url | quote }}
  - name: ENABLE_CUSTOM_ARTIFACTORY_SSL_VERIFICATION
    value: "{{ $root.Values.global.customArtifactory.enableSSLVerification }}"
  volumeMounts:
  - name: jdbc-lib-volume
    mountPath: /opt/pega/lib
  - name: {{ template "pegaVolumeCredentials" }}
    mountPath: "/opt/pega/secrets"
{{ if (eq (include "customArtifactorySSLVerificationEnabled" $root.root) "true") }}
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

{{- define "jdbcDownloadCmd" }}
  {{- $enableSSLVerification := "false" -}}
  {{- $customArtifactoryCert := "" }}
  {{- $customArtifactoryAuth := "" }}

  {{- $caAuthBasicName := "" -}}
  {{- $caAuthBasicPw := "" -}}
  {{- $caAuthApiHeaderName := "" -}}
  {{- $caAuthApiKey := "" -}}

  {{- if .root.Values.global.customArtifactory -}}
    {{- if .root.Values.global.customArtifactory.authentication -}}
      {{- if .root.Values.global.customArtifactory.authentication.basic -}}
        {{- $caAuthBasicName = .root.Values.global.customArtifactory.authentication.basic.username -}}
        {{- $caAuthBasicPw = .root.Values.global.customArtifactory.authentication.basic.password -}}
      {{- end -}}
      {{- if .root.Values.global.customArtifactory.authentication.apiKey -}}
        {{- $caAuthApiHeaderName = .root.Values.global.customArtifactory.authentication.apiKey.headerName -}}
        {{- $caAuthApiKey = .root.Values.global.customArtifactory.authentication.apiKey.value -}}
      {{- end -}}
    {{- end -}}

    {{- if .root.Values.global.customArtifactory.enableSSLVerification }}
      {{- $enableSSLVerification = "true" -}}
    {{- end }}

    {{- if or (ne $caAuthBasicName "") (ne $caAuthBasicPw "") -}}
      {{- if or (eq $caAuthBasicName "") (eq $caAuthBasicPw "") -}}
        {{- fail "Both username and password must be provided for basic authentication in global.customArtifactory.authentication.basic" | quote }}
      {{- else -}}
        {{- $customArtifactoryAuth := printf "-u %s:%s" $caAuthBasicName $caAuthBasicPw -}}
      {{- end -}}
    {{- end -}}
    {{- if eq $customArtifactoryAuth "" -}}
      {{- if or (ne $caAuthApiHeaderName "") (ne $caAuthApiKey "") -}}
        {{- if or (eq $caAuthApiHeaderName "") (eq $caAuthApiKey "") -}}
          {{- fail "Both header name and API key must be provided for API key authentication in global.customArtifactory.authentication.apiKey" | quote }}
        {{- else -}}
          {{- $customArtifactoryAuth = printf "-H \"%s:%s\"" $caAuthApiHeaderName $caAuthApiKey -}}
        {{- end -}}
      {{- end }}
    {{- end -}}

    {{- if .root.Values.global.customArtifactory.certificate }}
      {{- $certs := fromYaml .root.Values.global.customArtifactory.certificate -}}
      {{- $certName := "" }}
      {{- $certCount := 0 }}
      {{- range $key, $value := $certs }}
        {{- $certName = $key -}}
        {{- $certCount = add $certCount 1 -}}
      {{- end }}
      {{- if eq $certCount 1 }}
        {{- if regexMatch "^.+\\.(cer|pem|crt|der|cert|jks|p7b|p7c|key)$" $certName }}
          {{- $customArtifactoryCert = printf "--cacert /opt/pega/artifactory/cert/%s" $certName -}}
        {{- else -}}
          {{- fail "The certificate file provided in global.customArtifactory.certificate must have a valid certificate file extension such as .cer, .pem, .crt, .der, .cert, .jks, .p7b, .p7c or .key" | quote }}
        {{- end }}
      {{- end }}

    {{- end }}
  {{- end }}

  {{- $curlOptions := "" }}
  {{- if eq $enableSSLVerification "true" }}
    {{- $curlOptions = printf "-sSL --fail-with-body %s %s" $customArtifactoryAuth $customArtifactoryCert -}}
  {{- else }}
    {{- $curlOptions = printf "-ksSL --fail-with-body %s" $customArtifactoryAuth -}}
  {{- end }}

  {{- $url := .url -}}
  {{- $filename := base $url -}}

  {{- $curlCmd := printf "curl %s -o /opt/pega/lib/%s \"%s\" 2>&1 || exit $?" $curlOptions $filename $url }}
  {{- $curlCmd -}}
{{- end }}