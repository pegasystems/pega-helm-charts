### Embedded Stream with latest helm chart version
Starting from Infinity 24.2, support for embedded Stream is removed. As a best practice, update your Stream configuration to an external Kafka service.
To configure embedded Stream in Pega Platform ’24.1 and earlier using the Pega Helm chart version 3.25, perform the following steps.

#### Configure values.yaml
1. Add Stream tier details in the values.yaml file under Pega tiers section.
    #### Example for values.yaml and values-large.yaml 
    ```
    - name: "stream"
   # Create a stream tier for queue processing.  This tier deploys
   # as a stateful set to ensure durability of queued data. It may
   # be optionally exposed to the load balancer.
   # Note: Stream tier is deprecated. As a best practice, enable externalized Kafka service configuration under External Services.
   # When externalized Kafka service is enabled, remove the entire stream tier.
   nodeType: "Stream"

   # Pega requestor specific properties
   requestor:
   # Inactivity time after which requestor is passivated
   passivationTimeSec: 900

   service:
   port: 7003
   targetPort: 7003

   # If a nodeSelector is required for this or any tier, it may be specified here:
   # nodeSelector:
   #  disktype: ssd

   ingress:
   enabled: true
   # Enter the domain name to access web nodes via a load balancer.
   #  e.g. web.mypega.example.com
   domain: "YOUR_STREAM_NODE_DOMAIN"
   tls:
   # Enable TLS encryption
   enabled: true
   # secretName:
   # useManagedCertificate: false
   # ssl_annotation:

   livenessProbe:
   port: 8081

   # To configure an alternative user for your custom image, set value for runAsUser
   # To configure an alternative group for volume mounts, set value for fsGroup
   # See, https://github.com/pegasystems/pega-helm-charts/blob/master/charts/pega/README.md#security-context
   # securityContext:
   #   runAsUser: 9001
   #   fsGroup: 0

   # To specify security settings for a Container, include the securityContext field in the Container manifest
   # Security settings that you specify for a Container apply only to the pega container,
   # and they override settings made at the Pod level when there is overlap. Container settings
   # do not affect the Pod's Volumes.
   # See, https://kubernetes.io/docs/tasks/configure-pod-container/security-context/#set-the-security-context-for-a-container
   # containerSecurityContext:
   #   capabilities:
   #     add: ["SYS_TIME"]

   replicas: 2

   volumeClaimTemplate:
   resources:
   requests:
   storage: 5Gi

   # Set enabled to true to include a Pod Disruption Budget for this tier.
   # To enable this budget, specifiy either a pdb.minAvailable or pdb.maxUnavailable
   # value and comment out the other parameter.
   pdb:
   enabled: false
   minAvailable: 1
   # maxUnavailable: "50%"

   resources:
   requests:
   memory: "12Gi"
   cpu: 3
   limits:
   memory: "12Gi"
   cpu: 4
    ```
    
    #### Example for values-minimal.yaml
   For values-minimal.yaml, add the Stream nodeType to the minikube tier.
    ```
   # Specify the Pega tiers to deploy
    # For a minimal deployment, use a single tier to reduce resource consumption.
    # Note: Stream tier is deprecated. As a best practice, enable externalized Kafka service configuration under External Services.
    # configuration under External Services
    tier:
   - name: "minikube"
   nodeType: "Stream,BackgroundProcessing,WebUser,Search"

         service:
           httpEnabled: true
           port: 80
           targetPort: 8080
           # Without a load balancer, use a direct NodePort instead.
           serviceType: "NodePort"
           # To configure TLS between the ingress/load balancer and the backend, set the following:
           tls:
             enabled: false
             # To avoid entering the certificate values in plain text, configure the keystore, keystorepassword, cacertificate parameter
             # values in the External Secrets Manager, and enter the external secret name below
             # make sure the keys in the secret should be TOMCAT_KEYSTORE_CONTENT, TOMCAT_KEYSTORE_PASSWORD and ca.crt respectively
             external_secret_name: ""
             keystore:
             keystorepassword:
             port: 443
             targetPort: 8443
             # set the value of CA certificate here in case of baremetal/openshift deployments - CA certificate should be in base64 format
             # pass the certificateChainFile file if you are using certificateFile and certificateKeyFile
             cacertificate:
             # provide the SSL certificate and private key as a PEM format
             certificateFile:
             certificateKeyFile:
             # if you will deploy traefik addon chart and enable traefik, set enabled=true; otherwise leave the default setting.
             traefik:
               enabled: false
               # the SAN of the certificate present inside the container
               serverName: ""
               # set insecureSkipVerify=true, if the certificate verification has to be skipped
               insecureSkipVerify: false
      ```


2. Disable external Kafka service settings in the values.yaml file.
    ```
   # Stream (externalized Kafka service) settings.
    stream:
    # Beginning with Pega Platform '23, enabled by default; when disabled, your deployment does not use a"Kafka stream service" configuration.
    enabled: false
    # Provide externalized Kafka service broker urls.
    bootstrapServer: ""
    # Provide Security Protocol used to communicate with kafka brokers. Supported values are: PLAINTEXT, SSL, SASL_PLAINTEXT, SASL_SSL.
    securityProtocol: PLAINTEXT
    # If required, provide trustStore certificate file name
    # When using a trustStore certificate, you must also include a Kubernetes secret name, that contains the trustStore certificate,
    # in the global.certificatesSecrets parameter.
    # Pega deployments only support trustStores using the Java Key Store (.jks) format.
    trustStore: ""
    # If required provide trustStorePassword value in plain text.
    trustStorePassword: ""
    # If required, provide keyStore certificate file name
    # When using a keyStore certificate, you must also include a Kubernetes secret name, that contains the keyStore certificate,
    # in the global.certificatesSecrets parameter.
    # Pega deployments only support keyStores using the Java Key Store (.jks) format.
    keyStore: ""
    # If required, provide keyStore value in plain text.
    keyStorePassword: ""
    # If required, provide jaasConfig value in plain text.
    jaasConfig: ""
    # If required, provide a SASL mechanism**. Supported values are: PLAIN, SCRAM-SHA-256, SCRAM-SHA-512.
    saslMechanism: PLAIN
    # By default, topics originating from Pega Platform have the pega- prefix,
    # so that it is easy to distinguish them from topics created by other applications.
    # Pega supports customizing the name pattern for your Externalized Kafka configuration for each deployment.
    streamNamePattern: "pega-{stream.name}"
    # Your replicationFactor value cannot be more than the number of Kafka brokers. Pega recommended value is 3.
    replicationFactor: "3"
    # To avoid exposing trustStorePassword, keyStorePassword, and jaasConfig parameters, leave the values empty and
    # configure them using an External Secrets Manager, making sure you configure the keys in the secret in the order:
    # STREAM_TRUSTSTORE_PASSWORD, STREAM_KEYSTORE_PASSWORD and STREAM_JAAS_CONFIG.
    # Enter the external secret name below.
    external_secret_name: ""
   ```
