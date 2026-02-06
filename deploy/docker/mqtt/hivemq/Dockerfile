FROM hivemq/hivemq4:latest

USER root

# Download and install HiveMQ Prometheus extension
RUN apt-get update && \
    apt-get install -y wget unzip && \
    wget -q https://github.com/hivemq/hivemq-prometheus-extension/releases/download/4.1.0/hivemq-prometheus-extension-4.1.0.zip && \
    unzip -q hivemq-prometheus-extension-4.1.0.zip -d /opt/hivemq/extensions/ && \
    rm hivemq-prometheus-extension-4.1.0.zip && \
    apt-get remove -y wget unzip && \
    apt-get autoremove -y && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

USER 10000

# Extension will be auto-discovered by HiveMQ
