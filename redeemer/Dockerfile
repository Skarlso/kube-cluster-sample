FROM golang:1.13.0-alpine3.10

RUN apk --no-cache add git make vim
RUN go get -d -v github.com/Skarlso/kube-cluster-sample/...
WORKDIR /go/src/github.com/Skarlso/kube-cluster-sample/redeemer
RUN make

FROM alpine:latest
LABEL Author="Gergely Brautigam"
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /go/src/github.com/Skarlso/kube-cluster-sample/redeemer/build/redeemer .
CMD [ "/root/redeemer" ]
