ARG ENVD_VERSION

FROM ghcr.io/tensorchord/envd-from-scratch:${ENVD_VERSION} as envd

FROM moby/buildkit:v0.25.1-rootless
COPY --from=envd /usr/bin/envd /usr/bin/envd
COPY scripts/envd-daemonless.sh /envd-daemonless.sh

CMD [ "/envd-daemonless.sh" ]
