FROM alpine:latest

WORKDIR /root

# Copy binary
COPY ./.bin/app /root/.bin/app

# Copy configs
COPY ./.env /root/
COPY ./configs/main.yml /root/configs/main.yml

# Copy wait-db.sh
COPY ./scripts/database/wait-db.sh /root/

# Copy docs
COPY ./docs /root/docs

# Install psql-client
RUN apk add --no-cache postgresql-client

CMD ["sh", "-c", "sh wait-db.sh && ./.bin/app"]