FROM golang:1.13-alpine
RUN apk add -U git mercurial openssh ca-certificates gcc musl-dev
RUN go get -u github.com/golangci/golangci-lint/cmd/golangci-lint

CMD ["go", "run", "main.go"]