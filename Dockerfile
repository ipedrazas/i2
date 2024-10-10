# syntax=docker/dockerfile:1

ARG GO_VERSION=1.23
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION} AS build
WORKDIR /src

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x

ARG TARGETARCH

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    CGO_ENABLED=0 GOARCH=$TARGETARCH go build -o /bin/server .

FROM alpine:latest AS final


RUN --mount=type=cache,target=/var/cache/apk \
    apk --update add \
    ca-certificates \
    tzdata \
    openssh-client \
    && \
    update-ca-certificates


ARG UID=1000
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/home/appuser" \
    --uid "${UID}" \
    appuser
USER appuser

# Copy the executable from the "build" stage.
COPY --from=build /bin/server /bin/i2
# COPY --chown=appuser:appuser --chmod=0600 ./docker/ssh-config /home/appuser/.ssh/config
# COPY --chown=appuser:appuser --chmod=0600 ./docker/id_rsa /home/appuser/.ssh/id_rsa
# Expose the port that the application listens on.
EXPOSE 6001

# What the container should run when it is started.
ENTRYPOINT [ "/bin/i2" ]
