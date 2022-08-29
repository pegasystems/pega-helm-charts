default: es-prerequisite

NAMESPACE := <YOUR_NAMESPACE>
ELASTICSEARCH_IMAGE := docker.elastic.co/elasticsearch/elasticsearch:$(ELASTICSEARCH_VERSION)

es-prerequisite: es-secrets

es-secrets:
	docker run --name elastic-helm-charts-certs -i -w /tmp \
		$(ELASTICSEARCH_IMAGE) \
		/bin/sh -c " \
			elasticsearch-certutil ca --out /tmp/elastic-stack-ca.p12 --pass '' && \
			elasticsearch-certutil cert --name security-master --dns security-master --ca /tmp/elastic-stack-ca.p12 --pass '' --ca-pass '' --out /tmp/elastic-certificates.p12" && \
	docker cp elastic-helm-charts-certs:/tmp/elastic-certificates.p12 ./ && \
	docker rm -f elastic-helm-charts-certs && \
	openssl pkcs12 -nodes -passin pass:'' -in elastic-certificates.p12 -out elastic-certificate.pem && \
	openssl x509 -outform der -in elastic-certificate.pem -out elastic-certificate.crt && \
	kubectl create secret generic elastic-certificates --from-file=elastic-certificates.p12 --namespace=$(NAMESPACE) && \
	kubectl create secret generic elastic-certificate-pem --from-file=elastic-certificate.pem --namespace=$(NAMESPACE) && \
	kubectl create secret generic elastic-certificate-crt --from-file=elastic-certificate.crt --namespace=$(NAMESPACE) && \
	kubectl create secret generic srs-certificates --from-file=elastic-certificates.p12 --namespace=$(NAMESPACE)
	rm -f elastic-certificates.p12 elastic-certificate.pem elastic-certificate.crt elastic-stack-ca.p12

purge-es-secrets:
	kubectl delete secrets elastic-certificates elastic-certificate-pem elastic-certificate-crt --namespace=$(NAMESPACE) || true

external-es-secrets:
	kubectl create secret generic srs-certificates --from-file=$(PATH_TO_CERTIFICATE) --namespace=$(NAMESPACE)
