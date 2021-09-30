FROM golang:1.13.0-alpine3.10

RUN apk --no-cache add git make vim
RUN mkdir -p /data
#RUN git clone https://github.com/Skarlso/kube-cluster-sample /data
WORKDIR /data
COPY . .
RUN make

FROM alpine:latest
LABEL Author="Gergely Brautigam"
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /data/redeemer/build/redeemer .
CMD [ "/root/redeemer" ]
