FROM --platform=$BUILDPLATFORM cgr.dev/chainguard/wolfi-base AS builder
ARG VERSION=1.24
ARG APP_VERSION=edge
ARG TARGETOS TARGETARCH
ENV LANG=C.UTF-8
RUN apk update && \
  apk upgrade --purge --no-cache && \
  apk add --no-cache --purge --no-interactive \
  go-${VERSION} \
  ca-certificates && \
  update-ca-certificates

WORKDIR /app
COPY go.mod go.sum /app/
RUN go mod download

COPY . .
ENV GOOS=$TARGETOS
ENV GOARCH=$TARGETARCH
ENV CGO_ENABLED=0
ENV GOROOT=/usr/lib/go
RUN go build -o issue-syncer \
  -ldflags "-X github.com/alexdor/issue-syncer/cmd.root=${APP_VERSION}}" \
  -trimpath

FROM --platform=$BUILDPLATFORM cgr.dev/chainguard/wolfi-base
ENV LANG=C.UTF-8
COPY --from=builder /usr/share/ca-certificates/ /usr/share/ca-certificates/
COPY --from=builder /app/issue-syncer /usr/local/bin/issue-syncer

ENTRYPOINT ["issue-syncer"]
