# syntax=docker/dockerfile:1
FROM golang:1.26-bookworm AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 \
go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /bin/cone \
    ./cmd/cone

FROM debian:bookworm-slim

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
    git ca-certificates \
    && rm -rf /var/lib/apt/lists/*

RUN useradd -r -u 1000 -m -d /git-cone gitcone

COPY --from=builder /bin/cone /usr/local/bin/cone
RUN ln -s /usr/local/bin/cone /usr/local/bin/soft
RUN mkdir -p /git-cone/data /git-cone/hooks \
    && chown -R gitcone:gitcone /git-cone

ENV GIT_CONE_DATA_PATH=/git-cone/data

USER gitcone
WORKDIR /git-cone
VOLUME ["/git-cone/data"]
EXPOSE 23231 23232

ENTRYPOINT ["/usr/local/bin/cone"]
CMD ["serve"]
