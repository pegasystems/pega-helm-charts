{{- define "pega.config" -}}
{{- $arg := .mode -}}
# Node type specific configuration for {{ .name }}
kind: ConfigMap
apiVersion: v1
metadata:
  name: {{ .name }}
  namespace: {{ .root.Release.Namespace }}
data:

# Start of Pega Deployment Configuration

{{ if eq $arg "deploy-config" }}

{{- $prconfigPath := "config/deploy/prconfig.xml" }}
{{- $contextXMLTemplate := "config/deploy/context.xml.tmpl" }}
{{- $prlog4j2Path := "config/deploy/prlog4j2.xml" }}

{{- if .custom }}
{{- if .custom.prconfig }}
 # CUSTOM prconfig file to be used by {{ .name }}
  prconfig.xml: |-
{{ .custom.prconfig | indent 6 }}
{{ else if $prconfig := .root.Files.Glob $prconfigPath }}
 # prconfig file to be used by {{ .name }}
  prconfig.xml: |-
{{ .root.Files.Get $prconfigPath | indent 6 }}
{{- end }}
{{ else if $prconfig := .root.Files.Glob $prconfigPath }}
 # prconfig file to be used by {{ .name }}
  prconfig.xml: |-
{{ .root.Files.Get $prconfigPath | indent 6 }}
{{- end }}

{{ if $contextXML := .root.Files.Glob $contextXMLTemplate }}
  # contextXMLTemplate to be used by {{ .name }}
  context.xml.tmpl: |-
{{ .root.Files.Get $contextXMLTemplate | indent 6 }}
{{- end }}

{{- if .custom }}
{{- if .custom.context }}
 # CUSTOM context file to be used by {{ .name }}
  context.xml: |-
{{ .custom.context | indent 6 }}
{{- end }}
{{- end }}

  # prlog4j2 file to be used by {{ .name }}
  prlog4j2.xml: |-
{{ .root.Files.Get $prlog4j2Path | indent 6 }}

{{- end }}
# End of Pega Deployment Configuration

# Start of Pega Installer Configurations

{{ if eq $arg "installer-config" }}

{{- $prconfigTemplatePath := "config/installer/prconfig.xml.tmpl" }}
{{- $setupDatabasePath := "config/installer/setupDatabase.properties" }}
{{- $setupDatabasetemplatePath := "config/installer/setupDatabase.properties.tmpl" }}
{{- $prbootstraptemplatePath := "config/installer/prbootstrap.properties.tmpl" }}

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

{{- $prlog4j2Path := "config/installer/prlog4j2.xml" }}
  # prlog4j2 file to be used by {{ .name }}
  prlog4j2.xml: |-
{{ .root.Files.Get $prlog4j2Path | indent 6 }}

{{- $dbType := .dbType }}
{{- $postgresConfPath := "config/installer/postgres/postgres.conf" }}
{{- $oracledateConfPath := "config/installer/oracledate/oracledate.conf" }}
{{- $db2zosConfPath := "config/installer/db2zos/db2zos.conf" }}
{{- $mssqlConfPath := "config/installer/mssql/mssql.conf" }}
{{- $udbConfPath := "config/installer/udb/udb.conf" }}
{{- $zosPropertiesPath := "config/installer/db2zos/DB2SiteDependent.properties" }}

{{ if and (eq $dbType "postgres") ( $postgresConf := .root.Files.Glob $postgresConfPath ) }}
  postgres.conf: |-
{{ .root.Files.Get $postgresConfPath | indent 6 }}
{{- end }}

{{ if and (eq $dbType "oracledate") ( $oracledateConf := .root.Files.Glob $oracledateConfPath ) }}
  oracledate.conf: |-
{{ .root.Files.Get $oracledateConfPath | indent 6 }}
{{- end }}

{{ if and (eq $dbType "mssql") ( $mssqlConf := .root.Files.Glob $mssqlConfPath ) }}
  mssql.conf: |-
{{ .root.Files.Get $mssqlConfPath | indent 6 }}
{{- end }}

{{ if and (eq $dbType "db2zos") ( $db2zosConf := .root.Files.Glob $db2zosConfPath ) ( $db2zosProperties := .root.Files.Glob $zosPropertiesPath ) }}
  db2zos.conf: |-
{{ .root.Files.Get $db2zosConfPath | indent 6 }}
  DB2SiteDependent.properties: |-
{{ .root.Files.Get $zosPropertiesPath | indent 6 }}
{{- end }}

{{ if and (eq $dbType "udb") ( $udbConf := .root.Files.Glob $udbConfPath ) }}
  udb.conf: |-
{{ .root.Files.Get $udbConfPath | indent 6 }}
{{- end }}

{{- end }}
# End of Pega Installer Configurations
{{- end }}

