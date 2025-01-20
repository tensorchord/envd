FROM curlimages/curl:8.11.1 as builder
USER root
RUN curl --proto '=https' --tlsv1.2 -sSf https://starship.rs/install.sh | sh -s -- -y

FROM scratch as prod
COPY --from=builder /usr/local/bin/starship /usr/local/bin/starship
