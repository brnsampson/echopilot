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

COPY --from=0 /bin/containerpilot /bin/containerpilot
RUN chmod +x /bin/containerpilot

COPY --from=0 /go/bin/echopilot /usr/local/bin/echopilot
RUN chmod +x /usr/local/bin/echopilot

ENV ECHO_ADDR="0.0.0.0:8080"

ENTRYPOINT []

CMD ["/bin/containerpilot"]
