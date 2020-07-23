FROM ubuntu:20.04 AS builder

RUN \
        apt-get update && \
        apt-get clean && \
        apt-get install -f && \
        DEBIAN_FRONTEND=noninteractive \
                apt-get install -y -f --force-yes \
        		libreoffice \
                libreofficekit-dev \
        && \
        DEBIAN_FRONTEND=noninteractive \
                apt-get install --no-install-recommends -y -f --force-yes \
        		libvips \
                libvips-dev \
                curl \
        && \
        apt-get clean && \
        rm -rf /var/lib/apt/lists/

RUN \
        curl -O https://dl.google.com/go/go1.13.6.linux-amd64.tar.gz \
        && tar -C /usr/local -xzf go1.13.6.linux-amd64.tar.gz \
        && rm -rf ./go1.13.6.linux-amd64.tar.gz \
        && export PATH=$PATH:/usr/local/go/bin

RUN \
        apt-get update && \
        apt-get clean && \
        apt-get install -f && \
        DEBIAN_FRONTEND=noninteractive \
                apt-get install --no-install-recommends -y -f --force-yes \
        		gcc \
        && \
        apt-get clean && \
        rm -rf /var/lib/apt/lists/

ENV PATH="${PATH}:/usr/local/go/bin"

WORKDIR "/app"
COPY *.go *.mod *.sum /app/
RUN go build -tags extralibs





FROM ubuntu:19.04 AS worker
RUN \
        apt-get update && \
        apt-get clean && \
        apt-get install -f && \
        DEBIAN_FRONTEND=noninteractive \
                apt-get install -y -f --force-yes \
        		libreoffice \
                libreofficekit-dev \
        && \
        DEBIAN_FRONTEND=noninteractive \
                apt-get install --no-install-recommends -y -f --force-yes \
        		libvips \
        && \
        apt-get clean && \
        rm -rf /var/lib/apt/lists/

WORKDIR "/app"
COPY fonts /app/fonts/
COPY --from=builder /app/preview .

CMD ./preview