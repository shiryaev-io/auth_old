FROM golang:1.18.1-alpine

ENV APP_PATH=/app
ENV SERVICE_PATH=${APP_PATH}/cmd/service/main

RUN mkdir ${APP_PATH}
WORKDIR ${APP_PATH}

COPY . ./

RUN go mod download
WORKDIR ${SERVICE_PATH}
RUN go build -o server

EXPOSE ${SERVICE_PORT}

ENTRYPOINT [ "./server" ]