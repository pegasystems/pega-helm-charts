## Kafka cluster requirements
Configure your own managed Kafka infrastructure as per the below details.

### Deployment
Pega supports Client-managed cloud clients to configure an externalized Kafka configuration that connects to your organization's Kafka service infrastructure. These configuration options support both enterprise grade Kafka services offered by leading public cloud vendors or your a Kafka infrastructure that you manage across your enterprise.

Pega Platform deployments using Pega-provided Helm charts starting at version 2.2 or later provide Pega Helm chart settings that allow you configure the connection and authentication details required by your organization's Kafka service infrastructure. These latest, Kafka-specific Pega Helm chart enhancements provide a scalable Kafka configuration for your Pega applications running in your preferred Kubernetes environment while offering great flexibility in connecting to a Kafka service infrastructure using your company's preferred streaming policy and security profiles. To manage your externalized Kafka configuration in your deployment see [Kafka Helm charts](https://github.com/bitnami/charts/tree/master/bitnami/kafka).

#### Version
Pega recommends Apache Kafka versions 2.3.1 or later (Verified version 3.2.1)

### Configuration

Deployment Type | CPU     | Memory | Disk Space | Replicas
---         | ---     | ---    | ---        | ---
Development | 2 cores | 8Gi    | 100G*      | At least 2
Production  | 4 cores | 16Gi   | 200G*      | At least 3

* Disk Space depends on the required throughput (Number and size of messages) and retention period.
* In order to enable compression, it is enough to set `compression.type` in your kafka configuration.
* The above configuration can easily support up to 1000 kafka partitions; you can increase resources accordingly if your deployment requires more kafka partitions.
* Define appropriate quotas on network bandwidth and request rate if you want to share your kafka cluster across different environments.

#### Miscellaneous configuration
* message.max.bytes=5000000 
  This is the default maximum message size supported, if you want to increase this value then pass the following jvm arguments to pega tiers as well
  -Dstream.producer.max.request.size=<Max-message-size> -Dstream.producer.buffer.memory=<Max-message-size>
  See more details about JVM arguments [here](README.md#jvm-arguments)
* unclean.leader.election.enable=false
* auto.create.topics.enable=false

For best practices, see [this page.](https://docs.pega.com/decision-management/87/best-practices-stream-service-configuration)

### Security

Pega supports SSL for network traffic encryption an authentication for communicate with your organization's existing Kafka service. 

In order to secure, mount necessary certificates(trustStore and keyStore) during the Pega Platform deployment. For details, see [this section.](README.md#optional-support-for-providing-credentialscertificates-using-external-secrets-operator)

You may also securely pass settings like trustStorePassword,keyStorePassword, and jaasConfig through a secret in an external secret operator. For details, see [this section.](README.md#optional-support-for-providing-credentialscertificates-using-external-secrets-operator)

#### Permissions
To configure an externalized Kafka service connection using authentication and authorization profiles in Kafka Access control lists, your Pega profiles require following user permissions. To review configuration details, see [Kafka documentation for Authorization and ACLs](https://kafka.apache.org/documentation/#security_authz).

Principal |Resource Type  | Resource Name     | Operation | Permission Type | Patter Type
---         |---         | ---     | ---    | ---        | ---
User:\<user-name\> | TOPIC       | \<Prefix\> as in 'stream.streamNamePattern' | ALL    | ALLOW      | PREFIXED
User:\<user-name\> |TRANSACTIONAL_ID  | * | READ/WRITE   | ALLOW      | LITERAL
User:\<user-name\> |GROUP  | * | ALL   | ALLOW      | LITERAL
User:\<user-name\> |CLUSTER  | \<Cluster-Name\> | IDEMPOTENT_WRITE   |ALLOW      | LITERAL



