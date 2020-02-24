FROM golang:1.13

WORKDIR /go/src/github.com/brnsampson/echopilot

COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

ENV CONTAINERPILOT_VERSION 3.8.0

RUN curl -Lso /tmp/containerpilot.sha1.txt \
         "https://github.com/joyent/containerpilot/releases/download/${CONTAINERPILOT_VERSION}/containerpilot-${CONTAINERPILOT_VERSION}.sha1.txt" \
    && export CP_SHA1=$( cat /tmp/containerpilot.sha1.txt | grep containerpilot-${CONTAINERPILOT_VERSION}.tar.gz | awk '{print $1}' ) \
    && curl -Lso /tmp/containerpilot.tar.gz \
         "https://github.com/joyent/containerpilot/releases/download/${CONTAINERPILOT_VERSION}/containerpilot-${CONTAINERPILOT_VERSION}.tar.gz" \
    && echo "${CP_SHA1}  /tmp/containerpilot.tar.gz" | sha1sum -c \
    && tar zxf /tmp/containerpilot.tar.gz -C /bin \
    && rm /tmp/containerpilot.tar.gz \
    && rm /tmp/containerpilot.sha1.txt

FROM phusion/baseimage:latest

# COPY ContainerPilot configuration
ENV CONTAINERPILOT_PATH=/etc/containerpilot.json5
COPY containerpilot.json5 ${CONTAINERPILOT_PATH}
ENV CONTAINERPILOT=${CONTAINERPILOT_PATH}

COPY ./dist ./dist

COPY --from=0 /bin/containerpilot /bin/containerpilot
RUN chmod +x /bin/containerpilot

COPY --from=0 /go/bin/echopilot /usr/local/bin/echopilot
RUN chmod +x /usr/local/bin/echopilot

ENV ECHO_SERVER_CERT="/home/vagrant/cert.pem"
ENV ECHO_SERVER_KEY="/home/vagrant/key.pem"
ENV ECHO_GATEWAY_SKIP_VERIFY="false"
ENV ECHO_GRPC_ADDR="127.0.0.1:8080"
ENV ECHO_REST_ADDR="127.0.0.1:3000"
ENV ECHO_CLIENT_SKIP_VERIFY="true"

ENTRYPOINT []

CMD ["/bin/containerpilot"]
