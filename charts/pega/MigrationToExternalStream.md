### Switch from embedded Stream to externalized Kafka service
Beginning in 8.8, Pega is deprecating the use of a dedicated stream tier in Pega Platform deployments. Pega recommends you update your data streaming configuration in your deployments to use an externalized Kafka service. Use this article to migrate an embedded stream configuration to an externalized Kafka configuration in an existing deployment. This process requires that you use Pega Helm chart versions 2.2 or later so your deployment rely on configurations exclusively set in Pega-provided Helms and not in Pega Platform configuration files such as prconfig.

Kafka stores data on the filesystem. In deployments using embedded stream tier, the stream tier stores data on the volumes attached to the stream pods; in deployments using an externalized Kafka service, your application also stores data on the filesystem, but the storage configuration will be different from the volumes attached to the stream pods. Because of this difference, existing stream data will not be migrated.

Pega supports two processes to migrate a previously-existing Pega Platform deployment infrastructure using embedded stream nodes one using a Helm-chart-based, externalized Kafka configuration:
#### A non-production migration that could involve data loss
#### A migration with minimal or no data loss. Pega requires that clients migrating a production deployment use this process
See the appropriate section below for details. During a deployment update, you must not update Pega Platform software.

#### A non-production migration that could involve data loss
1. Edit pega chart
   
   1.1 Remove stream tier.
   
   1.2 Configure and enable externalized Kafka service.

2. Update deployment.
   
   2.1 Invoke the update process by using the `helm upgrade release --namespace mypega` command.

3. After restarts, new pods will connect to externalized Kafka service.

#### A migration with minimal or no data loss. Pega requires that clients migrating a production deployment use this process
1. Edit pega chart, 
   
   1.1 Set replica count to 0 for all the tiers producing stream data for e.g. Web.
   
   1.2 If a tier is producing as well as consuming then stop producer process.
   
2. Update deployment.
   
   2.1 Invoke the Update process by using the `helm upgrade release --namespace mypega` command.

3. Wait for all the consuming tiers for e.g. BackgroundProcessing, Batch, etc to process remaining stream data.
   
   3.1 To check the status
   
      3.1.1 Goto Admin Studio Page and check count against 'Ready to Process' for all the Queue Processors. Good to proceed further once count is zero for all the QPs.
   
      3.1.2 For other Data Flows that consume stream data, for each data flow
   1. Goto Data Flow Run, Open AllPartitionReport, check the Last IDs of all the partitions.
   2. Login to stream pod, goto to directory '/usr/local/tomcat/<embedded-kafka-version>/bin'
   3. Run command ./kafka-run-class.sh kafka.tools.GetOffsetShell --broker-list <IP of broker from the stream landing page>:9092 --topic <Topic Name>  --time -2
   4. Compare Last IDs from step ii. with offset numbers for each partition, Good to proceed further once all the offsets and IDs match

4. Edit pega chart again
   
   4.1 Restore replica count for all the producing tiers.
   
   4.2 Remove stream tier
   
   4.3 Configure and enable externalized Kafka service.

5. Update deployment.
   
   5.1 Invoke the Update process by using the `helm upgrade release --namespace mypega` command.

6. After restarts, both producers and consumers tiers will connect to externalized Kafka service.

7. If required, start producer processes stopped in step 1.2