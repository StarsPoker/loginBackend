FROM golang:latest

WORKDIR /app

COPY ./ /app

RUN go mod tidy

RUN go get github.com/githubnemo/CompileDaemon

EXPOSE 8079

ENTRYPOINT CompileDaemon --build="go build main.go" --command=./main