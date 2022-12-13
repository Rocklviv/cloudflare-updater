FROM --platform=${BUILDPLATFORM} golang:1.18-alpine as builder
ENV USER=appuser
ENV UID=10001
ARG BUILD_ARCH
ARG GOARM

ADD ./ /build
RUN adduser \
  --disabled-password \
  --gecos "" \
  --home "/nonexistent" \
  --shell "/sbin/nologin" \
  --no-create-home \
  --uid "${UID}" \
  "${USER}"
RUN apk --no-cache add ca-certificates

WORKDIR /build
RUN GOOS=linux GOARM=${GOARM} GOARCH=${BUILD_ARCH} go build -ldflags="-w -s" -o dns-checker .

FROM --platform=${BUILDPLATFORM} scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /build/dns-checker /opt/dns-checker
USER appuser:appuser

ENTRYPOINT ["/opt/dns-checker"]