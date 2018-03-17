FROM ubuntu

LABEL Author="Gergely Brautigam"

COPY build/.env /root
COPY build/frontend /root

EXPOSE 8081

CMD [ "/root/frontend" ]
