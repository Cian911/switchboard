FROM bitnami/minideb:jessie

ENV SWB_VERSION=0.2.4

RUN apt-get update && apt-get install -y curl

RUN curl -L https://github.com/Cian911/switchboard/releases/download/v${SWB_VERSION}/switchboard_${SWB_VERSION}_linux_64-bit.deb -o switchboard_${SWB_VERSION}_linux_64-bit.deb && \
    dpkg -i switchboard_${SWB_VERSION}_linux_64-bit.deb && \
    rm switchboard_${SWB_VERSION}_linux_64-bit.deb

ENTRYPOINT ["switchboard"]
