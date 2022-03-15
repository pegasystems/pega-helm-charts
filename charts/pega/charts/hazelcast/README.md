# Clustering Service Deployment (Draft)

The values.yaml in this hazelcast sub chart describes the additional deployment parameters used in the client server type deployment of Pega Platform with Hazelcast. 
In client server arrangement, Pega Platform acts as the client for the hazelcast server running on Kubernetes Pods. Pega highly recommends client server
form of deployment from Pega 8.6 onwards.


#### Managing Resources
You can optionally configure the memory and cpu limits for the hazelcast service using the below parameters. Default value is used in case an alternate value is not
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
`client.clusterName` | Cluster Name for Pega nodes to connect to for Client Server Deployment | `PRPC`
 

#### Server Parameters

Name                                           | Description                                           | Default value |
------------------------------------------------|-------------------------------------------------------|---|
`server.java_opts`                              | Parameter  for passing JVM Arguments| `-Xms820m -Xmx820m -XX:+HeapDumpOnOutOfMemoryError -XX:HeapDumpPath=/opt/hazelcast/logs/heapdump.hprof` |                                                                                          
`server.jmx_enabled` | Enables exposing the collected metrics over JMX if set to true, disables it otherwise | `true`
`server.health_monitoring_level` | Health monitoring log level. When SILENT, logs are printed only when values exceed some predefined threshold. When NOISY, logs are always printed periodically. Set OFF to turn off completely. | `OFF`
`server.operation_generic_thread_count` | Number of generic operation handler threads for each Hazelcast member. Its default value is the maximum of 2 and processor count / 2. | `""` 
`server.operation_thread_count`  | Number of partition based operation handler threads for each Hazelcast member. Its default value is the maximum of 2 and count of available processors. |  `""`
`server.io_thread_count` | Number of threads performing socket input and socket output. If, for example, the default value (3) is used, it means there are 3 threads performing input and 3 threads performing output (6 threads in total). | `""`
`server.event_thread_count` | Number of event handler threads | `""`
`server.max_join_seconds`  | Join timeout, maximum time to try to join before giving. | `""`
`server.group_name`  | Specifies the name of the cluster created by hazelcast nodes |  `PRPC`
`server.mancenter_url`  | URL of the Hazelcast Management center to which the hazelcast nodes can connect | `""`
`server.graceful_shutdown_max_wait_seconds` | Maximum wait in seconds during graceful shutdown. | `600`
`server.service_dns_timeout` | Custom time for how long the DNS Lookup is checked | `""`
`server.logging_level` | Can be used to set logging level for Hazelcast, available logging levels are OFF, FATAL, ERROR, WARN, INFO, DEBUG, TRACE and ALL. Invalid levels are assumed to be OFF| `info`
`server.diagnostics_enabled` | 	Specifies whether diagnostics tool is enabled or not for the cluster. | `true`
`server.diagnostics_metric_level` | - | `info`
`server.diagnostic_log_file_size_mb` | Size of each diagnostic file to be rolled. | `50`
`server.diagnostics_file_count` | Number of diagnostic log files to be rolled.   | `3` 


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