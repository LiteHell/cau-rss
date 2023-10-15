FROM golang:1.21rc3-alpine3.18 as base
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY cau_parser ./cau_parser
COPY server ./server

# To avoid tls error from swedu.cau.ac.kr
COPY swedu-cert.pem /usr/local/share/ca-certificates/swedu-cert.crt
RUN cat /usr/local/share/ca-certificates/swedu-cert.crt >> /etc/ssl/certs/ca-certificates.crt

COPY static ./static
COPY html ./html

COPY *.go ./

FROM base as build
RUN go build -v -o ./app ./

FROM build as test
# To check whether build is done without error, testing is performed after build stage
RUN ["go", "test" ,"-v", "./..."]

FROM build as deployment
CMD ["/app/app"]