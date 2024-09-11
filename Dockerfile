

ARG GO_VERSION=1.22
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION} AS build
# FROM golang:${GO_VERSION} AS build
WORKDIR /src


RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x

ARG TARGETARCH


RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    CGO_ENABLED=0 GOARCH=$TARGETARCH go build -o /bin/i2 .


FROM scratch AS prod

# Copy the executable from the "build" stage.
COPY --from=build /bin/i2 /bin/

# Expose the port that the application listens on.
EXPOSE 6001

# What the container should run when it is started.
ENTRYPOINT [ "/bin/i2" ]

