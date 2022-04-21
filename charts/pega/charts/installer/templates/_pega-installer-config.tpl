{{- define "pega.installer.config" -}}
{{- $arg := .mode -}}
# Node type specific configuration for {{ .name }}
kind: ConfigMap
apiVersion: v1
metadata:
  name: {{ .name }}
  namespace: {{ .root.Release.Namespace }}
data:
# Start of Pega Installer Configurations

{{ if eq $arg "installer-config" }}

{{- $prconfigTemplatePath := "config/prconfig.xml.tmpl" }}
{{- $setupDatabasePath := "config/setupDatabase.properties" }}
{{- $setupDatabasetemplatePath := "config/setupDatabase.properties.tmpl" }}
{{- $prbootstraptemplatePath := "config/prbootstrap.properties.tmpl" }}
{{- $prpcUtilsPropertiestemplatePath := "config/prpcUtils.properties.tmpl" }}
{{- $migrateSystempropertiestemplatePath := "config/migrateSystem.properties.tmpl" }}
{{- $custom_config := .root.Values.custom }}

{{- if $custom_config.configurations }}
{{ $custom_config.configurations  | toYaml | nindent 2 -}}
{{ else }}

{{ if $prconfigTemplate := .root.Files.Glob $prconfigTemplatePath }}
  # prconfigTemplate to be used by {{ .name }}
  prconfig.xml.tmpl: |-
{{ .root.Files.Get $prconfigTemplatePath | indent 6 }}
{{- end }}

{{ if $setupDatabase := .root.Files.Glob $setupDatabasePath }}
  # setupDatabase to be used by {{ .name }}
  setupDatabase.properties: |-
{{ .root.Files.Get $setupDatabasePath | indent 6 }}
{{- end }}

{{ if $setupDatabasetemplate := .root.Files.Glob $setupDatabasetemplatePath }}
  # setupDatabasetemplate to be used by {{ .name }}
  setupDatabase.properties.tmpl: |-
{{ .root.Files.Get $setupDatabasetemplatePath | indent 6 }}
{{- end }}

{{ if $prbootstraptemplate := .root.Files.Glob $prbootstraptemplatePath }}
  # prbootstraptemplate to be used by {{ .name }}
  prbootstrap.properties.tmpl: |-
{{ .root.Files.Get $prbootstraptemplatePath | indent 6 }}
{{- end }}

{{ if $prpcUtilsPropertiestemplate := .root.Files.Glob $prpcUtilsPropertiestemplatePath }}
  # prpcUtilsPropertiestemplate to be used by {{ .name }}
  prpcUtils.properties.tmpl: |-
{{ .root.Files.Get $prpcUtilsPropertiestemplatePath | indent 6 }}
{{- end }}

{{ if $migrateSystempropertiestemplate := .root.Files.Glob $migrateSystempropertiestemplatePath }}
  # migrateSystempropertiestemplate to be used by {{ .name }}
  migrateSystem.properties.tmpl: |-
{{ .root.Files.Get $migrateSystempropertiestemplatePath | indent 6 }}
{{- end }}

{{- $prlog4j2Path := "config/prlog4j2.xml" }}
  # prlog4j2 file to be used by {{ .name }}
  prlog4j2.xml: |-
{{ .root.Files.Get $prlog4j2Path | indent 6 }}
{{- end }}

{{- $dbType := .dbType }}
{{- $postgresConfPath := "config/postgres.conf" }}
{{- $oracledateConfPath := "config/oracledate.conf" }}
{{- $db2zosConfPath := "config/db2zos.conf" }}
{{- $mssqlConfPath := "config/mssql.conf" }}
{{- $udbConfPath := "config/udb.conf" }}
{{- $zosPropertiesPath := "config/DB2SiteDependent.properties" }}

{{ if and (eq $dbType "postgres") ( $postgresConf := .root.Files.Glob $postgresConfPath ) }}
  postgres.conf: |-
{{ include "customJdbcProps" .root | indent 6 }}
{{ .root.Files.Get $postgresConfPath | indent 6 }}
{{- end }}

{{ if and (eq $dbType "oracledate") ( $oracledateConf := .root.Files.Glob $oracledateConfPath ) }}
  oracledate.conf: |-
{{ include "customJdbcProps" .root | indent 6 }}
{{ .root.Files.Get $oracledateConfPath | indent 6 }}
{{- end }}

{{ if and (eq $dbType "mssql") ( $mssqlConf := .root.Files.Glob $mssqlConfPath ) }}
  mssql.conf: |-
{{ include "customJdbcProps" .root | indent 6 }}
{{ .root.Files.Get $mssqlConfPath | indent 6 }}
{{- end }}

{{ if and (eq $dbType "db2zos") ( $db2zosConf := .root.Files.Glob $db2zosConfPath ) ( $db2zosProperties := .root.Files.Glob $zosPropertiesPath ) }}
  db2zos.conf: |-
{{ include "customJdbcProps" .root | indent 6 }}
{{ .root.Files.Get $db2zosConfPath | indent 6 }}
  DB2SiteDependent.properties: |-
{{ .root.Files.Get $zosPropertiesPath | indent 6 }}
{{- end }}

{{ if and (eq $dbType "udb") ( $udbConf := .root.Files.Glob $udbConfPath ) }}
  udb.conf: |-
{{ include "customJdbcProps" .root | indent 6 }}
{{ .root.Files.Get $udbConfPath | indent 6 }}
{{- end }}

{{- end }}
# End of Pega Installer Configurations
{{- end }}

