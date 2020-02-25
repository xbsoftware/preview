FROM ubuntu:19.04

RUN \
        apt-get update && \
        apt-get clean && \
        apt-get install -f && \
        DEBIAN_FRONTEND=noninteractive \
                apt-get install -y -f --force-yes \
        		libreoffice \
                libreofficekit-dev \
        && \
        apt-get clean && \
        rm -rf /var/lib/apt/lists/

RUN \
        apt-get update && \
        apt-get clean && \
        apt-get install -f && \
        DEBIAN_FRONTEND=noninteractive \
                apt-get install --no-install-recommends -y -f --force-yes \
        		libvips \
                libvips-dev \
        && \
        apt-get clean && \
        rm -rf /var/lib/apt/lists/

RUN \
        apt-get update && \
        apt-get clean && \
        apt-get install -f && \
        DEBIAN_FRONTEND=noninteractive \
                apt-get install --no-install-recommends -y -f --force-yes \
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
COPY fonts /app/fonts/

RUN go build -tags extralibs
CMD ./preview