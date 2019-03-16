# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang:1.10-alpine3.8

# Copy the local package files to the container's workspace.
ADD . /go/src/twreporter.org/go-api/

COPY ./aws_credentials /root/.aws/credentials 

# Build the outyet command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
RUN cd /go/src/twreporter.org/go-api/ \
    && go build -o go-api; 

# Run the outyet command by default when the container starts.
ENTRYPOINT /go/src/twreporter.org/go-api/go-api

# Document that the service listens on port 8080.
EXPOSE 8080
