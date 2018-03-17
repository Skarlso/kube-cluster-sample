FROM ubuntu

LABEL Author="Gergely Brautigam"

COPY build/.env /root
COPY build/receiver /root

EXPOSE 4150

CMD [ "/root/receiver" ]
