FROM bufbuild/buf:1.15.0 as bufbuild

FROM golang:1.20.1-bullseye as gobuild

WORKDIR /go/src/github.com/brnsampson/echopilot

COPY . .

COPY --from=bufbuild /usr/local/bin/buf /usr/local/bin/buf

RUN apt update && apt install -y protobuf-compiler gzip

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest \
    && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest \
    && go install github.com/bufbuild/connect-go/cmd/protoc-gen-connect-go@latest

RUN  buf generate --config proto/buf.yaml --template proto/buf.gen.yaml proto/

RUN mkdir /web \
    && cp web/* /web

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


FROM debian:bullseye

# COPY ContainerPilot configuration
ENV CONTAINERPILOT_PATH=/etc/containerpilot.json5
COPY ./etc/containerpilot.json5 ${CONTAINERPILOT_PATH}
ENV CONTAINERPILOT=${CONTAINERPILOT_PATH}

COPY --from=gobuild /web ./web

COPY --from=gobuild /bin/containerpilot /bin/containerpilot
RUN chmod +x /bin/containerpilot

COPY --from=gobuild /go/bin/echopilot /usr/local/bin/echopilot
RUN chmod +x /usr/local/bin/echopilot

RUN mkdir -p /etc/echopilot
COPY ./etc/tls /etc/echopilot/tls

ENV ECHOPILOT_TLS_CERT="/etc/echopilot/tls/cert.pem"
ENV ECHOPILOT_TLS_KEY="/etc/echopilot/tls/key.pem"
ENV ECHOPILOT_TLS_SKIP_VERIFY="false"
ENV ECHOPILOT_TLS_ENABLED="false"
ENV ECHOPILOT_HOST="localhost"
ENV ECHOPILOT_BIND_IP="0.0.0.0"
ENV ECHOPILOT_PORT=3000
ENV ECHOPILOT_CLIENT_SKIP_VERIFY="true"

ENTRYPOINT []

CMD ["/bin/containerpilot"]
