# Stage 1: Build the Go application with dlv for debugging
FROM golang:alpine AS build

WORKDIR /root

COPY ./ ./

# Build application
RUN go build ./cmd/app/

# Stage 2: Create the final Docker image
FROM alpine:latest

WORKDIR /root

# Copy binary from Stage 1
COPY --from=build /root/app ./

# Copy configs from Stage 1
COPY --from=build /root/.env /root/
COPY --from=build /root/configs/main.yml /root/configs/main.yml

# Copy wait-db.sh
COPY --from=build /root/scripts/database/wait-db.sh /root/

# Install psql-client
RUN apk add --no-cache postgresql-client

# Set permissions
#RUN chmod +x /root/app
#RUN chmod +x /root/wait-db.sh