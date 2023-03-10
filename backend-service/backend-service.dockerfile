FROM alpine:latest 

RUN mkdir /app

WORKDIR /app

COPY backendApp /app

CMD [ "/app/backendApp" ]