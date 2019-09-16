## Configuration

The following table lists the configurable parameters of the Pega chart and their default values.

| Tier properties                                   | Default               | Description                                                                                                                                                         |
| --------------------------------------------------| ----------------------| --------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `global.tier.cpuRequest`                          | `200m`                | CPU request for each web node                                                                                                                                       |
| `global.tier.memRequest`                          | `"6Gi"`               | Memory request for each web node                                                                                                                                    |
| `global.tier.cpuLimit`                            | `2`                   | CPU limit for each web node                                                                                                                                         |
| `global.tier.memLimit`                            | `"8Gi"`               | Memory limit for each web node                                                                                                                                      |
| `global.tier.initialHeap`                         | `"4096m"`             | Initial heap size for the JVM                                                                                                                                       |
| `global.tier.maxHeap`                             | `"7168m"`             | Maximum heap size for the JVM                                                                                                                                       | 
| `global.tier.hpa.minReplicas`                     | `1`                   | Minimum number of replicas that HPA can scale-down                                                                                                                  |
| `global.tier.hpa.maxReplicas`                     | `5`                   | Maximum number of replicas that HPA can scale-up                                                                                                                    |
| `global.tier.hpa.targetAverageCPUUtilization`     | `700`                 |threshold value for average cpu utilization percentage (Recommended value is 700% of 200m ). HPA will scale up if average of all web nodes/pods cpu reaches 1.4c     |
| `global.tier.hpa.targetAverageMemoryUtilization`  | `85`                  |threshold value for average memory utilization percentage (Recommended value is 85% of 6Gi ).HPA will scale up if average of all web nodes/pods memory reaches 5.1Gi |


## Examples

Please refer [examples folder](/examples) to configure multiple tiers or to override default values while deploying Pega chart.