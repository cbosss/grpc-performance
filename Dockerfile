FROM golang:1.16-stretch AS proto

WORKDIR /proto

RUN apt-get update && \
    apt-get -y --no-install-recommends install unzip

# download pre-compiled protoc, validate checksum of binary
ENV PROTOC_VERSION=3.15.5
RUN curl -sL https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/protoc-${PROTOC_VERSION}-linux-x86_64.zip > tmp.zip && \
    echo "7358bf688ddad8d6ba430240b44e644c88c4678a21221987cd8d2a0dbf119beb  tmp.zip" | sha256sum --check && \
    unzip -jq tmp.zip bin/protoc -d /usr/local/bin && \
    unzip -q tmp.zip include/* -d /usr/local && \
    rm tmp.zip

# install compiler plugin (protoc-gen-go)
ADD go.mod .
ADD go.sum .

# uses version from go.mod
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
