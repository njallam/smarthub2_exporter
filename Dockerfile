FROM golang:alpine as builder
RUN apk update && apk add --no-cache git
WORKDIR $GOPATH/src/smarthub2_exporter
COPY . .
RUN go get -d -v
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags="-s -w" -o /go/bin/smarthub2_exporter

FROM scratch
COPY --from=builder /go/bin/smarthub2_exporter /go/bin/smarthub2_exporter
EXPOSE 9906
ENTRYPOINT ["/go/bin/smarthub2_exporter"]
