# Stage 1: Build the Go application with dlv for debugging
FROM golang:1.17.3-alpine3.14 AS build

WORKDIR /app

COPY ./ ./

# Install dependencies and build application
RUN go install github.com/go-delve/delve/cmd/dlv@latest
RUN go build -gcflags="-N -l" ./cmd/image-box-app/

# Stage 2: Create the final Docker image
FROM alpine:latest

WORKDIR /app

# Copy files from Stage 1
COPY --from=build /app/image-box-app /app/
COPY --from=build /app/scripts/database/wait-db.sh /app/
COPY --from=build /go/bin/dlv /usr/local/bin/

# Install psql-client
RUN apk add --no-cache postgresql-client

ENV PATH="/usr/local/bin:${PATH}"

CMD ["dlv", "exec", "./image-box-app"]