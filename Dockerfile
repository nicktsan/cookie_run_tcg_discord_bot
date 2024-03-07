FROM golang:1.22.1

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build

RUN go build -o /usr/src/app/mtbot

ENTRYPOINT ["/usr/src/app/mtbot"]