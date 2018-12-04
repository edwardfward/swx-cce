# reference https://hub.docker.com/_/golang/
# reference https://docs.docker.com/engine/reference/builder/

# use the latest version of golang
FROM golang:stretch

ENV GOBIN="/go/src/app"
ENV DB_URL=""
ENV DB_USERNAME:""
ENV DB_PASSWORD:""
ENV DB_PORT:""

WORKDIR /go/src/app

RUN apt-get update -y && apt-get upgrade -y && apt-get install ca-certificates git
RUN git clone --depth=1 --branch=master https://github.com/edwardfward/swx-cce .

RUN go install -v cmd/cce-server.go

CMD ["cce-server"]



