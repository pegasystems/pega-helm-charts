By default, the system saves garbage collection logs and heap dumps in the pod storage, which is ephemeral. When a pod crashes, all logs and heap dumps stored in the pod are lost. You can use one of the following options to persist the data:

### Persist garbage collection logs

Garbage collection (GC) logs are diagnostic tools for application memory management.

1. To enable GC logging, set the appropriate flag within the javaOpts parameter in values.yaml. For more information, see [JVM Arguments for better tunability](../charts/pega/RecommendedJVMArgs.md).

2. To collect and store GC logs permanently, configure the EFK stack (Elasticsearch-Fluentd-Kibana). For more information, see [EFK logging stack](../charts/addons/README.md#logging-with-elasticsearch-fluentd-kibana-efk).

### Persist heap dumps

When your application runs out of memory (OOM), capturing a heap dump can help in identifying the root cause.

1. To enable automatic heap dumps, set the appropriate flag within the `javaOpts` parameter in `values.yaml` file. For more information, see [JVM Arguments for better tunability](../charts/pega/RecommendedJVMArgs.md).

2. To persist heap dumps, configure a persistent storage using Persistent Volumes (PVs) and Persistent Volume Claims (PVCs). For more information on persistent volumes, see official [Kubernetes documentation](https://kubernetes.io/docs/concepts/storage/persistent-volumes/).

Include the following configuration snippet for each tier in `values.yaml` file.

Note: The system stores the heap dumps in a predefined location `/heapdumps`. Use the same location as the mount path in your `values.yaml` file.

```yaml
tier:
  - name: my-tier
    custom:
      volumes:
        - name: <volume name>
          persistentVolumeClaim:
            claimName: <persistent volume claim name>
      volumeMounts:
        - mountPath: /heapdumps
          name: <volume name>
```
