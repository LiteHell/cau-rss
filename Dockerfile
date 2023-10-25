FROM golang:alpine as base
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

FROM build as deployment
CMD ["/app/app"]

FROM base as test
ARG REDIS_ENABLED
ARG REDIS_ADDR
ENV REDIS_ENABLED=${REDIS_ENABLED:-false}
ENV REDIS_ADDR=${REDIS_ADDR:-127.0.0.1\:6379}
RUN echo REDIS_ENABLED: $REDIS_ENABLED
RUN echo REDIS_ADDR: $REDIS_ADDR
RUN ["go", "test" ,"-v", "./..."]

