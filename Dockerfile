FROM golang:latest

COPY ./ ./
ENV GOPATH=/

RUN chmod +x ./scripts/database/wait-db.sh
RUN apt-get update && apt-get install -y postgresql-client
RUN go build ./cmd/image-box-app/

CMD ["./image-box-app"]