FROM golang:1.17
WORKDIR /go/src/github.com/chrj/placetaytay/
COPY ./ /go/src/github.com/chrj/placetaytay/
RUN CGO_ENABLED=0 go install .

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /go/bin/placetaytay ./
CMD ["./placetaytay"]  