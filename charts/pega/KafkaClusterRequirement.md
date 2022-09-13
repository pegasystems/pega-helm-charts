## Kafka cluster requirements
Configure your own managed Kafka infrastructure as per the below details.

### Deployment
See [Kafka Helm charts](https://github.com/bitnami/charts/tree/master/bitnami/kafka) for deployments

#### Version
Pega recommends Apache Kafka versions 2.3.1 or later (Verified version 3.2.1)

#### Configuration

Deployment Type | CPU     | Memory | Disk Space | Replicas
---         | ---     | ---    | ---        | ---
Development | 2 cores | 8Gi    | 100G*      | At least 2
Production  | 4 cores | 16Gi   | 200G*      | At least 3

* Disk Space depends on the required throughput (Number and size of messages) and retention period.
* In order to enable compression, it is enough to set `compression.type` in your kafka configuration.
* Above configuration would support up to 1000 kafka partitions, increase resources accordingly if kafka partitions are more.
* Define appropriate quotas on network bandwidth and request rate if you want to share your kafka cluster across different environments.

##### Miscellaneous configuration
* message.max.bytes=5000000
* unclean.leader.election.enable=false
* auto.create.topics.enable=false

For best practices, see [this page.](https://docs.pega.com/decision-management/87/best-practices-stream-service-configuration)

#### Security

Pega supports SSL for encryption of traffic as well as authentication to communicate with external Kafka. 

In order to secure, mount necessary certificates(trustStore and keyStore) during the Pega Platform deployment. For details, see [this section.](README.md#optional-support-for-providing-credentialscertificates-using-external-secrets-operator)

You may also securely pass settings like trustStorePassword,keyStorePassword, and jaasConfig through a secret in an external secret operator. For details, see [this section.](README.md#optional-support-for-providing-credentialscertificates-using-external-secrets-operator)

### Kafka Permissions
If external kafka is configured with authentication and authorization through Kafka Access control lists then Pega requires following user permissions

Principal |Resource Type  | Resource Name     | Operation | Permission Type | Patter Type
---         |---         | ---     | ---    | ---        | ---
User:<user-name> | TOPIC       | <Prefix> as in 'stream.streamNamePattern' | ALL    | ALLOW      | PREFIXED
User:<user-name> |TRANSACTIONAL_ID  | * | READ/WRITE   | ALLOW      | LITERAL
User:<user-name> |GROUP  | * | ALL   | ALLOW      | LITERAL
User:<user-name> |CLUSTER  | <Cluster-Name> | IDEMPOTENT_WRITE   |ALLOW      | LITERAL



