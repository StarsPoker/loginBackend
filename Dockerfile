FROM golang:latest

WORKDIR /app

COPY ./ /app

RUN go mod tidy \
    && go get github.com/githubnemo/CompileDaemon \
    && go install github.com/githubnemo/CompileDaemon

EXPOSE 8079

ENTRYPOINT CompileDaemon --build="go build main.go" --command=./main