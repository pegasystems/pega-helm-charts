### Persist garbage collection logs

Garbage collection logs are like diagnostic tools for application's memory management. When a pod crashes, all the logs stored within the pod are lost.

#### Enable GC logging: 
Set the appropriate flag (mentioned [here](https://github.com/pegasystems/pega-helm-charts/blob/master/charts/pega/RecommendedJVMArgs.md)) within the `javaOpts` parameter in `values.yaml` file.

#### Persist logs with EFK Stack:
Configure the EFK stack (Elasticsearch-Fluentd-Kibana) to collect and store these GC logs permanently. Refer to the Pega documentation for EFK setup: [EFK logging stack](https://github.com/pegasystems/pega-helm-charts/blob/master/charts/addons/README.md#logging-with-elasticsearch-fluentd-kibana-efk).

### Persist heap dumps

When your application runs out of memory (OOM), capturing a heap dump can help in identifying the root cause

#### Enable Automatic Heap Dumps:
Set the appropriate flag (mentioned [here](https://github.com/pegasystems/pega-helm-charts/blob/master/charts/pega/RecommendedJVMArgs.md)) within the `javaOpts` parameter in `values.yaml` file.

#### Persist heap dumps using persistent volumes:
Configure a persistent storage using Persistent Volumes (PVs) and Persistent Volume Claims (PVCs). Refer this [document](https://kubernetes.io/docs/concepts/storage/persistent-volumes/) to create and use persistent volumes.

Include the following configuration snippet for each tier in `values.yaml` file.

```yaml
tier:
  - name: my-tier
    custom:
      volumes:
        - name: <volume name>
          persistentVolumeClaim:
            climName: <persistent volume claim name>
      volumeMounts:
        - mountPath: <mount path>
          name: <volume name>
```
