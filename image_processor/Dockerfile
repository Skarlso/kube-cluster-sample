FROM alpine
LABEL Author="Gergely Brautigam"
RUN apk add -u ca-certificates
COPY ./build/linux/amd64/processor /app/

WORKDIR /app/
ENTRYPOINT [ "/app/processor" ]
