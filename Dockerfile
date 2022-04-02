FROM golang:1.18 AS go-build
ENV GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64

WORKDIR /build
COPY . .
RUN go mod download
RUN go build -o app ./cmd

FROM alpine:3.10
COPY --from=go-build /build/app /app
COPY --from=go-build /build/views /views
EXPOSE 80
CMD /app