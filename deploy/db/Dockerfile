FROM alpine

WORKDIR /app

RUN apk add --update --no-cache postgresql-client
RUN apk add --update --no-cache curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz && \
    chmod +x ./migrate


COPY . . 


CMD ["sh", "./deploy/db/migrations.sh"]

