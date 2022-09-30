FROM golang:alpine3.14 as base
RUN apk update && apk add bash inotify-tools && apk add git

WORKDIR /app

COPY ./ /app

FROM base as dev

ENV CGO_ENABLED 0
COPY startScript.sh /build/startScript.sh

RUN git clone https://github.com/go-delve/delve.git && \
    cd delve && \
    go install github.com/go-delve/delve/cmd/dlv

RUN go mod tidy 

RUN go build -o /server -gcflags -N -gcflags -l

EXPOSE 8079
EXPOSE 40000

ENTRYPOINT sh startScript.sh

FROM base as prod

RUN go mod tidy
RUN go build -o /server
EXPOSE 8079

CMD ["/server"]
