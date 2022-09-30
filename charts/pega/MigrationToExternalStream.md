### Switch from embedded Stream to External Kafka
Pega recommends deployment of stream with external kafka. This means existing deployment may switch from embedded stream to external kafka.

Kafka stores data on the filesystem.
Stream tier stores data on the volumes attached to the stream pods.
External kafka also stores data on the filesystem, but their storage will be different from the volumes attached to the stream pods.
Hence, existing stream data will not be migrated.

The switch can be done in two ways.
Switch requires a deployment upgrade, please do not perform Pega Platform upgrade during the switch.

#### Switch non-production environments with a potential data loss.
1. Edit pega chart
   
   1.1 Remove stream tier.
   
   1.2 Configure and enable external kafka.

2. Upgrade deployment.
   
   2.1 Invoke the upgrade process by using the `helm upgrade release --namespace mypega` command.

3. After restarts, new pods will connect to external kafka.

#### Switch with minimal or no data loss, recommended for production environments.
1. Edit pega chart, 
   
   1.1 Set replica count to 0 for all the tiers producing stream data for e.g. Web.
   
   1.2 If a tier is producing as well as consuming then stop producer process.
   
2. Upgrade deployment.
   
   2.1 Invoke the upgrade process by using the `helm upgrade release --namespace mypega` command.

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
   
   4.3 Configure and enable external kafka.

5. Upgrade deployment.
   
   5.1 Invoke the upgrade process by using the `helm upgrade release --namespace mypega` command.

6. After restarts, both producers and consumers tiers will connect to external kafka.

7. If required, start producer processes stopped in step 1.2