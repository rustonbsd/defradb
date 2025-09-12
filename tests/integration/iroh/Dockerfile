FROM defra:latest AS defrabin

FROM docker.io/golang:1.24 AS runtime

WORKDIR /app

RUN  echo "DEFRA_KEYRING_SECRET=$(head -c 32 /dev/urandom | base32)" > /app/secret.env

COPY start_script_* /usr/local/bin/

COPY --from=defrabin /defradb /usr/local/bin/defradb
COPY --from=defrabin /lib/libwasmer.so /lib/libwasmer.so

RUN apt-get update && apt-get install -y jq

RUN chmod +x /usr/local/bin/defradb \
    /usr/local/bin/start_script_a.sh \
    /usr/local/bin/start_script_b.sh

FROM runtime AS nodea
ENTRYPOINT ["/usr/local/bin/start_script_a.sh"]

FROM runtime AS nodeb
ENTRYPOINT ["/usr/local/bin/start_script_b.sh"]