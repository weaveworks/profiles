# Build the manager binary
FROM golang:1.16-alpine as builder

RUN apk add --no-cache gcc pkgconfig libc-dev binutils-gold musl~=1.2 libgit2-dev~=1.1

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY main.go main.go
COPY api/ api/
COPY controllers/ controllers/
COPY pkg/ pkg/

# This is automatically set
ARG TARGETARCH

# Build
RUN CGO_ENABLED=1 GOOS=linux GOARCH=${TARGETARCH} GO111MODULE=on go build -a -o manager main.go

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM alpine:3.13

RUN apk add --no-cache ca-certificates tini libgit2~=1.1 musl~=1.2

WORKDIR /
COPY --from=builder /workspace/manager .
USER 65532:65532

ENTRYPOINT ["/manager"]
