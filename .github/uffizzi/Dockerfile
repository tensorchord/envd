FROM uffizzi/ttyd:golang1.18-alpine

COPY envd /usr/bin/envd

WORKDIR /

ENTRYPOINT ["tini", "--"]
CMD     ["ttyd", "/bin/zsh"]