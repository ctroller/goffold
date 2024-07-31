# Stage 1: Node for prettier in VSCode
FROM node:22-bookworm AS node

# Stage 2: Install Tailwind
FROM debian:bookworm-slim AS tailwind
RUN apt-get update && apt-get install -y curl && \
    curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/download/v3.4.6/tailwindcss-linux-x64 && \
    chmod +x tailwindcss-linux-x64

# Stage 3: Install act
FROM debian:bookworm-slim AS act
RUN apt-get update && apt-get install -y curl sudo && \
    curl --proto '=https' --tlsv1.2 -sSf https://raw.githubusercontent.com/nektos/act/master/install.sh | sudo bash -s -- -b /usr/local/bin

# Stage 4: Install Go tools
FROM golang:1.22.5-bookworm AS go-tools
RUN go install golang.org/x/tools/gopls@latest && \
    go install github.com/go-delve/delve/cmd/dlv@latest && \
    go install golang.org/x/lint/golint@latest && \
    go install honnef.co/go/tools/cmd/staticcheck@latest && \
    go install golang.org/x/tools/cmd/goimports@latest && \
    go install github.com/golang/mock/mockgen@latest && \
    go install golang.org/x/tools/cmd/godoc@latest && \
    go install golang.org/x/perf/cmd/benchstat@latest

# Final stage: Main image
FROM mcr.microsoft.com/devcontainers/go:1-1.22-bookworm AS main

# Copy Prettier from node stage
COPY --from=node /usr/local /usr/local

# Copy Tailwind from tailwind stage
COPY --from=tailwind /tailwindcss-linux-x64 /usr/local/bin/tailwindcss

# Copy act from act stage
COPY --from=act /usr/local/bin/act /usr/local/bin/act

# Copy Go tools from go-tools stage
COPY --from=go-tools /go/bin/* /usr/local/bin/

# Set environment variables for Go
ENV CGO_ENABLED=0

# Update and install dependencies
RUN apt-get update && apt-get upgrade -y && \
    apt-get install -y git curl && \
    # Clean up
    apt-get autoremove -y && apt-get clean -y && \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

# Set the working directory
WORKDIR /workspace

# Make prettier work in vscode with go templates. Absolute disgusting hackfix, but alas
# https://discourse.gohugo.io/t/formatter-for-go-templates/38403
RUN npm install -g prettier prettier-plugin-go-template

USER vscode

# Keep the container running
CMD ["sleep", "infinity"]