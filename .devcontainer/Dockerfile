FROM mcr.microsoft.com/devcontainers/go:1.25

# Installation des outils de build
RUN echo 'deb [trusted=yes] https://repo.goreleaser.com/apt/ /' | tee /etc/apt/sources.list.d/goreleaser.list \
    && apt-get update && apt-get install goreleaser \
    && rm -f /etc/apt/sources.list.d/yarn.list \
    && apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get clean -y && rm -rf /var/lib/apt/lists/*
