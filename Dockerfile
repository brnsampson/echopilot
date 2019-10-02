FROM golang:1.13

WORKDIR /go/src/github.com/brnsampson/echopilot

COPY . .

RUN go get -d -v ./...
RUN go install -v ./...


FROM python:3

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

# COPY ContainerPilot configuration
ENV CONTAINERPILOT_PATH=/etc/containerpilot.json5
COPY containerpilot.json5 ${CONTAINERPILOT_PATH}
ENV CONTAINERPILOT=${CONTAINERPILOT_PATH}

# Currently the prestart and prestop scripts are python, so we have to do this
RUN apt-get update && apt-get install -y libsystemd-dev
COPY requirements.txt ./
RUN pip install --no-cache-dir -r requirements.txt

COPY ./bin/* /usr/local/bin/
RUN chmod +x /usr/local/bin/*

COPY --from=0 /go/bin/echopilot /usr/local/bin/echopilot
RUN chmod +x /usr/local/bin/echopilot

ENV ECHO_ADDR="0.0.0.0:8080"

ENTRYPOINT []

CMD ["/bin/containerpilot"]
