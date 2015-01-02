# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

# Build the outyet command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
RUN go get github.com/mitchellh/goamz/...
RUN go get github.com/zerklabs/auburn/log
RUN go install github.com/zerklabs/up53

# Run the outyet command by default when the container starts.
ENTRYPOINT /go/bin/up53
