# reference https://hub.docker.com/_/golang/
# reference https://docs.docker.com/engine/reference/builder/

# use the latest version of golang
FROM golang:stretch

ARG APP_REPO=""

ENV GOBIN="/go/src/app"
ENV APP_ENV_REPO=https://${APP_REPO}
ENV DB_URL=""
ENV DB_USERNAME=""
ENV DB_PASSWORD=""
ENV DB_PORT=""

RUN git clone --depth=1 --branch=master ${APP_ENV_REPO} /go/src/app
WORKDIR /go/src/app
RUN ls -al cmd/

RUN go install -v cmd/cceapp.go

CMD ["./cceapp"]



