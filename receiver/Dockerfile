FROM golang:1.10.0-alpine3.7

RUN apk --no-cache add git make vim
ADD version /root/version
RUN go get -d -v github.com/Skarlso/kube-cluster-sample/...
WORKDIR /go/src/github.com/Skarlso/kube-cluster-sample/receiver
RUN make

FROM alpine:latest
LABEL Author="Gergely Brautigam"
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /go/src/github.com/Skarlso/kube-cluster-sample/receiver/build/receiver .
EXPOSE 8000

CMD [ "/root/receiver" ]
