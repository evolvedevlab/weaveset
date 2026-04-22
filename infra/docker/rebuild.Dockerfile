FROM alpine:latest

RUN apk add --no-cache \
    ca-certificates \
    curl \
    tar \
    bash \
    libc6-compat \
    libstdc++ \
    libgcc

# hugo version to download
ENV HUGO_VERSION=latest

RUN set -eux; \
    ARCH=$(uname -m); \
    case "$ARCH" in \
        x86_64) ARCH=64bit ;; \
        aarch64) ARCH=ARM64 ;; \
        *) echo "Unsupported arch: $ARCH"; exit 1 ;; \
    esac; \
    \
    if [ "$HUGO_VERSION" = "latest" ]; then \
        HUGO_VERSION=$(curl -s https://api.github.com/repos/gohugoio/hugo/releases/latest | grep tag_name | cut -d '"' -f 4); \
    fi; \
    \
    curl -L -o hugo.tar.gz \
        "https://github.com/gohugoio/hugo/releases/download/${HUGO_VERSION}/hugo_extended_${HUGO_VERSION#v}_Linux-${ARCH}.tar.gz"; \
    \
    tar -xzf hugo.tar.gz; \
    mv hugo /usr/local/bin/hugo; \
    chmod +x /usr/local/bin/hugo; \
    rm -f hugo.tar.gz LICENSE README.md

# Create application directory
RUN mkdir -p /app /app/site

WORKDIR /app

COPY ./script/rebuild.sh /app/script/rebuild.sh

# Default trigger file
ENV FILEPATH=site/content/list/.changed

RUN chmod +x /app/script/rebuild.sh

CMD ["/bin/bash", "/app/script/rebuild.sh"]
