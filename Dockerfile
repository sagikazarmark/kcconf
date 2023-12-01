FROM --platform=$BUILDPLATFORM golang:1.21.3-alpine3.18@sha256:926f7f7e1ab8509b4e91d5ec6d5916ebb45155b0c8920291ba9f361d65385806 AS builder

RUN apk add --update --no-cache ca-certificates git

ARG TARGETOS
ARG TARGETARCH
ARG TARGETPLATFORM

WORKDIR /usr/local/src/kcconf

ARG GOPROXY

ENV CGO_ENABLED=0
ENV GOOS=$TARGETOS GOARCH=$TARGETARCH

COPY go.* ./
RUN go mod download

COPY . .

RUN go build -o /usr/local/bin/kcconf .


FROM alpine:3.18.5@sha256:34871e7290500828b39e22294660bee86d966bc0017544e848dd9a255cdf59e0

RUN apk add --update --no-cache ca-certificates tzdata

COPY --from=builder /usr/local/bin/kcconf /usr/local/bin/

EXPOSE 8080

CMD kcconf
