{{- define "pega.installer.environment.config" -}}
data:
  # Temporary password for administrator@pega.com that is used to install Pega Platform
  ADMIN_PASSWORD: {{ .Values.adminPassword }}
  # Database Type for installation
  DB_TYPE: {{ .Values.global.jdbc.dbType }}
  # JDBC URL of the DB where Pega is installed
  JDBC_URL: {{ .Values.global.jdbc.url }}
  # Class name of the DB's JDBC driver
  JDBC_CLASS: {{ .Values.global.jdbc.driverClass }}
{{- if .Values.global.jdbc.driverUri }}
  # URI that the JDBC driver can be downloaded from
  JDBC_DRIVER_URI: {{ .Values.global.jdbc.driverUri }}
{{- end }}
  # Rules schema of the Pega installation
  RULES_SCHEMA: {{ .Values.global.jdbc.rulesSchema }}
  # Data schema of the Pega installation
  DATA_SCHEMA: {{ .Values.global.jdbc.dataSchema }}
{{- if .Values.global.jdbc.customerDataSchema }}
  # CustomerData schema of the Pega installation
  CUSTOMERDATA_SCHEMA: {{ .Values.global.jdbc.customerDataSchema }}
{{- end }}
{{- if .Values.global.jdbc.connectionProperties }}
  # JDBC custom connection properties
  JDBC_CUSTOM_CONNECTION: {{ .Values.global.jdbc.connectionProperties }}
{{- end }}
{{- if .Values.multitenant_system }}
  # Whether this is a Multitenant System ('true' if yes, 'false' if no)
  MT_SYSTEM: {{ .Values.multitenant_system | quote}}
{{- end }}
  # UDF generation will be skipped if this property is set to true
  BYPASS_UDF_GENERATION: {{ .Values.bypassUdfGeneration | quote}}
  # Z/OS SITE-SPECIFIC PROPERTIES FILE
  ZOS_PROPERTIES: {{ .Values.zos.zosProperties }}
{{- if .Values.zos.db2zosUdfWlm }}
  # Specify the workload manager to load UDFs into db2zos
  DB2ZOS_UDF_WLM: {{ .Values.zos.db2zosUdfWlm}}
{{- end }}
{{- if .Values.distributionKitURL }}
  # Distribution kit URL
  DISTRIBUTION_KIT_URL: {{ .Values.distributionKitURL }}
{{- end }}
{{- if .Values.bypassLoadEngineClasses }}
  # Bypass loading engine classes into database during installation
  BYPASS_LOAD_ENGINE_CLASSES: {{ .Values.bypassLoadEngineClasses | quote }}
{{- end }}
{{- if .Values.bypassLoadAssembledClasses }}
  # Bypass loading assembly classes into database during installation
  BYPASS_LOAD_ASSEMBLED_CLASSES: {{ .Values.bypassLoadAssembledClasses | quote }}
{{- end }}
  # Enable ssl verification for jdbc driver download
  ENABLE_CUSTOM_ARTIFACTORY_SSL_VERIFICATION: {{ .Values.global.customArtifactory.enableSSLVerification | quote }}
{{- if .Values.advancedSettings }}
  ADVANCED_SETTINGS: |-
{{- range .Values.advancedSettings }}
    {{ .name }}={{ .value }}
{{- end }}
{{ end }}
{{- end -}}