FROM golang:1.24-alpine3.17 AS build

COPY conf/ /go/src/conf/
COPY internal/ /go/src/internal/
COPY pkg/ /go/src/pkg/
COPY go.mod go.sum *.go /go/src/

WORKDIR /go/src/
RUN go env -w GO111MODULE=on \
  && go env -w GOPROXY=https://goproxy.cn,direct \
  && go env -w GOOS=linux \
  && go env -w GOARCH=amd64
RUN go mod tidy
RUN go build -o offer_save

FROM alpine:3.17

COPY --from=build /go/src/offer_save /app/offer_save

RUN chmod +x /app/offer_save

EXPOSE 8080
ENTRYPOINT [ "/app/offer_save" ]