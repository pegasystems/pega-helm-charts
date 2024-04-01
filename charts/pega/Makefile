default: test

test:secrets

NAMESPACE := <YOUR_NAMESPACE>
FIPS_ENABLED := <FIPS_MODE>

CLUSTERING_SERVICE_IMAGE := cloudservices-docker-dev-local.bin.pega.io/platform/clustering-service:2.0.0-SNAPSHOT
ENC_KEYSTORE_PASSWORD := mystorePwd
ENC_TRUSTSTORE_PASSWORD := mystorePwd
ENCRYPTION_KEYSTORE_NAME := cluster-keystore
ENCRYPTION_TRUSTSTORE_NAME := cluster-truststore
ALIAS := myalias

ifeq ($(FIPS_ENABLED), true)
secrets:
	docker run --name hazelcast-helm-charts-certs -i -w /tmp \
    		$(CLUSTERING_SERVICE_IMAGE) \
            		/bin/sh -c " \
            		         keytool -genkeypair -alias $(ALIAS) -storetype BCFKS -keyalg RSA -keysize 2048 -validity 365 -keypass $(ENC_KEYSTORE_PASSWORD) -keystore $(ENCRYPTION_KEYSTORE_NAME).BCFKS -storepass $(ENC_TRUSTSTORE_PASSWORD) -providerclass org.bouncycastle.jcajce.provider.BouncyCastleFipsProvider -providerpath '/opt/hazelcast/lib/bc-fips-1.0.2.4.jar' && \
            		         keytool -export -alias $(ALIAS) -storetype BCFKS -file server.crt -keystore $(ENCRYPTION_KEYSTORE_NAME).BCFKS -storepass $(ENC_TRUSTSTORE_PASSWORD) -providerclass sun.security.pkcs11.SunPKCS11 -providerclass org.bouncycastle.jcajce.provider.BouncyCastleFipsProvider -providerpath '/opt/hazelcast/lib/bc-fips-1.0.2.4.jar' &&\
            		         keytool -import -alias $(ALIAS) -storetype BCFKS -file server.crt -keystore $(ENCRYPTION_TRUSTSTORE_NAME).BCFKS -storepass $(ENC_TRUSTSTORE_PASSWORD) -providerclass sun.security.pkcs11.SunPKCS11 -providerclass org.bouncycastle.jcajce.provider.BouncyCastleFipsProvider -providerpath '/opt/hazelcast/lib/bc-fips-1.0.2.4.jar'" &&\
                             docker cp hazelcast-helm-charts-certs:/tmp/$(ENCRYPTION_KEYSTORE_NAME).BCFKS ./ &&\
                             docker cp hazelcast-helm-charts-certs:/tmp/$(ENCRYPTION_TRUSTSTORE_NAME).BCFKS ./ &&\
                             kubectl create secret generic hz-encryption-secrets --from-literal=HZ_ENC_KEYSTORE_PASSWORD=$(ENC_KEYSTORE_PASSWORD) --from-literal=HZ_ENC_TRUSTSTORE_PASSWORD=$(ENC_TRUSTSTORE_PASSWORD) --from-file=$(ENCRYPTION_KEYSTORE_NAME).BCFKS --from-file=$(ENCRYPTION_TRUSTSTORE_NAME).BCFKS --namespace=$(NAMESPACE) &&\
                             rm -f $(ENCRYPTION_KEYSTORE_NAME).BCFKS  $(ENCRYPTION_TRUSTSTORE_NAME).BCFKS
else
secrets:
	docker run --name hazelcast-helm-charts-certs -i -w /tmp \
    		$(CLUSTERING_SERVICE_IMAGE) \
            		/bin/sh -c " \
            		         keytool -genkey -alias $(ALIAS) -keyalg RSA -keysize 2048 -keypass $(ENC_KEYSTORE_PASSWORD) -keystore $(ENCRYPTION_KEYSTORE_NAME).jks -validity 9999 -storepass $(ENC_TRUSTSTORE_PASSWORD) && \
                             keytool -export -alias $(ALIAS) -file server.crt -keystore $(ENCRYPTION_KEYSTORE_NAME).jks -validity 9999 -storepass $(ENC_TRUSTSTORE_PASSWORD) && \
                             keytool -import -alias $(ALIAS) -file server.crt -keystore $(ENCRYPTION_TRUSTSTORE_NAME).jks -storepass $(ENC_TRUSTSTORE_PASSWORD)" && \
            		docker cp hazelcast-helm-charts-certs:/tmp/$(ENCRYPTION_KEYSTORE_NAME).jks ./ &&\
            		docker cp hazelcast-helm-charts-certs:/tmp/$(ENCRYPTION_TRUSTSTORE_NAME).jks ./ &&\
                    docker rm -f hazelcast-helm-charts-certs &&\
                    kubectl create secret generic hz-encryption-secrets --from-literal=HZ_ENC_KEYSTORE_PASSWORD=$(ENC_KEYSTORE_PASSWORD) --from-literal=HZ_ENC_TRUSTSTORE_PASSWORD=$(ENC_TRUSTSTORE_PASSWORD) --from-file=$(ENCRYPTION_KEYSTORE_NAME).jks --from-file=$(ENCRYPTION_TRUSTSTORE_NAME).jks --namespace=$(NAMESPACE) &&\
                   rm -f $(ENCRYPTION_KEYSTORE_NAME).jks  $(ENCRYPTION_TRUSTSTORE_NAME).jks

endif