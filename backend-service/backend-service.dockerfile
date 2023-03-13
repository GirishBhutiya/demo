#FROM alpine:latest 
#FROM golang:1.20-alpine AS build
#RUN apk add libc6-compat

#WORKDIR /app

#COPY backendApp .

#CMD [ "/app/backendApp" ]



#Build stage
FROM golang:1.20.1-alpine3.17 AS builder
WORKDIR /app
COPY backendApp .
#RUN go build -o main main.go

#Run stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/backendApp ./

EXPOSE 8080
CMD [ "/app/backendApp" ]