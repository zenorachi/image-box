FROM golang:latest

COPY ./ ./

RUN go build ./cmd/image-box-app/

CMD ["./image-box-app"]