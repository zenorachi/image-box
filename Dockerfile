# Stage 1: Build the Go application with dlv for debugging
FROM golang:alpine AS build

WORKDIR /app

COPY ./ ./

# Build application
RUN go build ./cmd/image-box-app/

# Stage 2: Create the final Docker image
FROM alpine:latest

WORKDIR /app

# Copy binary from Stage 1
COPY --from=build /app/image-box-app /app/

# Copy configs from Stage 1
COPY --from=build /app/.env /app/
COPY --from=build /app/configs/main.yml /app/configs/main.yml

# Copy wait-db.sh
COPY --from=build /app/scripts/database/wait-db.sh /app/

# Install psql-client
RUN apk add --no-cache postgresql-client

# Set permissions
#RUN chmod +x /app/image-box-app
#RUN chmod +x /app/wait-db.sh