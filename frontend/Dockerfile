FROM golang:1.13.0-alpine3.10

ADD version /root/version
RUN apk --no-cache add git make vim
RUN mkdir -p /frontend
WORKDIR /frontend
COPY . .
RUN make

FROM alpine:latest
LABEL Author="Gergely Brautigam"
ARG port=8081
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /frontend/build/frontend .
COPY --from=0 /frontend/index.html .

EXPOSE ${port}
CMD [ "/root/frontend" ]
