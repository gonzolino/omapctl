FROM golang:1.26.4-trixie AS builder

RUN apt-get update && apt-get install -y --no-install-recommends \
    golang \
    gcc \
    libcephfs-dev \
    librbd-dev \
    librados-dev \
    libc-bin \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 go build -o /bin/omapctl .

FROM ubuntu:26.04

RUN apt-get update && apt-get install -y --no-install-recommends \
    ceph-common \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /bin/omapctl /usr/local/bin/omapctl

ENTRYPOINT ["omapctl"]
