Garbage collection logs and heap dumps are written to the pod storage which is ephemeral. When a pod crashes, all the logs and heap dumps stored in the pod are lost. To persist them we can use the following options:

### Persist garbage collection logs

Garbage collection logs are like diagnostic tools for application's memory management. 

#### 1. To enable GC logging: 
Set the appropriate flag (mentioned [here](../charts/pega/RecommendedJVMArgs.md)) within the `javaOpts` parameter in `values.yaml` file.

#### 2. Persist logs with EFK Stack:
Configure the EFK stack (Elasticsearch-Fluentd-Kibana) to collect and store these GC logs permanently. Refer to the Pega documentation for EFK setup: [EFK logging stack](../charts/addons/README.md#logging-with-elasticsearch-fluentd-kibana-efk).

### Persist heap dumps

When your application runs out of memory (OOM), capturing a heap dump can help in identifying the root cause

#### 1. To enable automatic heap dumps:
Set the appropriate flag (mentioned [here](../charts/pega/RecommendedJVMArgs.md)) within the `javaOpts` parameter in `values.yaml` file.

#### 2. Persist heap dumps using persistent volumes:
Configure a persistent storage using Persistent Volumes (PVs) and Persistent Volume Claims (PVCs). Refer this [document](https://kubernetes.io/docs/concepts/storage/persistent-volumes/) to create and use persistent volumes.

Include the following configuration snippet for each tier in `values.yaml` file.

Please note that the heap dumps are stored in a predefined location `/heapdumps`. Use the same location as the mount path in your `values.yaml`.

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
