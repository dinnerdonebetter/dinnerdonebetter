FROM otel/opentelemetry-collector-contrib:0.43.0

COPY environments/dev/config_files/opentelemetry/config.yaml /etc/config.yaml

EXPOSE 4317

CMD ["--config=/etc/config.yaml"]