FROM golang:1.10.0-alpine3.7

ADD version /root/version
RUN apk --no-cache add git make vim
RUN go get -d -v github.com/Skarlso/kube-cluster-sample/...
WORKDIR /go/src/github.com/Skarlso/kube-cluster-sample/frontend
RUN make

FROM alpine:latest
LABEL Author="Gergely Brautigam"
ARG port=8081
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /go/src/github.com/Skarlso/kube-cluster-sample/frontend/build/frontend .
COPY --from=0 /go/src/github.com/Skarlso/kube-cluster-sample/frontend/index.html .

EXPOSE ${port}
CMD [ "/root/frontend" ]
