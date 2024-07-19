FROM golang:1.22.5-alpine

RUN apk add --no-cache curl git bash ca-certificates \
    && apk add --no-cache --virtual .build-deps gcc musl-dev \
    && apk add --no-cache --virtual .build-deps g++ make 


WORKDIR /app

COPY ../go.mod ../go.sum ./

RUN go mod download

COPY ../ .

RUN go build -o /find-home .

COPY ../docker/crontab /etc/crontabs/root

CMD ["sh", "-c", "crond -f -L /var/log/cron.log"]

