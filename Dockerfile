FROM python:3

WORKDIR /usr/src/app

RUN apt-get update && apt-get install -y libsystemd-dev

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

COPY requirements.txt ./
RUN pip install --no-cache-dir -r requirements.txt

# COPY ContainerPilot configuration
ENV CONTAINERPILOT_PATH=/etc/containerpilot.json5
COPY containerpilot.json5 ${CONTAINERPILOT_PATH}
ENV CONTAINERPILOT=${CONTAINERPILOT_PATH}

COPY ./bin/* /usr/local/bin/
RUN chmod +x /usr/local/bin/* 


ENTRYPOINT []

CMD ["/bin/containerpilot"]
