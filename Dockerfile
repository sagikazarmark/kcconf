FROM --platform=$BUILDPLATFORM golang:1.21.1-alpine3.18@sha256:96634e55b363cb93d39f78fb18aa64abc7f96d372c176660d7b8b6118939d97b AS builder

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


FROM alpine:3.18.4@sha256:eece025e432126ce23f223450a0326fbebde39cdf496a85d8c016293fc851978

RUN apk add --update --no-cache ca-certificates tzdata

COPY --from=builder /usr/local/bin/kcconf /usr/local/bin/

EXPOSE 8080

CMD kcconf
