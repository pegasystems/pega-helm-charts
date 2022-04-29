# Clustering Service Deployment

The Hazelcast clustering service is implemented using a Pega-provided Docker image that contains the Hazelcast clustering service image, which you must make available in your deployment Docker image repository.
This Hazelcast subchart defines the additional deployment parameters the clustering service uses to define dynamic features for the client-server deployment model using Hazelcast.
In the latter deployment model, Pega Platform uses Hazelcast in Client configuration to connect to a cluster of PODs running Hazelcast in server configuration. 
**For Kubernetes based deployments client-server is the default and the recommended form of deployment.**

**Using Clustering service for client-server form of deployment is only supported from Pega Platform 8.6 or later.**

#### Managing Resources
You can configure the memory and cpu limits for the Hazelcast service using the below parameters. Default value is used in case an alternate value is not
specified.

Name                                           | Description                                           | Default value |
------------------------------------------------|-------------------------------------------------------|---|
`resources.requests.memory` | Initial Memory request for PODS running the Hazelcast service | `1Gi` |
`resources.requests.cpu` | Initial CPU request for PODS running the Hazelcast service | `1`|
`resources.limits.memory` | Memory Limit for PODS running the Hazelcast service | `1Gi`|
`resources.limits.cpu` | CPU Limit for PODS running the Hazelcast service | `2`|

#### Client Parameter


Name                                           | Description                                           | Default value |
------------------------------------------------|-------------------------------------------------------|---|
`client.clusterName` | Cluster Name for Pega nodes to connect to for client-server Deployment | `PRPC`
 

#### Server Parameters

Name                                           | Description                                           | Default value |
------------------------------------------------|-------------------------------------------------------|---|
`server.java_opts`                              | Parameter  for passing JVM Arguments| `-Xms820m -Xmx820m -XX:+HeapDumpOnOutOfMemoryError -XX:HeapDumpPath=/opt/hazelcast/logs/heapdump.hprof` |                                                                                          
`server.jmx_enabled` | Enables exposing the collected metrics over JMX if set to true, disables it otherwise | `true`
`server.health_monitoring_level` | Health monitoring log level. When SILENT, logs are printed only when values exceed some predefined threshold. When NOISY, logs are always printed periodically. Set OFF to turn off completely. | `OFF`
`server.operation_generic_thread_count` | Number of generic operation handler threads for each Hazelcast member. The default value is the maximum of 2 and (available processor count / 2). | `""` 
`server.operation_thread_count`  | Number of partition based operation handler threads for each Hazelcast member. The default value is the maximum of 2 and the number of available processors. |  `""`
`server.io_thread_count` | Number of threads performing socket input and socket output. For example, If a default value of (3) is used, it implies there are 3 threads performing input and 3 threads performing output (6 threads in total). | `""`
`server.event_thread_count` | Number of event handler threads | `""`
`server.max_join_seconds`  | Join timeout, maximum time to try to join before giving. | `""`
`server.group_name`  | Specifies the name of the cluster created by Hazelcast nodes |  `PRPC`
`server.mancenter_url`  | URL of the Hazelcast Management center to which the Hazelcast nodes can connect | `""`
`server.graceful_shutdown_max_wait_seconds` | Maximum wait in seconds during graceful shutdown. | `600`
`server.service_dns_timeout` | Custom time for how long the DNS Lookup is checked | `""`
`server.logging_level` | Can be used to set logging level for Hazelcast, available logging levels are OFF, FATAL, ERROR, WARN, INFO, DEBUG, TRACE and ALL. Invalid levels are assumed to be OFF| `info`
`server.diagnostics_enabled` | 	Specifies whether diagnostics tool is enabled or not for the cluster. | `true`
`server.diagnostics_metric_level` | The level of information saved in the log file is set to info by default. To change the level of information saved in the log file, change the value of the setting to the level that you want to use. | `info`
`server.diagnostic_log_file_size_mb` | The maximum size of each diagnostic file, default value used is 50 MB | `50`
`server.diagnostics_file_count` | The maximum number of diagnostic files that the system keeps. Default value is 3  | `3` 


#### Example

```yaml
image: "YOUR_HAZELCAST_IMAGE:TAG"
imagePullPolicy: "Always"
replicas: 3
enabled: true
username: ""
password: ""

resources:
  requests:
    memory: "1Gi"
    cpu: 1
  limits:
    memory: "1Gi"
    cpu: 2
    
client:
  clusterName: "PRPC"
  
server:
  java_opts: "-Xms820m -Xmx820m -XX:+HeapDumpOnOutOfMemoryError -XX:HeapDumpPath=/opt/hazelcast/logs/heapdump.hprof
  -XX:+UseParallelGC -Xlog:gc*,gc+phases=debug:file=/opt/hazelcast/logs/gc.log:time,pid,tags:filecount=5,filesize=3m"
  jmx_enabled: "true"
  health_monitoring_level: "OFF"
  operation_generic_thread_count: ""
  operation_thread_count: ""
  io_thread_count: ""
  event_thread_count: ""
  max_join_seconds: ""
  group_name: "PRPC"
  mancenter_url: ""
  graceful_shutdown_max_wait_seconds: "600"
  service_dns_timeout: ""
  logging_level: "info"
  diagnostics_enabled: "true"
  diagnostics_metric_level: "info"
  diagnostic_log_file_size_mb: "50"
  diagnostics_file_count: "3"
```