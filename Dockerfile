FROM golang:1.21rc3-alpine3.18
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY cau_parser ./cau_parser
COPY server ./server

COPY *.go ./
RUN go build -v -o ./app ./

COPY static ./static
COPY html ./html

CMD ["/app/app"]
