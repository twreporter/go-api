# Define global user name
ARG server_user=goapi

# Start from a Alpine Linux image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang:1.12.6-alpine3.10 As build

RUN apk add --update --no-cache \
	tzdata \
	ca-certificates \
        git \
 && update-ca-certificates

ENV GO111MODULE on

# Module cache pre-warm
WORKDIR /go/cache

COPY go.mod .
COPY go.sum .
RUN go mod download
RUN go mod verify

WORKDIR /go/src/twreporter.org/go-api

# Copy the local package files to the container's workspace.
COPY . .

ENV DOCKERIZE_VERSION v0.6.1

# Download mysql health check docker
RUN wget https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && tar -C /usr/local/bin -xzvf dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && rm dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz

# Build the migrate tool
RUN MIGRATE_PATH=$(cat go.mod | grep migrate | awk '{print $1}') && \
    go build -o /go/bin/migrate -tags 'mysql' ${MIGRATE_PATH}/cmd/migrate 
   
# Inherit global user argument
ARG server_user

# Add the user for running go-api
RUN adduser -D -g '' ${server_user}

# Install
RUN go install

# Minimize image size by only using the required binary
FROM alpine:3.9

COPY --from=build /go/bin /usr/local/bin /usr/local/bin/
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build /etc/passwd /etc/passwd

ARG server_user

WORKDIR /home/${server_user}

COPY ./aws_credentials /home/${server_user}/.aws/credentials 
COPY ./entrypoint.sh /home/${server_user}/entrypoint.sh
COPY ./migrations /home/${server_user}/migrations/
COPY ./template /home/${server_user}/template/
RUN chmod +x entrypoint.sh

ENV MIGRATION_DIR /home/${server_user}/migrations/

# Specify the user for running go-api
USER ${server_user}

# Run the outyet command by default when the container starts.
ENTRYPOINT ["./entrypoint.sh"]

CMD ["go-api"]

# Document that the service listens on port 8080.
EXPOSE 8080
