FROM golang:latest

COPY ./ ./
ENV GOPATH=/

RUN go build ./cmd/image-box-app/

CMD ["./image-box-app"]