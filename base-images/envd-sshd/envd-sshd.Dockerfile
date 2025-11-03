FROM scratch
ARG TARGETPLATFORM
COPY $TARGETPLATFORM/envd-sshd /usr/bin/envd-sshd
