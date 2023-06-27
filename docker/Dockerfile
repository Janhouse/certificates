FROM golang:alpine AS builder

# Build Delve
#RUN go install github.com/go-delve/delve/cmd/dlv@latest

WORKDIR /src
COPY . .

RUN apk add --no-cache curl git make libcap
RUN make V=1 bin/step-ca
RUN setcap CAP_NET_BIND_SERVICE=+eip bin/step-ca

FROM smallstep/step-kms-plugin:cloud AS kms

FROM smallstep/step-cli:latest

COPY --from=builder /src/bin/step-ca /usr/local/bin/step-ca
COPY --from=kms /usr/local/bin/step-kms-plugin /usr/local/bin/step-kms-plugin
#COPY --from=builder /go/bin/dlv /
USER step

ENV CONFIGPATH="/home/step/config/ca.json"
ENV PWDPATH="/home/step/secrets/password"

VOLUME ["/home/step"]
STOPSIGNAL SIGTERM
HEALTHCHECK CMD step ca health 2>/dev/null | grep "^ok" >/dev/null

COPY docker/entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/bin/bash", "/entrypoint.sh"]
CMD exec /usr/local/bin/step-ca --password-file $PWDPATH $CONFIGPATH
#CMD exec /dlv --listen=:40000 --headless=true --api-version=2 --accept-multiclient exec /usr/local/bin/step-ca -- --password-file $PWDPATH $CONFIGPATH