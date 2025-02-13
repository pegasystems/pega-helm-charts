default: secrets

NAMESPACE := <YOUR_NAMESPACE>
CLUSTERING_SERVICE_IMAGE := <CLUSTERING_SERVICE_IMAGE>
ENC_KEYSTORE_PASSWORD := <ENC_KEYSTORE_PASSWORD>
ENC_TRUSTSTORE_PASSWORD := <ENC_TRUSTSTORE_PASSWORD>
HIGHLY_SECURE_CRYPTO_MODE_ENABLED := <HIGHLY_SECURE_CRYPTO_MODE_ENABLED>
ALIAS := myalias

ifeq ($(HIGHLY_SECURE_CRYPTO_MODE_ENABLED), true)
secrets:
	docker run --name hazelcast-helm-charts-certs -i -w /tmp \
    		$(CLUSTERING_SERVICE_IMAGE) \
            		/bin/sh -c " \
            		         ./certs.sh HIGHLY_SECURE_CRYPTO_MODE_ENABLED $(ENC_KEYSTORE_PASSWORD) cluster-keystore $(ENC_TRUSTSTORE_PASSWORD) cluster-truststore $(ALIAS)" &&\
                             docker cp hazelcast-helm-charts-certs:/tmp/cluster-keystore.jks ./ &&\
                             docker cp hazelcast-helm-charts-certs:/tmp/cluster-truststore.jks ./ &&\
                             docker rm -f hazelcast-helm-charts-certs &&\
                             kubectl create secret generic hz-encryption-secrets --from-literal=HZ_SSL_KEYSTORE_PASSWORD=$(ENC_KEYSTORE_PASSWORD) --from-literal=HZ_SSL_TRUSTSTORE_PASSWORD=$(ENC_TRUSTSTORE_PASSWORD) --from-file=cluster-keystore.jks --from-file=cluster-truststore.jks --namespace=$(NAMESPACE)
else
secrets:
	docker run --name hazelcast-helm-charts-certs -i -w /tmp \
    		$(CLUSTERING_SERVICE_IMAGE) \
            		/bin/sh -c " \
            		         ./certs.sh SSL_MODE_ENABLED $(ENC_KEYSTORE_PASSWORD) cluster-keystore $(ENC_TRUSTSTORE_PASSWORD) cluster-truststore $(ALIAS)" && \
            		docker cp hazelcast-helm-charts-certs:/tmp/cluster-keystore.jks ./ &&\
            		docker cp hazelcast-helm-charts-certs:/tmp/cluster-truststore.jks ./ &&\
                    docker rm -f hazelcast-helm-charts-certs &&\
                    kubectl create secret generic hz-encryption-secrets --from-literal=HZ_SSL_KEYSTORE_PASSWORD=$(ENC_KEYSTORE_PASSWORD) --from-literal=HZ_SSL_TRUSTSTORE_PASSWORD=$(ENC_TRUSTSTORE_PASSWORD) --from-file=cluster-keystore.jks --from-file=cluster-truststore.jks --namespace=$(NAMESPACE)

endif