FROM golang:1-alpine

WORKDIR /go/src/app

COPY go.mod go.sum ./

RUN go mod download

COPY . .
RUN go install -v ./...

CMD ["container-manager"]