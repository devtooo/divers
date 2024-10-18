apiVersion: apps/v1
kind: Deployment
metadata:
  name: kafka-cli
  namespace: <namespace>
spec:
  replicas: 1
  selector:
    matchLabels:
      app: kafka-cli
  template:
    metadata:
      labels:
        app: kafka-cli
    spec:
      containers:
      - name: kafka-cli
        image: confluentinc/cp-kafka:latest
        command: ["/bin/bash", "-c", "sleep 3600"]  # Keep the pod running for manual CLI commands
        volumeMounts:
        - name: kafka-certificates
          mountPath: /opt/kafka/certs  # Path to certificates inside the container
        env:
        - name: KAFKA_OPTS
          value: "-Djavax.net.ssl.trustStore=/opt/kafka/certs/ca.crt -Djavax.net.ssl.keyStore=/opt/kafka/certs/client.crt -Djavax.net.ssl.keyPassword=/opt/kafka/certs/client.key"
      volumes:
      - name: kafka-certificates
        secret:
          secretName: kafka-client-certificates
